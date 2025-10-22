package models

type MissionBadge struct {
	ID        int64 `db:"id"`
	MissionID int64 `db:"mission_id"`
	BadgeID   int64 `db:"badge_id"`
}
