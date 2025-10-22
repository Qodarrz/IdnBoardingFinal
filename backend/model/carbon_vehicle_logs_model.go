package models

import "time"


type Coordinates struct {
    Lat float64 `json:"lat" validate:"required"`
    Lon float64 `json:"lon" validate:"required"`
}

type CarbonVehicleLog struct {
	ID              int64     `db:"id"`
	VehicleID       int64     `db:"vehicle_id"`
	StartLat        float64   `db:"start_lat"`
	StartLon        float64   `db:"start_lon"`
	EndLat          float64   `db:"end_lat"`
	EndLon          float64   `db:"end_lon"`
	DistanceKm      float64   `db:"distance_km"`
	DurationMinutes int       `db:"duration_minutes"`
	CarbonEmission  float64   `db:"carbon_emission_g"`
	LoggedAt        time.Time `db:"logged_at"`
}

type CarbonVehicleWithLog struct {
    ID          int64               `db:"id"`
    UserID      int64               `db:"user_id"`
    VehicleType string              `db:"vehicle_type"`
    FuelType    string              `db:"fuel_type"`
    Name        string              `db:"name"`
    LatestLog   *CarbonVehicleLog   `db:"-"`
}



func (CarbonVehicleLog) TableName() string {
	return "carbon_vehicle_logs"
}
