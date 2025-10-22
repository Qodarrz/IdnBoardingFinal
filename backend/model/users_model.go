package models

import "time"

// Di model/user.go
type User struct {
	ID              int64          `json:"id" db:"id"`
	Username        string         `json:"username" db:"username"`
	Email           string         `json:"email" db:"email"`
	Password        string         `json:"password" db:"password"`
	Role            string         `json:"role" db:"role"`
	GoogleID        *string        `json:"google_id,omitempty" db:"google_id"` // Ubah menjadi string
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty" db:"email_verified_at"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	Profile         *UserProfile   `json:"profile,omitempty"`
}