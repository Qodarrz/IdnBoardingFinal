package models

import "time"

type UserProfile struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	FullName  *string    `db:"full_name"`
	AvatarURL *string    `db:"avatar_url"`
	Birthdate *time.Time `db:"birthdate"`
	Gender    *string    `db:"gender"`
	CreatedAt time.Time  `db:"created_at"`
}
