// @title           Task Manager API
// @version         1.0
// @description     REST API для менеджера задач с аутентификацией и управлением задачами.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer {token}"
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daniilgit/task-manager-api/internal/config"
	"github.com/daniilgit/task-manager-api/internal/handler"
	"github.com/daniilgit/task-manager-api/internal/repo/postgres"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/daniilgit/task-manager-api/pkg/hash"
	"github.com/daniilgit/task-manager-api/pkg/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	// Config.
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Logger.
	logLevel := new(slog.LevelVar)
	switch cfg.LogLevel {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelInfo)
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(log)

	// Database.
	db, err := sqlx.Connect("pgx", cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Info("connected to database")

	// Dependencies.
	userRepo := postgres.NewUserRepo(db)
	taskRepo := postgres.NewTaskRepo(db)

	hasher := hash.NewBcryptHasher(12)

	jwtManager := jwt.NewManager(
		cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL,
		"task-manager-api",
	)
	tokenAdapter := jwt.NewManagerAdapter(jwtManager)

	authUC := usecase.NewAuthUseCase(userRepo, hasher, tokenAdapter, log)
	taskUC := usecase.NewTaskUseCase(taskRepo, log)

	// HTTP server.
	router := handler.NewRouter(authUC, taskUC, tokenAdapter, log, db.PingContext)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("server started", slog.Int("port", cfg.Server.Port))
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	case <-ctx.Done():
		stop()
		log.Info("shutting down server...")
		srv.SetKeepAlivesEnabled(false)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		log.Info("server stopped gracefully")
	}

	return nil
}
