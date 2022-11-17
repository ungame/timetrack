package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func New() *sql.DB {
	cfg := getDefaultConfig()
	conn, err := sql.Open("mysql", cfg.Source())
	if err != nil {
		log.Panicln(err)
	}
	conn.SetConnMaxLifetime(time.Minute * 5)
	conn.SetMaxIdleConns(25)
	conn.SetMaxOpenConns(25)
	return conn
}
