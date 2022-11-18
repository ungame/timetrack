package db

import (
	"database/sql"
	"fmt"
	"github.com/ungame/timetrack/ioext"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const sqlFile = "infra/timetrack/sql/init.sql"

type Migration interface {
	Up(conn *sql.DB)
}

type migration struct {
	name  string
	query string
}

func NewMigration(name, query string) Migration {
	return &migration{
		name:  name,
		query: query,
	}
}

func (m *migration) Up(conn *sql.DB) {
	m.createTable(conn)

	if m.hasMigration(conn) {
		return
	}

	m.migrate(conn)
}

func (m *migration) createTable(conn *sql.DB) {
	query := "CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY)"
	if _, err := conn.Exec(query); err != nil {
		log.Panicln(err.Error())
	}
}

func (m *migration) hasMigration(conn *sql.DB) bool {
	query := "SELECT COUNT(*) FROM migrations WHERE name='" + m.name + "'"
	var count int64
	err := conn.QueryRow(query).Scan(&count)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Panicln(err)
	}
	return count > 0
}

func (m *migration) migrate(conn *sql.DB) {
	log.Println("initializing transaction...")
	tx, err := conn.Begin()
	result, err := tx.Exec(m.query)
	if err != nil {
		log.Panicln(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("error on get affected rows:", err.Error())
	} else {
		log.Println("query executed. rows affected:", rows)
	}

	query := "INSERT INTO migrations (name) VALUES ('" + m.name + "')"
	result, err = tx.Exec(query)
	if err != nil {
		log.Println("error on exec second query:", err.Error())
		if err = tx.Rollback(); err != nil {
			log.Panicln(err)
		} else {
			log.Println("rollback successfully!")
		}
	}

	rows, err = result.RowsAffected()
	if err != nil {
		log.Println("error on get affected rows:", err.Error())
	} else {
		log.Println("migrations updated. rows affected:", rows)
	}

	if err = tx.Commit(); err != nil {
		log.Panicln(err.Error())
	}

	log.Println("transaction completed successfully.")
}

func GetSqliteSeed() (string, string) {
	var (
		rootPath = filepath.Join(getCurrentPath(), "../")
		fullPath = rootPath + "/" + sqlFile
	)

	file, err := os.Open(fullPath)
	if err != nil {
		log.Panicln(err)
	}
	defer ioext.Close(file)

	content, err := io.ReadAll(file)
	if err != nil {
		log.Panicln(err)
	}

	return fullPath, convertToSqLiteQuery(string(content))
}

func convertToSqLiteQuery(query string) string {
	query = strings.ReplaceAll(query, "BIGINT PRIMARY KEY", "INTEGER PRIMARY KEY")
	query = strings.ReplaceAll(query, "INT PRIMARY KEY", "INTEGER PRIMARY KEY")
	query = strings.ReplaceAll(query, "AUTO_INCREMENT", "AUTOINCREMENT")
	query = strings.ReplaceAll(query, "ENGINE = INNODB", "")
	query = strings.ReplaceAll(query, "DEFAULT CHARSET = UTF8", "")
	fmt.Println(query)
	return query
}
