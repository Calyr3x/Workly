package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func InitDB() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return errors.New("DB_URL is not set")
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return errors.New("Failed to connect to database: " + err.Error())
	}

	return nil
}
