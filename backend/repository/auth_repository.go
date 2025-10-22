package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
	models "github.com/Qodarrz/fiber-app/model"
)

type UserRepositoryInterface interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, userID int64) (*model.User, error)
	Create(ctx context.Context, user *model.User, profile *model.UserProfile) error
	Update(ctx context.Context, user *model.User) error
	VerifyEmailByToken(ctx context.Context, userID int64) error
	SaveResetPasswordToken(ctx context.Context, userID int64, token string) error
	FindByResetPasswordToken(ctx context.Context, token string) (*model.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	UpdateGoogleID(ctx context.Context, userID int64, googleID string) error
	CreateOrUpdateWithOAuth(ctx context.Context, user *model.User, profile *model.UserProfile) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, role, created_at
	          FROM users WHERE email = $1 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, userID int64) (*models.User, error) {
	user := &models.User{}

	query := `
        SELECT u.id, u.username, u.email, u.password, u.role, 
               u.google_id, u.email_verified_at, u.created_at,
               p.id AS profile_id, p.user_id, p.full_name, p.avatar_url, 
               p.birthdate, p.gender, p.created_at AS profile_created_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.id = $1;
    `

	var (
		googleID         sql.NullString
		emailVerifiedAt  sql.NullTime
		profileID        sql.NullInt64
		profileUserID    sql.NullInt64
		fullName         sql.NullString
		avatarURL        sql.NullString
		birthdate        sql.NullTime
		gender           sql.NullString
		profileCreatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &googleID, &emailVerifiedAt, &user.CreatedAt,
		&profileID, &profileUserID, &fullName, &avatarURL,
		&birthdate, &gender, &profileCreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Handle nullable fields
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}

	// Jika profile exists
	if profileID.Valid {
		user.Profile = &models.UserProfile{
			ID:        profileID.Int64,
			UserID:    profileUserID.Int64,
			FullName:  getStringPtrFromNullString(fullName),
			AvatarURL: getStringPtrFromNullString(avatarURL),
			Birthdate: getTimePtrFromNullTime(birthdate),
			Gender:    getStringPtrFromNullString(gender),
			CreatedAt: profileCreatedAt.Time,
		}
	}

	return user, nil
}

// Helper functions untuk nullable fields
func getStringPtrFromNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func getTimePtrFromNullTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
func (r *userRepository) Create(ctx context.Context, user *model.User, profile *model.UserProfile) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userQuery := `
		INSERT INTO users (username, email, password, role, google_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, userQuery,
		user.Username, user.Email, user.Password, user.Role, user.GoogleID, user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	// Insert ke tabel user_profiles
	profileQuery := `
		INSERT INTO user_profiles (user_id, created_at)
		VALUES ($1, $2)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, profileQuery, user.ID, time.Now()).Scan(&profile.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET username = $1, email = $2, password = $3, role = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.Role, user.ID)
	return err
}

func (r *userRepository) VerifyEmailByToken(ctx context.Context, userID int64) error {
	query := `UPDATE users SET email_verified_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (r *userRepository) SaveResetPasswordToken(ctx context.Context, userID int64, token string) error {
	query := `INSERT INTO reset_password_token (user_id, token) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, userID, token)
	return err
}

func (r *userRepository) FindByResetPasswordToken(ctx context.Context, token string) (*model.User, error) {
	query := `SELECT u.id, u.username, u.email, u.password, u.role, u.created_at
			  FROM users u
			  JOIN reset_password_token r ON u.id = r.user_id
			  WHERE r.token = $1`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Di repository/user_repository.go
func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, role, google_id, created_at
	          FROM users WHERE google_id = $1 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.GoogleID, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateGoogleID(ctx context.Context, userID int64, googleID string) error {
	query := `UPDATE users SET google_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, googleID, userID)
	return err
}

func (r *userRepository) CreateOrUpdateWithOAuth(ctx context.Context, user *model.User, profile *model.UserProfile) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah user sudah ada berdasarkan email
	existingUser, err := r.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		// Update Google ID untuk user yang sudah ada
		updateQuery := `UPDATE users SET google_id = $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateQuery, user.GoogleID, existingUser.ID)
		if err != nil {
			return err
		}
		user.ID = existingUser.ID
	} else {
		// Buat user baru
		userQuery := `
			INSERT INTO users (username, email, password, role, google_id, email_verified_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, userQuery,
			user.Username, user.Email, user.Password, user.Role,
			user.GoogleID, time.Now(), user.CreatedAt, // Email langsung verified untuk OAuth
		).Scan(&user.ID)
		if err != nil {
			return err
		}

		// Buat profile baru
		profileQuery := `
			INSERT INTO user_profiles (user_id, full_name, avatar_url, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, profileQuery,
			user.ID, profile.FullName, profile.AvatarURL, time.Now(),
		).Scan(&profile.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
