package authservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/models"
	"github.com/wafi11/workspaces/pkg/proto"
	"github.com/wafi11/workspaces/pkg/utils"
)

type Repository struct {
	db    *sqlx.DB
	redis *redis.Client
	conf  *config.Config
}

func NewRepository(db *sqlx.DB, redis *redis.Client, conf *config.Config) *Repository {
	return &Repository{
		db:    db,
		redis: redis,
		conf:  conf,
	}
}

func (repo *Repository) Register(c context.Context, req *RegisterRequest, provider string) (*RegisterResponse, error) {
	if provider == UserProvidersLocal && req.Password == "" {
		return nil, fmt.Errorf("password must be required")
	}

	tx, err := repo.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var userId, role, username string
	var hashedPassword string
	userId = uuid.New().String()

	if provider == UserProvidersLocal {
		hashedPassword, err = utils.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to process password: %w", err)
		}
	}

	query := `
		INSERT INTO users (id, username, email, password, terminal_url,role,avatar_url)
		VALUES ($1, $2, $3, $4, $5,'user',$6)
		RETURNING id, username,role
	`

	err = tx.QueryRowContext(c, query, userId, req.Username, req.Email, hashedPassword, utils.GenerateUrl(userId, models.DOMAIN), req.AvatarURL).Scan(&userId, &username, &role)
	if err != nil {
		return nil, fmt.Errorf("username or email already registered: %w", err)
	}

	insertProvider := `
			INSERT INTO providers (id, user_id, name, provider_id)
			VALUES ($1, $2, $3, $4)
		`
	_, err = tx.ExecContext(c, insertProvider,
		uuid.New().String(), userId, provider, req.ProviderId,
	)

	if err != nil {
		log.Printf("[register]  failed register providers : %s", err.Error())
		return nil, fmt.Errorf("invalid credentials")
	}

	quotaQuery := `
		INSERT INTO user_quotas (id, user_id, max_workspaces, max_storage_gb, max_ram_mb, max_cpu_cores)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.ExecContext(c, quotaQuery, uuid.New().String(), userId, models.MaxQuota, models.MaxStorage-models.StorTerminalLimitGi, models.MaxRam-models.MemTerminalLimitMi, models.MaxCpu-models.CpuTerminalLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create user quota: %w", err)
	}
	sessionId := uuid.New().String()

	token, err := repo.GenerateToken(c, GenerateTokenReq{
		UserID:    userId,
		Role:      role,
		SessionID: sessionId,
	})

	sessionQuery := `
		INSERT INTO sessions (id, user_id, is_active, user_agent, ip_address, refresh_token)
		VALUES ($1, $2, true, $3, $4,$5)
	`
	_, err = tx.ExecContext(c, sessionQuery, sessionId, userId, "", "", token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	var templateId string

	queryTemplateId := `
		SELECT id FROM templates where category = 'terminal' and is_public = false
	`

	err = tx.QueryRowContext(c, queryTemplateId).Scan(&templateId)
	if err != nil {
		log.Printf("[register] template not found : %s", err.Error())
		return nil, fmt.Errorf("failed to register")
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	cpuReq := fmt.Sprintf("%.2f", models.MaxCpu)
	memReq := fmt.Sprintf("%dMi", models.MaxRam)
	storReq := fmt.Sprintf("%dGi", models.MaxStorage)

	messagebroker.PublishEvent(c, repo.redis, &proto.WorkspaceEnvelope{
		Payload: &proto.WorkspaceEnvelope_Create{
			Create: &proto.CreateWorkspaceEvent{
				Identity: &proto.WorkspaceIdentity{
					UserId:    userId,
					Username:  username,
					Namespace: GenerateNamespace(userId),
					Password:  req.Password,
				},
				TemplateId: templateId,
				Resources: &proto.ResourceSpec{
					CpuRequest:          cpuReq,
					CpuLimit:            cpuReq,
					MemoryRequest:       memReq,
					MemoryLimit:         memReq,
					StorageRequest:      storReq,
					StorageLimit:        storReq,
					CpuTerminalLimit:    fmt.Sprintf("%.2f", models.CpuTerminalLimit),
					MemoryTerminalLimit: fmt.Sprintf("%dMi", models.MemTerminalLimitMi),
				},
				EnvVars: map[string]string{
					"WS_TOKEN":         token.AccessToken,
					"WS_REFRESH_TOKEN": token.RefreshToken,
					"WS_API_URL":       "http://192.168.1.49:8080",
				},
				Timezone: "UTC",
				Replicas: 1,
			},
		},
	})

	return &RegisterResponse{
		UserId:       userId,
		Message:      "Successfully created user",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (repo *Repository) Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	var id, password, role, username string

	err := repo.db.QueryRowContext(c, `
		SELECT 
			u.id, 
			u.username, 
			u.password, 
			u.role
		FROM users u 
		LEFT JOIN providers p on p.user_id = u.id
		WHERE u.email = $1 and p.name = 'local'
	`, req.Email).Scan(&id, &username, &password, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if !utils.VerifyPassword(password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	sessionId := uuid.New().String()
	token, err := repo.GenerateToken(c, GenerateTokenReq{
		UserID:    id,
		Role:      role,
		SessionID: sessionId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	_, err = repo.db.ExecContext(c, `
		INSERT INTO sessions (id, user_id, is_active, user_agent, ip_address, refresh_token)
		VALUES ($1, $2, true, $3, $4, $5)
	`, sessionId, id, userAgent, ipAddress, token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err = repo.redis.Set(c, refreshTokenKey(sessionId), token.RefreshToken, 24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	if err = repo.redis.Set(c, sessionKey(sessionId), true, 1*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	return &LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserId:       id,
		Role:         role,
		SessionId:    sessionId,
	}, nil
}
func (repo *Repository) RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	claims, err := config.ValidationToken(req.RefreshToken, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
	}

	// coba dari redis dulu, fallback ke DB
	storedToken, err := repo.redis.Get(c, refreshTokenKey(claims.SessionID)).Result()
	if err != nil || storedToken == "" {
		err = repo.db.QueryRowContext(c,
			`SELECT refresh_token FROM sessions WHERE id = $1`,
			claims.SessionID,
		).Scan(&storedToken)
		if err != nil {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
	}

	if storedToken != req.RefreshToken {
		return nil, fmt.Errorf("refresh token mismatch")
	}

	token, err := repo.GenerateToken(c, GenerateTokenReq{
		UserID:    claims.UserID,
		Role:      claims.Role,
		SessionID: claims.SessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	_, err = repo.db.ExecContext(c,
		`UPDATE sessions SET refresh_token = $1, updated_at = NOW() WHERE id = $2`,
		token.RefreshToken, claims.SessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token in db: %w", err)
	}

	if err = repo.redis.Set(c, refreshTokenKey(claims.SessionID), token.RefreshToken, 24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (repo *Repository) Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// nonaktifkan session di database
	query := `
		UPDATE sessions SET is_active = false WHERE id = $1
	`
	result, err := repo.db.ExecContext(c, query, req.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to invalidate session: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("session not found")
	}

	// ambil user_id dari session untuk hapus redis
	var userId string
	err = repo.db.QueryRowContext(c,
		`SELECT user_id FROM sessions WHERE id = $1`, req.SessionId,
	).Scan(&userId)

	// hapus refresh token dari redis
	if err == nil {
		repo.redis.Del(c, refreshTokenKey(req.SessionId))
		repo.redis.Del(c, sessionKey(req.SessionId))
	}

	return &LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}

func (repo *Repository) Validate(ctx context.Context, token string) (bool, error) {
	claims, err := config.ValidationToken(token, repo.conf)
	if err != nil {
		return false, fmt.Errorf("invalid or expired refresh token: %w", err)
	}
	val, err := repo.redis.Get(ctx, refreshTokenKey(claims.SessionID)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val != "", nil
}
