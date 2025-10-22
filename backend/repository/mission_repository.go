// repository/mission.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Qodarrz/fiber-app/dto"
	model "github.com/Qodarrz/fiber-app/model"
)

type MissionRepositoryInterface interface {
	Create(ctx context.Context, mission *model.Mission) error
	FindByID(ctx context.Context, id int64) (*model.Mission, error)
	FindAll(ctx context.Context, page, limit int) ([]*model.Mission, error)
	FindActiveMissions(ctx context.Context) ([]*model.Mission, error)
	FindUserMissions(ctx context.Context, userID int64) ([]*model.UserMission, error)
	CreateUserMission(ctx context.Context, userMission *model.UserMission) error
	UpdateUserMission(ctx context.Context, userMission *model.UserMission) error
	FindUserMissionByID(ctx context.Context, userID, missionID int64) (*model.UserMission, error)
	GetAllMissionProgress(ctx context.Context, userID int64) ([]dto.MissionProgressResponse, error)
}

type missionRepository struct {
	db *sql.DB
}

func NewMissionRepository(db *sql.DB) MissionRepositoryInterface {
	return &missionRepository{db: db}
}

// repository/mission.go (bagian Create)
func (r *missionRepository) Create(ctx context.Context, mission *model.Mission) error {
	query := `
		INSERT INTO missions 
		    (title, description, mission_type, criteria_type, points_reward, 
		     gives_badge, badge_id, target_value, expired_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var badgeID interface{}
	if mission.BadgeID.Valid {
		badgeID = mission.BadgeID.Int64
	} else {
		badgeID = nil
	}

	var criteriaType interface{}
	if mission.CriteriaType != "" {
		criteriaType = string(mission.CriteriaType)
	} else {
		criteriaType = nil
	}

	var expiredAt interface{}
	if mission.ExpiredAt.Valid {
		expiredAt = mission.ExpiredAt.Time
	} else {
		expiredAt = nil
	}

	// langsung QueryRowContext tanpa prepare
	return r.db.QueryRowContext(ctx, query,
		mission.Title, mission.Description, mission.MissionType, criteriaType,
		mission.PointsReward, mission.GivesBadge, badgeID, mission.TargetValue,
		expiredAt, mission.CreatedAt,
	).Scan(&mission.ID)
}

func (r *missionRepository) FindByID(ctx context.Context, id int64) (*model.Mission, error) {
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward, gives_badge,
		       badge_id, target_value, expired_at, created_at
		FROM missions
		WHERE id = $1
	`

	mission := &model.Mission{}
	var badgeID sql.NullInt64
	var criteriaType sql.NullString
	var expiredAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&mission.ID, &mission.Title, &mission.Description, &mission.MissionType,
		&criteriaType, &mission.PointsReward, &mission.GivesBadge, &badgeID,
		&mission.TargetValue, &expiredAt, &mission.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	mission.BadgeID = badgeID
	mission.ExpiredAt = expiredAt

	if criteriaType.Valid {
		mission.CriteriaType = model.MissionCriteriaType(criteriaType.String)
	}

	return mission, nil
}

