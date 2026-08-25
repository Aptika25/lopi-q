package user

import "context"

type User struct {
	Id          int32    `json:"id"`
	Nip         string   `json:"nip"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Jabatan     string   `json:"jabatan"`
	UnitKerja   string   `json:"unit_kerja"`
	Permissions []string `json:"permissions"`
	TotpEnabled bool     `json:"totp_enabled"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   string   `json:"created_at"`
}

type GetProfileRequest struct {
	UserId int32 `json:"user_id"`
}

type ListUsersRequest struct{}

type ListUsersResponse struct {
	Success bool    `json:"success"`
	Users   []*User `json:"users"`
	Error   string  `json:"error"`
}

type CreateUserRequest struct {
	Nip         string   `json:"nip"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Jabatan     string   `json:"jabatan"`
	UnitKerja   string   `json:"unit_kerja"`
	Permissions []string `json:"permissions"`
	Password    string   `json:"password"`
}

type UpdateUserRequest struct {
	Id          int32    `json:"id"`
	Nip         string   `json:"nip"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Jabatan     string   `json:"jabatan"`
	UnitKerja   string   `json:"unit_kerja"`
	Permissions []string `json:"permissions"`
	IsActive    bool     `json:"is_active"`
	Password    string   `json:"password"`
}

type ToggleActiveRequest struct {
	Id       int32 `json:"id"`
	IsActive bool  `json:"is_active"`
}

type Reset2FARequest struct {
	Id int32 `json:"id"`
}

type Reset2FAResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type UserResponse struct {
	Success bool   `json:"success"`
	User    *User  `json:"user"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type UnimplementedUserServiceServer struct{}

type UserServiceServer interface {
	GetProfile(context.Context, *GetProfileRequest) (*UserResponse, error)
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*UserResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UserResponse, error)
	ToggleActive(context.Context, *ToggleActiveRequest) (*UserResponse, error)
	Reset2FA(context.Context, *Reset2FARequest) (*Reset2FAResponse, error)
}

func RegisterUserServiceServer(s interface{}, srv UserServiceServer) {
	_ = s
	_ = srv
}
