package notificationservices

import (
	"context"
	"database/sql"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateNotifications(ctx context.Context, req *NotificationRequest, tx *sql.Tx) (*NotificationResponse, error) {
	return s.repo.CreateNotifications(ctx, req,tx)
}
func (s *Service) UpdateReadNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error {
	return s.repo.UpdateReadNotification(ctx,notificationID,userID,tx)
}
func (s *Service) UpdateReadAllNotification(ctx context.Context, userID string, tx *sql.Tx) error {
	return s.repo.UpdateReadAllNotification(ctx,userID,tx)
}
func (s *Service) DeleteNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error {
	return s.repo.DeleteNotification(ctx,notificationID,userID,tx)
}
func (s *Service) GetReceivedNotifications(ctx context.Context, userID string) ([]Notification, error) {
	return s.repo.GetReceivedNotifications(ctx,userID)
}
func (s *Service) GetSentNotifications(ctx context.Context, userID string) ([]Notification, error) {
	return s.repo.GetSentNotifications(ctx,userID)
}
func (s *Service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx,userID)
}