func (r *missionRepository) FindAll(ctx context.Context, page, limit int) ([]*model.Mission, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward, gives_badge,
		       badge_id, target_value, expired_at, created_at
		FROM missions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission
	for rows.Next() {
		mission := &model.Mission{}
		var badgeID sql.NullInt64
		var criteriaType sql.NullString
		var expiredAt sql.NullTime

		err := rows.Scan(
			&mission.ID, &mission.Title, &mission.Description, &mission.MissionType,
			&criteriaType, &mission.PointsReward, &mission.GivesBadge, &badgeID,
			&mission.TargetValue, &expiredAt, &mission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		mission.BadgeID = badgeID
		mission.ExpiredAt = expiredAt

		if criteriaType.Valid {
			mission.CriteriaType = model.MissionCriteriaType(criteriaType.String)
		}

		missions = append(missions, mission)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return missions, nil
}

func (r *missionRepository) FindActiveMissions(ctx context.Context) ([]*model.Mission, error) {
	query := `
		SELECT id, title, description, mission_type, criteria_type, points_reward, gives_badge,
		       badge_id, target_value, expired_at, created_at
		FROM missions
		WHERE expired_at IS NULL OR expired_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission
	for rows.Next() {
		mission := &model.Mission{}
		var badgeID sql.NullInt64
		var criteriaType sql.NullString
		var expiredAt sql.NullTime

		err := rows.Scan(
			&mission.ID, &mission.Title, &mission.Description, &mission.MissionType,
			&criteriaType, &mission.PointsReward, &mission.GivesBadge, &badgeID,
			&mission.TargetValue, &expiredAt, &mission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		mission.BadgeID = badgeID
		mission.ExpiredAt = expiredAt

		if criteriaType.Valid {
			mission.CriteriaType = model.MissionCriteriaType(criteriaType.String)
		}

		missions = append(missions, mission)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return missions, nil
}

func (r *missionRepository) FindUserMissions(ctx context.Context, userID int64) ([]*model.UserMission, error) {
	query := `
		SELECT um.id, um.user_id, um.mission_id, um.completed_at, um.created_at,
		       m.title, m.description, m.mission_type, m.criteria_type, m.points_reward, 
		       m.gives_badge, m.badge_id, m.target_value, m.expired_at, m.created_at as mission_created_at
		FROM user_missions um
		JOIN missions m ON um.mission_id = m.id
		WHERE um.user_id = $1
		ORDER BY um.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userMissions []*model.UserMission
	for rows.Next() {
		userMission := &model.UserMission{Mission: &model.Mission{}}
		var completedAt sql.NullTime
		var badgeID sql.NullInt64
		var criteriaType sql.NullString
		var expiredAt sql.NullTime

		err := rows.Scan(
			&userMission.ID, &userMission.UserID, &userMission.MissionID, &completedAt, &userMission.CreatedAt,
			&userMission.Mission.Title, &userMission.Mission.Description, &userMission.Mission.MissionType,
			&criteriaType, &userMission.Mission.PointsReward, &userMission.Mission.GivesBadge, &badgeID,
			&userMission.Mission.TargetValue, &expiredAt, &userMission.Mission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		userMission.CompletedAt = completedAt
		userMission.Mission.ID = userMission.MissionID
		userMission.Mission.BadgeID = badgeID
		userMission.Mission.ExpiredAt = expiredAt

		if criteriaType.Valid {
			userMission.Mission.CriteriaType = model.MissionCriteriaType(criteriaType.String)
		}

		userMissions = append(userMissions, userMission)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userMissions, nil
}

func (r *missionRepository) CreateUserMission(ctx context.Context, userMission *model.UserMission) error {
	query := `
		INSERT INTO user_missions (user_id, mission_id, completed_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var completedAt interface{}
	if userMission.CompletedAt.Valid {
		completedAt = userMission.CompletedAt.Time
	} else {
		completedAt = nil
	}

	err := r.db.QueryRowContext(ctx, query,
		userMission.UserID, userMission.MissionID, completedAt, userMission.CreatedAt,
	).Scan(&userMission.ID)

	return err
}

func (r *missionRepository) UpdateUserMission(ctx context.Context, userMission *model.UserMission) error {
	query := `
		UPDATE user_missions
		SET completed_at = $1
		WHERE id = $2
	`

	var completedAt interface{}
	if userMission.CompletedAt.Valid {
		completedAt = userMission.CompletedAt.Time
	} else {
		completedAt = nil
	}

	_, err := r.db.ExecContext(ctx, query, completedAt, userMission.ID)
	return err
}

func (r *missionRepository) FindUserMissionByID(ctx context.Context, userID, missionID int64) (*model.UserMission, error) {
	query := `
		SELECT id, user_id, mission_id, completed_at, created_at
		FROM user_missions
		WHERE user_id = $1 AND mission_id = $2
	`

	userMission := &model.UserMission{}
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, missionID).Scan(
		&userMission.ID, &userMission.UserID, &userMission.MissionID,
		&completedAt, &userMission.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userMission.CompletedAt = completedAt
	return userMission, nil
}

func (r *missionRepository) GetAllMissionProgress(ctx context.Context, userID int64) ([]dto.MissionProgressResponse, error) {
	query := `
		SELECT m.id, m.title, m.description, m.target_value, m.points_reward, m.gives_badge, m.badge_id,
		       COALESCE(ump.progress_value, 0) AS progress_value,
		       CASE 
		          WHEN ump.progress_value >= m.target_value AND m.target_value IS NOT NULL THEN true
		          ELSE false
		       END AS is_completed
		FROM missions m
		LEFT JOIN user_mission_progress ump ON m.id = ump.mission_id AND ump.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.MissionProgressResponse
	for rows.Next() {
		var res dto.MissionProgressResponse
		if err := rows.Scan(
			&res.MissionID,
			&res.Title,
			&res.Description,
			&res.TargetValue,
			&res.PointsReward,
			&res.GivesBadge,
			&res.BadgeID,
			&res.Progress,
			&res.IsCompleted,
		); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}
