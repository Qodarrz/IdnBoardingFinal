package models

import (
	"database/sql"
	"time"
)

type UserMission struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	MissionID   int64        `json:"mission_id"`
	CompletedAt sql.NullTime `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at"`
	Mission     *Mission     `json:"mission,omitempty"`
}