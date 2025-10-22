package models

import "time"

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	TotalPoints int       `json:"total_points"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
