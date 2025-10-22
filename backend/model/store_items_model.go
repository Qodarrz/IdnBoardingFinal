package models

import "time"

type StoreItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string   `json:"description,omitempty"` // pakai pointer biar bisa null
	PricePoints int       `json:"price_points"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	ImageURL    string   `json:"image_url,omitempty"` // bisa null
	CreatedAt   time.Time `json:"created_at"`
}