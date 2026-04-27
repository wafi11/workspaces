package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wafi11/workspaces/core/auth-services/config"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/internal"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
)

var globalDB *sqlx.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	os.Setenv("JWT_SECRET", "test-secret")

	// 1. Jalankan Container Postgres SEKALI SAJA untuk semua test
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("workspaces_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	// 2. Hubungkan ke DB ke variabel global
	globalDB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// 3. Jalankan Migrasi awal
	globalDB.MustExec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            username TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            role TEXT NOT NULL DEFAULT 'user',
            created_at TIMESTAMPTZ DEFAULT now()
        );
        CREATE TABLE IF NOT EXISTS sessions (
            id UUID PRIMARY KEY,
            user_id UUID NOT NULL REFERENCES users(id),
            refresh_token TEXT NOT NULL,
            created_at TIMESTAMPTZ DEFAULT now()
        );
    `)

	code := m.Run()

	globalDB.Close()
	pgContainer.Terminate(ctx)

	os.Exit(code)
}

func clearData(t *testing.T) {
	_, err := globalDB.Exec("TRUNCATE TABLE sessions, users CASCADE")
	require.NoError(t, err)
}

func TestRepository_Register_Success(t *testing.T) {
	clearData(t)
	repo := internal.NewRepository(nil, globalDB)

	resp, err := repo.Register(context.Background(), &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	})

	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Successfully Register", resp.Message)
}

func TestRepository_Register_PasswordHashed(t *testing.T) {
	clearData(t)
	repo := internal.NewRepository(nil, globalDB)

	_, err := repo.Register(context.Background(), &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	})
	require.NoError(t, err)

	var storedPassword string
	err = globalDB.Get(&storedPassword, "SELECT password FROM users WHERE email = $1", "wafi@example.com")
	assert.NoError(t, err)
	assert.NotEqual(t, "Secret123", storedPassword)
}

func TestRepository_Register_DefaultRole(t *testing.T) {
	clearData(t)
	repo := internal.NewRepository(nil, globalDB)

	repo.Register(context.Background(), &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	})

	var role string
	err := globalDB.Get(&role, "SELECT role FROM users WHERE email = $1", "wafi@example.com")
	if err != nil {
		log.Printf("role not found")
		return
	}
	assert.Equal(t, "user", role)
}

func TestRepository_Register_EmailDuplicate(t *testing.T) {
	clearData(t)
	repo := internal.NewRepository(nil, globalDB)
	req := &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	}

	repo.Register(context.Background(), req)
	resp, err := repo.Register(context.Background(), req)

	assert.ErrorIs(t, err, pkg.ErrEmailAlreadyExist)
	assert.Nil(t, resp)
}

func TestRepository_Login_Success(t *testing.T) {
	clearData(t)
	repo := internal.NewRepository(&config.Config{
		JWT_SECRET: "KAKLAKMAKAKAKAMAKMK",
	}, globalDB)

	// Persiapan user
	_, err := repo.Register(context.Background(), &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	})
	require.NoError(t, err)

	resp, err := repo.Login(context.Background(), &v1.LoginRequest{
		Email:    "wafi@example.com",
		Password: "Secret123",
	})

	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestRepository_Login_WrongPassword(t *testing.T) {
	clearData(t)
	repo := internal.NewRepository(nil, globalDB)

	repo.Register(context.Background(), &v1.RegisterRequest{
		Email:    "wafi@example.com",
		Password: "Secret123",
		Username: "wafi11",
	})

	resp, err := repo.Login(context.Background(), &v1.LoginRequest{
		Email:    "wafi@example.com",
		Password: "WrongPass1",
	})

	assert.ErrorIs(t, err, pkg.ErrInvalidCredentials)
	assert.Nil(t, resp)
}
