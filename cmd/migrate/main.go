package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/daniilgit/task-manager-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}

	if err := goose.RunContext(context.Background(), command, db.DB, "migrations"); err != nil {
		return fmt.Errorf("goose %s: %w", command, err)
	}

	fmt.Printf("goose %s: done\n", command)
	return nil
}
