package notificationservices

import (
	"context"
	"database/sql"
)

type NotificationRepository interface {
	CreateNotifications(ctx context.Context, req *NotificationRequest, tx *sql.Tx) (*NotificationResponse, error)
	UpdateReadNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error
	UpdateReadAllNotification(ctx context.Context, userID string, tx *sql.Tx) error
	DeleteNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error
	GetReceivedNotifications(ctx context.Context, userID string) ([]Notification, error)
	GetSentNotifications(ctx context.Context, userID string) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}

type NotificationService interface {
	CreateNotifications(ctx context.Context, req *NotificationRequest, tx *sql.Tx) (*NotificationResponse, error)
	UpdateReadNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error
	UpdateReadAllNotification(ctx context.Context, userID string, tx *sql.Tx) error
	DeleteNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error
	GetReceivedNotifications(ctx context.Context, userID string) ([]Notification, error)
	GetSentNotifications(ctx context.Context, userID string) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}