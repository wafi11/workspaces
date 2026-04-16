package services

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	workspacev1 "github.com/wafi11/workspace-operator/pkg/proto"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
)

type Repository struct {
	redisClient *redis.Client
	db          *sqlx.DB
	minioClient *minio.Client
	dynClient   dynamic.Interface
	mapper      meta.RESTMapper
}

func NewRepository(
	redis *redis.Client,
	db *sqlx.DB,
	minioClient *minio.Client,
	dynClient dynamic.Interface,
	mapper meta.RESTMapper,
) *Repository {
	return &Repository{
		redisClient: redis,
		db:          db,
		minioClient: minioClient,
		dynClient:   dynClient,
		mapper:      mapper,
	}
}

// ─── Operator ────────────────────────────────────────────

func StartOperator(
	ctx context.Context,
	jobQueue <-chan *workspacev1.WorkspaceEnvelope,
	k8sClient IK8SClient,
	r *Repository,
) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[operator] shutting down")
				return

			case envelope := <-jobQueue:
				if err := dispatch(ctx, envelope, k8sClient, r); err != nil {
					log.Printf("[operator] error: %v", err)

					// publish UpdateStatus FAILED balik ke redis
					// supaya repository bisa rollback DB
					publishFailed(ctx, r.redisClient, envelope, err)
				}
			}
		}
	}()
}

func dispatch(
	ctx context.Context,
	envelope *workspacev1.WorkspaceEnvelope,
	k8sClient IK8SClient,
	r *Repository,
) error {
log.Printf("[operator] envelope type=%T payload=%+v", envelope.Payload, envelope.Payload)
	switch p := envelope.Payload.(type) {
		
	case *workspacev1.WorkspaceEnvelope_Create:
		log.Printf("[operator] create workspace_id=%s", p.Create.Identity.WorkspaceId)
		return handleCreate(ctx, p.Create, k8sClient, r)

	case *workspacev1.WorkspaceEnvelope_Add:
		log.Printf("[operator] add pod workspace_id=%s", p.Add.Identity.WorkspaceId)
		return handleAdd(ctx, p.Add, r)

	case *workspacev1.WorkspaceEnvelope_Delete:
		log.Printf("[operator] delete workspace_id=%s", p.Delete.Identity.WorkspaceId)
		return handleDelete(ctx, p.Delete, k8sClient)

	case *workspacev1.WorkspaceEnvelope_Stop:
		log.Printf("[operator] stop workspace_id=%s", p.Stop.Identity.WorkspaceId)
		return handleStop(ctx, p.Stop, k8sClient)

	case *workspacev1.WorkspaceEnvelope_Start:
		log.Printf("[operator] start workspace_id=%s", p.Start.Identity.WorkspaceId)
		return handleStart(ctx, p.Start, k8sClient)

	case *workspacev1.WorkspaceEnvelope_CreatePort:
		log.Printf("[operator] create port workspace_name=%s", p.CreatePort.WorkspaceName)
		return  handleCreatePort(ctx,p.CreatePort,k8sClient)
	
	case *workspacev1.WorkspaceEnvelope_DeletePort:
		log.Printf("[operator] delete port workspace_name=%s", p.DeletePort.WorkspaceName)
		return handleDeletePort(ctx,p.DeletePort,k8sClient)

	default:
		return fmt.Errorf("unknown payload type: %T", envelope.Payload)
	}
}

// ─── Handlers ────────────────────────────────────────────

func handleCreate(ctx context.Context, e *workspacev1.CreateWorkspaceEvent, k8sClient IK8SClient, r *Repository) error {
	id := e.Identity

	if err := k8sClient.CreateNamespace(ctx, id.UserId); err != nil {
		return fmt.Errorf("create namespace: %w", err)
	}

	if err := k8sClient.CreateResourceQuota(ctx, id.UserId, toQuotaConfig(e.Resources)); err != nil {
		return fmt.Errorf("create quota: %w", err)
	}

	if err := k8sClient.SetupRBAC(ctx, id.UserId); err != nil {
		return fmt.Errorf("setup rbac: %w", err)
	}

	if err := r.ExecuteDeployment(ctx, e.TemplateId, toDeployParams(e.Identity, e.Resources, e.EnvVars)); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}

	log.Printf("[operator] workspace %s provisioned successfully", id.WorkspaceId)
	return nil
}

