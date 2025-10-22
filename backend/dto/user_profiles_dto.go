package dto

import "time"

// RegisterDTO digunakan untuk menerima data registrasi user + profile
type ProfilesDTO struct {
	// UserProfile fields
	FullName  *string    `json:"full_name,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
}
