// dto/badge.go
package dto

import (
	"time"
)

type CreateBadgeDTO struct {
	Name           string `json:"name" validate:"required"`
	ImageURL       string `json:"image_url" validate:"required,url"`
	Description    string `json:"description"`
	RequiredPoints int    `json:"required_points" validate:"required,min=0"`
}

type BadgeResponseDTO struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	ImageURL       string    `json:"image_url"`
	Description    string    `json:"description"`
	RequiredPoints int       `json:"required_points"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserBadgeResponseDTO struct {
	ID         int64             `json:"id"`
	UserID     int64             `json:"user_id"`
	Badge      BadgeResponseDTO  `json:"badge"`
	RedeemedAt *time.Time        `json:"redeemed_at,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}