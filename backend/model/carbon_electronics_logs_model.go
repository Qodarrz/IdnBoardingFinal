package models

import "time"

type CarbonElectronicLog struct {
	ID             int64     `db:"id"`
	DeviceID       int64     `db:"device_id"`
	DurationHours  float64   `db:"duration_hours"`
	CarbonEmission float64   `db:"carbon_emission_g"`
	LoggedAt       time.Time `db:"logged_at" json:"LoggedAt"`
}

func (CarbonElectronicLog) TableName() string {
	return "carbon_electronics_logs"
}
