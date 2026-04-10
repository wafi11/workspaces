package services

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
)

type Repository struct {
	redisClient *redis.Client
	db          *sqlx.DB
	minioClient *minio.Client
	jobQueue    chan<- WorkspaceJob
	dynClient   dynamic.Interface
	mapper      meta.RESTMapper
}

func NewRepository(redis *redis.Client, jobQueue chan<- WorkspaceJob, db *sqlx.DB, minioClient *minio.Client, dynClient dynamic.Interface, mapper meta.RESTMapper) *Repository {
	return &Repository{
		redisClient: redis,
		jobQueue:    jobQueue,
		db:          db,
		minioClient: minioClient,
		dynClient:   dynClient,
		mapper:      mapper,
	}
}

func StartOperator(ctx context.Context, jobQueue <-chan WorkspaceJob, k8sClient IK8SClient, r *Repository) {
	go func() {
		for {
			select {
			case job := <-jobQueue:
				log.Printf("[operator] received job action=%q workspaceId=%s\n", job.Action, job.WorkspaceId) // ← tambah ini

				switch job.Action {
				case JobCreate:
					if err := handleCreate(ctx, job, k8sClient, r); err != nil {
						log.Printf("failed to create workspace %s: %v", job.WorkspaceId, err)
					}
				case JobAdd:
					if err := handleAdd(ctx, job, r); err != nil {
						log.Printf("failed to create services %s: %v", job.WorkspaceId, err)
					}
				case JobDelete:
					if err := handleDelete(ctx, job, k8sClient); err != nil {
						log.Printf("failed to delete workspace %s: %v", job.WorkspaceId, err)
					}
				default:
					log.Printf("[operator] unknown action: %q", job.Action) // ← dan ini
				}

			case <-ctx.Done():
				log.Println("operator shutting down")
				return
			}
		}
	}()
}

func handleAdd(ctx context.Context, job WorkspaceJob, r *Repository) error {
	log.Printf("DEBUG JOB: CPUReq: %s, CPULimit: %s, MemReq: %s, MemLimit: %s",
		job.CPURequest, job.CPULimit, job.MemoryRequest, job.MemoryLimit)
	if err := r.ExecuteDeployment(ctx, job.TemplateId, DeployParams{
		Namespace:  generateNamespace(job.UserId),
		User:       &job.UserId,
		Name:       job.Name,
		Image:      &job.Image,
		Replicas:   1,
		RunAsUser:  1000,
		RunAsGroup: 1000,
		FsGroup:    1000,
		Password:   "password123",
		CPURequest: job.CPURequest,
		MemRequest: job.MemoryRequest,
		CPULimit:   job.CPULimit,
		MemLimit:   job.MemoryLimit,
		Username:   job.Username,
		Domain:     "wfdnstore.online",
	}); err != nil {
		return fmt.Errorf("failed deployment: %w", err)
	}

	return nil

}

func handleCreate(ctx context.Context, job WorkspaceJob, k8sClient IK8SClient, r *Repository) error {

	if err := k8sClient.CreateNamespace(ctx, job.UserId); err != nil {
		return (fmt.Errorf("failed to create namespace: %w", err))
	}
	ensureValue := func(val, fallback string) string {
		if val == "" || val == "<nil>" {
			return fallback
		}
		return val
	}

	if err := k8sClient.CreateResourceQuota(ctx, job.UserId, QuotaConfig{
		CPULimit:      ensureValue(job.CPURequest, "1"),
		MemoryLimit:   ensureValue(job.MemoryLimit, "1024Mi"),
		StorageLimit:  ensureValue(job.StorageLimit, "5Gi"),
		CPURequest:    ensureValue(job.CPURequest, "1"),
		MemoryRequest: ensureValue(job.MemoryRequest, "1024Mi"),
	}); err != nil {
		return (fmt.Errorf("failed to create resource quota: %w", err))
	}
	err := k8sClient.SetupRBAC(ctx, job.UserId)
	if err != nil {
		return (fmt.Errorf("failed to setup rbac: %w", err))
	}
	dbName := getEnvString(job.EnvVars, "DB_NAME")

	// 4. deploy template
	if err := r.ExecuteDeployment(ctx, job.TemplateId, DeployParams{
		WS_TOKEN:         getEnvString(job.EnvVars, "WS_TOKEN"),
		WS_REFRESH_TOKEN: getEnvString(job.EnvVars, "WS_REFRESH_TOKEN"),
		WS_API_URL:       getEnvString(job.EnvVars, "WS_API_URL"),
		Namespace:        generateNamespace(job.UserId),
		DB_NAME:          &dbName,
		User:             &job.UserId,
		Name:             job.UserId,
		StorageClass:     "nfs",
		StorageSize:      job.StorageRequest,
		Replicas:         1,
		RunAsUser:        1000,
		RunAsGroup:       1000,
		FsGroup:          1000,
		Password:         "password123",
		CPULimit:         "0.25",
		MemLimit:         "128Mi",
		Username:         job.Username,
		CPURequest:       "0.10",
		MemRequest:       "100Mi",
		Domain:           "wfdnstore.online",
	}); err != nil {
		return (fmt.Errorf("failed deployment: %w", err))
	}

	log.Printf("workspace %s provisioned successfully", job.WorkspaceId)
	return nil
}

func handleDelete(ctx context.Context, job WorkspaceJob, k8sClient IK8SClient) error {
	// hapus namespace → otomatis hapus semua resource di dalamnya
	if err := k8sClient.DeleteNamespace(ctx, job.Namespace); err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	log.Printf("workspace %s deleted successfully", job.WorkspaceId)
	return nil
}
