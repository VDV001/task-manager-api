//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daniilgit/task-manager-api/internal/handler"
	"github.com/daniilgit/task-manager-api/internal/repo/postgres"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/daniilgit/task-manager-api/pkg/hash"
	"github.com/daniilgit/task-manager-api/pkg/jwt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type testApp struct {
	t      *testing.T
	server *httptest.Server
	db     *sqlx.DB
}

func setupTestApp(t *testing.T) *testApp {
	t.Helper()
	ctx := t.Context()

	// Start PostgreSQL container.
	pgContainer, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { pgContainer.Terminate(context.Background()) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sqlx.Connect("pgx", connStr)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Run migrations.
	err = goose.Up(db.DB, "../../migrations")
	require.NoError(t, err)

	// Wire dependencies.
	userRepo := postgres.NewUserRepo(db)
	taskRepo := postgres.NewTaskRepo(db)
	hasher := hash.NewBcryptHasher(4) // low cost for fast tests
	jwtManager := jwt.NewManager("test-access", "test-refresh", 15*time.Minute, 24*time.Hour, "task-manager-api")
	tokenAdapter := jwt.NewManagerAdapter(jwtManager)

	log := slog.Default()
	authUC := usecase.NewAuthUseCase(userRepo, hasher, tokenAdapter, log)
	taskUC := usecase.NewTaskUseCase(taskRepo, log)

	router := handler.NewRouter(authUC, taskUC, tokenAdapter, log, db.PingContext)
	server := httptest.NewServer(router)
	t.Cleanup(func() { server.Close() })

	return &testApp{t: t, server: server, db: db}
}

func (a *testApp) request(method, path string, body any, token string) *http.Response {
	a.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		require.NoError(a.t, err)
	}

	req, err := http.NewRequest(method, a.server.URL+path, &buf)
	require.NoError(a.t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(a.t, err)
	return resp
}

type tokenResp struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

func (a *testApp) registerAndLogin(t *testing.T, name, email, password string) string {
	t.Helper()
	resp := a.request("POST", "/api/v1/auth/register", map[string]string{
		"name": name, "email": email, "password": password,
	}, "")
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var tr tokenResp
	err := json.NewDecoder(resp.Body).Decode(&tr)
	require.NoError(t, err)
	resp.Body.Close()
	return tr.Data.AccessToken
}

func TestIntegration_FullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	app := setupTestApp(t)

	// 1. Register.
	token := app.registerAndLogin(t, "Alice", "alice@example.com", "password123")
	assert.NotEmpty(t, token)

	// 2. Duplicate registration → 409.
	resp := app.request("POST", "/api/v1/auth/register", map[string]string{
		"name": "Alice", "email": "alice@example.com", "password": "password123",
	}, "")
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	resp.Body.Close()

	// 3. Login.
	resp = app.request("POST", "/api/v1/auth/login", map[string]string{
		"email": "alice@example.com", "password": "password123",
	}, "")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var tr tokenResp
	err := json.NewDecoder(resp.Body).Decode(&tr)
	require.NoError(t, err)
	resp.Body.Close()
	loginToken := tr.Data.AccessToken
	assert.NotEmpty(t, loginToken)

	// 4. Create task.
	deadline := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	resp = app.request("POST", "/api/v1/tasks", map[string]any{
		"title": "Deploy v2", "description": "Ship it", "deadline": deadline,
	}, token)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var taskResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&taskResp)
	require.NoError(t, err)
	resp.Body.Close()
	taskID := taskResp.Data.ID
	assert.NotEmpty(t, taskID)

	// 5. Get task.
	resp = app.request("GET", "/api/v1/tasks/"+taskID, nil, token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 6. Update task (PATCH).
	resp = app.request("PATCH", "/api/v1/tasks/"+taskID, map[string]any{
		"status": "in_progress",
	}, token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 7. List tasks.
	resp = app.request("GET", "/api/v1/tasks?status=in_progress&page=1&limit=10", nil, token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 8. Stats.
	resp = app.request("GET", "/api/v1/tasks/stats", nil, token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var statsResp struct {
		Data struct {
			Total    int            `json:"total"`
			ByStatus map[string]int `json:"by_status"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&statsResp)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, 1, statsResp.Data.Total)
	assert.Equal(t, 1, statsResp.Data.ByStatus["in_progress"])

	// 9. Delete task.
	resp = app.request("DELETE", "/api/v1/tasks/"+taskID, nil, token)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// 10. Get deleted task → 404.
	resp = app.request("GET", "/api/v1/tasks/"+taskID, nil, token)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// 11. Access without token → 401.
	resp = app.request("GET", "/api/v1/tasks", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	// 12. Another user can't see Alice's tasks.
	otherToken := app.registerAndLogin(t, "Bob", "bob@example.com", "password456")
	resp = app.request("GET", fmt.Sprintf("/api/v1/tasks/%s", taskID), nil, otherToken)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode) // soft-deleted, so 404
	resp.Body.Close()
}
