package models

import "time"

type CarbonElectronic struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	DeviceName string   `json:"device_name"`
	DeviceType string   `json:"device_type"`
	PowerWatts int      `json:"power_watts"`
	CreatedAt  time.Time `json:"created_at"`
}