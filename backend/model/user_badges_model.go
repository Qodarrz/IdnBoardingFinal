package models

import (
	"database/sql"
	"time"
)

type UserBadge struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	BadgeID    int64      `json:"badge_id"`
	RedeemedAt sql.NullTime   `json:"redeemed_at"`
	CreatedAt  time.Time  `json:"created_at"`
	Badge      *Badge     `json:"badge,omitempty"`
}