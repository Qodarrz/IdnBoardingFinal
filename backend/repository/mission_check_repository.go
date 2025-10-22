package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
)

// =========================
// Repository Interface
// =========================
type CheckMissionRepositoryInterface interface {
	FindByID(ctx context.Context, id int64) (*model.Mission, error)
	FindActiveMissions(ctx context.Context) ([]*model.Mission, error)
	FindMissionsByType(ctx context.Context, missionType model.MissionType) ([]*model.Mission, error)
	FindMissionsByCriteriaType(ctx context.Context, criteriaType model.MissionCriteriaType) ([]*model.Mission, error)
	AssignMissionToUser(ctx context.Context, userID, missionID int64) error
	GetMissionProgress(ctx context.Context, userID, missionID int64) (float64, error)
	UpdateMissionProgress(ctx context.Context, userID, missionID int64, progress float64) error
	MarkMissionCompleted(ctx context.Context, userID, missionID int64) error
	HasUserCompletedMission(ctx context.Context, userID, missionID int64) (bool, error)
	CheckMission(ctx context.Context, userID int64, mission *model.Mission) (bool, error)
	CheckAllUserMissions(ctx context.Context, userID int64) error
	CheckUserMissionsByType(ctx context.Context, userID int64, missionType model.MissionType) error
	CheckUserMissionsByCriteriaType(ctx context.Context, userID int64, criteriaType model.MissionCriteriaType) error
}

type checkMissionRepository struct {
	db *sql.DB
}

func CheckMissionRepository(db *sql.DB) CheckMissionRepositoryInterface {
	return &checkMissionRepository{db: db}
}

// =========================
// Query Functions
// =========================
func (r *checkMissionRepository) FindByID(ctx context.Context, id int64) (*model.Mission, error) {
	mission := &model.Mission{}
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward, 
		       gives_badge, badge_id, target_value, created_at, expired_at
		FROM missions
		WHERE id = $1
	`

	var criteriaType sql.NullString
	var badgeID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&mission.ID, &mission.Title, &mission.Description, &mission.MissionType,
		&criteriaType, &mission.PointsReward, &mission.GivesBadge, &badgeID,
		&mission.TargetValue, &mission.CreatedAt, &mission.ExpiredAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Handle nullable fields
	if criteriaType.Valid {
		mission.CriteriaType = model.MissionCriteriaType(criteriaType.String)
	}
	if badgeID.Valid {
		mission.BadgeID = badgeID
	}

	return mission, nil
}

func (r *checkMissionRepository) FindActiveMissions(ctx context.Context) ([]*model.Mission, error) {
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward,
		       gives_badge, badge_id, target_value, created_at, expired_at
		FROM missions
		WHERE (expired_at IS NULL OR expired_at > $1)
  AND mission_type IN ('carbon_reduction', 'streak', 'activity', 'custom')
	`
	rows, err := r.db.QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission
	for rows.Next() {
		m := &model.Mission{}
		var criteriaType sql.NullString
		var badgeID sql.NullInt64

		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.MissionType, &criteriaType,
			&m.PointsReward, &m.GivesBadge, &badgeID,
			&m.TargetValue, &m.CreatedAt, &m.ExpiredAt,
		); err != nil {
			return nil, err
		}

		// Handle nullable fields
		if criteriaType.Valid {
			m.CriteriaType = model.MissionCriteriaType(criteriaType.String)
		}
		if badgeID.Valid {
			m.BadgeID = badgeID
		}

		missions = append(missions, m)
	}
	return missions, nil
}

