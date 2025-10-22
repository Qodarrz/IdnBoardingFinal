package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type ValidationRepositoryInterface interface {
	IsUnique(ctx context.Context, table, column, value string, id uint32) (bool, error)
}

type validationRepository struct {
	db *sql.DB
}

func InitValidationRepository(db *sql.DB) ValidationRepositoryInterface {
	return &validationRepository{db: db}
}

func (r *validationRepository) IsUnique(ctx context.Context, table, column, value string, id uint32) (bool, error) {
	var query string
	var err error
	var result int8

	if id != 0 {
		query = fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? AND id != ? LIMIT 1", table, column)
		err = r.db.QueryRowContext(ctx, query, value, id).Scan(&result)
	} else {
		query = fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? LIMIT 1", table, column)
		err = r.db.QueryRowContext(ctx, query, value).Scan(&result)
	}

	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, nil
}
