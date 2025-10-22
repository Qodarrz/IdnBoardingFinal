package models

import "time"

type ActivityLog struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Activity  string    `db:"activity" json:"activity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
