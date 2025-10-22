// repository/user_badge.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
)

type UserBadgeRepositoryInterface interface {
	FindByUserID(ctx context.Context, userID int64) ([]*model.UserBadge, error)
	FindByUserAndBadgeID(ctx context.Context, userID, badgeID int64) (*model.UserBadge, error)
	AssignBadge(ctx context.Context, userID, badgeID int64) error
	RedeemBadge(ctx context.Context, userID, badgeID int64) error
	GetUserEarnedBadges(ctx context.Context, userID int64) ([]*model.UserBadge, error)
	GetUserRedeemedBadges(ctx context.Context, userID int64) ([]*model.UserBadge, error)
}

type userBadgeRepository struct {
	db *sql.DB
}

func NewUserBadgeRepository(db *sql.DB) UserBadgeRepositoryInterface {
	return &userBadgeRepository{db: db}
}

func (r *userBadgeRepository) FindByUserID(ctx context.Context, userID int64) ([]*model.UserBadge, error) {
	query := `
		SELECT ub.id, ub.user_id, ub.badge_id, ub.redeemed_at, ub.created_at,
		       b.name, b.image_url, b.description, b.required_points, b.created_at as badge_created_at
		FROM user_badges ub
		JOIN badges b ON ub.badge_id = b.id
		WHERE ub.user_id = ?
		ORDER BY ub.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges []*model.UserBadge
	for rows.Next() {
		userBadge := &model.UserBadge{Badge: &model.Badge{}}
		var redeemedAt sql.NullTime

		err := rows.Scan(
			&userBadge.ID, &userBadge.UserID, &userBadge.BadgeID, &redeemedAt, &userBadge.CreatedAt,
			&userBadge.Badge.Name, &userBadge.Badge.ImageURL, &userBadge.Badge.Description,
			&userBadge.Badge.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		userBadge.RedeemedAt = redeemedAt
		userBadge.Badge.ID = userBadge.BadgeID

		userBadges = append(userBadges, userBadge)
	}

	return userBadges, nil
}

func (r *userBadgeRepository) FindByUserAndBadgeID(ctx context.Context, userID, badgeID int64) (*model.UserBadge, error) {
	query := `
		SELECT id, user_id, badge_id, redeemed_at, created_at
		FROM user_badges
		WHERE user_id = ? AND badge_id = ?
	`

	userBadge := &model.UserBadge{}
	var redeemedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, badgeID).Scan(
		&userBadge.ID, &userBadge.UserID, &userBadge.BadgeID, &redeemedAt, &userBadge.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userBadge.RedeemedAt = redeemedAt
	return userBadge, nil
}

func (r *userBadgeRepository) AssignBadge(ctx context.Context, userID, badgeID int64) error {
	// Check if badge already assigned
	existing, err := r.FindByUserAndBadgeID(ctx, userID, badgeID)
	if err != nil {
		return err
	}

	if existing != nil {
		return errors.New("badge already assigned to user")
	}

	query := `
		INSERT INTO user_badges (user_id, badge_id, created_at)
		VALUES (?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query, userID, badgeID, time.Now())
	return err
}

func (r *userBadgeRepository) RedeemBadge(ctx context.Context, userID, badgeID int64) error {
	userBadge, err := r.FindByUserAndBadgeID(ctx, userID, badgeID)
	if err != nil {
		return err
	}

	if userBadge == nil {
		return errors.New("badge not assigned to user")
	}

	if userBadge.RedeemedAt.Valid {
		return errors.New("badge already redeemed")
	}

	query := `UPDATE user_badges SET redeemed_at = ? WHERE id = ?`
	_, err = r.db.ExecContext(ctx, query, time.Now(), userBadge.ID)
	return err
}

func (r *userBadgeRepository) GetUserEarnedBadges(ctx context.Context, userID int64) ([]*model.UserBadge, error) {
	query := `
		SELECT ub.id, ub.user_id, ub.badge_id, ub.redeemed_at, ub.created_at,
		       b.name, b.image_url, b.description, b.required_points, b.created_at as badge_created_at
		FROM user_badges ub
		JOIN badges b ON ub.badge_id = b.id
		WHERE ub.user_id = ?
		ORDER BY ub.created_at DESC
	`

	return r.queryUserBadges(ctx, query, userID)
}

func (r *userBadgeRepository) GetUserRedeemedBadges(ctx context.Context, userID int64) ([]*model.UserBadge, error) {
	query := `
		SELECT ub.id, ub.user_id, ub.badge_id, ub.redeemed_at, ub.created_at,
		       b.name, b.image_url, b.description, b.required_points, b.created_at as badge_created_at
		FROM user_badges ub
		JOIN badges b ON ub.badge_id = b.id
		WHERE ub.user_id = ? AND ub.redeemed_at IS NOT NULL
		ORDER BY ub.redeemed_at DESC
	`

	return r.queryUserBadges(ctx, query, userID)
}

func (r *userBadgeRepository) queryUserBadges(ctx context.Context, query string, userID int64) ([]*model.UserBadge, error) {
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges []*model.UserBadge
	for rows.Next() {
		userBadge := &model.UserBadge{Badge: &model.Badge{}}
		var redeemedAt sql.NullTime

		err := rows.Scan(
			&userBadge.ID, &userBadge.UserID, &userBadge.BadgeID, &redeemedAt, &userBadge.CreatedAt,
			&userBadge.Badge.Name, &userBadge.Badge.ImageURL, &userBadge.Badge.Description,
			&userBadge.Badge.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		userBadge.RedeemedAt = redeemedAt
		userBadge.Badge.ID = userBadge.BadgeID

		userBadges = append(userBadges, userBadge)
	}

	return userBadges, nil
}
