// service/mission.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto "github.com/Qodarrz/fiber-app/dto"
	model "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
)

type MissionServiceInterface interface {
	CreateMission(ctx context.Context, req *dto.CreateMissionDTO) (*dto.MissionResponseDTO, error)
	CreateMissionWithBadge(ctx context.Context, req *dto.CreateMissionWithBadgeDTO) (*dto.MissionWithBadgeResponseDTO, error)
	GetMissionByID(ctx context.Context, id int64) (*dto.MissionResponseDTO, error)
	GetAllMissions(ctx context.Context, page, limit int) ([]*dto.MissionResponseDTO, error)
	GetActiveMissions(ctx context.Context) ([]*dto.MissionResponseDTO, error)
	GetUserMissions(ctx context.Context, userID int64) ([]*dto.UserMissionResponseDTO, error)
	CheckMissionCompletion(ctx context.Context, userID, missionID int64) (bool, error)
}

type missionService struct {
	missionRepo     repository.MissionRepositoryInterface
	userMissionRepo repository.MissionRepositoryInterface
	badgeRepo       repository.BadgeRepositoryInterface
}


func NewMissionService(
	missionRepo repository.MissionRepositoryInterface,
	userMissionRepo repository.MissionRepositoryInterface,
	badgeRepo repository.BadgeRepositoryInterface,
) MissionServiceInterface {
	return &missionService{
		missionRepo:     missionRepo,
		userMissionRepo: userMissionRepo,
		badgeRepo:       badgeRepo,
	}		
}

func (s *missionService) CreateMissionWithBadge(ctx context.Context, req *dto.CreateMissionWithBadgeDTO) (*dto.MissionWithBadgeResponseDTO, error) {
	// Validasi: Jika gives_badge true, pastikan data badge lengkap
	if req.GivesBadge {
		if req.BadgeName == "" {
			return nil, errors.New("badge_name is required when gives_badge is true")
		}
		if req.BadgeImageURL == "" {
			return nil, errors.New("badge_image_url is required when gives_badge is true")
		}
	}

	var badgeID *int64
	var badgeResponse *dto.BadgeResponseDTO

	// Jika mission memberikan badge, buat badge terlebih dahulu
	if req.GivesBadge {
		// Buat badge
		badge := &model.Badge{
			Name:           req.BadgeName,
			ImageURL:       req.BadgeImageURL,
			Description:    req.BadgeDescription,
			CreatedAt:      time.Now(),
		}

		err := s.badgeRepo.Create(ctx, badge)
		if err != nil {
			return nil, fmt.Errorf("failed to create badge: %w", err)
		}

		badgeID = &badge.ID
		badgeResponse = s.badgeToDTO(badge)
	}

	mission := &model.Mission{
		Title:        req.Title,
		Description:  req.Description,
		MissionType:  model.MissionType(req.MissionType),
		CriteriaType: model.MissionCriteriaType(req.CriteriaType),
		PointsReward: req.PointsReward,
		GivesBadge:   req.GivesBadge,
		TargetValue:  req.TargetValue,
		CreatedAt:    time.Now(),
	}

	if badgeID != nil {
		mission.BadgeID = model.NewNullInt64(*badgeID)
	}

	if req.ExpiredAt != nil {
		mission.ExpiredAt = model.NewNullTime(*req.ExpiredAt)
	}

	err := s.missionRepo.Create(ctx, mission)
	if err != nil {
		return nil, fmt.Errorf("failed to create mission: %w", err)
	}

	response := &dto.MissionWithBadgeResponseDTO{
		Mission: *s.missionToDTO(mission),
	}

	if badgeResponse != nil {
		response.Badge = *badgeResponse
	}

	return response, nil
}

// Tambahkan helper function untuk badge
func (s *missionService) badgeToDTO(badge *model.Badge) *dto.BadgeResponseDTO {
	return &dto.BadgeResponseDTO{
		ID:             badge.ID,
		Name:           badge.Name,
		ImageURL:       badge.ImageURL,
		Description:    badge.Description,
		CreatedAt:      badge.CreatedAt,
	}
}

