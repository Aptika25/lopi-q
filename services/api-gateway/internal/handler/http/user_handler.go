package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"api-gateway/internal/client"
	"api-gateway/internal/middleware"
	userProto "proto/user"
)

type UserHTTPHandler struct {
	userSvc userProto.UserServiceServer
}

func NewUserHTTPHandler(userSvc userProto.UserServiceServer) *UserHTTPHandler {
	return &UserHTTPHandler{userSvc: userSvc}
}

func (h *UserHTTPHandler) HandleAdminActivityLogs(w http.ResponseWriter, r *http.Request) {
	if stub, ok := h.userSvc.(*client.UserClientDirectStub); ok {
		logs, err := stub.GetActivityLogs(r.Context())
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "logs": logs})
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "logs": []interface{}{}})
}

func (h *UserHTTPHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	res, _ := h.userSvc.GetProfile(r.Context(), &userProto.GetProfileRequest{UserId: int32(userID)})
	if !res.Success {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": res.Error})
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "user": res.User})
}

func (h *UserHTTPHandler) HandleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		res, _ := h.userSvc.ListUsers(r.Context(), &userProto.ListUsersRequest{})
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "users": res.Users})
		return
	}
	if r.Method == "POST" {
		var body struct {
			NIP         string   `json:"nip"`
			Email       string   `json:"email"`
			Name        string   `json:"name"`
			Role        string   `json:"role"`
			Jabatan     string   `json:"jabatan"`
			UnitKerja   string   `json:"unit_kerja"`
			Permissions []string `json:"permissions"`
			Password    string   `json:"password"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)

		// Validate required fields
		if body.Email == "" || body.Name == "" || body.Password == "" {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Email, nama, dan password wajib diisi."})
			return
		}
		if body.Role == "" {
			body.Role = "admin"
		}
		if body.Permissions == nil {
			body.Permissions = []string{}
		}

		res, err := h.userSvc.CreateUser(r.Context(), &userProto.CreateUserRequest{
			Nip:         body.NIP,
			Email:       body.Email,
			Name:        body.Name,
			Role:        body.Role,
			Jabatan:     body.Jabatan,
			UnitKerja:   body.UnitKerja,
			Permissions: body.Permissions,
			Password:    body.Password,
		})
		if err != nil || res == nil {
			middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": false, "error": "Gagal menambahkan user."})
			return
		}
		if res.Success {
			client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "CREATE_USER", fmt.Sprintf("Admin menambahkan pengguna baru: %s (NIP: %s, Role: %s)", body.Name, body.NIP, body.Role), client.GetClientIP(r), r.UserAgent())
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": res.Success, "user": res.User, "message": res.Message, "error": res.Error})
		return
	}
}

func (h *UserHTTPHandler) HandleAdminUserDetail(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		middleware.RespondJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Not found"})
		return
	}

	idStr := pathParts[3]
	var userID int32
	_, _ = fmt.Sscanf(idStr, "%d", &userID)

	if strings.HasSuffix(r.URL.Path, "/toggle-active") && r.Method == "PUT" {
		var body struct {
			IsActive bool `json:"is_active"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)

		res, _ := h.userSvc.ToggleActive(r.Context(), &userProto.ToggleActiveRequest{Id: userID, IsActive: body.IsActive})
		client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "TOGGLE_USER_ACTIVE", fmt.Sprintf("Admin mengubah status akses user ID #%d menjadi: %v", userID, body.IsActive), client.GetClientIP(r), r.UserAgent())
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": res.Success, "user": res.User})
		return
	}

	if strings.HasSuffix(r.URL.Path, "/reset-2fa") && r.Method == "POST" {
		res, _ := h.userSvc.Reset2FA(r.Context(), &userProto.Reset2FARequest{Id: userID})
		client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "RESET_USER_2FA", fmt.Sprintf("Admin mereset 2FA Google Authenticator untuk user ID #%d", userID), client.GetClientIP(r), r.UserAgent())
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": res.Success, "message": res.Message})
		return
	}

	if r.Method == "PUT" {
		var body struct {
			NIP         string   `json:"nip"`
			Email       string   `json:"email"`
			Name        string   `json:"name"`
			Role        string   `json:"role"`
			Jabatan     string   `json:"jabatan"`
			UnitKerja   string   `json:"unit_kerja"`
			Permissions []string `json:"permissions"`
			Password    string   `json:"password"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)

		if body.Permissions == nil {
			body.Permissions = []string{}
		}

		res, _ := h.userSvc.UpdateUser(r.Context(), &userProto.UpdateUserRequest{
			Id:          userID,
			Nip:         body.NIP,
			Email:       body.Email,
			Name:        body.Name,
			Role:        body.Role,
			Jabatan:     body.Jabatan,
			UnitKerja:   body.UnitKerja,
			Permissions: body.Permissions,
			Password:    body.Password,
		})
		if res != nil && res.Success {
			client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "UPDATE_USER", fmt.Sprintf("Admin memperbarui profil pengguna: %s (NIP: %s, Role: %s)", body.Name, body.NIP, body.Role), client.GetClientIP(r), r.UserAgent())
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": res.Success, "user": res.User})
		return
	}
}
