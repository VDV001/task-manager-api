package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	authUC *usecase.AuthUseCase,
	taskUC *usecase.TaskUseCase,
	tokens usecase.TokenManager,
	log *slog.Logger,
	dbPing func(ctx context.Context) error,
	corsOrigins []string,
	buildInfo map[string]string,
) http.Handler {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Request-Id", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(c.Handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(StructuredLogger(log))
	r.Use(middleware.Recoverer)
	r.Use(maxBytesReader(1 << 20)) // 1 MB

	validate := validator.New()

	authHandler := NewAuthHandler(authUC, validate)
	taskHandler := NewTaskHandler(taskUC, validate)

	authLimiter := NewRateLimiter(1, 5) // 1 req/sec, burst of 5

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes.
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Use(authLimiter.Middleware)
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
		})

		// Protected routes.
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(tokens))

			r.Route("/tasks", func(r chi.Router) {
				r.With(middleware.AllowContentType("application/json")).Post("/", taskHandler.Create)
				r.Get("/", taskHandler.List)
				r.Get("/stats", taskHandler.Stats)
				r.Get("/{id}", taskHandler.GetByID)
				r.With(middleware.AllowContentType("application/json")).Patch("/{id}", taskHandler.Update)
				r.Delete("/{id}", taskHandler.Delete)
			})
		})
	})

	// Swagger UI.
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Version info.
	r.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(buildInfo)
	})

	// Liveness — always OK if process is running.
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Readiness — checks database connectivity.
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := dbPing(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}

// maxBytesReader limits the size of incoming request bodies.
func maxBytesReader(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
