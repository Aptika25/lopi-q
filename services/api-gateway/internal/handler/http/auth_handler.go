package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"api-gateway/internal/client"
	"api-gateway/internal/middleware"
	authProto "proto/auth"
)

type AuthHTTPHandler struct {
	authSvc authProto.AuthServiceServer
}

func NewAuthHTTPHandler(authSvc authProto.AuthServiceServer) *AuthHTTPHandler {
	return &AuthHTTPHandler{authSvc: authSvc}
}

func (h *AuthHTTPHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Identifier string `json:"identifier"`
		NIP        string `json:"nip"`
		Email      string `json:"email"`
		Password   string `json:"password"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	id := body.Identifier
	if id == "" {
		if body.Email != "" {
			id = body.Email
		} else {
			id = body.NIP
		}
	}

	res, err := h.authSvc.Login(r.Context(), &authProto.LoginRequest{
		Identifier: id,
		Password:   body.Password,
	})
	if err != nil || res == nil || !res.Success {
		errMsg := "NIP / Email Dinas atau password salah."
		if res != nil && res.Error != "" {
			errMsg = res.Error
		}
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": errMsg})
		return
	}
	if res.OtpRequired {
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success":            true,
			"otp_required":       true,
			"otp_setup_required": res.OtpSetupRequired,
			"temp_token":         res.TempToken,
		})
		return
	}
	client.RecordActivityLog(int(res.UserId), res.Nip, res.Name, "LOGIN_SUCCESS", "Login awal pengguna ke sistem LOPI-Q", client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   res.Token,
		"user": map[string]interface{}{
			"id":     res.UserId,
			"nip":    res.Nip,
			"name":   res.Name,
			"role":   res.Role,
			"email":  id,
		},
	})
}

func (h *AuthHTTPHandler) HandleVerify2FA(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TempToken string `json:"temp_token"`
		Code      string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	res, err := h.authSvc.Verify2FA(r.Context(), &authProto.Verify2FARequest{
		TempToken: body.TempToken,
		Code:      body.Code,
	})
	if err != nil || res == nil || !res.Success {
		errMsg := "Kode verifikasi salah atau expired."
		if res != nil && res.Error != "" {
			errMsg = res.Error
		}
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": errMsg})
		return
	}
	client.RecordActivityLog(int(res.UserId), res.Nip, res.Name, "VERIFY_2FA_SUCCESS", "Berhasil verifikasi 6-digit OTP Google Authenticator saat login", client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   res.Token,
		"user": map[string]interface{}{
			"id":   res.UserId,
			"nip":  res.Nip,
			"name": res.Name,
			"role": res.Role,
		},
	})
}

func (h *AuthHTTPHandler) HandleSetup2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)
	authHeader := r.Header.Get("Authorization")
	tempToken := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tempToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	res, err := h.authSvc.Setup2FA(r.Context(), &authProto.Setup2FARequest{
		UserId:    int32(userID),
		TempToken: tempToken,
	})
	if err != nil || res == nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal setup 2FA"})
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":     res.Success,
		"secret":      res.Secret,
		"otpauth_url": res.OtpauthUrl,
		"qr_code":     res.QrCode,
		"error":       res.Error,
	})
}

func (h *AuthHTTPHandler) HandleEnable2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)
	var body struct {
		Code      string `json:"code"`
		Secret    string `json:"secret"`
		TempToken string `json:"temp_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	authHeader := r.Header.Get("Authorization")
	tempToken := body.TempToken
	if tempToken == "" && strings.HasPrefix(authHeader, "Bearer ") {
		tempToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	res, err := h.authSvc.Enable2FA(r.Context(), &authProto.Enable2FARequest{
		UserId:    int32(userID),
		Code:      body.Code,
		Secret:    body.Secret,
		TempToken: tempToken,
	})
	if err != nil || res == nil || !res.Success {
		errMsg := "Kode verifikasi 2FA tidak cocok."
		if res != nil && res.Error != "" {
			errMsg = res.Error
		}
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": errMsg})
		return
	}
	client.RecordActivityLog(userID, "", "", "ENABLE_2FA", "Pengguna berhasil mengaktifkan 2FA Google Authenticator", client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":      true,
		"backup_codes": res.BackupCodes,
	})
}

func (h *AuthHTTPHandler) HandleDisable2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)
	res, _ := h.authSvc.Disable2FA(r.Context(), &authProto.Disable2FARequest{UserId: int32(userID)})
	client.RecordActivityLog(userID, "", "", "DISABLE_2FA", "Pengguna menonaktifkan 2FA Google Authenticator", client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": res.Success})
}

func (h *AuthHTTPHandler) HandleSelfReset2FA(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TempToken  string `json:"temp_token"`
		BackupCode string `json:"backup_code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	res, err := h.authSvc.SelfReset2FA(r.Context(), &authProto.SelfReset2FARequest{
		TempToken:  body.TempToken,
		BackupCode: body.BackupCode,
	})
	if err != nil || res == nil || !res.Success {
		errMsg := "Kode backup pemulihan tidak valid."
		if res != nil && res.Error != "" {
			errMsg = res.Error
		}
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": errMsg})
		return
	}
	client.RecordActivityLog(0, "", "", "SELF_RESET_2FA", "Pengguna melakukan reset 2FA via Kode Pemulihan Backup", client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}
