package template

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/middlewares"
	templateservice "github.com/wafi11/workspaces/services/template-service"
)

func NewTemplateRouter(c *echo.Echo, db *sqlx.DB, redis *redis.Client, minioClient *minio.Client, conf *config.Config) {
	repo := templateservice.NewRepository(db, redis, minioClient)
	svc := templateservice.NewService(repo)
	handler := NewHandler(svc)

	protected := c.Group("/api/v1/templates",middlewares.AuthMiddleware(conf))
	protected_admin := protected.Group("", middlewares.AdminMiddleware())
	protected_admin.POST("", handler.CreateTemplate)
	protected.GET("", handler.GetListTemplates)
	protected.GET("/workspace/form", handler.FinTemplateWorkspaceForm)
	protected.GET("/:template_id", handler.GetTemplateDetails)
	protected_admin.GET("/:template_id/form", handler.FindTemplateDetailsForm)
	protected_admin.PUT("/:id", handler.UpdateTemplate)
	protected_admin.DELETE("/:id", handler.DeleteTemplate)

	// template variables
	protected.GET("/:template_id/variables", handler.GetTemplateVariables)
	protected_admin.POST("/:template_id/variables", handler.CreateTemplateVariable)
	protected_admin.PUT("/variables/:id", handler.UpdateTemplateVariable)
	protected_admin.DELETE("/variables/:id", handler.DeleteTemplateVariable)

	// template addons
	protected_admin.POST("/:template_id/add-ons", handler.CreateTemplateAddon)
	protected_admin.GET("/:template_id/add-ons", handler.GetTemplateAddons)
	protected_admin.PUT("/add-ons/:id", handler.UpdateTemplateAddon)
	protected_admin.DELETE("/add-ons/:id", handler.DeleteTemplateAddon)

	// template files
	protected_admin.POST("/:template_id/files", handler.CreateTemplateFiles)
	protected_admin.GET("/:template_id/files", handler.GetTemplateFiles)
	protected_admin.PUT("/files/:id", handler.UpdateTemplateFiles)
	protected_admin.DELETE("/files/:id", handler.DeleteTemplateFiles)

}
