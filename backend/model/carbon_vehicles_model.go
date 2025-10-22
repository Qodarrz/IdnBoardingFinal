package models

import "time"

type VehicleType string
type FuelType string

const (
    VehicleCar          VehicleType = "car"
    VehicleMotorcycle   VehicleType = "motorcycle"
    VehicleBicycle      VehicleType = "bicycle"
    VehiclePublicTrans  VehicleType = "public_transport"
    VehicleWalk         VehicleType = "walk"
)

const (
    FuelPetrol   FuelType = "petrol"
    FuelDiesel   FuelType = "diesel"
    FuelElectric FuelType = "electric"
    FuelNone     FuelType = "none"
)


type CarbonVehicle struct {
    ID          int64       `json:"id"`
    UserID      int64       `json:"user_id"`
    VehicleType VehicleType `json:"vehicle_type"`
    FuelType    FuelType    `json:"fuel_type"`
    Name        string      `json:"name"`
    CreatedAt   time.Time   `json:"created_at"`
}
