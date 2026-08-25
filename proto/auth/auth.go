package auth

import "context"

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Success          bool   `json:"success"`
	OtpRequired      bool   `json:"otp_required"`
	OtpSetupRequired bool   `json:"otp_setup_required"`
	TempToken        string `json:"temp_token"`
	Token            string `json:"token"`
	UserId           int32  `json:"user_id"`
	Role             string `json:"role"`
	Nip              string `json:"nip"`
	Name             string `json:"name"`
	Error            string `json:"error"`
}

type Verify2FARequest struct {
	TempToken string `json:"temp_token"`
	Code      string `json:"code"`
}

type Verify2FAResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	UserId  int32  `json:"user_id"`
	Role    string `json:"role"`
	Nip     string `json:"nip"`
	Name    string `json:"name"`
	Error   string `json:"error"`
}

type Setup2FARequest struct {
	UserId    int32  `json:"user_id"`
	TempToken string `json:"temp_token"`
}

type Setup2FAResponse struct {
	Success    bool   `json:"success"`
	Secret     string `json:"secret"`
	OtpauthUrl string `json:"otpauth_url"`
	QrCode     string `json:"qr_code"`
	Error      string `json:"error"`
}

type Enable2FARequest struct {
	UserId    int32  `json:"user_id"`
	Code      string `json:"code"`
	Secret    string `json:"secret"`
	TempToken string `json:"temp_token"`
}

type Enable2FAResponse struct {
	Success     bool     `json:"success"`
	BackupCodes []string `json:"backup_codes"`
	Error       string   `json:"error"`
}

type Disable2FARequest struct {
	UserId int32 `json:"user_id"`
}

type Disable2FAResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SelfReset2FARequest struct {
	TempToken  string `json:"temp_token"`
	BackupCode string `json:"backup_code"`
}

type SelfReset2FAResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type UnimplementedAuthServiceServer struct{}

func (s *UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}

func (s *UnimplementedAuthServiceServer) Verify2FA(context.Context, *Verify2FARequest) (*Verify2FAResponse, error) {
	return nil, nil
}

func (s *UnimplementedAuthServiceServer) Setup2FA(context.Context, *Setup2FARequest) (*Setup2FAResponse, error) {
	return nil, nil
}

func (s *UnimplementedAuthServiceServer) Enable2FA(context.Context, *Enable2FARequest) (*Enable2FAResponse, error) {
	return nil, nil
}

func (s *UnimplementedAuthServiceServer) Disable2FA(context.Context, *Disable2FARequest) (*Disable2FAResponse, error) {
	return nil, nil
}

func (s *UnimplementedAuthServiceServer) SelfReset2FA(context.Context, *SelfReset2FARequest) (*SelfReset2FAResponse, error) {
	return nil, nil
}

type AuthServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Verify2FA(context.Context, *Verify2FARequest) (*Verify2FAResponse, error)
	Setup2FA(context.Context, *Setup2FARequest) (*Setup2FAResponse, error)
	Enable2FA(context.Context, *Enable2FARequest) (*Enable2FAResponse, error)
	Disable2FA(context.Context, *Disable2FARequest) (*Disable2FAResponse, error)
	SelfReset2FA(context.Context, *SelfReset2FARequest) (*SelfReset2FAResponse, error)
}

func RegisterAuthServiceServer(s interface{}, srv AuthServiceServer) {
	_ = s
	_ = srv
}
