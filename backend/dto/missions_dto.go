// dto/mission.go
package dto

import (
	"time"
)

type MissionType string

const (
	MissionTypeCarbonReduction MissionType = "carbon_reduction"
	MissionTypeStreak          MissionType = "streak"
	MissionTypeActivity        MissionType = "activity"
	MissionTypeCustom          MissionType = "custom"
)

type CriteriaType string

const (
	// Vehicles
	CriteriaCar             CriteriaType = "car"
	CriteriaMotorcycle      CriteriaType = "motorcycle"
	CriteriaBicycle         CriteriaType = "bicycle"
	CriteriaPublicTransport CriteriaType = "public_transport"
	CriteriaWalk            CriteriaType = "walk"

	CriteriaLaptop         CriteriaType = "laptop"
	CriteriaDesktop        CriteriaType = "desktop"
	CriteriaTV             CriteriaType = "tv"
	CriteriaAC             CriteriaType = "ac"
	CriteriaFridge         CriteriaType = "fridge"
	CriteriaFan            CriteriaType = "fan"
	CriteriaWashingMachine CriteriaType = "washing_machine"
	CriteriaOther          CriteriaType = "other"
)

type CreateMissionDTO struct {
	Title            string        `json:"title" validate:"required"`
	Description      string        `json:"description"`
	MissionType      MissionType   `json:"mission_type" validate:"required,oneof=carbon_reduction streak activity custom"`
	CriteriaType     *CriteriaType `json:"criteria_type,omitempty"` // ✅ baru ditambahkan
	PointsReward     int           `json:"points_reward" validate:"required,min=0"`
	GivesBadge       bool          `json:"gives_badge"`
	BadgeID          *int64        `json:"badge_id"`
	CarbonReductionG *float64      `json:"carbon_reduction_g"`
	TargetValue      float64       `json:"target_value" validate:"required"`
	ExpiredAt        *time.Time    `json:"expired_at"`
}

type MissionResponseDTO struct {
	ID               int64         `json:"id"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	MissionType      MissionType   `json:"mission_type"`
	CriteriaType     *CriteriaType `json:"criteria_type,omitempty"`
	PointsReward     int           `json:"points_reward"`
	GivesBadge       bool          `json:"gives_badge"`
	BadgeID          *int64        `json:"badge_id,omitempty"`
	CarbonReductionG *float64      `json:"carbon_reduction_g,omitempty"`
	TargetValue      interface{}   `json:"target_value,omitempty"`
	ExpiredAt        *time.Time    `json:"expired_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
}

type UserMissionResponseDTO struct {
	ID          int64              `json:"id"`
	UserID      int64              `json:"user_id"`
	Mission     MissionResponseDTO `json:"mission"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

type CreateMissionWithBadgeDTO struct {
	Title            string       `json:"title" validate:"required"`
	Description      string       `json:"description"`
	MissionType      MissionType  `json:"mission_type" validate:"required,oneof=carbon_reduction streak activity custom"`
	CriteriaType     CriteriaType `json:"criteria_type,omitempty"`
	PointsReward     int          `json:"points_reward" validate:"required,min=0"`
	GivesBadge       bool         `json:"gives_badge" validate:"required"`
	BadgeName        string       `json:"badge_name,omitempty"`
	BadgeImageURL    string       `json:"badge_image_url,omitempty"`
	BadgeDescription string       `json:"badge_description,omitempty"`
	TargetValue      float64      `json:"target_value" validate:"required"`
	ExpiredAt        *time.Time   `json:"expired_at"`
}

type MissionWithBadgeResponseDTO struct {
	Mission MissionResponseDTO `json:"mission"`
	Badge   BadgeResponseDTO   `json:"badge,omitempty"`
}

type MissionProgressResponse struct {
	MissionID    int64   `json:"mission_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	TargetValue  float64 `json:"target_value"`
	Progress     float64 `json:"progress"`
	IsCompleted  bool    `json:"is_completed"`
	PointsReward float64 `json:"points_reward"`
	GivesBadge   bool    `json:"gives_badge"`
	BadgeID      *int64  `json:"badge_id,omitempty"`
}
