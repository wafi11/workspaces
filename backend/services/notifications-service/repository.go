package notificationservices

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	DB *sql.DB
	Redis *redis.Client
}


func NewRepository(Db *sql.DB,redis *redis.Client) *Repository {
	return &Repository{
		DB: Db,
		Redis: redis,
	}
}

func (repo *Repository) CreateNotifications(ctx context.Context, req *NotificationRequest, tx *sql.Tx) (*NotificationResponse, error) {
    var notification_id string

    query := `
        INSERT INTO notifications (
            user_id, retreived_id, type, title,
            message, metadata, is_read, created_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        RETURNING id
    `

    var err error
    if tx != nil {
        err = tx.QueryRowContext(ctx, query,
            req.UserID, req.RetreivedID, req.NotificationType,
            req.Title, req.Message, req.Metadata, req.IsRead, time.Now(),
        ).Scan(&notification_id)
    } else {
        err = repo.DB.QueryRowContext(ctx, query,
            req.UserID, req.RetreivedID, req.NotificationType,
            req.Title, req.Message, req.Metadata, req.IsRead, time.Now(),
        ).Scan(&notification_id)
    }

    if err != nil {
        log.Printf("[create_notification] err: %s", err.Error())
        return nil, fmt.Errorf("failed to create notification: %w", err)
    }

    return &NotificationResponse{
        Notification: Notification{
            ID:               notification_id,
            UserID:           req.UserID,
            RetreivedID:      req.RetreivedID,
            NotificationType: req.NotificationType,
            Title:            req.Title,
            Message:          req.Message,
            Metadata:         req.Metadata,
            IsRead:           req.IsRead,
            CreatedAt:        time.Now(),
        },
        Message: "Successfully create notifications",
    }, nil
}

func (repo *Repository) UpdateReadNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error {
	query := `
		UPDATE notifications
		SET is_read = true, updated_at = $1
		WHERE id = $2 AND user_id = $3
	`

	var err error
	var result sql.Result

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, time.Now(), notificationID, userID)
	} else {
		result, err = repo.DB.ExecContext(ctx, query, time.Now(), notificationID, userID)
	}

	if err != nil {
		log.Printf("[update_read_notification] err: %s", err.Error())
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("notification not found or unauthorized")
	}

	return nil
}

func (repo *Repository) UpdateReadAllNotification(ctx context.Context, userID string, tx *sql.Tx) error {
	query := `
		UPDATE notifications
		SET is_read = true, updated_at = $1
		WHERE user_id = $2 AND is_read = false
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, time.Now(), userID)
	} else {
		_, err = repo.DB.ExecContext(ctx, query, time.Now(), userID)
	}

	if err != nil {
		log.Printf("[update_read_all_notification] err: %s", err.Error())
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

func (repo *Repository) DeleteNotification(ctx context.Context, notificationID, userID string, tx *sql.Tx) error {
	query := `
		UPDATE notifications
		SET deleted_at_retreived = $1
		WHERE id = $2 AND user_id = $3 AND deleted_at_retreived IS NULL
	`

	var err error
	var result sql.Result

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, time.Now(), notificationID, userID)
	} else {
		result, err = repo.DB.ExecContext(ctx, query, time.Now(), notificationID, userID)
	}

	if err != nil {
		log.Printf("[delete_notification] err: %s", err.Error())
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("notification not found or already deleted")
	}

	return nil
}

func (repo *Repository) GetReceivedNotifications(ctx context.Context, userID string) ([]Notification, error) {
	query := `
		SELECT 
			id, user_id, retreived_id, type, title,
			message, metadata, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		  AND deleted_at_received IS NULL
		ORDER BY created_at DESC
	`

	rows, err := repo.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("[get_received_notifications] err: %s", err.Error())
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID, &n.UserID, &n.RetreivedID, &n.NotificationType,
			&n.Title, &n.Message, &n.Metadata, &n.IsRead, &n.CreatedAt,
		)
		if err != nil {
			log.Printf("[get_received_notifications] scan err: %s", err.Error())
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return notifications, nil
}

func (repo *Repository) GetSentNotifications(ctx context.Context, userID string) ([]Notification, error) {
	query := `
		SELECT 
			id, user_id, retreived_id, type, title,
			message, metadata, is_read, created_at
		FROM notifications
		WHERE retreived_id = $1
		  AND deleted_at_retreived IS NULL
		ORDER BY created_at DESC
	`

	rows, err := repo.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("[get_sent_notifications] err: %s", err.Error())
		return nil, fmt.Errorf("failed to get sent notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID, &n.UserID, &n.RetreivedID, &n.NotificationType,
			&n.Title, &n.Message, &n.Metadata, &n.IsRead, &n.CreatedAt,
		)
		if err != nil {
			log.Printf("[get_sent_notifications] scan err: %s", err.Error())
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return notifications, nil
}

func (repo *Repository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM notifications
		WHERE user_id = $1
		  AND is_read = false
		  AND deleted_at_retreived IS NULL
	`

	var count int
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		log.Printf("[get_unread_count] err: %s", err.Error())
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}
