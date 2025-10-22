package models

import "time"

type PointTransaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Amount        int       `json:"amount"`
	Direction     string    `json:"direction"`
	Source        string    `json:"source"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   int64     `json:"reference_id"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}
