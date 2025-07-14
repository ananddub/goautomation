package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"goautomation/db/sqlc"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB      *sql.DB
	Queries *sqlc.Queries
	Context context.Context
}

// NewDatabase creates a new database connection and applies migrations
func NewDatabase(dbPath string) (*Database, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := sqlc.New(db)

	return &Database{
		DB:      db,
		Queries: queries,
		Context: context.Background(),
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
