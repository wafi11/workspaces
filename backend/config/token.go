package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string


const (
	RoleUser UserRole = "user"
	RoleAdmin UserRole = "admin"
)
type TokenRequest struct {
	UserID    string
	SessionID string
	Username  string
	Exp       int
	Role string
	TokenName string
}

type TokenWorkspaceRequest struct {
	UserID  string
	Exp int
	AcessLevel string

}

type WorkspaceClaims struct {
	UserID   string `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	SessionID string `json:"session_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(c context.Context, req *TokenRequest, conf *Config) (string, error) {
	if conf == nil {
		log.Printf("config not found")
		return "", fmt.Errorf("config not found")
	}
	claims := Claims{
		UserID:   req.UserID,
		Username: req.Username,
		SessionID: req.SessionID,
		Role: req.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(req.Exp) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   req.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(conf.JWT_SECRET))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}


func GenerateTokenWorkspaces(ctx context.Context,req *TokenWorkspaceRequest,conf *Config) (string,error){
	if conf == nil {
		log.Printf("config not found")
		return "", fmt.Errorf("config not found")
	}

	claims := WorkspaceClaims{
	UserID:   req.UserID,
	Role: req.AcessLevel,
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(req.Exp) * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   req.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(conf.JWT_SECRET))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}


func ValidateTokenWorkspace(tokenStr string,conf *Config) (*WorkspaceClaims,error) {
	if conf == nil {
		log.Printf("[validate token workspaces] config not found")
		return nil,fmt.Errorf("config not found")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &WorkspaceClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JWT_SECRET), nil 
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*WorkspaceClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ValidationToken(tokenStr string, conf *Config) (*Claims, error) {

	if conf == nil {
		log.Printf("config not found")
		return nil, fmt.Errorf("config not found")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JWT_SECRET), nil 
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
