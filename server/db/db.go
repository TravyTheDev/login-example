package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func NewSqlStorage() (*sql.DB, error) {
	dbUrl := os.Getenv("DB_URL")
	dbType := "sqlite"
	db, err := sql.Open(dbType, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return db, nil
}
