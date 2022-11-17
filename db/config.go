package db

import "fmt"

const (
	defaultUser     = "root"
	defaultPassword = "root"
	defaultHost     = "localhost"
	defaultPort     = 3306
	defaultDatabase = "timetrack"
)

type config struct {
	user   string
	pass   string
	host   string
	port   int
	dbName string
}

func (c *config) Source() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.user, c.pass, c.host, c.port, c.dbName)
}

func getDefaultConfig() *config {
	return &config{
		user:   defaultUser,
		pass:   defaultPassword,
		host:   defaultHost,
		port:   defaultPort,
		dbName: defaultDatabase,
	}
}
