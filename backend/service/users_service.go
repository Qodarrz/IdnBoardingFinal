// service/user_profile_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	dto "github.com/Qodarrz/fiber-app/dto"
	models "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
)

type UserProfileServiceInterface interface {
	GetProfile(ctx context.Context, userID int64) (*dto.UserWithProfileResponseDTO, error)
	UpdateProfile(ctx context.Context, userID int64, req *dto.UserProfileUpdateDTO) (*dto.UserWithProfileResponseDTO, error)
}

type UserProfileService struct {
	userRepo        repository.UserRepositoryInterface
	userProfileRepo repository.UserProfileRepositoryInterface
	activityRepo    repository.ActivityRepositoryInterface
}

func NewUserProfileService(
	userProfileRepo repository.UserProfileRepositoryInterface,
	activityRepo repository.ActivityRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *UserProfileService {
	return &UserProfileService{
		userProfileRepo: userProfileRepo,
		activityRepo:    activityRepo,
		userRepo:       userRepo,
	}
}

func (s *UserProfileService) GetProfile(ctx context.Context, userID int64) (*dto.UserWithProfileResponseDTO, error) {
	user, profile, err := s.userProfileRepo.FindByUID(ctx, userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if profile == nil {
		profile = &models.UserProfile{
			UserID: userID,
		}
		if err := s.userProfileRepo.Create(ctx, profile); err != nil {
			return nil, errors.New("gagal membuat profile")
		}
	}

	response := &dto.UserWithProfileResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Profile: dto.UserProfileResponseDTO{
			ID:        profile.ID,
			UserID:    profile.UserID,
			FullName:  profile.FullName,
			AvatarURL: profile.AvatarURL,
			Birthdate: profile.Birthdate,
			Gender:    profile.Gender,
			CreatedAt: profile.CreatedAt,
		},
	}

	return response, nil

}

func (s *UserProfileService) UpdateProfile(ctx context.Context, userID int64, req *dto.UserProfileUpdateDTO) (*dto.UserWithProfileResponseDTO, error) {
	log.Printf("UpdateProfile called with userID: %d, req: %+v", userID, req)

	// Ambil user dari table users
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found") // or use a custom error type
	}

	// Now safely access user fields
	// Your existing code that uses user.Profile, etc.
	// Ambil profile dari table user_profiles
	profile, err := s.userProfileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data profile: %w", err)
	}

	// Jika profile belum ada, buat baru
	if profile == nil {
		profile = &models.UserProfile{
			UserID: userID,
		}
		if err := s.userProfileRepo.Create(ctx, profile); err != nil {
			return nil, errors.New("gagal membuat profile")
		}
	}

	// Update fields profile yang di-provide (hanya jika tidak nil)
	if req.FullName != nil {
		profile.FullName = req.FullName // Dereference jika perlu
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = req.AvatarURL // Dereference jika perlu
	}
	if req.Birthdate != nil {
		profile.Birthdate = req.Birthdate // Dereference jika perlu
	}
	if req.Gender != nil {
		profile.Gender = req.Gender // Dereference jika perlu
	}

	// Update profile di repo
	if err := s.userProfileRepo.Update(ctx, profile); err != nil {
		return nil, errors.New("gagal mengupdate profile")
	}

	// Log activity
	if s.activityRepo != nil {
		s.activityRepo.LogActivity(ctx, userID, "User update profile")
	}

	// Build response
	response := &dto.UserWithProfileResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Profile: dto.UserProfileResponseDTO{
			ID:        profile.ID,
			UserID:    profile.UserID,
			FullName:  profile.FullName,
			AvatarURL: profile.AvatarURL,
			Birthdate: profile.Birthdate,
			Gender:    profile.Gender,
			CreatedAt: profile.CreatedAt,
		},
	}

	return response, nil
}
