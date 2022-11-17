package repository

import (
	"context"
	"database/sql"
	"github.com/ungame/timetrack/app/models"
	"github.com/ungame/timetrack/ioext"
	"github.com/ungame/timetrack/queries"
	"time"
)

const (
	createActivityQuery       = "insert into activities (category_id, description, status, started_at, updated_at) values (?, ?, ?, ?, ?)"
	updateActivityQuery       = "update activities set category_id = ?, description = ?, status = ?, updated_at = ?, finished_at = ? where id = ?"
	deleteActivityQuery       = "delete from activities where id = ?"
	defaultPrepareStmtTimeout = time.Second * 30
)

type ActivitiesRepository interface {
	Create(ctx context.Context, activity *models.Activity) (*models.Activity, error)
	Get(ctx context.Context, id int64) (*models.Activity, error)
	GetAll(ctx context.Context) ([]*models.Activity, error)
	GetByStatus(ctx context.Context, status models.ActivityStatus) ([]*models.Activity, error)
	Update(ctx context.Context, activity *models.Activity) (*models.Activity, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Close()
}

type activitiesRepository struct {
	conn       *sql.DB
	createStmt *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
}

func NewActivitiesRepository(conn *sql.DB) ActivitiesRepository {
	ctx, cancel := context.WithTimeout(context.Background(), defaultPrepareStmtTimeout)
	defer cancel()

	return &activitiesRepository{
		conn:       conn,
		createStmt: queries.MustPrepare(ctx, conn, createActivityQuery),
		updateStmt: queries.MustPrepare(ctx, conn, updateActivityQuery),
		deleteStmt: queries.MustPrepare(ctx, conn, deleteActivityQuery),
	}
}

func (r *activitiesRepository) Create(ctx context.Context, activity *models.Activity) (*models.Activity, error) {
	result, err := r.createStmt.ExecContext(
		ctx,
		activity.CategoryID,
		activity.Description,
		activity.Status,
		activity.StartedAt,
		activity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	activity.ID, err = result.LastInsertId()
	return activity, err
}

func (r *activitiesRepository) Get(ctx context.Context, id int64) (*models.Activity, error) {
	var (
		query    = "select * from activities where id = ?"
		row      = r.conn.QueryRowContext(ctx, query, id)
		activity = new(models.Activity)
	)
	err := row.Scan(
		&activity.ID,
		&activity.CategoryID,
		&activity.Description,
		&activity.Status,
		&activity.StartedAt,
		&activity.UpdatedAt,
		&activity.FinishedAt,
	)
	return activity, err
}

func (r *activitiesRepository) GetAll(ctx context.Context) ([]*models.Activity, error) {
	query := "select * from activities"
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer ioext.Close(rows)
	activities := make([]*models.Activity, 0, 10)
	for rows.Next() {
		activity := new(models.Activity)
		err = rows.Scan(
			&activity.ID,
			&activity.CategoryID,
			&activity.Description,
			&activity.Status,
			&activity.StartedAt,
			&activity.UpdatedAt,
			&activity.FinishedAt,
		)
		if err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	return activities, err
}

func (r *activitiesRepository) Update(ctx context.Context, activity *models.Activity) (*models.Activity, error) {
	_, err := r.updateStmt.ExecContext(
		ctx,
		activity.CategoryID,
		activity.Description,
		activity.Status,
		activity.UpdatedAt,
		activity.FinishedAt,
		activity.ID,
	)
	return activity, err
}

func (r *activitiesRepository) Delete(ctx context.Context, id int64) (int64, error) {
	result, err := r.deleteStmt.ExecContext(ctx, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *activitiesRepository) GetByStatus(ctx context.Context, status models.ActivityStatus) ([]*models.Activity, error) {
	query := "select * from activities where status = ?"
	rows, err := r.conn.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer ioext.Close(rows)
	activities := make([]*models.Activity, 0, 10)
	for rows.Next() {
		activity := new(models.Activity)
		err = rows.Scan(
			&activity.ID,
			&activity.CategoryID,
			&activity.Description,
			&activity.Status,
			&activity.StartedAt,
			&activity.UpdatedAt,
			&activity.FinishedAt,
		)
		if err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	return activities, err
}

func (r *activitiesRepository) Close() {
	ioext.Close(r.createStmt)
	ioext.Close(r.updateStmt)
	ioext.Close(r.deleteStmt)
}
