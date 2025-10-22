// dto/user_profile.go
package dto

import (
	"time"
)

// dto/user_profile.go
type UserProfileUpdateDTO struct {
	FullName  *string    `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	AvatarURL *string    `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Gender    *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
}

type UserProfileResponseDTO struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	FullName  *string     `json:"full_name"`
	AvatarURL *string     `json:"avatar_url"`
	Birthdate *time.Time  `json:"birthdate"`
	Gender    *string     `json:"gender"`	
	CreatedAt time.Time  `json:"created_at"`
}

type UserWithProfileResponseDTO struct {
	ID       int64                  `json:"id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Role     string                 `json:"role"`
	Profile  UserProfileResponseDTO `json:"profile"`
}