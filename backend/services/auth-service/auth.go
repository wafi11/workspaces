package authservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (repo *Repository) Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	tx, err := repo.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var userId,role, username string
	userId = uuid.New().String()

	query := `
		INSERT INTO users (id, username, email, password, terminal_url,role)
		VALUES ($1, $2, $3, $4, $5,'user')
		RETURNING id, username,role
	`
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	err = tx.QueryRowContext(c, query, userId, req.Username, req.Email, hashedPassword, utils.GenerateUrl(userId, models.DOMAIN)).Scan(&userId, &username,&role)
	if err != nil {
		return nil, fmt.Errorf("username or email already registered: %w", err)
	}

	const (
		maxQuota   = 5
		maxStorage = 20
		maxRam     = 4096
		maxCpu     = 4.0
		
		cpuTerminalLimit    = 0.5
		memTerminalLimitMi  = 512
		storTerminalLimitGi = 1
	)

	quotaQuery := `
		INSERT INTO user_quotas (id, user_id, max_workspaces, max_storage_gb, max_ram_mb, max_cpu_cores)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_,err = tx.ExecContext(c, quotaQuery, uuid.New().String(), userId, maxQuota, maxStorage-storTerminalLimitGi, maxRam-memTerminalLimitMi, maxCpu-cpuTerminalLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create user quota: %w", err)
	}
	sessionId := uuid.New().String()

		// Generate tokens setelah commit (di luar transaction)
	accessToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    userId,
		Username:  username,
		Role: role,
		Exp:       1,
		SessionID: sessionId,
		TokenName: "access_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    userId,
		Username:  username,
		Role: role,
		SessionID: sessionId,
		Exp:       24,
		TokenName: "refresh_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	sessionQuery := `
		INSERT INTO sessions (id, user_id, is_active, user_agent, ip_address,refresh_token)
		VALUES ($1, $2, true, $3, $4,$5)
	`
	_, err = repo.db.ExecContext(c, sessionQuery, sessionId, userId, "","",refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}


	cpuReq := fmt.Sprintf("%.2f", maxCpu)
	memReq := fmt.Sprintf("%dMi", maxRam)
	storReq := fmt.Sprintf("%dGi", maxStorage)

	messagebroker.PublishEvent(c, repo.redis, &proto.WorkspaceEnvelope{
		Payload: &proto.WorkspaceEnvelope_Create{
			Create: &proto.CreateWorkspaceEvent{
				Identity: &proto.WorkspaceIdentity{
					UserId:    userId,
					Username:  username,
					Namespace: GenerateNamespace(userId),
					Password:  req.Password,
				},
				TemplateId: "a7fda0ee-092c-40dc-be9f-8917784764b2",
				Resources: &proto.ResourceSpec{
					CpuRequest:          cpuReq,
					CpuLimit:            cpuReq,
					MemoryRequest:       memReq,
					MemoryLimit:         memReq,
					StorageRequest:      storReq,
					StorageLimit:        storReq,
					CpuTerminalLimit:    fmt.Sprintf("%.2f", cpuTerminalLimit),
					MemoryTerminalLimit: fmt.Sprintf("%dMi", memTerminalLimitMi),
				},
				EnvVars: map[string]string{
					"WS_TOKEN":         accessToken,
					"WS_REFRESH_TOKEN": refreshToken,
					"WS_API_URL":       "http://192.168.1.49:8080",
				},
				Timezone: "UTC",
				Replicas: 1,
			},
		},
	})

	return &RegisterResponse{
		UserId:  userId,
		Message: "Successfully created user",
	}, nil
}

func (repo *Repository) Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	var id, password,role, username string

	query := `
		SELECT 
			id, 
			username, 
			password,
			role
		FROM users 
		WHERE email = $1
	`

	err := repo.db.QueryRowContext(c, query, req.Email).Scan(&id, &username, &password,&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// log.Printf("login error : ")
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if !utils.VerifyPassword(password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	// generate access token
	sessionId := uuid.New().String()
	accessToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    id,
		Role: role,
		Username:  username,
		Exp:       1,
		TokenName: "access_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// generate refresh token
	refreshToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    id,
		Username:  username,
		Exp:       24,
		Role: role,
		TokenName: "refresh_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// simpan session ke database
	sessionQuery := `
		INSERT INTO sessions (id, user_id, is_active, user_agent, ip_address,refresh_token)
		VALUES ($1, $2, true, $3, $4,$5)
	`
	_, err = repo.db.ExecContext(c, sessionQuery, sessionId, id, userAgent, ipAddress,refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	redisKey := fmt.Sprintf("refresh_token:%s", id)
	err = repo.redis.Set(c, redisKey, refreshToken, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	repo.redis.Set(c, "session:"+accessToken,id, 1*time.Hour)


	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       id,
		Role: role,
		SessionId:    sessionId,
	}, nil

}

func (repo *Repository) RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
    claims, err := config.ValidationToken(req.RefreshToken, repo.conf)
    if err != nil {
        return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
    }

    // pakai session_id dari claims
    redisKey := fmt.Sprintf("refresh_token:%s", claims.SessionID)
    storedToken, err := repo.redis.Get(c, redisKey).Result()

    if err != nil || storedToken == "" {
        // fallback ke DB by session_id
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

    newAccessToken, err := config.GenerateToken(c, &config.TokenRequest{
        UserID:    claims.UserID,
        Username:  claims.Username,
        Exp:       1,
        Role:      claims.Role,
        SessionID: claims.SessionID,
        TokenName: "access_token",
    }, repo.conf)
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }

    newRefreshToken, err := config.GenerateToken(c, &config.TokenRequest{
        UserID:    claims.UserID,
        Username:  claims.Username,
        Exp:       24,
        Role:      claims.Role,
        SessionID: claims.SessionID,
        TokenName: "refresh_token",
    }, repo.conf)
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }

    // update DB by session_id
    _, err = repo.db.ExecContext(c,
        `UPDATE sessions SET refresh_token = $1, updated_at = NOW() WHERE id = $2`,
        newRefreshToken, claims.SessionID,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to update refresh token in db: %w", err)
    }

    // update redis
    err = repo.redis.Set(c, redisKey, newRefreshToken, 24*time.Hour).Err()
    if err != nil {
        return nil, fmt.Errorf("failed to store refresh token: %w", err)
    }

    return &RefreshTokenResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
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
		redisKey := fmt.Sprintf("refresh_token:%s", userId)
		repo.redis.Del(c, redisKey)
	}

	return &LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}


func (repo *Repository) Validate(ctx context.Context, token string) (bool, error) {
    val, err := repo.redis.Get(ctx, "session:"+token).Result()
    if err == redis.Nil {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return val != "", nil
}