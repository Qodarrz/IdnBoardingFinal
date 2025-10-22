package models

import "time"

type Points struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	TotalPoints int       `json:"total_points"`
	CreatedAt   time.Time `json:"created_at"`
}

