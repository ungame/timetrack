package models

import (
	"database/sql"
	"github.com/ungame/timetrack/pointer"
	"github.com/ungame/timetrack/types"
	"time"
)

type Category struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

func (c *Category) SetCreatedAt(createdAt time.Time) {
	c.CreatedAt = sql.NullTime{
		Time:  createdAt,
		Valid: !createdAt.IsZero(),
	}
}

func (c *Category) SetUpdatedAt(updatedAt time.Time) {
	c.UpdatedAt = sql.NullTime{
		Time:  updatedAt,
		Valid: !updatedAt.IsZero(),
	}
}

func (c *Category) Type() *types.Category {
	return &types.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   pointer.New(c.CreatedAt.Time),
		UpdatedAt:   pointer.New(c.UpdatedAt.Time),
	}
}
