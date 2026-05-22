package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL DEFAULT '',
			country TEXT NOT NULL DEFAULT ''
		);
	`)
	if err != nil {
		return fmt.Errorf("creating authors table issues: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL DEFAULT '',
			isbn CHAR(20) UNIQUE NOT NULL DEFAULT '',
			year INTEGER NOT NULL DEFAULT 0,
			authors_id INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (authors_id) REFERENCES authors(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("creating books table issues: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS readers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL DEFAULT '',
			email TEXT UNIQUE NOT NULL DEFAULT '',
			phone TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("creating readers table issues: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS borrowing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id INTEGER NOT NULL DEFAULT 0,
			readers_id INTEGER NOT NULL DEFAULT 0,
			borrow_date TEXT NOT NULL DEFAULT '',
			return_date TEXT,
			FOREIGN KEY (book_id) REFERENCES books(id),
			FOREIGN KEY (readers_id) REFERENCES readers(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("creating borrowing table issues: %w", err)
	}
	return nil
}

func InitDb(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open connection issues: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection issues: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("creating tables issues: %w", err)
	}

	log.Println("DB is initialized")
	return db, nil
}
