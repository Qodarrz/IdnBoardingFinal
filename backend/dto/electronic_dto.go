// dto/electronics.go
package dto

import "time"

type DeviceType string

const (
	DeviceTypeKulkas     DeviceType = "kulkas"
	DeviceTypeLampu      DeviceType = "lampu"
	DeviceTypeMesinCuci  DeviceType = "mesin cuci"
	DeviceTypeTV         DeviceType = "tv"
	DeviceTypeKomputer   DeviceType = "komputer"
	DeviceTypeLaptop     DeviceType = "laptop"
	DeviceTypeSmartphone DeviceType = "smartphone"
	DeviceTypeMicrowave  DeviceType = "microwave"
	DeviceTypeFan        DeviceType = "fan"
	DeviceTypeAC         DeviceType = "ac"
)

func (dt DeviceType) IsValid() bool {
	switch dt {
	case DeviceTypeKulkas,
		DeviceTypeLampu,
		DeviceTypeMesinCuci,
		DeviceTypeTV,
		DeviceTypeKomputer,
		DeviceTypeLaptop,
		DeviceTypeSmartphone,
		DeviceTypeMicrowave,
		DeviceTypeFan,
		DeviceTypeAC:
		return true
	}
	return false
}

type CreateElectronicDTO struct {
	DeviceName string `json:"device_name" validate:"required"`
	DeviceType string `json:"device_type" validate:"required,oneof=kulkas lampu 'mesin cuci' tv komputer laptop smartphone microwave fan ac"`
	PowerWatts int    `json:"power_watts" validate:"required,gt=0"`
}

type AddElectronicsLogDTO struct {
	DeviceID      *int64    `json:"device_id,omitempty"`   // Optional jika membuat device baru
	DeviceName    string    `json:"device_name,omitempty"` // Required jika DeviceID tidak provided
	DeviceType    string    `json:"device_type,omitempty"` // Required jika DeviceID tidak providedx
	PowerWatts    int       `json:"power_watts,omitempty"` // Required jika DeviceID tidak provided
	DurationHours float64   `json:"duration_hours" validate:"required,gt=0"`
	LoggedAt      time.Time `json:"logged_at" validate:"required"`
}

type EditElectronicDTO struct {
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	PowerWatts int    `json:"power_watts"`
}
