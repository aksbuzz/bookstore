package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func New(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn is empty")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Instance() *sql.DB {
	return db.db
}

func (db *DB) Close() error {
	return db.db.Close()
}
