package dto

import "time"

// User
type UserDTO struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	FullName    string     `json:"full_name"`
	AvatarURL   string     `json:"avatar_url"`
	Birthdate   *time.Time `json:"birthdate"`
	Gender      string     `json:"gender"`
	TotalPoints int64      `json:"total_points"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Vehicles
type VehicleDTO struct {
	ID          int64   `json:"id"`
	VehicleType string  `json:"vehicle_type"`
	FuelType    string  `json:"fuel_type"`
	Name        string  `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	TotalLogs   int     `json:"total_logs"`
	TotalCarbon float64 `json:"total_carbon_emission_g"`
}

// Electronics
type ElectronicDTO struct {
	ID          int64   `json:"id"`
	DeviceName  string  `json:"device_name"`
	DeviceType  string  `json:"device_type"`
	PowerWatts  int     `json:"power_watts"`
	CreatedAt   time.Time `json:"created_at"`
	TotalLogs   int     `json:"total_logs"`
	TotalCarbon float64 `json:"total_carbon_emission_g"`
}

// Missions
type MissionDTO struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	MissionType   string     `json:"mission_type"`
	PointsReward  int        `json:"points_reward"`
	TargetValue   float64    `json:"target_value"`
	ProgressValue float64    `json:"progress_value"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Badges
type BadgeDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	RedeemedAt  time.Time `json:"redeemed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Point Transactions
type PointTransactionDTO struct {
	ID            int64     `json:"id"`
	Amount        int       `json:"amount"`
	Direction     string    `json:"direction"`
	Source        string    `json:"source"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   int64     `json:"reference_id"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

// Activity Logs
type ActivityLogDTO struct {
	ID        int64     `json:"id"`
	Activity  string    `json:"activity"`
	CreatedAt time.Time `json:"created_at"`
}

// Orders


// Leaderboard
type LeaderboardEntryDTO struct {
	User            UserSimpleDTO `json:"user"`
	TotalPoints     int64         `json:"total_points"`
	CompletedMissions int         `json:"completed_missions"`
	Score           float64       `json:"score"`
	CarbonReduction float64       `json:"carbon_reduction_g"`
}

// Aggregate UserCustomData
type UserCustomDataDTO struct {
	User              UserDTO                `json:"user"`
	Vehicles          []VehicleDTO           `json:"vehicles"`
	Electronics       []ElectronicDTO        `json:"electronics"`
	Missions          []MissionDTO           `json:"missions"`
	Badges            []BadgeDTO             `json:"badges"`
	PointTransactions []PointTransactionDTO  `json:"point_transactions"`
	ActivityLogs      []ActivityLogDTO       `json:"activity_logs"`
	Orders            []OrderDTO             `json:"orders"`
}
