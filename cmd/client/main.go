package main

import (
	"flag"
	"github.com/ungame/timetrack/app/cli"
)

var baseURL string

func init() {
	flag.StringVar(&baseURL, "base_url", "http://localhost:15555", "set base url server")
	flag.Parse()
}

func main() {
	c := cli.New(baseURL)
	c.Run()
}