func (r *checkMissionRepository) FindMissionsByType(ctx context.Context, missionType model.MissionType) ([]*model.Mission, error) {
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward, 
		       gives_badge, badge_id, target_value, created_at, expired_at
		FROM missions
		WHERE mission_type = $1 AND (expired_at IS NULL OR expired_at > $2)
	`
	rows, err := r.db.QueryContext(ctx, query, missionType, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission
	for rows.Next() {
		m := &model.Mission{}
		var criteriaType sql.NullString
		var badgeID sql.NullInt64

		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.MissionType, &criteriaType,
			&m.PointsReward, &m.GivesBadge, &badgeID,
			&m.TargetValue, &m.CreatedAt, &m.ExpiredAt,
		); err != nil {
			return nil, err
		}

		if criteriaType.Valid {
			m.CriteriaType = model.MissionCriteriaType(criteriaType.String)
		}
		if badgeID.Valid {
			m.BadgeID = badgeID
		}

		missions = append(missions, m)
	}
	return missions, nil
}

func (r *checkMissionRepository) FindMissionsByCriteriaType(ctx context.Context, criteriaType model.MissionCriteriaType) ([]*model.Mission, error) {
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward, 
		       gives_badge, badge_id, target_value, created_at, expired_at
		FROM missions
		WHERE criteria_type = $1 AND (expired_at IS NULL OR expired_at > $2)
	`
	rows, err := r.db.QueryContext(ctx, query, criteriaType, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission
	for rows.Next() {
		m := &model.Mission{}
		var badgeID sql.NullInt64

		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.MissionType, &m.CriteriaType,
			&m.PointsReward, &m.GivesBadge, &badgeID,
			&m.TargetValue, &m.CreatedAt, &m.ExpiredAt,
		); err != nil {
			return nil, err
		}

		if badgeID.Valid {
			m.BadgeID = badgeID
		}

		missions = append(missions, m)
	}
	return missions, nil
}

func (r *checkMissionRepository) AssignMissionToUser(ctx context.Context, userID, missionID int64) error {
	query := `
		INSERT INTO user_missions(user_id, mission_id, created_at) 
		VALUES ($1, $2, $3)	
		ON CONFLICT (user_id, mission_id) 
		DO UPDATE SET created_at = EXCLUDED.created_at
	`
	_, err := r.db.ExecContext(ctx, query, userID, missionID, time.Now())
	return err
}

func (r *checkMissionRepository) GetMissionProgress(ctx context.Context, userID, missionID int64) (float64, error) {
	var progress float64
	query := `SELECT progress_value FROM user_mission_progress WHERE user_id = $1 AND mission_id = $2`
	err := r.db.QueryRowContext(ctx, query, userID, missionID).Scan(&progress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return progress, nil
}

func (r *checkMissionRepository) UpdateMissionProgress(ctx context.Context, userID, missionID int64, progress float64) error {
	query := `
		INSERT INTO user_mission_progress(user_id, mission_id, progress_value, last_updated)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, mission_id) 
		DO UPDATE SET progress_value = EXCLUDED.progress_value, last_updated = EXCLUDED.last_updated
	`
	_, err := r.db.ExecContext(ctx, query, userID, missionID, progress, time.Now())
	return err
}

func (r *checkMissionRepository) MarkMissionCompleted(ctx context.Context, userID, missionID int64) error {
	query := `UPDATE user_missions SET completed_at = $1 WHERE user_id = $2 AND mission_id = $3`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID, missionID)
	return err
}

func (r *checkMissionRepository) HasUserCompletedMission(ctx context.Context, userID, missionID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_missions WHERE user_id = $1 AND mission_id = $2 AND completed_at IS NOT NULL)`
	err := r.db.QueryRowContext(ctx, query, userID, missionID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// =========================
// Mission Progress Calculation Functions
// =========================
func (r *checkMissionRepository) calculateMissionProgress(ctx context.Context, userID int64, mission *model.Mission) (float64, error) {
	switch mission.MissionType {
	case model.MissionTypeStreak:
		return r.calculateLoginStreakProgress(ctx, userID)
	case model.MissionTypeCarbonReduction:
		return r.calculateCarbonReductionProgress(ctx, userID, mission.CriteriaType)
	case model.MissionTypeActivity:
		return r.calculateActivityCountProgress(ctx, userID, mission.CriteriaType)
	case model.MissionTypeCustom:
		return r.calculateCustomMissionProgress(ctx, userID, mission.CriteriaType)
	default:
		return 0, fmt.Errorf("unknown mission type: %s", mission.MissionType)
	}
}

func (r *checkMissionRepository) calculateLoginStreakProgress(ctx context.Context, userID int64) (float64, error) {
	var streak float64
	query := `
		SELECT COUNT(DISTINCT DATE(created_at)) 
		FROM activity_logs
		WHERE user_id = $1 AND activity LIKE '%login%'
		AND created_at >= NOW() - INTERVAL '30 days'
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&streak)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return streak, nil
}

