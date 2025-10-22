// dto/user_custom_endpoint.go
package dto

import (
	"time"
)

// UserCustomDataResponseDTO untuk endpoint custom data user
// UserCustomDataResponseDTO untuk endpoint custom data user
type UserCustomDataResponseDTO struct {
	User          UserDetailResponseDTO   `json:"user"`
	Vehicles      []CustomVehicleDTO      `json:"vehicles,omitempty"`
	Electronics   []CustomElectronicDTO   `json:"electronics,omitempty"`
	Missions      []CustomMissionProgressDTO `json:"missions,omitempty"`
	Badges        []CustomBadgeDTO        `json:"badges,omitempty"`
	PointHistory  []CustomPointTransactionDTO `json:"point_history,omitempty"`
	ActivityLogs  []CustomActivityLogDTO  `json:"activity_logs,omitempty"`
	Orders        []CustomOrderDTO        `json:"orders,omitempty"`
	MonthlyVehicleCarbon    []MonthlyCarbonDTO    `json:"monthly_vehicle_carbon,omitempty"`    // Tambahan baru
	MonthlyElectronicCarbon []MonthlyCarbonDTO    `json:"monthly_electronic_carbon,omitempty"` // Tambahan baru
}

// MonthlyCarbonDTO untuk data karbon bulanan
type MonthlyCarbonDTO struct {
	Month        time.Time `json:"month"`
	TotalCarbon  float64   `json:"total_carbon_emission_g"`
}


type UserDetailResponseDTO struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	FullName    string     `json:"full_name"`
	AvatarURL   string     `json:"avatar_url"`
	Birthdate   *time.Time `json:"birthdate"`
	Gender      string     `json:"gender"`
	TotalPoints int        `json:"total_points"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CustomVehicleDTO struct {
	ID           int64     `json:"id"`
	VehicleType  string    `json:"vehicle_type"`
	FuelType     string    `json:"fuel_type"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	TotalLogs    int       `json:"total_logs"`
	TotalCarbon  float64   `json:"total_carbon_emission_g"`
}

type CustomElectronicDTO struct {
	ID           int64     `json:"id"`
	DeviceName   string    `json:"device_name"`
	DeviceType   string    `json:"device_type"`
	PowerWatts   int       `json:"power_watts"`
	CreatedAt    time.Time `json:"created_at"`
	TotalLogs    int       `json:"total_logs"`
	TotalCarbon  float64   `json:"total_carbon_emission_g"`
}

type CustomMissionProgressDTO struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	MissionType  string     `json:"mission_type"`
	PointsReward int        `json:"points_reward"`
	TargetValue  float64    `json:"target_value"`
	Progress     float64    `json:"progress"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CustomBadgeDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	RedeemedAt  time.Time `json:"redeemed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type CustomPointTransactionDTO struct {
	ID            int64     `json:"id"`
	Amount        int       `json:"amount"`
	Direction     string    `json:"direction"`
	Source        string    `json:"source"`
	ReferenceType string    `json:"reference_type,omitempty"`
	ReferenceID   int64     `json:"reference_id,omitempty"`
	Note          string    `json:"note,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CustomActivityLogDTO struct {
	ID        int64     `json:"id"`
	Activity  string    `json:"activity"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomOrderDTO struct {
	ID           int64              `json:"id"`
	TotalPoints  int                `json:"total_points"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	Items        []CustomOrderItemDTO `json:"items,omitempty"`
}

type CustomOrderItemDTO struct {
	ID              int64  `json:"id"`
	ItemName        string `json:"item_name"`
	Qty             int    `json:"qty"`
	PriceEachPoints int    `json:"price_each_points"`
}

// LeaderboardRequestDTO untuk request leaderboard
type LeaderboardRequestDTO struct {
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	TimeRange string `query:"time_range" validate:"omitempty,oneof=day week month all"`
}

type LeaderboardItemDTO struct {
	Rank              int           `json:"rank"`
	User              UserSimpleDTO `json:"user"`
	TotalPoints       int           `json:"total_points"`
	CompletedMissions int           `json:"completed_missions"`
	Score             float64       `json:"score"` // total_points * completed_missions
	CarbonReduction   float64       `json:"carbon_reduction_g,omitempty"`
}

type LeaderboardResponseDTO struct {
	Leaderboard  []LeaderboardItemDTO `json:"leaderboard"`
	Pagination   PaginationDTO        `json:"pagination"`
	TimeRange    string               `json:"time_range"`
}

type UserSimpleDTO struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

type PaginationDTO struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}