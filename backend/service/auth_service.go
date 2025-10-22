package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	dto "github.com/Qodarrz/fiber-app/dto"
	helpers "github.com/Qodarrz/fiber-app/helper"
	models "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
	"golang.org/x/crypto/bcrypt"
)

// Interface
type AuthServiceInterface interface {
	Register(ctx context.Context, req *dto.RegisterDTO) (*models.User, error)
	Login(ctx context.Context, req *dto.LoginDTO) (*models.User, string, error)
	GetProfile(ctx context.Context, userID int64) (*models.User, error)
	VerifyEmail(ctx context.Context, userID int64) error
	RequestResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	LoginWithGoogle(ctx context.Context, googleUser *dto.GoogleUserDTO) (*models.User, string, error)
	LinkGoogleAccount(ctx context.Context, userID int64, googleUser *dto.GoogleUserDTO) error
	UnlinkGoogleAccount(ctx context.Context, userID int64) error
}

// Implementasi
type AuthService struct {
	userRepo     repository.UserRepositoryInterface
	activityRepo repository.ActivityRepositoryInterface
	missionRepo  repository.CheckMissionRepositoryInterface
}

func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	activityRepo repository.ActivityRepositoryInterface,
	missionRepo repository.CheckMissionRepositoryInterface,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		activityRepo: activityRepo,
		missionRepo:  missionRepo,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterDTO) (*models.User, error) {
	exists, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if exists != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPass, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPass,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.GoogleID != nil {
		user.GoogleID = req.GoogleID
	}

	profile := &models.UserProfile{
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user, profile); err != nil {
		return nil, err
	}

	token, err := helpers.GenerateEmailVerificationToken(fmt.Sprint(user.ID), user.Email)
	if err != nil {
		return nil, err
	}

	if err := helpers.SendEmailVerification(user.Email, token); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginDTO) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, "", errors.New("email atau password salah")
	}

	if !checkPassword(user.Password, req.Password) {
		return nil, "", errors.New("email atau password salah")
	}

	token, err := helpers.GenerateJWT(fmt.Sprint(user.ID))
	if err != nil {
		return nil, "", err
	}

	msg := fmt.Sprintf("User %d login berhasil", user.ID)
	if err := s.activityRepo.LogActivity(ctx, user.ID, msg); err != nil {
		fmt.Printf("gagal simpan activity log: %v\n", err)
	}

	go func() {
		bgCtx := context.Background()

		if err := s.missionRepo.CheckAllUserMissions(bgCtx, user.ID); err != nil {
			fmt.Printf("Gagal check missions setelah login: %v\n", err)
		} else {
			fmt.Printf("Success check missions untuk user %d setelah login\n", user.ID)
		}
	}()

	return user, token, nil
}
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return user, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, userID int64) error {
	return s.userRepo.VerifyEmailByToken(ctx, userID)
}

func (s *AuthService) RequestResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("email tidak ditemukan")
	}

	token := helpers.GenerateRandomToken(7)
	err = s.userRepo.SaveResetPasswordToken(ctx, user.ID, token)
	if err != nil {
		return fmt.Errorf("gagal menyimpan token reset password: %w", err)
	}

	return helpers.SendTokenForgotEmail(email, token)
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.userRepo.FindByResetPasswordToken(ctx, token)
	if err != nil || user == nil {
		return errors.New("token tidak valid atau sudah kadaluarsa")
	}

	hashedPass, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPass
	return s.userRepo.Update(ctx, user)
}

// Update password
func (s *AuthService) UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}

	if !checkPassword(user.Password, oldPassword) {
		return errors.New("password lama salah")
	}

	hashedPass, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPass
	return s.userRepo.Update(ctx, user)
}

// Di service/auth_service.go
func (s *AuthService) LoginWithGoogle(ctx context.Context, googleUser *dto.GoogleUserDTO) (*models.User, string, error) {
	// Cek apakah user sudah ada berdasarkan Google ID
	user, err := s.userRepo.FindByGoogleID(ctx, googleUser.ID)
	if err != nil {
		return nil, "", err
	}

	// Jika tidak ditemukan berdasarkan Google ID, cari berdasarkan email
	if user == nil {
		user, err = s.userRepo.FindByEmail(ctx, googleUser.Email)
		if err != nil {
			return nil, "", err
		}

		if user != nil {
			// User ditemukan berdasarkan email, update Google ID
			err = s.userRepo.UpdateGoogleID(ctx, user.ID, googleUser.ID)
			if err != nil {
				return nil, "", err
			}
			user.GoogleID = &googleUser.ID
		} else {
			// Buat user baru untuk OAuth
			hashedPass, _ := hashPassword("") // Password kosong untuk OAuth user

			user = &models.User{
				Username:  googleUser.Name,
				Email:     googleUser.Email,
				Password:  hashedPass,
				Role:      "user",
				GoogleID:  &googleUser.ID,
				CreatedAt: time.Now(),
			}

			profile := &models.UserProfile{
				FullName:  &googleUser.Name,
				AvatarURL: &googleUser.Picture,
				CreatedAt: time.Now(),
			}

			if err := s.userRepo.CreateOrUpdateWithOAuth(ctx, user, profile); err != nil {
				return nil, "", err
			}
		}
	}

	// Generate JWT token
	token, err := helpers.GenerateJWT(fmt.Sprint(user.ID))
	if err != nil {
		return nil, "", err
	}

	// Log activity
	msg := fmt.Sprintf("User %d login dengan Google OAuth", user.ID)
	if err := s.activityRepo.LogActivity(ctx, user.ID, msg); err != nil {
		fmt.Printf("gagal simpan activity log: %v\n", err)
	}

	// Check missions
	go func() {
		bgCtx := context.Background()
		if err := s.missionRepo.CheckAllUserMissions(bgCtx, user.ID); err != nil {
			fmt.Printf("Gagal check missions setelah login OAuth: %v\n", err)
		} else {
			fmt.Printf("Success check missions untuk user %d setelah login OAuth\n", user.ID)
		}
	}()

	return user, token, nil
}

func (s *AuthService) LinkGoogleAccount(ctx context.Context, userID int64, googleUser *dto.GoogleUserDTO) error {
	// Cek apakah Google ID sudah terpakai oleh user lain
	existingUser, err := s.userRepo.FindByGoogleID(ctx, googleUser.ID)
	if err != nil {
		return err
	}

	if existingUser != nil && existingUser.ID != userID {
		return errors.New("akun Google ini sudah terhubung dengan user lain")
	}

	// Update Google ID untuk user
	return s.userRepo.UpdateGoogleID(ctx, userID, googleUser.ID)
}

func (s *AuthService) UnlinkGoogleAccount(ctx context.Context, userID int64) error {
	// Set Google ID menjadi NULL untuk user
	return s.userRepo.UpdateGoogleID(ctx, userID, "")
}
