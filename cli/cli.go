package cli

import (
	"flag"
	"fmt"
	"os"
)

var baseURL string

func init() {
	flag.StringVar(&baseURL, "base_url", "http://localhost:15555", "set base url server")
	flag.Parse()
}

type CommandLine struct {
}

func (c *CommandLine) Run() {

	if len(os.Args) < 2 {
		c.Usage()
		os.Exit(0)
	}

	var (
		listName    string
		description string
		category    int
	)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCmd.StringVar(&listName, "n", "", "name list")

	startActivityCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startActivityCmd.StringVar(&description, "d", "", "set description")
	startActivityCmd.IntVar(&category, "c", 0, "set category id")

	switch os.Args[1] {
	case "list":
		if err := listCmd.Parse(os.Args[2:]); err != nil {
			c.Usage()
			return
		}

	case "start":
		if err := startActivityCmd.Parse(os.Args[2:]); err != nil {
			c.Usage()
			return
		}

	default:
		c.Usage()
		os.Exit(0)
	}

	if listCmd.Parsed() {
		c.List(listName)
	}

	if startActivityCmd.Parsed() {
		CreateActivity(description, int64(category))
	}
}

func (c *CommandLine) Usage() {
	fmt.Println("-- Usage:")
	fmt.Println(" ")
	fmt.Println("     list  -n LIST_NAME")
	fmt.Println("     start -d DESCRIPTION -c CATEGORY_ID")
}

func (c *CommandLine) List(name string) {
	switch name {
	case "categories":
		ListCategories()
	case "activities":
		ListActivities()
	default:
		c.Usage()
	}
}
