package repository

import (
	"context"
	"database/sql"
	"github.com/ungame/timetrack/app/models"
	"github.com/ungame/timetrack/ioext"
)

type CategoriesRepository interface {
	Get(ctx context.Context, id int64) (*models.Category, error)
	GetAll(ctx context.Context) ([]*models.Category, error)
}

type categoriesRepository struct {
	conn *sql.DB
}

func NewCategoriesRepository(conn *sql.DB) CategoriesRepository {
	return &categoriesRepository{conn: conn}
}

func (r *categoriesRepository) Get(ctx context.Context, id int64) (*models.Category, error) {
	var (
		query    = "select * from categories where id = ?"
		row      = r.conn.QueryRowContext(ctx, query, id)
		category = new(models.Category)
	)
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	return category, err
}

func (r *categoriesRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	query := "select * from categories"
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer ioext.Close(rows)
	categories := make([]*models.Category, 0, 5)
	for rows.Next() {
		category := new(models.Category)
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
