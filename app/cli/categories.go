package cli

import (
	"encoding/json"
	"fmt"
	"github.com/ungame/timetrack/ioext"
	"github.com/ungame/timetrack/timeext"
	"github.com/ungame/timetrack/types"
	"io"
	"net/http"
)

func PrintCategory(c *types.Category) {
	fmt.Println("-- Category")
	fmt.Println("     ID:         ", c.ID)
	fmt.Println("     Name:       ", c.Name)
	fmt.Println("     Description:", c.Description)
	if c.CreatedAt != nil {
		fmt.Println("     Created:", c.CreatedAt.Local().Format(timeext.DateTimeFormat))
	}
	if c.UpdatedAt != nil {
		fmt.Println("     Updated:", c.UpdatedAt.Local().Format(timeext.DateTimeFormat))
	}
	fmt.Println("")
}

func (c *CommandLine) GetCategory(id int64) (category *types.Category) {
	uri := fmt.Sprintf("%s/categories/%d", c.baseURL, id)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if res != nil && res.Body != nil {
		defer ioext.Close(res.Body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal(body, &category)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func (c *CommandLine) ListCategories() {
	uri := fmt.Sprintf("%s/categories", c.baseURL)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if res != nil && res.Body != nil {
		defer ioext.Close(res.Body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var categories []*types.Category
	err = json.Unmarshal(body, &categories)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, category := range categories {
		PrintCategory(category)
	}
}
