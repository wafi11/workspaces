package template

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	templateservice "github.com/wafi11/workspaces/services/template-service"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
)

func NewTemplateRouter(c *echo.Echo, db *sqlx.DB, redis *redis.Client, minioClient *minio.Client, dynClient dynamic.Interface, mapper meta.RESTMapper) {
	repo := templateservice.NewRepository(db, redis, minioClient, dynClient, mapper)
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
	protected.POST("/:template_id/addons", handler.CreateTemplateAddon)
	protected.GET("/:template_id/addons", handler.GetTemplateAddons)
	protected.PUT("/addons/:id", handler.UpdateTemplateAddon)
	protected.DELETE("/addons/:id", handler.DeleteTemplateAddon)

	// template files
	protected.POST("/:template_id/files", handler.CreateTemplateFiles)
	protected.GET("/:template_id/files", handler.GetTemplateFiles)
	protected.PUT("/files/:id", handler.UpdateTemplateFiles)
	protected.DELETE("/files/:id", handler.DeleteTemplateFiles)

}
