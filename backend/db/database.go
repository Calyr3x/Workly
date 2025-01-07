package db

import (
	"database/sql"
	"errors"
	"fmt"
	"workly/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg config.DBConfig) error {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return errors.New("Failed to connect to database:" + err.Error())
	}

	return nil
}
