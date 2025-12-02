// Package database holds all the database logic for
// create a new driver connection and insert webhooks
package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{conn: db}

	if err := database.init(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *Database) init() error {
	query := `
	CREATE TABLE IF NOT EXISTS webhooks (
		id TEXT PRIMARY KEY,
		endpoint TEXT NOT NULL,
		method TEXT NOT NULL,
		headers TEXT NOT NULL,
		body TEXT,
		ip TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_endpoint ON webhooks(endpoint);
	CREATE INDEX IF NOT EXISTS idx_created_at ON webhooks(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_method ON webhooks(method);
	`

	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) Exec(query string, args ...any) error {
	_, err := db.conn.Exec(query, args...)
	return err
}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}
