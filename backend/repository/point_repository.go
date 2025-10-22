// repository/points.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
)

type PointsRepositoryInterface interface {
	GetUserPoints(ctx context.Context, userID int64) (*model.Points, error)
	AddPoints(ctx context.Context, userID int64, amount int, source string, referenceID int64) error
	DeductPoints(ctx context.Context, userID int64, amount int, source string, referenceID int64) error
}

type pointsRepository struct {
	db *sql.DB
}

func NewPointsRepository(db *sql.DB) PointsRepositoryInterface {
	return &pointsRepository{db: db}
}

func (r *pointsRepository) GetUserPoints(ctx context.Context, userID int64) (*model.Points, error) {
	points := &model.Points{}
	query := `SELECT id, user_id, total_points, created_at FROM points WHERE user_id = $1`
	
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&points.ID, &points.UserID, &points.TotalPoints, &points.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create new points record if it doesn't exist
			return r.createUserPoints(ctx, userID)
		}
		return nil, err
	}
	return points, nil
}

func (r *pointsRepository) createUserPoints(ctx context.Context, userID int64) (*model.Points, error) {
	points := &model.Points{
		UserID:      userID,
		TotalPoints: 0,
		CreatedAt:   time.Now(),
	}

	// Gunakan ON CONFLICT DO NOTHING agar kalau sudah ada, tidak insert lagi
	query := `
		INSERT INTO points (user_id, total_points, created_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, userID, points.TotalPoints, points.CreatedAt).Scan(&points.ID)
	if err != nil {
		// Kalau DO NOTHING maka RETURNING tidak jalan → harus ambil existing row
		if errors.Is(err, sql.ErrNoRows) {
			return r.GetUserPoints(ctx, userID)
		}
		return nil, err
	}

	return points, nil
}


func (r *pointsRepository) AddPoints(ctx context.Context, userID int64, amount int, source string, referenceID int64) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update points
	query := `UPDATE points SET total_points = total_points + $1 WHERE user_id = $2`
	_, err = tx.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	// Record transaction
	transactionQuery := `INSERT INTO point_transactions 
		(user_id, amount, direction, source, reference_type, reference_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = tx.ExecContext(ctx, transactionQuery, 
		userID, amount, "in", source, "order", referenceID, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *pointsRepository) DeductPoints(ctx context.Context, userID int64, amount int, source string, referenceID int64) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if user has enough points
	var currentPoints int
	checkQuery := `SELECT total_points FROM points WHERE user_id = $1`
	err = tx.QueryRowContext(ctx, checkQuery, userID).Scan(&currentPoints)
	if err != nil {
		return err
	}
	if currentPoints < amount {
		return errors.New("insufficient points")
	}

	// Update points
	updateQuery := `UPDATE points SET total_points = total_points - $1 WHERE user_id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, amount, userID)
	if err != nil {
		return err
	}

	// Record transaction
	transactionQuery := `INSERT INTO point_transactions 
		(user_id, amount, direction, source, reference_type, reference_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = tx.ExecContext(ctx, transactionQuery, 
		userID, amount, "out", source, "order", referenceID, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}