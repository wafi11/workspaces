package template

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	templateservice "github.com/wafi11/workspaces/services/template-service"
)

func NewTemplateRouter(c *echo.Echo, db *sqlx.DB, redis *redis.Client, minioClient *minio.Client) {
	repo := templateservice.NewRepository(db, redis, minioClient)
	svc := templateservice.NewService(repo)
	handler := NewHandler(svc)

	protected := c.Group("/api/v1/templates")
	protected.POST("", handler.CreateTemplate)
	protected.GET("", handler.GetListTemplates)
	protected.GET("/:template_id", handler.GetTemplateDetails)
	protected.PUT("/:id", handler.UpdateTemplate)
	protected.DELETE("/:id", handler.DeleteTemplate)


	// template variables
	protected.POST("/:template_id/variables", handler.CreateTemplateVariable)
	protected.GET("/:template_id/variables", handler.GetTemplateVariables)
	protected.PUT("/variables/:id", handler.UpdateTemplateVariable)
	protected.DELETE("/variables/:id", handler.DeleteTemplateVariable)

	// template addons
	protected.POST("/:template_id/add-ons", handler.CreateTemplateAddon)
	protected.GET("/:template_id/add-ons", handler.GetTemplateAddons)
	protected.PUT("/add-ons/:id", handler.UpdateTemplateAddon)
	protected.DELETE("/add-ons/:id", handler.DeleteTemplateAddon)

	// template files
	protected.POST("/:template_id/files", handler.CreateTemplateFiles)
	protected.GET("/:template_id/files", handler.GetTemplateFiles)
	protected.PUT("/files/:id", handler.UpdateTemplateFiles)
	protected.DELETE("/files/:id", handler.DeleteTemplateFiles)

}
