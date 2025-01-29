package store

import (
	"database/sql"
	"fmt"

	"maas/internal/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// NewDB creates a new database connection.
func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
