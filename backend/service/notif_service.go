// service/badge_service.go
package service

import (
	"context"

	model "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
)

type NotificationService interface {
	GetNotificationsByUserID(ctx context.Context, userID int64) ([]*model.Notification, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) GetNotificationsByUserID(ctx context.Context, userID int64) ([]*model.Notification, error) {
	return s.notificationRepo.GetByUserID(ctx, userID)
}
