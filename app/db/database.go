package db

import (
	"database/sql"
	"fmt"
	"log"

	"workly/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg config.DBConfig) {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
}
