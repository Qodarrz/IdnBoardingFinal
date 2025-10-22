package models

import "time"

type Notification struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    Type      string    `json:"type"`       // general, mission, store, etc.
    IsRead    bool      `json:"is_read"`
    CreatedAt time.Time `json:"created_at"`
    ReadAt    *time.Time `json:"read_at,omitempty"`
}
