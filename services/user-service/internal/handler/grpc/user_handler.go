package grpc

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"proto/user"
	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserHandler struct {
	user.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) toProtoUser(u *model.User) *user.User {
	perms := u.Permissions
	if perms == nil {
		perms = []string{}
	}
	return &user.User{
		Id:          int32(u.ID),
		Nip:         u.NIP,
		Email:       u.Email,
		Name:        u.Name,
		Role:        u.Role,
		Jabatan:     u.Jabatan,
		UnitKerja:   u.UnitKerja,
		Permissions: perms,
		TotpEnabled: u.TotpEnabled,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserHandler) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.UserResponse, error) {
	u := h.repo.FindByID(int(req.UserId))
	if u == nil {
		return &user.UserResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}
	return &user.UserResponse{Success: true, User: h.toProtoUser(u)}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	users := h.repo.ListAll()
	result := make([]*user.User, 0, len(users))
	for i := range users {
		result = append(result, h.toProtoUser(&users[i]))
	}
	return &user.ListUsersResponse{Success: true, Users: result}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	perms := req.Permissions
	if perms == nil {
		perms = []string{}
	}
	u := h.repo.Create(model.User{
		NIP:          req.Nip,
		Email:        req.Email,
		Name:         req.Name,
		Role:         req.Role,
		Jabatan:      req.Jabatan,
		UnitKerja:    req.UnitKerja,
		Permissions:  perms,
		PasswordHash: string(hash),
		IsActive:     true,
	})
	return &user.UserResponse{Success: true, User: h.toProtoUser(&u), Message: "User berhasil dibuat."}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	u := model.User{
		ID:           int(req.Id),
		Email:        req.Email,
		Name:         req.Name,
		Role:         req.Role,
		Jabatan:      req.Jabatan,
		UnitKerja:    req.UnitKerja,
		Permissions:  req.Permissions,
		PasswordHash: req.Password,
	}
	ok := h.repo.Update(&u)
	if !ok {
		return &user.UserResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}
	return &user.UserResponse{Success: true, User: h.toProtoUser(&u), Message: "Data user diperbarui."}, nil
}

func (h *UserHandler) ToggleActive(ctx context.Context, req *user.ToggleActiveRequest) (*user.UserResponse, error) {
	u := h.repo.ToggleActive(int(req.Id), req.IsActive)
	if u == nil {
		return &user.UserResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}
	return &user.UserResponse{Success: true, User: h.toProtoUser(u)}, nil
}

func (h *UserHandler) Reset2FA(ctx context.Context, req *user.Reset2FARequest) (*user.Reset2FAResponse, error) {
	ok := h.repo.Reset2FA(int(req.Id))
	if !ok {
		return &user.Reset2FAResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}
	return &user.Reset2FAResponse{Success: true, Message: "2FA berhasil direset."}, nil
}
