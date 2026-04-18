package workspaceservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/wafi11/workspaces/pkg/proto"
	"github.com/wafi11/workspaces/pkg/websocket"
)

type Service struct {
	repo     WorkspaceRepository
	jobQueue <-chan *proto.WorkspaceEnvelope
	hub      *websocket.Hub
}

func NewService(repo WorkspaceRepository, jobQueue <-chan *proto.WorkspaceEnvelope, hub *websocket.Hub) *Service {
	return &Service{repo: repo, jobQueue: jobQueue, hub: hub}
}

func (s *Service) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error) {
	if err := validateCreateWorkspace(req); err != nil {
		return nil, err
	}
	return s.repo.CreateWorkspace(ctx, req, username)
}

func (s *Service) ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error) {
	if err := validateListWorkspaces(req); err != nil {
		return nil, err
	}
	return s.repo.ListWorkspacesByUserId(ctx, req)
}

func (s *Service) ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error) {
	data, err := s.repo.ListWorkspaces(ctx, limit, offset, status)
	if err != nil {
		log.Printf("ERROR ListWorkspaces: %v", err)
		return nil, errors.New("failed to fetch workspaces")
	}
	return &ListWorkspacesResponse{
		Workspaces: data.Workspaces,
	}, nil
}

func (s *Service) ListWorkspaceForm(ctx context.Context, userId string) ([]ListWorkspaceForm, error) {
	return s.repo.ListWorkspaceForm(ctx, userId)
}

func (s *Service) GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error) {
	if err := validateGetWorkspace(req, true); err != nil {
		return nil, err
	}
	return s.repo.GetWorkspace(ctx, req)
}

func (s *Service) DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error) {
	if err := validateDeleteWorkspace(req); err != nil {
		return nil, err
	}
	return s.repo.DeleteWorkspace(ctx, req)
}


func (s *Service) UpdateWorkspaceStatus(ctx context.Context, workspaceId,userId string, status string) error {
	log.Printf("[consumer] sending to userID=%s clients=%v", userId, s.hub)
	s.hub.SendToUser(userId, map[string]string{
		"type":         fmt.Sprintf("workspace.%s",status),
		"workspace_id": workspaceId,
		"status":       status,
	})

	return s.repo.UpdateWorkspaceStatus(ctx, workspaceId, status)
}

func (s *Service)  ListWorkspacePorts(ctx context.Context, workspaceID string) ([]WorkspacePort, error){
	return s.repo.ListWorkspacePorts(ctx,workspaceID)
}


func (s *Service) CreateWorkspacePort(ctx context.Context, workspaceID string, port int, userID string)(*WorkspacePort, error){
	
	return s.repo.CreateWorkspacePort(ctx,workspaceID,port,userID)
}
func (s *Service) ValidateWorkspaceOwner(ctx context.Context, workspaceID, userID string) error{
	return s.repo.ValidateWorkspaceOwner(ctx,workspaceID,userID)
}

func (s *Service) DeleteWorkspacePort(ctx context.Context, workspaceID string, port int) error{
	return s.repo.DeleteWorkspacePort(ctx,workspaceID,port)
}

func (s *Service)  GetWorkspacePort(ctx context.Context, workspaceID string, port int) (*WorkspacePort, error){
	return s.repo.GetWorkspacePort(ctx,workspaceID,port)
}


func (s *Service)  AddCollaborators(c context.Context, req WorkspaceCollaborator) (*WorkspaceCollaboratorResponse, error){
	return s.repo.AddCollaborators(c,req)
}

func (s *Service)  UpdateCollaborator(c context.Context, req UpdateCollaboratorRequest) error {
	return s.repo.UpdateCollaborator(c,req)
}

func (s *Service) RemoveCollaborator(c context.Context, req RemoveCollaboratorRequest) error {
	return s.repo.RemoveCollaborator(c,req)
}


func (s *Service)  GetCollaboratedWorkspaces(ctx context.Context, userID string) ([]CollaboratedWorkspace, error){
	return s.repo.GetCollaboratedWorkspaces(ctx,userID)
}



func (s *Service) AcceptOrDeniedInvitationCollborator(c context.Context,types,notificationID string) error{
	return s.repo.AcceptOrDeniedInvitationCollborator(c,types,notificationID)
}

func (s *Service) StartEventConsumer(ctx context.Context) {
	go func() {
		for {
			select {
			case event := <-s.jobQueue:
				update, ok := event.Payload.(*proto.WorkspaceEnvelope_Update)
				if !ok {
					log.Println("[consumer] bukan update event, skip")
					continue
				}

				wsID := update.Update.WorkspaceId
				userID := update.Update.UserId
				status := update.Update.Status

				if err := s.repo.UpdateWorkspaceStatus(ctx, wsID, string(ConvertWorkspaceStatus(status))); err != nil {
					log.Printf("[consumer] failed update DB: %v", err)
					continue
				}

				if err := s.repo.CreateWorkspaceSessions(ctx, CreateWorkspaceSessions{
					WorkspaceId: wsID,
					UserId:      userID,
					Status:      string(ConvertWorkspaceStatus(status)),
				}); err != nil {
					log.Printf("[consumer] failed start sessions: %v", err)
					continue
				}

				log.Printf("[consumer] sending to userID=%s clients=%v", userID, s.hub)
				s.hub.SendToUser(userID, map[string]any{
					"type":         "workspace.update",
					"workspace_id": wsID,
					"status":       string(ConvertWorkspaceStatus(status)),
				})

			case <-ctx.Done():
				return
			}
		}
	}()
}
