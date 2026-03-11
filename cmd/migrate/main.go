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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer db.Close()

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose set dialect: %v", err)
	}

	if err := goose.RunContext(context.Background(), command, db.DB, "migrations"); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}

	fmt.Printf("goose %s: done\n", command)
}
