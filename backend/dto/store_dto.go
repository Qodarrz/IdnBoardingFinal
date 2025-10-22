// dto/store.go
package dto

import (
	"time"
)

type StoreItemDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	PricePoints int     `json:"price_points" validate:"required,min=1"`
	Stock       int     `json:"stock" validate:"min=0"`
	Status      string  `json:"status"`
	ImageURL    string  `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateStoreItemDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	PricePoints int    `json:"price_points" validate:"required,min=1"`
	Stock       int    `json:"stock" validate:"min=0"`
	ImageURL    string `json:"image_url"`
}

type UpdateStoreItemDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PricePoints int    `json:"price_points" validate:"min=1"`
	Stock       int    `json:"stock" validate:"min=0"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
}

type OrderItemDTO struct {
	ItemID int `json:"item_id" validate:"required"`
	Qty    int `json:"qty" validate:"required,min=1"`
}

type CreateOrderDTO struct {
	Items []OrderItemDTO `json:"items" validate:"required,min=1,dive"`
}

type OrderDTO struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	TotalPoints int            `json:"total_points"`
	Status      string         `json:"status"`
	Items       []OrderItemResponseDTO `json:"items"`
	CreatedAt   time.Time      `json:"created_at"`
}

type OrderItemResponseDTO struct {
	ID              int64  `json:"id"`
	ItemID          int64  `json:"item_id"`
	ItemName        string `json:"item_name"`
	Qty             int    `json:"qty"`
	PriceEachPoints int    `json:"price_each_points"`
	TotalPoints     int    `json:"total_points"`
}

type OrderResponseDTO struct {
	Order     OrderDTO `json:"order"`
	RemainingPoints int `json:"remaining_points"`
}