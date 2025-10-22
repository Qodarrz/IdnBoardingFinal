package models

import "time"	

type Badge struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	ImageURL       string    `json:"image_url"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type BadgeWithOwnership struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	ImageURL    string     `json:"image_url"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	IsOwned     bool       `json:"is_owned"`
	RedeemedAt  *time.Time `json:"redeemed_at,omitempty"`
}



