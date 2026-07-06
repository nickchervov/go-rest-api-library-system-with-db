package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func InitDb(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open connection issues: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection issues: %w", err)
	}
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("create db instance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://internal/db/migrations", "sqlite", driver)
	if err != nil {
		return nil, fmt.Errorf("create migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("up migration: %w", err)
	}

	log.Println("DB is initialized with migrations")
	return db, nil
}
