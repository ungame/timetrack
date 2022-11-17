package db

import (
	"context"
	"github.com/ungame/timetrack/ioext"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	conn := New()
	defer ioext.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := conn.PingContext(ctx)
	if err != nil {
		t.Error("unexpected error on to ping database:", err.Error())
	}
}
