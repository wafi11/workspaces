package notifications

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/middlewares"
	notificationservices "github.com/wafi11/workspaces/services/notifications-service"
)

func NewNotificationRoutesfunc(c *echo.Echo, db *sqlx.DB, redis *redis.Client, conf *config.Config) {

	repo := notificationservices.NewRepository(db.DB,redis)
	svc := notificationservices.NewService(repo)
	h := NewNotificationHandler(svc)


	notif := c.Group("/api/v1/notifications",middlewares.AuthMiddleware(conf))

	notif.GET("/retreived",h.FindAllRetrived)
	notif.GET("/received",h.FindAllReceived)
	notif.GET("/unread-count",h.FindAllUnreadCount)

}