package service

import (
	"context"

	"github.com/Qodarrz/fiber-app/dto"
	"github.com/Qodarrz/fiber-app/repository"
)

type UserCustomEndpointServiceInterface interface {
	GetUserCustomData(ctx context.Context, userID int64) (*dto.UserCustomDataResponseDTO, error)
	GetLeaderboard(ctx context.Context, req *dto.LeaderboardRequestDTO) (*dto.LeaderboardResponseDTO, error)
	GetMissionProgress(ctx context.Context, userID, missionID int64) (float64, error)
	GetAllMissionProgress(ctx context.Context, userID int64) ([]dto.MissionProgressResponse, error)
}

type userCustomEndpointService struct {
	userCustomRepo     repository.UserCustomEndpointRepoInterface
	checkMissionRepo   repository.CheckMissionRepositoryInterface
	missionRepo        repository.MissionRepositoryInterface
}

func NewUserCustomEndpointService(
	userCustomRepo repository.UserCustomEndpointRepoInterface,
	checkMissionRepo repository.CheckMissionRepositoryInterface,
	missionRepo repository.MissionRepositoryInterface,
) UserCustomEndpointServiceInterface {
	return &userCustomEndpointService{
		userCustomRepo:   userCustomRepo,
		checkMissionRepo: checkMissionRepo,
		missionRepo:      missionRepo,	
	}
}

