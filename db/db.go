package db

import (
	"database/sql"
	"fmt"

	"github.com/dnc-data-mcp/config"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
