package queries

import (
	"context"
	"database/sql"
	"log"
)

func MustPrepare(ctx context.Context, conn *sql.DB, query string) *sql.Stmt {
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Panicln("Error on prepare query:", err.Error())
	}

	return stmt
}
