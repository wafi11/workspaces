package models

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
type UserQuota struct {
	ID               string  `db:"id" json:"id"`
	UserID           string  `db:"user_id" json:"userId"`
	MaxWorkspaces    int     `db:"max_workspaces" json:"maxWorkspaces"`
	MaxStorageGB     int     `db:"max_storage_gb" json:"maxStorageGb"`
	MaxRamMB         int     `db:"max_ram_mb" json:"maxRamMb"`
	MaxCpuCores      float64 `db:"max_cpu_cores" json:"maxCpuCores"`
	UsedWorkspaces   int     `db:"used_workspaces" json:"used_workspaces"`
	UsedCpuCores     float64 `db:"used_cpu_cores" json:"used_cpu_cores"`
	UsedRamCores     float64 `db:"used_ram_mb" json:"used_ram_mb"`
	UsedStorageCores float64 `db:"used_storage_gb" json:"used_storage_gb"`
}
