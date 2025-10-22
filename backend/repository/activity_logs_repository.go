package repository

import (
	"context"
	"database/sql"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
)

type ActivityRepositoryInterface interface {
	LogActivity(ctx context.Context, userID int64, activity string) error
	GetUserLogs(ctx context.Context, userID int64, limit int) ([]*model.ActivityLog, error)
}

type activityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) ActivityRepositoryInterface {
	return &activityRepository{db: db}
}

// Simpan log aktivitas user
func (r *activityRepository) LogActivity(ctx context.Context, userID int64, activity string) error {
	query := `
		INSERT INTO activity_logs (user_id, activity, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, activity, time.Now())
	return err
}

// Ambil log aktivitas user (misalnya untuk history)
func (r *activityRepository) GetUserLogs(ctx context.Context, userID int64, limit int) ([]*model.ActivityLog, error) {
	query := `
		SELECT id, user_id, activity, created_at
		FROM activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*model.ActivityLog
	for rows.Next() {
		l := &model.ActivityLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Activity, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	
	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return logs, nil
}