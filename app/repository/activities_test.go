package repository

import (
	"context"
	"fmt"
	"github.com/ungame/timetrack/app/models"
	"github.com/ungame/timetrack/db"
	"github.com/ungame/timetrack/ioext"
	"github.com/ungame/timetrack/queries"
	"github.com/ungame/timetrack/timeext"
	"sort"
	"testing"
	"time"
)

func TestActivitiesRepository_FilterByPeriod(t *testing.T) {

	var (
		ctx      = context.Background()
		conn     = db.New()
		repo     = NewActivitiesRepository(conn)
		testTime = time.Date(2022, time.November, 01, 12, 0, 0, 0, time.UTC)
		err      error
	)

	defer ioext.Close(conn)
	defer repo.Close()

	finishedActivity := &models.Activity{
		CategoryID:  5,
		Description: "YESTERDAY",
		Status:      models.Finished,
	}

	yesterday := testTime.AddDate(0, 0, -1)
	finishedAt := yesterday.Add(time.Hour)
	finishedActivity.SetStartedAt(yesterday)
	finishedActivity.SetUpdatedAt(finishedAt)
	finishedActivity.SetFinishedAt(finishedAt)

	finishedActivity, err = repo.Create(ctx, finishedActivity)
	if err != nil {
		t.Error(err)
	}

	// try cleanup tests activities from db
	defer func() { _, _ = repo.Delete(ctx, finishedActivity.ID) }()

	startedActivity := &models.Activity{
		CategoryID:  1,
		Description: "TODAY",
		Status:      models.Started,
	}
	startedActivity.SetStartedAt(testTime)
	startedActivity.SetUpdatedAt(testTime)

	todayPeriod := _TestPeriod{
		start: timeext.GetStartOfDayFrom(testTime, testTime.Location()),
		end:   testTime,
	}

	startedActivity, err = repo.Create(ctx, startedActivity)
	if err != nil {
		t.Error(err)
	}

	// try cleanup tests activities from db
	defer func() { _, _ = repo.Delete(ctx, startedActivity.ID) }()

	// starting tests
	items, err := repo.FilterByPeriod(ctx, &todayPeriod, queries.Desc, 10)
	if err != nil {
		t.Error(err)
	}

	// order by id, pre-requisite to use sort.Search
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	index := sort.Search(len(items), func(i int) bool {
		return items[i].ID >= startedActivity.ID
	})

	if items[index].ID != startedActivity.ID {
		t.Errorf("unexpected item: \nexpected=%+v \ngot=%+v", startedActivity, items[index])
	}

	// ensuring all items started after todayPeriod
	for _, item := range items {
		if item.StartedAt.Time.Before(todayPeriod.start) {
			t.Errorf("started_at field must be greater or equal = %s", todayPeriod.start.Format(timeext.DateTimeFormat))
		}
	}

	// ensuring that started item yesterday is not present in items
	index = sort.Search(len(items), func(i int) bool {
		return items[i].ID >= finishedActivity.ID
	})

	// index is equals items length when item is not found
	if index != len(items) && items[index].ID == finishedActivity.ID {
		t.Errorf("unexpected item: %+v", items[index])
	}

	// start testing yesterday items only
	yesterdayPeriod := _TestPeriod{
		start: timeext.GetStartOfDayFrom(yesterday, testTime.Location()),
		end:   timeext.GetEndOfDayFrom(yesterday, testTime.Location()),
	}

	items, err = repo.FilterByPeriod(ctx, &yesterdayPeriod, queries.Desc, 10)
	if err != nil {
		t.Error(err)
	}

	// order by id, pre-requisite to use sort.Search
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	index = sort.Search(len(items), func(i int) bool {
		return items[i].ID >= finishedActivity.ID
	})
	if items[index].ID != finishedActivity.ID {
		t.Errorf("unexpected item: \nexpected=%+v \ngot=%+v", finishedActivity, items[index])
	}

	// ensuring all items started yesterday
	for _, item := range items {
		if item.StartedAt.Time.Before(yesterdayPeriod.start) {
			t.Errorf("started_at field must be greater or equal = %s", yesterdayPeriod.start.Format(timeext.DateTimeFormat))
		}
		if item.StartedAt.Time.After(yesterdayPeriod.end) {
			t.Errorf("started_at field must be less or equal = %s", yesterdayPeriod.end.Format(timeext.DateTimeFormat))
		}
	}

	// ensuring that started item today is not present in items
	index = sort.Search(len(items), func(i int) bool {
		return items[i].ID >= startedActivity.ID
	})

	// index is equals items length when item is not found
	if index != len(items) && items[index].ID == startedActivity.ID {
		t.Errorf("unexpected item: %+v", items[index])
	}
}

type _TestPeriod struct {
	start time.Time
	end   time.Time
}

func (t *_TestPeriod) Range(_ *time.Location) (time.Time, time.Time) {
	return t.start, t.end
}

func (t *_TestPeriod) Print() {
	fmt.Println("Start:", t.start.Format(timeext.DateTimeFormat))
	fmt.Println("End:", t.end.Format(timeext.DateTimeFormat))
}
