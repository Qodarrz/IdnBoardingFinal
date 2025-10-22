package models

import "time"

type OrderItem struct {
	ID              int64     `json:"id"`
	OrderID         int64     `json:"order_id"`
	ItemID          int64     `json:"item_id"`
	Qty             int       `json:"qty"`
	PriceEachPoints int       `json:"price_each_points"`
	CreatedAt       time.Time `json:"created_at"`
}