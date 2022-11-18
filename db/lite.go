package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Lite(fileStorage FileStorage, migrations ...Migration) *sql.DB {
	conn, err := sql.Open("sqlite3", fileStorage.Create())
	if err != nil {
		log.Panicln(err.Error())
	}

	for _, migration := range migrations {
		migration.Up(conn)
	}

	return conn
}
