package types

import "time"

type Activity struct {
	ID          int64      `json:"id"`
	CategoryID  int64      `json:"category_id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}
