package cli

import (
	"flag"
	"fmt"
	"github.com/ungame/timetrack/types"
	"os"
)

type CommandLine struct {
	baseURL string
}

func New(baseURL string) *CommandLine {
	return &CommandLine{baseURL: baseURL}
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
		period      string
		order       string
		limit       int
		activityID  int64
	)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCmd.StringVar(&listName, "n", "", "name list")
	listCmd.StringVar(&period, "p", "", "list by period")
	listCmd.StringVar(&order, "o", "", "order items")
	listCmd.IntVar(&limit, "l", 0, "limit items")

	startActivityCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startActivityCmd.StringVar(&description, "d", "", "set description")
	startActivityCmd.IntVar(&category, "c", 0, "set category id")

	finishActivityCmd := flag.NewFlagSet("finish", flag.ExitOnError)
	finishActivityCmd.Int64Var(&activityID, "id", 0, "activity id to finish")

	deleteActivityCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteActivityCmd.Int64Var(&activityID, "id", 0, "activity id to delete")

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

	case "finish":
		if err := finishActivityCmd.Parse(os.Args[2:]); err != nil {
			c.Usage()
			return
		}

	case "delete":
		if err := deleteActivityCmd.Parse(os.Args[2:]); err != nil {
			c.Usage()
			return
		}

	default:
		c.Usage()
		os.Exit(0)
	}

	if listCmd.Parsed() {
		hasFilter := period != "" || order != "" || limit > 0
		if hasFilter {

			if limit == 0 {
				fmt.Println("[WARN] Items limit = 0")
			}

			c.List(listName, &types.PeriodFilter{
				PeriodName: period,
				OrderBy:    order,
				Limit:      limit,
			})
		} else {
			c.List(listName, nil)
		}
	}

	if startActivityCmd.Parsed() {
		c.StartActivity(description, int64(category))
	}

	if finishActivityCmd.Parsed() {
		c.FinishActivity(activityID)
	}

	if deleteActivityCmd.Parsed() {
		c.DeleteActivity(activityID)
	}
}

func (c *CommandLine) Usage() {
	fmt.Println("-- Usage:")
	fmt.Println(" ")
	fmt.Println("     list  -n LIST_NAME -p PERIOD -o ORDER -l LIMIT")
	fmt.Println("              LIST_NAME (required): [categories, activities]")
	fmt.Println("              PERIOD (optional):    [today,yesterday,weekly,monthly]")
	fmt.Println("              ORDER  (optional):    [asc,desc]")
	fmt.Println("              LIMIT  (optional):    must be a number greater than 0")
	fmt.Println("")
	fmt.Println("     start -d DESCRIPTION -c CATEGORY_ID")
	fmt.Println("              DESCRIPTION (optional)")
	fmt.Println("              CATEGORY_ID (required): must be an existing category")
	fmt.Println("")
	fmt.Println("     finish -id ACTIVITY_ID")
	fmt.Println("                ACTIVITY_ID (required)")
}

func (c *CommandLine) List(name string, filter *types.PeriodFilter) {
	switch name {
	case "categories":
		c.ListCategories()
	case "activities":
		if filter != nil {
			c.FilterActivities(filter)
		} else {
			c.ListActivities()
		}
	default:
		c.Usage()
	}
}
