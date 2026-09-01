package grpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"auth-service/internal/repository"
	"pkg/jwt"
	"proto/auth"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	repo *repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func generateBackupCodes(count int) []string {
	var codes []string
	for i := 0; i < count; i++ {
		bytes := make([]byte, 4)
		_, _ = rand.Read(bytes)
		codes = append(codes, hex.EncodeToString(bytes))
	}
	return codes
}

func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	u := h.repo.FindByIdentifier(req.Identifier)
	if u == nil {
		return &auth.LoginResponse{Success: false, Error: "NIP / Email Dinas atau password salah."}, nil
	}

	if !u.IsActive {
		return &auth.LoginResponse{Success: false, Error: "Akun Anda sedang non-aktif."}, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		return &auth.LoginResponse{Success: false, Error: "NIP / Email Dinas atau password salah."}, nil
	}

	tempToken, _ := jwt.GenerateTempToken(u.ID, 15*time.Minute)

	if u.TotpEnabled || u.TotpSecret != "" {
		return &auth.LoginResponse{
			Success:          true,
			OtpRequired:      true,
			OtpSetupRequired: false,
			TempToken:        tempToken,
		}, nil
	}

	// Mandatory 2FA Setup for accounts without 2FA configured yet
	return &auth.LoginResponse{
		Success:          true,
		OtpRequired:      true,
		OtpSetupRequired: true,
		TempToken:        tempToken,
	}, nil
}

func (h *AuthHandler) Verify2FA(ctx context.Context, req *auth.Verify2FARequest) (*auth.Verify2FAResponse, error) {
	claims, err := jwt.ValidateToken(req.TempToken)
	if err != nil || !claims.IsTemp {
		return &auth.Verify2FAResponse{Success: false, Error: "Sesi verifikasi expired. Silakan login ulang."}, nil
	}

	user := h.repo.FindByID(claims.UserID)
	if user == nil {
		return &auth.Verify2FAResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}

	cleanCode := strings.TrimSpace(req.Code)
	isValid := false

	if user.TotpSecret != "" && totp.Validate(cleanCode, user.TotpSecret) {
		isValid = true
	} else if len(cleanCode) == 8 {
		// Check backup codes
		for i, bCode := range user.BackupCodes {
			if bCode == cleanCode {
				isValid = true
				// Remove used backup code
				user.BackupCodes = append(user.BackupCodes[:i], user.BackupCodes[i+1:]...)
				h.repo.Save(user)
				break
			}
		}
	}

	if !isValid {
		return &auth.Verify2FAResponse{Success: false, Error: "Kode verifikasi 2FA atau Kode Backup tidak valid."}, nil
	}

	token, _ := jwt.GenerateToken(user.ID, user.NIP, user.Role, 24*time.Hour)
	return &auth.Verify2FAResponse{
		Success: true,
		Token:   token,
		UserId:  int32(user.ID),
		Role:    user.Role,
		Nip:     user.NIP,
		Name:    user.Name,
	}, nil
}

func (h *AuthHandler) Setup2FA(ctx context.Context, req *auth.Setup2FARequest) (*auth.Setup2FAResponse, error) {
	userID := int(req.UserId)
	if userID == 0 && req.TempToken != "" {
		if claims, err := jwt.ValidateToken(req.TempToken); err == nil && claims.IsTemp {
			userID = claims.UserID
		}
	}

	user := h.repo.FindByID(userID)
	if user == nil {
		return &auth.Setup2FAResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}

	accountIdentifier := user.Email
	if accountIdentifier == "" {
		accountIdentifier = user.NIP
	}

	secret := user.TotpSecret
	if secret == "" {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "LOPI-Q BULUKUMBA",
			AccountName: accountIdentifier,
		})
		if err != nil {
			return &auth.Setup2FAResponse{Success: false, Error: "Gagal generate secret 2FA."}, nil
		}
		secret = key.Secret()
		user.TotpSecret = secret
		h.repo.Save(user)
	}

	escapedAccount := url.QueryEscape(accountIdentifier)
	otpUrl := fmt.Sprintf("otpauth://totp/LOPI-Q%%20BULUKUMBA:%s?secret=%s&issuer=LOPI-Q%%20BULUKUMBA", escapedAccount, secret)

	return &auth.Setup2FAResponse{
		Success:    true,
		Secret:     secret,
		OtpauthUrl: otpUrl,
	}, nil
}

func (h *AuthHandler) Enable2FA(ctx context.Context, req *auth.Enable2FARequest) (*auth.Enable2FAResponse, error) {
	userID := int(req.UserId)
	if userID == 0 && req.TempToken != "" {
		if claims, err := jwt.ValidateToken(req.TempToken); err == nil && claims.IsTemp {
			userID = claims.UserID
		}
	}

	user := h.repo.FindByID(userID)
	secretToUse := req.Secret
	if secretToUse == "" {
		secretToUse = user.TotpSecret
	}

	if user == nil || secretToUse == "" {
		return &auth.Enable2FAResponse{Success: false, Error: "Setup 2FA belum diinisiasi."}, nil
	}

	if !totp.Validate(strings.TrimSpace(req.Code), secretToUse) {
		return &auth.Enable2FAResponse{Success: false, Error: "Kode verifikasi 2FA tidak cocok."}, nil
	}

	backupCodes := generateBackupCodes(4)
	user.TotpSecret = secretToUse
	user.TotpEnabled = true
	user.BackupCodes = backupCodes
	h.repo.Save(user)

	return &auth.Enable2FAResponse{
		Success:     true,
		BackupCodes: backupCodes,
	}, nil
}

func (h *AuthHandler) Disable2FA(ctx context.Context, req *auth.Disable2FARequest) (*auth.Disable2FAResponse, error) {
	user := h.repo.FindByID(int(req.UserId))
	if user == nil {
		return &auth.Disable2FAResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}

	user.TotpEnabled = false
	user.TotpSecret = ""
	user.BackupCodes = nil
	h.repo.Save(user)

	return &auth.Disable2FAResponse{Success: true}, nil
}

func (h *AuthHandler) SelfReset2FA(ctx context.Context, req *auth.SelfReset2FARequest) (*auth.SelfReset2FAResponse, error) {
	claims, err := jwt.ValidateToken(req.TempToken)
	if err != nil || !claims.IsTemp {
		return &auth.SelfReset2FAResponse{Success: false, Error: "Sesi expired."}, nil
	}

	user := h.repo.FindByID(claims.UserID)
	if user == nil {
		return &auth.SelfReset2FAResponse{Success: false, Error: "User tidak ditemukan."}, nil
	}

	cleanCode := strings.TrimSpace(req.BackupCode)
	isValid := false

	for i, bCode := range user.BackupCodes {
		if bCode == cleanCode {
			isValid = true
			user.BackupCodes = append(user.BackupCodes[:i], user.BackupCodes[i+1:]...)
			break
		}
	}

	if !isValid {
		return &auth.SelfReset2FAResponse{Success: false, Error: "Kode backup pemulihan tidak valid."}, nil
	}

	user.TotpEnabled = false
	user.TotpSecret = ""
	h.repo.Save(user)

	return &auth.SelfReset2FAResponse{Success: true}, nil
}