func (r *checkMissionRepository) calculateCarbonReductionProgress(ctx context.Context, userID int64, criteriaType model.MissionCriteriaType) (float64, error) {
	var totalCarbon float64

	if criteriaType == "" {
		query := `
			SELECT COALESCE(SUM(carbon_emission_g), 0) 
			FROM (
				SELECT carbon_emission_g FROM carbon_vehicle_logs cvl
				JOIN carbon_vehicles cv ON cvl.vehicle_id = cv.id
				WHERE cv.user_id = $1
				UNION ALL
				SELECT carbon_emission_g FROM carbon_electronics_logs cel
				JOIN carbon_electronics ce ON cel.device_id = ce.id
				WHERE ce.user_id = $2
			) AS emissions
		`
		err := r.db.QueryRowContext(ctx, query, userID, userID).Scan(&totalCarbon)
		if err != nil {
			return 0, err
		}
	} else {
		// Carbon reduction berdasarkan criteria type
		switch criteriaType {
		case model.CriteriaCar, model.CriteriaMotorcycle, model.CriteriaBicycle,
			model.CriteriaPublicTransport, model.CriteriaWalk:
			// Carbon dari kendaraan tertentu
			query := `
				SELECT COALESCE(SUM(carbon_emission_g), 0) 
				FROM carbon_vehicle_logs cvl
				JOIN carbon_vehicles cv ON cvl.vehicle_id = cv.id
				WHERE cv.user_id = $1 AND cv.vehicle_type = $2
			`
			err := r.db.QueryRowContext(ctx, query, userID, string(criteriaType)).Scan(&totalCarbon)
			if err != nil {
				return 0, err
			}

		case model.CriteriaLaptop, model.CriteriaDesktop, model.CriteriaTV,
			model.CriteriaAC, model.CriteriaFridge, model.CriteriaFan,
			model.CriteriaWashingMachine, model.CriteriaOther:
			// Carbon dari elektronik tertentu
			query := `
				SELECT COALESCE(SUM(carbon_emission_g), 0) 
				FROM carbon_electronics_logs cel
				JOIN carbon_electronics ce ON cel.device_id = ce.id
				WHERE ce.user_id = $1 AND ce.device_type = $2
			`
			err := r.db.QueryRowContext(ctx, query, userID, string(criteriaType)).Scan(&totalCarbon)
			if err != nil {
				return 0, err
			}
		}
	}

	return totalCarbon, nil
}

