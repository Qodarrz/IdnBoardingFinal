// repository/badge.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
)

type BadgeRepositoryInterface interface {
	FindByID(ctx context.Context, id int64) (*model.Badge, error)
	FindAll(ctx context.Context, page, limit int) ([]*model.Badge, error)
	FindAllWithOwnership(ctx context.Context, userID int64, page, limit int) ([]*model.BadgeWithOwnership, error)
	Create(ctx context.Context, badge *model.Badge) error
}

type badgeRepository struct {
	db *sql.DB
}

func NewBadgeRepository(db *sql.DB) BadgeRepositoryInterface {
	return &badgeRepository{db: db}
}

func (r *badgeRepository) FindByID(ctx context.Context, id int64) (*model.Badge, error) {
	query := `
		SELECT id, name, image_url, description, created_at
		FROM badges
		WHERE id = $1
	`

	badge := &model.Badge{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&badge.ID, &badge.Name, &badge.ImageURL, &badge.Description,
		&badge.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return badge, nil
}

func (r *badgeRepository) FindAll(ctx context.Context, page, limit int) ([]*model.Badge, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, name, image_url, description, created_at
		FROM badges
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []*model.Badge
	for rows.Next() {
		badge := &model.Badge{}
		err := rows.Scan(
			&badge.ID, &badge.Name, &badge.ImageURL, &badge.Description,
			&badge.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return badges, nil
}

func (r *badgeRepository) FindAllWithOwnership(ctx context.Context, userID int64, page, limit int) ([]*model.BadgeWithOwnership, error) {
	offset := (page - 1) * limit
	query := `
	SELECT 
    b.id,
    b.name,
    b.image_url,
    b.description,
    b.created_at,
    CASE WHEN ub.id IS NOT NULL THEN true ELSE false END AS is_owned,
    ub.redeemed_at
FROM badges b
LEFT JOIN user_badges ub 
    ON b.id = ub.badge_id AND ub.user_id = $1
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3

`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println("DEBUG >>> userID =", userID, "page =", page, "limit =", limit, "offset =", offset)


	var badges []*model.BadgeWithOwnership
	for rows.Next() {
		badge := &model.BadgeWithOwnership{}
		var redeemedAt sql.NullTime

		err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.ImageURL,
			&badge.Description,
			&badge.CreatedAt,
			&badge.IsOwned,
			&redeemedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable redeemed_at
		if redeemedAt.Valid {
			badge.RedeemedAt = &redeemedAt.Time
		}

		badges = append(badges, badge)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return badges, nil
}

func (r *badgeRepository) Create(ctx context.Context, badge *model.Badge) error {
	query := `
		INSERT INTO badges (name, image_url, description, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		badge.Name, badge.ImageURL, badge.Description, time.Now(),
	).Scan(&badge.ID)

	return err
}
