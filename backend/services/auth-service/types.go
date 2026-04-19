package authservices

import "context"

type IServices interface {
	Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error)
	Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error)
	RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
	Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Validate(c context.Context,req string) (bool,error)
	CreatePAT(c context.Context,req *CreatePATRequest) (*CreatePATResponse,error)
	DeletePAT(c context.Context,PatId,userId string) error
	 GetAllPAT(c context.Context,userID string) ([]Pat,error)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Role 		 string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"user_id"`
	SessionId    string `json:"session_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	SessionId string `json:"session_id"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type ValidateTokenResponse struct {
	Valid    bool
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