func (s *missionService) CreateMission(ctx context.Context, req *dto.CreateMissionDTO) (*dto.MissionResponseDTO, error) {

	mission := &model.Mission{
		Title:        req.Title,
		Description:  req.Description,
		MissionType:  model.MissionType(req.MissionType),
		PointsReward: req.PointsReward,
		GivesBadge:   req.GivesBadge,
		TargetValue:  req.TargetValue,
		CreatedAt:    time.Now(),
	}

	// Set optional fields
	if req.BadgeID != nil {
		mission.BadgeID = model.NewNullInt64(*req.BadgeID)
	}

	if req.CarbonReductionG != nil {
		mission.CarbonReductionG = model.NewNullFloat64(*req.CarbonReductionG)
	}

	if req.ExpiredAt != nil {
		mission.ExpiredAt = model.NewNullTime(*req.ExpiredAt)
	}

	err := s.missionRepo.Create(ctx, mission)
	if err != nil {
		return nil, err
	}

	return s.missionToDTO(mission), nil
}

func (s *missionService) GetMissionByID(ctx context.Context, id int64) (*dto.MissionResponseDTO, error) {
	mission, err := s.missionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if mission == nil {
		return nil, errors.New("mission not found")
	}

	return s.missionToDTO(mission), nil
}

func (s *missionService) GetAllMissions(ctx context.Context, page, limit int) ([]*dto.MissionResponseDTO, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	missions, err := s.missionRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	var result []*dto.MissionResponseDTO
	for _, mission := range missions {
		result = append(result, s.missionToDTO(mission))
	}

	return result, nil
}

func (s *missionService) GetActiveMissions(ctx context.Context) ([]*dto.MissionResponseDTO, error) {
	missions, err := s.missionRepo.FindActiveMissions(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.MissionResponseDTO
	for _, mission := range missions {
		result = append(result, s.missionToDTO(mission))
	}

	return result, nil
}

func (s *missionService) GetUserMissions(ctx context.Context, userID int64) ([]*dto.UserMissionResponseDTO, error) {
	userMissions, err := s.userMissionRepo.FindUserMissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []*dto.UserMissionResponseDTO
	for _, userMission := range userMissions {
		missionDTO := s.missionToDTO(userMission.Mission)

		var completedAt *time.Time
		if userMission.CompletedAt.Valid {
			completedAt = &userMission.CompletedAt.Time
		}

		result = append(result, &dto.UserMissionResponseDTO{
			ID:          userMission.ID,
			UserID:      userMission.UserID,
			Mission:     *missionDTO,
			CompletedAt: completedAt,
			CreatedAt:   userMission.CreatedAt,
		})
	}

	return result, nil
}

func (s *missionService) CheckMissionCompletion(ctx context.Context, userID, missionID int64) (bool, error) {
	userMission, err := s.userMissionRepo.FindUserMissionByID(ctx, userID, missionID)
	if err != nil {
		return false, err
	}

	return userMission != nil && userMission.CompletedAt.Valid, nil
}

func (s *missionService) missionToDTO(mission *model.Mission) *dto.MissionResponseDTO {
	dto := &dto.MissionResponseDTO{
		ID:           mission.ID,
		Title:        mission.Title,
		Description:  mission.Description,
		MissionType:  dto.MissionType(mission.MissionType),
		PointsReward: mission.PointsReward,
		GivesBadge:   mission.GivesBadge,
		TargetValue:  mission.TargetValue,
		CreatedAt:    mission.CreatedAt,
	}

	if mission.BadgeID.Valid {
		dto.BadgeID = &mission.BadgeID.Int64
	}

	if mission.CarbonReductionG.Valid {
		dto.CarbonReductionG = &mission.CarbonReductionG.Float64
	}

	if mission.ExpiredAt.Valid {
		dto.ExpiredAt = &mission.ExpiredAt.Time
	}

	return dto
}