func (r *checkMissionRepository) calculateActivityCountProgress(ctx context.Context, userID int64, criteriaType model.MissionCriteriaType) (float64, error) {
	var count float64

	if criteriaType == "" {
		// Total semua aktivitas
		query := `SELECT COUNT(*) FROM activity_logs WHERE user_id = $1`
		err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		// Aktivitas berdasarkan criteria type
		query := `SELECT COUNT(*) FROM activity_logs WHERE user_id = $1 AND activity = $2`
		err := r.db.QueryRowContext(ctx, query, userID, string(criteriaType)).Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (r *checkMissionRepository) calculateCustomMissionProgress(ctx context.Context, userID int64, criteriaType model.MissionCriteriaType) (float64, error) {
	switch {
	case criteriaType == model.CriteriaCar || criteriaType == model.CriteriaMotorcycle ||
		criteriaType == model.CriteriaBicycle || criteriaType == model.CriteriaPublicTransport ||
		criteriaType == model.CriteriaWalk:
		// Hitung jarak kendaraan tertentu
		var totalDistance float64
		query := `
			SELECT COALESCE(SUM(distance_km), 0) 
			FROM carbon_vehicle_logs cvl
			JOIN carbon_vehicles cv ON cvl.vehicle_id = cv.id
			WHERE cv.user_id = $1 AND cv.vehicle_type = $2
		`
		err := r.db.QueryRowContext(ctx, query, userID, string(criteriaType)).Scan(&totalDistance)
		if err != nil {
			return 0, err
		}
		return totalDistance, nil

	case criteriaType == model.CriteriaLaptop || criteriaType == model.CriteriaDesktop ||
		criteriaType == model.CriteriaTV || criteriaType == model.CriteriaAC ||
		criteriaType == model.CriteriaFridge || criteriaType == model.CriteriaFan ||
		criteriaType == model.CriteriaWashingMachine || criteriaType == model.CriteriaOther:
		// Hitung jam penggunaan elektronik tertentu
		var totalHours float64
		query := `
			SELECT COALESCE(SUM(duration_hours), 0) 
			FROM carbon_electronics_logs cel
			JOIN carbon_electronics ce ON cel.device_id = ce.id
			WHERE ce.user_id = $1 AND ce.device_type = $2
		`
		err := r.db.QueryRowContext(ctx, query, userID, string(criteriaType)).Scan(&totalHours)
		if err != nil {
			return 0, err
		}
		return totalHours, nil

	default:
		// Default: hitung points earned
		var pointsEarned float64
		query := `
			SELECT COALESCE(SUM(amount), 0) 
			FROM point_transactions 
			WHERE user_id = $1 AND direction = 'in'
		`
		err := r.db.QueryRowContext(ctx, query, userID).Scan(&pointsEarned)
		if err != nil {
			return 0, err
		}
		return pointsEarned, nil
	}
}

// =========================
// Main Mission Checking Logic
// =========================
func (r *checkMissionRepository) CheckMission(ctx context.Context, userID int64, mission *model.Mission) (bool, error) {
	// Hitung progress berdasarkan mission type dan criteria type
	progress, err := r.calculateMissionProgress(ctx, userID, mission)
	if err != nil {
		return false, err
	}

	// Update progress di user_mission_progress
	if err := r.UpdateMissionProgress(ctx, userID, mission.ID, progress); err != nil {
		return false, err
	}

	// Check jika mission completed
	if progress >= mission.TargetValue {
		completed, err := r.HasUserCompletedMission(ctx, userID, mission.ID)
		if err != nil {
			return false, err
		}

		if !completed {
			return r.completeMission(ctx, userID, mission)
		}
		return true, nil
	}

	return false, nil
}

func (r *checkMissionRepository) completeMission(ctx context.Context, userID int64, mission *model.Mission) (bool, error) {
	// Assign mission ke user
	if err := r.AssignMissionToUser(ctx, userID, mission.ID); err != nil {
		return false, err
	}

	// Mark as completed
	if err := r.MarkMissionCompleted(ctx, userID, mission.ID); err != nil {
		return false, err
	}

	// Give points reward
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO points (user_id, total_points, created_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) 
		DO UPDATE SET total_points = points.total_points + $4
	`, userID, mission.PointsReward, time.Now(), mission.PointsReward)
	if err != nil {
		return false, err
	}

	// Record point transaction
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO point_transactions (user_id, amount, direction, source, reference_type, reference_id, created_at)
		VALUES ($1, $2, 'in', 'mission_completion', 'missions', $3, $4)
	`, userID, mission.PointsReward, mission.ID, time.Now())
	if err != nil {
		return false, err
	}

	// Give badge jika applicable
	if mission.GivesBadge && mission.BadgeID.Valid {
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO user_badges (user_id, badge_id, redeemed_at, created_at) 
			VALUES ($1, $2, NULL, $3)
			ON CONFLICT (user_id, badge_id) 
			DO NOTHING
		`, userID, mission.BadgeID.Int64, time.Now())
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *checkMissionRepository) CheckAllUserMissions(ctx context.Context, userID int64) error {
	// Dapatkan semua mission aktif
	missions, err := r.FindActiveMissions(ctx)
	if err != nil {
		return err
	}

	// Check setiap mission
	for _, mission := range missions {
		completed, err := r.CheckMission(ctx, userID, mission)
		if err != nil {
			fmt.Printf("Gagal check mission %d: %v\n", mission.ID, err)
			continue
		}

		if completed {
			fmt.Printf("User %d complaeted mission: %s\n", userID, mission.Title)
		}
	}

	return nil
}

func (r *checkMissionRepository) CheckUserMissionsByType(ctx context.Context, userID int64, missionType model.MissionType) error {
	missions, err := r.FindMissionsByType(ctx, missionType)
	if err != nil {
		return err
	}

	for _, mission := range missions {
		completed, err := r.CheckMission(ctx, userID, mission)
		if err != nil {
			fmt.Printf("Gagal check missison %d: %v\n", mission.ID, err)
			continue
		}

		if completed {
			fmt.Printf("User %d completed mission: %s\n", userID, mission.Title)
		}
	}

	return nil
}

func (r *checkMissionRepository) CheckUserMissionsByCriteriaType(ctx context.Context, userID int64, criteriaType model.MissionCriteriaType) error {
	missions, err := r.FindMissionsByCriteriaType(ctx, criteriaType)
	if err != nil {
		return err
	}

	for _, mission := range missions {
		completed, err := r.CheckMission(ctx, userID, mission)
		if err != nil {
			fmt.Printf("Gagal check misssion %d: %v\n", mission.ID, err)
			continue
		}

		if completed {
			fmt.Printf("User %d completed mission: %s\n", userID, mission.Title)	
		}
	}

	return nil
}
