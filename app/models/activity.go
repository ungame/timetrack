package models

import (
	"database/sql"
	"github.com/ungame/timetrack/pointer"
	"github.com/ungame/timetrack/types"
	"time"
)

type ActivityStatus string

func (a ActivityStatus) String() string {
	if a == Started {
		return "STARTED"
	}
	return "FINISHED"
}

const (
	Started  ActivityStatus = "1"
	Finished ActivityStatus = "0"
)

type Activity struct {
	ID          int64
	CategoryID  int64
	Description string
	Status      ActivityStatus
	StartedAt   sql.NullTime
	UpdatedAt   sql.NullTime
	FinishedAt  sql.NullTime
}

func (a *Activity) SetStartedAt(startedAt time.Time) {
	a.StartedAt = sql.NullTime{
		Time:  startedAt,
		Valid: !startedAt.IsZero(),
	}
}

func (a *Activity) SetUpdatedAt(updatedAt time.Time) {
	a.UpdatedAt = sql.NullTime{
		Time:  updatedAt,
		Valid: !updatedAt.IsZero(),
	}
}

func (a *Activity) SetFinishedAt(finishedAt time.Time) {
	a.FinishedAt = sql.NullTime{
		Time:  finishedAt,
		Valid: !finishedAt.IsZero(),
	}
}

func (a *Activity) Type() *types.Activity {
	activity := &types.Activity{
		ID:          a.ID,
		CategoryID:  a.CategoryID,
		Description: a.Description,
		Status:      a.Status.String(),
		StartedAt:   pointer.New(a.StartedAt.Time),
		UpdatedAt:   pointer.New(a.UpdatedAt.Time),
	}
	if a.FinishedAt.Valid {
		activity.FinishedAt = pointer.New(a.FinishedAt.Time)
	}
	return activity
}
