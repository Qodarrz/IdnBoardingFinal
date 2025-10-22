// service/badge_service.go
package service

import (
	"context"

	model "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
)

type BadgeService interface {
	GetBadgesWithOwnership(ctx context.Context, userID int64, page, limit int) ([]*model.BadgeWithOwnership, error)
}

type badgeService struct {
	badgeRepo repository.BadgeRepositoryInterface
}

func NewBadgeService(badgeRepo repository.BadgeRepositoryInterface) BadgeService {
	return &badgeService{badgeRepo: badgeRepo}
}

func (s *badgeService) GetBadgesWithOwnership(ctx context.Context, userID int64, page, limit int) ([]*model.BadgeWithOwnership, error) {
	return s.badgeRepo.FindAllWithOwnership(ctx, userID, page, limit)
}