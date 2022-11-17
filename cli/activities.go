package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ungame/timetrack/ioext"
	"github.com/ungame/timetrack/timeext"
	"github.com/ungame/timetrack/types"
	"io"
	"log"
	"net/http"
)

func PrintActivity(a *types.Activity) {

	category := GetCategory(a.CategoryID)

	fmt.Println("-- Activity")
	fmt.Println("     ID:         ", a.ID)
	fmt.Printf("     Category:    %s (ID=%d)\n", category.Name, a.CategoryID)
	fmt.Println("     Description:", a.Description)
	fmt.Println("     Status:     ", a.Status)
	if a.StartedAt != nil {
		fmt.Println("     Started:    ", a.StartedAt.Local().Format(timeext.DateTimeFormat))
	}
	if a.UpdatedAt != nil {
		fmt.Println("     Updated:    ", a.UpdatedAt.Local().Format(timeext.DateTimeFormat))
	}
	if a.FinishedAt != nil {
		fmt.Println("     Finished:   ", a.FinishedAt.Local().Format(timeext.DateTimeFormat))
		fmt.Println("     Duration:   ", a.FinishedAt.Sub(*a.StartedAt).String())
	}
	fmt.Println("")
}

func CreateActivity(description string, category int64) {

	input := &types.Activity{
		CategoryID:  category,
		Description: description,
	}

	payload, err := json.Marshal(input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	uri := fmt.Sprintf("%s/activities", baseURL)

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err.Error())
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

	var activity types.Activity
	err = json.Unmarshal(body, &activity)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	PrintActivity(&activity)
}

func ListActivities() {

	uri := fmt.Sprintf("%s/activities", baseURL)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Println(err.Error())
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

	var activities []*types.Activity
	err = json.Unmarshal(body, &activities)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, activity := range activities {
		PrintActivity(activity)
	}
}
