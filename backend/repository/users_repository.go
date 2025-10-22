// repository/user_profile_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	model "github.com/Qodarrz/fiber-app/model"
)

type UserProfileRepositoryInterface interface {
	FindByUserID(ctx context.Context, userID int64) (*model.UserProfile, error)
	FindByUID(ctx context.Context, userID int64) (*model.User, *model.UserProfile, error)
	Update(ctx context.Context, profile *model.UserProfile) error
	Create(ctx context.Context, profile *model.UserProfile) error
}

type userProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) UserProfileRepositoryInterface {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID int64) (*model.UserProfile, error) {
	profile := &model.UserProfile{}
	query := `
		SELECT id, user_id, full_name, avatar_url, birthdate, gender, created_at
		FROM user_profiles 
		WHERE user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.AvatarURL,
		&profile.Birthdate,
		&profile.Gender,
		&profile.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}

func strOrNil(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func timeOrNil(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
func (r *userProfileRepository) Update(ctx context.Context, profile *model.UserProfile) error {
	if profile == nil {
		return errors.New("profile is nil")
	}

	log.Printf("Repository Update called with profile: %+v", profile)

	query := `
	UPDATE user_profiles
	SET
		full_name = COALESCE($1, full_name),
		avatar_url = COALESCE($2, avatar_url),
		birthdate = COALESCE($3, birthdate),
		gender = COALESCE($4, gender)
	WHERE user_id = $5
	RETURNING id, created_at
	`

	// Helper functions untuk handle nil pointers
	strOrNil := func(s *string) interface{} {
		if s == nil {
			return nil
		}
		return *s
	}

	timeOrNil := func(t *time.Time) interface{} {
		if t == nil {
			return nil
		}
		return *t
	}

	err := r.db.QueryRowContext(ctx, query,
		strOrNil(profile.FullName),
		strOrNil(profile.AvatarURL),
		timeOrNil(profile.Birthdate),
		strOrNil(profile.Gender),
		profile.UserID,
	).Scan(&profile.ID, &profile.CreatedAt)

	if err != nil {
		log.Printf("Repository Update error: %v", err)
		return fmt.Errorf("failed to update profile: %w", err)
	}

	log.Printf("Repository Update successful: ID=%d, CreatedAt=%v", profile.ID, profile.CreatedAt)
	return nil
}

func (r *userProfileRepository) Create(ctx context.Context, profile *model.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, full_name, avatar_url, birthdate, gender, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		profile.UserID,
		profile.FullName,
		profile.AvatarURL,
		profile.Birthdate,
		profile.Gender,
		time.Now(),
	).Scan(&profile.ID)
}

func (r *userProfileRepository) FindByUID(ctx context.Context, userID int64) (*model.User, *model.UserProfile, error) {
	user := &model.User{}
	profile := &model.UserProfile{}

	query := `
        SELECT u.id, u.username, u.email, u.role,
               p.id, p.user_id, p.full_name, p.avatar_url, p.birthdate, p.gender, p.created_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.id = $1
    `
	var (
		profileID   sql.NullInt64
		profileUser sql.NullInt64
		fullName    sql.NullString
		avatarURL   sql.NullString
		birthdate   sql.NullTime
		gender      sql.NullString
		createdAt   sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role,
		&profileID, &profileUser, &fullName, &avatarURL, &birthdate, &gender, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	// Jika profileID NULL, profil tidak ada
	if !profileID.Valid {
		return user, nil, nil
	}

	// Assign ke struct pointer
	profile.ID = profileID.Int64
	profile.UserID = profileUser.Int64
	if fullName.Valid {
		profile.FullName = &fullName.String
	}
	if avatarURL.Valid {
		profile.AvatarURL = &avatarURL.String
	}
	if birthdate.Valid {
		profile.Birthdate = &birthdate.Time
	}
	if gender.Valid {
		profile.Gender = &gender.String
	}
	if createdAt.Valid {
		profile.CreatedAt = createdAt.Time
	}

	return user, profile, nil
}
