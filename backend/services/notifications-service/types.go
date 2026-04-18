package notificationservices

import (
	"encoding/json"
	"time"
)

type NotificationType string

const (
	InvitationCollaborator NotificationType = "INVITATION_COLLABORATOR"
)

type Notification struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	RetreivedID  	 string  `json:"retreived_id"`
	NotificationType NotificationType `json:"notification_type"`
	Title            string           `json:"title"`
	Message          string           `json:"message"`
	Metadata         json.RawMessage         `json:"metadata"`
IsRead bool `json:"is_read"`
	CreatedAt        time.Time  `json:"created_at"`
}



type NotificationRequest struct {
	UserID  string  `json:"user_id"`
	NotificationType NotificationType `json:"notification_type"`
	Title            string           `json:"title"`
	RetreivedID  	 string  `json:"retreived_id"`
	Message          string           `json:"message"`
	Metadata         json.RawMessage         `json:"metadata"`
IsRead bool `json:"is_read"`
}


type NotificationResponse struct {
	Notification
	Message string `json:"message"`
}

type ReadNotificationRequest struct {
	UserID string `json:"user_id"`
}


type ReadNotificationResponse struct {
	Message string `json:"message"`
}

type DeleteMessageNotificationRequest struct {
	NotificationID string `json:"notification_id"`
}
