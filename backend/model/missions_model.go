// model/mission.go
package models

import (
	"database/sql"
	"time"
)

type MissionType string

const (
	MissionTypeCarbonReduction MissionType = "carbon_reduction"
	MissionTypeStreak          MissionType = "streak"
	MissionTypeActivity        MissionType = "activity"
	MissionTypeCustom          MissionType = "custom"
)

type MissionCriteriaType string

const (
	// Vehicles
	CriteriaCar             MissionCriteriaType = "car"
	CriteriaMotorcycle      MissionCriteriaType = "motorcycle"
	CriteriaBicycle         MissionCriteriaType = "bicycle"
	CriteriaPublicTransport MissionCriteriaType = "public_transport"
	CriteriaWalk            MissionCriteriaType = "walk"

	// Electronics
	CriteriaLaptop         MissionCriteriaType = "laptop"
	CriteriaDesktop        MissionCriteriaType = "desktop"
	CriteriaTV             MissionCriteriaType = "tv"
	CriteriaAC             MissionCriteriaType = "ac"
	CriteriaFridge         MissionCriteriaType = "fridge"
	CriteriaFan            MissionCriteriaType = "fan"
	CriteriaWashingMachine MissionCriteriaType = "washing_machine"
	CriteriaOther          MissionCriteriaType = "other"
)

type Mission struct {
	ID               int64           `json:"id"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	MissionType      MissionType     `json:"mission_type"`
	CriteriaType     MissionCriteriaType  `json:"criteria_type"` // ✅ ditambahkan
	PointsReward     int             `json:"points_reward"`
	GivesBadge       bool            `json:"gives_badge"`
	BadgeID          sql.NullInt64   `json:"badge_id"`
	CarbonReductionG sql.NullFloat64 `json:"carbon_reduction_g"`
	TargetValue      float64         `json:"target_value"`
	ExpiredAt        sql.NullTime    `json:"expired_at"`
	CreatedAt        time.Time       `json:"created_at"`
}

// Membuat sql.NullInt64 dari int64
func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// Membuat sql.NullFloat64 dari float64
func NewNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

// Membuat sql.NullTime dari time.Time
func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

// Membuat sql.NullString dari string
func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