func handleAdd(ctx context.Context, e *workspacev1.AddPodEvent, r *Repository) error {
	res := e.Resources
	log.Printf("[operator] add pod CPUReq=%s CPULimit=%s MemReq=%s MemLimit=%s",
		res.CpuRequest, res.CpuLimit, res.MemoryRequest, res.MemoryLimit)

	if err := r.ExecuteDeployment(ctx, e.TemplateId, DeployParams{
		WsID:       e.Identity.WorkspaceId,
		DB_USER: &e.AddOns.DbUser,
		DB_NAME: &e.AddOns.DbName,
		DB_PASSWORD: &e.AddOns.DbPassword,
		Namespace:  generateNamespace(e.Identity.UserId),
		User:       &e.Identity.UserId,
		Name:       e.Identity.Name,
		Image:      &e.Image,
		Password:   e.Identity.Password,
		Replicas:   int(e.Replicas),
		RunAsUser:  1000,
		RunAsGroup: 1000,
		FsGroup:    1000,
		CPURequest: res.CpuRequest,
		CPULimit:   res.CpuLimit,
		MemRequest: res.MemoryRequest,
		MemLimit:   res.MemoryLimit,
		Username:   e.Identity.Username,
		Domain:     "wfdnstore.online",
	}); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	return nil
}

func handleDelete(ctx context.Context, e *workspacev1.DeleteWorkspaceEvent, k8sClient IK8SClient) error {
	if err := k8sClient.DeleteNamespace(ctx, e.Identity.Namespace); err != nil {
		return fmt.Errorf("delete namespace: %w", err)
	}
	log.Printf("[operator] workspace %s deleted", e.Identity.WorkspaceId)
	return nil
}

func handleStop(ctx context.Context, e *workspacev1.StopWorkspaceEvent, k8sClient IK8SClient) error {
	if err := k8sClient.StopScalling(ctx,e.Identity.Namespace,e.Identity.Name); err != nil {
		return fmt.Errorf("scale down: %w", err)
	}
	log.Printf("[operator] workspace %s stopped", e.Identity.WorkspaceId)
	return nil
}

func handleStart(ctx context.Context, e *workspacev1.StartWorkspaceEvent, k8sClient IK8SClient) error {
	if err := k8sClient.StartScalling(ctx,e.Identity.Namespace,e.Identity.Name); err != nil {
		return fmt.Errorf("scale up: %w", err)
	}
	log.Printf("[operator] workspace %s started", e.Identity.WorkspaceId)
	return nil
}


func handleCreatePort(ctx context.Context, e *workspacev1.CreatePort, k8sClient IK8SClient) error{
	if err := k8sClient.CreatePort(ctx,e.UserId,e.WorkspaceName,int(e.Port)); err != nil {
		return fmt.Errorf("failed to create port : %s",err.Error())
	}

	serviceName := fmt.Sprintf("%s-%d-svc", e.WorkspaceName, e.Port)

	if err := k8sClient.ExposeToIngress(ctx,e.UserId,e.WorkspaceName,serviceName,e.Domain,e.Port); err != nil {
		return fmt.Errorf("failed to expose port : %s",err.Error())
	}

	log.Printf("[operator] update  port %s started", e.WorkspaceName)

	return  nil
}

func handleDeletePort(ctx context.Context,e *workspacev1.DeletePort, k8sClient IK8SClient) error {
	if err := k8sClient.DeletePort(ctx,e.UserId,e.WorkspaceName,int(e.Port)); err != nil {
		return fmt.Errorf("failed to delete port : %s",err.Error())
	}

	domain :=  fmt.Sprintf("%d-%s-%s.wfdnstore.online",e.Port,e.WorkspaceName,e.UserId)

	if err :=  k8sClient.RemoveFromIngress(ctx,e.UserId,e.WorkspaceName,domain); err != nil {
		return fmt.Errorf("failed to remove ingress : %s",err.Error())
	}

	log.Printf("[operator] delete port %s started", e.WorkspaceName)
	return nil
}