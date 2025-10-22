package repository

import (
	"context"
	"database/sql"

	model "github.com/Qodarrz/fiber-app/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	GetByUserID(ctx context.Context, userID int64) ([]*model.Notification, error)
}


type notificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(ctx context.Context, notif *model.Notification) error {
	query := `INSERT INTO notifications (user_id, title, message, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, notif.UserID, notif.Title, notif.Message, notif.CreatedAt).Scan(&notif.ID)
	return err
}

func (r *notificationRepo) GetByUserID(ctx context.Context, userID int64) ([]*model.Notification, error) {
	query := `SELECT id, user_id, title, message, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		notif := &model.Notification{}
		err := rows.Scan(&notif.ID, &notif.UserID, &notif.Title, &notif.Message, &notif.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}

