package main

import (
	"github.com/ungame/timetrack/it/internal"
	"time"
)

func main() {
	internal.Load(time.Now().AddDate(0, -2, 1))
}
