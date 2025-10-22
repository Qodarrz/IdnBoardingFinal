// model/user_custom.go
package models

import "time"

type UserCustomData struct {
	User             User
	UserProfile      UserProfile
	Vehicles         []VehicleWithCarbon
	Electronics      []ElectronicWithCarbon
	Missions         []MissionProgress
	Badges           []UserBadge
	PointTransactions []PointTransaction
	Points           Points
	ActivityLogs     []ActivityLog
	Orders           []Order
}

type VehicleWithCarbon struct {
	ID          int64     `json:"id"`
	VehicleType string    `json:"vehicle_type"`
	FuelType    string    `json:"fuel_type"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	TotalLogs   int       `json:"total_logs"`
	TotalCarbon float64   `json:"total_carbon_emission_g"`
}

type ElectronicWithCarbon struct {
	ID          int64     `json:"id"`
	DeviceName  string    `json:"device_name"`
	DeviceType  string    `json:"device_type"`
	PowerWatts  int       `json:"power_watts"`
	CreatedAt   time.Time `json:"created_at"`
	TotalLogs   int       `json:"total_logs"`
	TotalCarbon float64   `json:"total_carbon_emission_g"`
}

type MissionProgress struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	MissionType  string     `json:"mission_type"`
	PointsReward int        `json:"points_reward"`
	TargetValue  float64    `json:"target_value"`
	ProgressValue float64   `json:"progress_value"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type LeaderboardEntry struct {
	User               UserListItem
	TotalPoints        int     `json:"total_points"`
	CompletedMissions  int     `json:"completed_missions"`
	Score              float64 `json:"score"`
	CarbonReduction    float64 `json:"carbon_reduction_g"`
}

type UserListItem struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}