func (s *userCustomEndpointService) GetUserCustomData(ctx context.Context, userID int64) (*dto.UserCustomDataResponseDTO, error) {
	data, err := s.userCustomRepo.GetUserCustomData(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.mapToDTO(data), nil
}

func (s *userCustomEndpointService) GetLeaderboard(ctx context.Context, req *dto.LeaderboardRequestDTO) (*dto.LeaderboardResponseDTO, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.TimeRange == "" {
		req.TimeRange = "all"
	}

	entries, total, err := s.userCustomRepo.GetLeaderboard(ctx, req.Page, req.Limit, req.TimeRange)
	if err != nil {
		return nil, err
	}

	response := &dto.LeaderboardResponseDTO{
		Leaderboard: s.mapLeaderboardToDTO(entries, req.Page, req.Limit),
		Pagination: dto.PaginationDTO{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: (total + req.Limit - 1) / req.Limit,
		},
		TimeRange: req.TimeRange,
	}

	return response, nil
}


func (s *userCustomEndpointService) mapToDTO(data *repository.UserCustomData) *dto.UserCustomDataResponseDTO {
	response := &dto.UserCustomDataResponseDTO{
		User: dto.UserDetailResponseDTO{
			ID:        data.User.ID,
			Username:  data.User.Username,
			Email:     data.User.Email,
			Role:      data.User.Role,
			FullName:  data.User.FullName,
			AvatarURL: data.User.AvatarURL,
			Birthdate: data.User.Birthdate,
			Gender:    data.User.Gender,
			CreatedAt: data.User.CreatedAt,
			TotalPoints: data.UserPoints.TotalPoints,
		},
	}

	for _, v := range data.Vehicles {
		response.Vehicles = append(response.Vehicles, dto.CustomVehicleDTO{
			ID:           v.ID,
			VehicleType:  v.VehicleType,
			FuelType:     v.FuelType,
			Name:         v.Name,
			CreatedAt:    v.CreatedAt,
			TotalLogs:    v.TotalLogs,
			TotalCarbon:  v.TotalCarbon,
		})
	}

	// Map electronics
	for _, e := range data.Electronics {
		response.Electronics = append(response.Electronics, dto.CustomElectronicDTO{
			ID:           e.ID,
			DeviceName:   e.DeviceName,
			DeviceType:   e.DeviceType,
			PowerWatts:   e.PowerWatts,
			CreatedAt:    e.CreatedAt,
			TotalLogs:    e.TotalLogs,
			TotalCarbon:  e.TotalCarbon,
		})
	}

	// Map missions
	for _, m := range data.Missions {
		response.Missions = append(response.Missions, dto.CustomMissionProgressDTO{
			ID:           m.ID,
			Title:        m.Title,
			Description:  m.Description,
			MissionType:  m.MissionType,
			PointsReward: m.PointsReward,
			TargetValue:  m.TargetValue,
			Progress:     m.ProgressValue,
			CompletedAt:  m.CompletedAt,
			CreatedAt:    m.CreatedAt,
		})
	}

	// Map badges
	for _, b := range data.Badges {
		response.Badges = append(response.Badges, dto.CustomBadgeDTO{
			ID:          b.ID,
			Name:        b.Name,
			ImageURL:    b.ImageURL,
			Description: b.Description,
			RedeemedAt:  b.RedeemedAt,
			CreatedAt:   b.CreatedAt,
		})
	}

	// Map point transactions
	for _, pt := range data.PointTransactions {
		response.PointHistory = append(response.PointHistory, dto.CustomPointTransactionDTO{
			ID:            pt.ID,
			Amount:        pt.Amount,
			Direction:     pt.Direction,
			Source:        pt.Source,
			ReferenceType: pt.ReferenceType,
			ReferenceID:   pt.ReferenceID,
			Note:          pt.Note,
			CreatedAt:     pt.CreatedAt,
		})
	}

	// Map activity logs
	for _, al := range data.ActivityLogs {
		response.ActivityLogs = append(response.ActivityLogs, dto.CustomActivityLogDTO{
			ID:        al.ID,
			Activity:  al.Activity,
			CreatedAt: al.CreatedAt,
		})
	}

	// Map orders
	for _, o := range data.Orders {
		orderDTO := dto.CustomOrderDTO{
			ID:          o.ID,
			TotalPoints: o.TotalPoints,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
		}

		for _, item := range o.Items {
			orderDTO.Items = append(orderDTO.Items, dto.CustomOrderItemDTO{
				ID:              item.ID,
				ItemName:        item.ItemName,
				Qty:             item.Qty,
				PriceEachPoints: item.PriceEachPoints,
			})
		}

		response.Orders = append(response.Orders, orderDTO)
	}

	for _, mvc := range data.MonthlyVehicleCarbon {
		response.MonthlyVehicleCarbon = append(response.MonthlyVehicleCarbon, dto.MonthlyCarbonDTO{
			Month:       mvc.Month,
			TotalCarbon: mvc.TotalCarbon,
		})
	}

	// Map monthly electronic carbon
	for _, mec := range data.MonthlyElectronicCarbon {
		response.MonthlyElectronicCarbon = append(response.MonthlyElectronicCarbon, dto.MonthlyCarbonDTO{
			Month:       mec.Month,
			TotalCarbon: mec.TotalCarbon,
		})
	}

	return response
}


func (s *userCustomEndpointService) mapLeaderboardToDTO(entries []repository.LeaderboardEntry, page, limit int) []dto.LeaderboardItemDTO {
	var leaderboard []dto.LeaderboardItemDTO

	for i, entry := range entries {
		rank := (page-1)*limit + i + 1
		leaderboard = append(leaderboard, dto.LeaderboardItemDTO{
			Rank:              rank,
			User: dto.UserSimpleDTO{
				ID:        entry.User.ID,
				Username:  entry.User.Username,
				FullName:  entry.User.FullName,
				AvatarURL: entry.User.AvatarURL,
			},
			TotalPoints:       entry.TotalPoints,
			CompletedMissions: entry.CompletedMissions,
			Score:            entry.Score,
			CarbonReduction:  entry.CarbonReduction,
		})
	}

	return leaderboard
}

func (s *userCustomEndpointService) GetMissionProgress(ctx context.Context, userID, missionID int64) (float64, error) {
	progress, err := s.checkMissionRepo.GetMissionProgress(ctx, userID, missionID)
	if err != nil {
		return 0, err
	}
	return progress, nil
}

func (s *userCustomEndpointService) GetAllMissionProgress(ctx context.Context, userID int64) ([]dto.MissionProgressResponse, error) {
	return s.missionRepo.GetAllMissionProgress(ctx, userID)
}