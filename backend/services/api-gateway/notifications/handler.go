package notifications

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	notificationservices "github.com/wafi11/workspaces/services/notifications-service"
)

type NotificationHandler struct {
	svc *notificationservices.Service
}


func NewNotificationHandler(svc *notificationservices.Service) *NotificationHandler{
	return &NotificationHandler{
		svc: svc,
	}
}


func (h *NotificationHandler) FindAllRetrived(c echo.Context) error{

	user_id :=  c.Get("user_id").(string)
	
	data,err := h.svc.GetSentNotifications(c.Request().Context(),user_id)
	if err != nil {
		return response.Error(c,http.StatusBadRequest,"failed to find notifications",nil)
	}

	return response.Success(c,http.StatusOK,"Successfully retreived notifications",data)
}


func (h *NotificationHandler) FindAllReceived(c echo.Context) error{

	user_id :=  c.Get("user_id").(string)
	
	data,err := h.svc.GetReceivedNotifications(c.Request().Context(),user_id)
	if err != nil {
		return response.Error(c,http.StatusBadRequest,"failed to find notifications",nil)
	}

	return response.Success(c,http.StatusOK,"Successfully received notifications",data)
}



func (h *NotificationHandler) FindAllUnreadCount(c echo.Context) error{

	user_id :=  c.Get("user_id").(string)
	
	data,err := h.svc.GetUnreadCount(c.Request().Context(),user_id)
	if err != nil {
		return response.Error(c,http.StatusBadRequest,"failed to find notifications",nil)
	}

	return response.Success(c,http.StatusOK,"Successfully retreived count notifications",data)
}