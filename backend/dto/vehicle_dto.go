// dto/vehicle.go
package dto

type CreateVehicleDTO struct {
	VehicleType string `json:"vehicle_type" validate:"required"`
	FuelType    string `json:"fuel_type" validate:"required"`
	Name        string `json:"name" validate:"required"`
}

type Coordinates struct {
	Lat float64 `json:"lat" validate:"required"`
	Lon float64 `json:"lon" validate:"required"`
}

type AddVehicleLogDTO struct {
	VehicleID       *int64  `json:"vehicle_id,omitempty"`
	VehicleType     string  `json:"vehicle_type,omitempty"`
	FuelType        string  `json:"fuel_type,omitempty"`
	VehicleName     string  `json:"vehicle_name,omitempty"`
	StartLat        float64 `json:"start_lat" validate:"required"`
	StartLon        float64 `json:"start_lon" validate:"required"`
	EndLat          float64 `json:"end_lat" validate:"required"`
	EndLon          float64 `json:"end_lon" validate:"required"`
	DistanceKm      float64 `json:"distance_km" validate:"required,gt=0"`
	DurationMinutes int     `json:"duration_minutes" validate:"required,gt=0"`
}

type EditVehicleDTO struct {
	VehicleType string `json:"vehicle_type"`
	FuelType    string `json:"fuel_type"`
	Name        string `json:"name"`
}
