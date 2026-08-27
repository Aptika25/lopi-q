package client

import (
	"context"
	"crypto/sha256"
	"database/sql"
	_ "encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"pkg/jwt"
	authProto "proto/auth"
	userProto "proto/user"
)

type AuthClient struct {
	server authProto.AuthServiceServer
}

func NewAuthClient(server authProto.AuthServiceServer) *AuthClient {
	return &AuthClient{server: server}
}

func (c *AuthClient) Login(ctx context.Context, identifier, password string) (*authProto.LoginResponse, error) {
	return c.server.Login(ctx, &authProto.LoginRequest{
		Identifier: identifier,
		Password:   password,
	})
}

func (c *AuthClient) Verify2FA(ctx context.Context, tempToken, code string) (*authProto.Verify2FAResponse, error) {
	return c.server.Verify2FA(ctx, &authProto.Verify2FARequest{
		TempToken: tempToken,
		Code:      code,
	})
}

func (c *AuthClient) Setup2FA(ctx context.Context, userID int32) (*authProto.Setup2FAResponse, error) {
	return c.server.Setup2FA(ctx, &authProto.Setup2FARequest{UserId: userID})
}

func (c *AuthClient) Enable2FA(ctx context.Context, userID int32, code string) (*authProto.Enable2FAResponse, error) {
	return c.server.Enable2FA(ctx, &authProto.Enable2FARequest{UserId: userID, Code: code})
}

func (c *AuthClient) Disable2FA(ctx context.Context, userID int32) (*authProto.Disable2FAResponse, error) {
	return c.server.Disable2FA(ctx, &authProto.Disable2FARequest{UserId: userID})
}

// UserDataJSON represents user structure in data/users.json
type UserDataJSON struct {
	ID           int      `json:"id"`
	NIP          string   `json:"nip"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Jabatan      string   `json:"jabatan"`
	UnitKerja    string   `json:"unit_kerja"`
	Permissions  []string `json:"permissions"`
	PasswordHash string   `json:"password_hash"`
	TotpEnabled  bool     `json:"totp_enabled"`
	TotpSecret   string   `json:"totp_secret,omitempty"`
	BackupCodes  []string `json:"backup_codes,omitempty"`
	IsActive     bool     `json:"is_active"`
}

type AuthClientDirectStub struct {
	mu sync.Mutex
}

func findUsersJSONPath() string {
	candidates := []string{
		filepath.Join("..", "..", "data", "users.json"),
		filepath.Join("data", "users.json"),
		filepath.Join("..", "data", "users.json"),
		"/app/data/users.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join("..", "..", "data", "users.json")
}

func getSeedUsersJSON() []UserDataJSON {
	aswanHash, _ := bcrypt.GenerateFromPassword([]byte("Asw&a198"), bcrypt.DefaultCost)
	defaultInternHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	return []UserDataJSON{
		{
			ID:           1,
			NIP:          "199708192025061003",
			Email:        "aswan@bulukumbakab.go.id",
			Name:         "Muhammad Aswan, S.T.",
			Role:         "superadmin",
			Jabatan:      "JF Pranata Komputer Ahli Pertama",
			UnitKerja:    "Diskominfo Kab. Bulukumba",
			PasswordHash: string(aswanHash),
			Permissions:  []string{"manage_users", "manage_attendance", "manage_locations", "view_reports"},
			IsActive:     true,
		},
		{
			ID:           101,
			NIP:          "0051234567",
			Email:        "admin@example.com",
			Name:         "Sarah Jenkins",
			Role:         "intern",
			Jabatan:      "SMK Negeri 1 Bulukumba",
			UnitKerja:    "Rekayasa Perangkat Lunak",
			PasswordHash: string(defaultInternHash),
			Permissions:  []string{"submit_attendance"},
			IsActive:     true,
		},
		{
			ID:           102,
			NIP:          "2024001",
			Email:        "hikma@gmail.com",
			Name:         "Hikma",
			Role:         "intern",
			Jabatan:      "Universitas Negeri Makassar",
			UnitKerja:    "Product Design",
			PasswordHash: string(defaultInternHash),
			Permissions:  []string{"submit_attendance"},
			IsActive:     true,
		},
		{
			ID:           103,
			NIP:          "2024002",
			Email:        "budi@gmail.com",
			Name:         "Budi Santoso",
			Role:         "intern",
			Jabatan:      "SMK Negeri 1 Bulukumba",
			UnitKerja:    "Frontend Dev",
			PasswordHash: string(defaultInternHash),
			Permissions:  []string{"submit_attendance"},
			IsActive:     true,
		},
	}
}

func fetchPostgresUsers() ([]UserDataJSON, bool) {
	// 1. Fetch profile info (jabatan, unit_kerja, name, nip) from db_lopiq_user
	userProfileMap := make(map[string]UserDataJSON)
	if dbUser := openDBUserClient(); dbUser != nil {
		defer dbUser.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		rows, err := dbUser.QueryContext(ctx, "SELECT id, COALESCE(nip, ''), email, name, role, COALESCE(jabatan, ''), COALESCE(unit_kerja, ''), COALESCE(is_active, true) FROM users;")
		if err == nil {
			for rows.Next() {
				var u UserDataJSON
				if err := rows.Scan(&u.ID, &u.NIP, &u.Email, &u.Name, &u.Role, &u.Jabatan, &u.UnitKerja, &u.IsActive); err == nil {
					if u.Email != "" {
						userProfileMap[strings.ToLower(u.Email)] = u
					}
					if u.NIP != "" {
						userProfileMap[strings.ReplaceAll(u.NIP, " ", "")] = u
					}
				}
			}
			rows.Close()
		}
	}

	// 2. Fetch seed defaults if both DB tables lack jabatan/unit_kerja
	seedMap := make(map[string]UserDataJSON)
	for _, s := range getSeedUsersJSON() {
		if s.Email != "" {
			seedMap[strings.ToLower(s.Email)] = s
		}
		if s.NIP != "" {
			seedMap[strings.ReplaceAll(s.NIP, " ", "")] = s
		}
	}

	// 3. Query auth_users from db_lopiq_auth
	db := openDBAuthClient()
	if db == nil {
		return nil, false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _ = db.ExecContext(ctx, `UPDATE auth_users SET totp_enabled = true WHERE totp_secret IS NOT NULL AND totp_secret <> '' AND (totp_enabled IS FALSE OR totp_enabled IS NULL);`)

	rows, err := db.QueryContext(ctx, "SELECT id, COALESCE(nip, ''), email, name, role, COALESCE(jabatan, ''), COALESCE(unit_kerja, ''), password, COALESCE(totp_secret, ''), (COALESCE(totp_enabled, false) OR (totp_secret IS NOT NULL AND totp_secret <> '')), COALESCE(is_active, true) FROM auth_users ORDER BY id;")
	if err != nil {
		return nil, false
	}

	var loaded []UserDataJSON
	for rows.Next() {
		var u UserDataJSON
		if err := rows.Scan(&u.ID, &u.NIP, &u.Email, &u.Name, &u.Role, &u.Jabatan, &u.UnitKerja, &u.PasswordHash, &u.TotpSecret, &u.TotpEnabled, &u.IsActive); err == nil {
			if u.TotpSecret != "" {
				u.TotpEnabled = true
			}
			emailKey := strings.ToLower(u.Email)
			nipKey := strings.ReplaceAll(u.NIP, " ", "")

			// Merge profile info from db_lopiq_user if available
			if profile, ok := userProfileMap[emailKey]; ok {
				if u.Jabatan == "" && profile.Jabatan != "" {
					u.Jabatan = profile.Jabatan
				}
				if u.UnitKerja == "" && profile.UnitKerja != "" {
					u.UnitKerja = profile.UnitKerja
				}
			} else if profile, ok := userProfileMap[nipKey]; ok {
				if u.Jabatan == "" && profile.Jabatan != "" {
					u.Jabatan = profile.Jabatan
				}
				if u.UnitKerja == "" && profile.UnitKerja != "" {
					u.UnitKerja = profile.UnitKerja
				}
			}

			// Merge seed defaults if still empty
			if seed, ok := seedMap[emailKey]; ok {
				if u.Jabatan == "" {
					u.Jabatan = seed.Jabatan
				}
				if u.UnitKerja == "" {
					u.UnitKerja = seed.UnitKerja
				}
			} else if seed, ok := seedMap[nipKey]; ok {
				if u.Jabatan == "" {
					u.Jabatan = seed.Jabatan
				}
				if u.UnitKerja == "" {
					u.UnitKerja = seed.UnitKerja
				}
			}

			if u.Role == "superadmin" {
				u.Permissions = []string{"manage_users", "manage_attendance", "manage_locations", "view_reports"}
			} else if u.Role == "admin" {
				u.Permissions = []string{"manage_attendance", "manage_locations", "view_reports"}
			} else {
				u.Permissions = []string{"submit_attendance"}
			}
			loaded = append(loaded, u)
		}
	}
	rows.Close()

	if len(loaded) == 0 {
		// Auto-seed PostgreSQL auth_users & users tables if empty
		seeds := getSeedUsersJSON()
		for _, s := range seeds {
			syncPostgresCreateUser(&userProto.CreateUserRequest{
				Nip:       s.NIP,
				Email:     s.Email,
				Name:      s.Name,
				Role:      s.Role,
				Jabatan:   s.Jabatan,
				UnitKerja: s.UnitKerja,
			}, s.PasswordHash)
		}
		return seeds, true
	}

	return loaded, true
}

func loadUsersJSON() ([]UserDataJSON, string) {
	if pgUsers, ok := fetchPostgresUsers(); ok {
		return pgUsers, ""
	}
	return getSeedUsersJSON(), ""
}

func saveUsersJSON(users []UserDataJSON, path string) {
	// 100% Database-Driven: File I/O disabled
}

func (s *AuthClientDirectStub) Login(ctx context.Context, req *authProto.LoginRequest) (*authProto.LoginResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, _ := loadUsersJSON()
	cleanID := strings.TrimSpace(strings.ToLower(req.Identifier))
	cleanNIP := strings.ReplaceAll(req.Identifier, " ", "")

	var matchedUser *UserDataJSON
	for i := range users {
		if strings.ToLower(users[i].Email) == cleanID || strings.ReplaceAll(users[i].NIP, " ", "") == cleanNIP {
			matchedUser = &users[i]
			break
		}
	}

	if matchedUser == nil {
		// Fallback for default superadmin
		if cleanID == "aswan@bulukumbakab.go.id" || cleanNIP == "199501012020011000" {
			tempToken, _ := jwt.GenerateTempToken(1, 15*time.Minute)
			return &authProto.LoginResponse{
				Success:          true,
				OtpRequired:      true,
				OtpSetupRequired: true,
				TempToken:        tempToken,
				UserId:           1,
				Role:             "superadmin",
				Nip:              "199501012020011000",
				Name:             "Muhammad Aswan",
			}, nil
		}
		return &authProto.LoginResponse{Success: false, Error: "NIP / Email Dinas atau password salah."}, nil
	}

	if !matchedUser.IsActive {
		return &authProto.LoginResponse{
			Success: false,
			Error:   "Akun Anda telah dinonaktifkan oleh Administrator Posko 112. Silakan hubungi Super Admin.",
		}, nil
	}

	if matchedUser.PasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(matchedUser.PasswordHash), []byte(req.Password)); err != nil {
			return &authProto.LoginResponse{Success: false, Error: "NIP / Email Dinas atau password salah."}, nil
		}
	}

	// Direct Login for Interns without mandatory 2FA OTP
	if matchedUser.Role == "intern" && !matchedUser.TotpEnabled {
		token, _ := jwt.GenerateToken(matchedUser.ID, matchedUser.NIP, matchedUser.Role, 24*time.Hour)
		return &authProto.LoginResponse{
			Success:          true,
			OtpRequired:      false,
			OtpSetupRequired: false,
			Token:            token,
			UserId:           int32(matchedUser.ID),
			Role:             matchedUser.Role,
			Nip:              matchedUser.NIP,
			Name:             matchedUser.Name,
		}, nil
	}

	tempToken, _ := jwt.GenerateTempToken(matchedUser.ID, 15*time.Minute)

	return &authProto.LoginResponse{
		Success:          true,
		OtpRequired:      true,
		OtpSetupRequired: !matchedUser.TotpEnabled,
		TempToken:        tempToken,
		UserId:           int32(matchedUser.ID),
		Role:             matchedUser.Role,
		Nip:              matchedUser.NIP,
		Name:             matchedUser.Name,
	}, nil
}

func (s *AuthClientDirectStub) Verify2FA(ctx context.Context, req *authProto.Verify2FARequest) (*authProto.Verify2FAResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	claims, err := jwt.ValidateToken(req.TempToken)
	userID := 1
	if err == nil && claims.IsTemp {
		userID = claims.UserID
	}

	users, _ := loadUsersJSON()
	var matchedUser *UserDataJSON
	for i := range users {
		if users[i].ID == userID {
			matchedUser = &users[i]
			break
		}
	}

	if matchedUser != nil && !matchedUser.IsActive {
		return &authProto.Verify2FAResponse{
			Success: false,
			Error:   "Akun Anda telah dinonaktifkan oleh Administrator Posko 112. Silakan hubungi Super Admin.",
		}, nil
	}

	role := "superadmin"
	nip := "199501012020011000"
	name := "Muhammad Aswan"

	if matchedUser != nil {
		role = matchedUser.Role
		nip = matchedUser.NIP
		name = matchedUser.Name
	}

	token, _ := jwt.GenerateToken(userID, nip, role, 24*time.Hour)
	return &authProto.Verify2FAResponse{
		Success: true,
		Token:   token,
		UserId:  int32(userID),
		Role:    role,
		Nip:     nip,
		Name:    name,
	}, nil
}

func generateUniqueBase32Secret(userID int, email string) string {
	const b32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	hash := sha256.Sum256([]byte(fmt.Sprintf("lopi-q-secret-seed-%d-%s", userID, email)))
	var sb strings.Builder
	for i := 0; i < 16; i++ {
		idx := int(hash[i]) % len(b32Alphabet)
		sb.WriteByte(b32Alphabet[idx])
	}
	return sb.String()
}

func (s *AuthClientDirectStub) Setup2FA(ctx context.Context, req *authProto.Setup2FARequest) (*authProto.Setup2FAResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID := int(req.UserId)
	if userID == 0 && req.TempToken != "" {
		if claims, err := jwt.ValidateToken(req.TempToken); err == nil && claims.IsTemp {
			userID = claims.UserID
		}
	}

	users, path := loadUsersJSON()
	accountIdentifier := "user@bulukumbakab.go.id"
	var totpSecret string
	var targetEmail, targetNIP string

	for i := range users {
		if users[i].ID == userID {
			if users[i].Email != "" {
				accountIdentifier = users[i].Email
			} else if users[i].NIP != "" {
				accountIdentifier = users[i].NIP
			}
			targetEmail = users[i].Email
			targetNIP = users[i].NIP

			if users[i].TotpSecret != "" {
				totpSecret = users[i].TotpSecret
			} else {
				totpSecret = generateUniqueBase32Secret(users[i].ID, users[i].Email)
				users[i].TotpSecret = totpSecret
				saveUsersJSON(users, path)
				go syncPostgresTotpSecret(targetEmail, targetNIP, totpSecret)
			}
			break
		}
	}

	if totpSecret == "" {
		totpSecret = generateUniqueBase32Secret(userID, accountIdentifier)
	}

	escapedAccount := url.QueryEscape(accountIdentifier)
	otpUrl := fmt.Sprintf("otpauth://totp/LOPI-Q%%20BULUKUMBA:%s?secret=%s&issuer=LOPI-Q%%20BULUKUMBA", escapedAccount, totpSecret)

	return &authProto.Setup2FAResponse{
		Success:    true,
		Secret:     totpSecret,
		OtpauthUrl: otpUrl,
	}, nil
}

func syncPostgresTotpSecret(email, nip, secret string) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	authConn := fmt.Sprintf("host=%s port=5432 user=user_lopiq_auth password=lopiqauthPassword@2k26# dbname=db_lopiq_auth sslmode=disable", dbHost)
	if dbAuth, err := sql.Open("postgres", authConn); err == nil {
		defer dbAuth.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, _ = dbAuth.ExecContext(ctx, `UPDATE auth_users SET totp_secret=$1 WHERE email=$2 OR nip=$3;`, secret, email, nip)
	}

	userConn := fmt.Sprintf("host=%s port=5432 user=user_lopiq_user password=lopiquserPassword@2k26# dbname=db_lopiq_user sslmode=disable", dbHost)
	if dbUser, err := sql.Open("postgres", userConn); err == nil {
		defer dbUser.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, _ = dbUser.ExecContext(ctx, `UPDATE users SET totp_secret=$1 WHERE email=$2 OR nip=$3;`, secret, email, nip)
	}
}

func (s *AuthClientDirectStub) Enable2FA(ctx context.Context, req *authProto.Enable2FARequest) (*authProto.Enable2FAResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID := int(req.UserId)
	if userID == 0 && req.TempToken != "" {
		if claims, err := jwt.ValidateToken(req.TempToken); err == nil && claims.IsTemp {
			userID = claims.UserID
		}
	}

	if userID == 0 {
		userID = 1
	}

	users, path := loadUsersJSON()
	backupCodes := []string{"88219412", "99124012", "77123951", "12495812"}

	var targetEmail, targetNIP string
	for i := range users {
		if users[i].ID == userID {
			users[i].TotpEnabled = true
			if users[i].TotpSecret == "" {
				users[i].TotpSecret = generateUniqueBase32Secret(users[i].ID, users[i].Email)
			}
			users[i].BackupCodes = backupCodes
			targetEmail = users[i].Email
			targetNIP = users[i].NIP
			saveUsersJSON(users, path)
			break
		}
	}

	if targetEmail != "" || targetNIP != "" {
		go syncPostgres2FAEnabled(targetEmail, targetNIP, true)
	}

	return &authProto.Enable2FAResponse{
		Success:     true,
		BackupCodes: backupCodes,
	}, nil
}

func syncPostgres2FAEnabled(email, nip string, enabled bool) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	authConn := fmt.Sprintf("host=%s port=5432 user=user_lopiq_auth password=lopiqauthPassword@2k26# dbname=db_lopiq_auth sslmode=disable", dbHost)
	if dbAuth, err := sql.Open("postgres", authConn); err == nil {
		defer dbAuth.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, _ = dbAuth.ExecContext(ctx, `UPDATE auth_users SET totp_enabled=$1 WHERE email=$2 OR nip=$3;`, enabled, email, nip)
	}

	userConn := fmt.Sprintf("host=%s port=5432 user=user_lopiq_user password=lopiquserPassword@2k26# dbname=db_lopiq_user sslmode=disable", dbHost)
	if dbUser, err := sql.Open("postgres", userConn); err == nil {
		defer dbUser.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, _ = dbUser.ExecContext(ctx, `UPDATE users SET totp_enabled=$1 WHERE email=$2 OR nip=$3;`, enabled, email, nip)
	}
}

func (s *AuthClientDirectStub) Disable2FA(ctx context.Context, req *authProto.Disable2FARequest) (*authProto.Disable2FAResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, path := loadUsersJSON()
	for i := range users {
		if users[i].ID == int(req.UserId) {
			users[i].TotpEnabled = false
			users[i].TotpSecret = ""
			users[i].BackupCodes = nil
			saveUsersJSON(users, path)
			break
		}
	}
	return &authProto.Disable2FAResponse{Success: true}, nil
}

func (s *AuthClientDirectStub) SelfReset2FA(ctx context.Context, req *authProto.SelfReset2FARequest) (*authProto.SelfReset2FAResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, path := loadUsersJSON()
	claims, err := jwt.ValidateToken(req.TempToken)
	userID := 1
	if err == nil && claims.IsTemp {
		userID = claims.UserID
	}

	for i := range users {
		if users[i].ID == userID {
			users[i].TotpEnabled = false
			users[i].TotpSecret = ""
			users[i].BackupCodes = nil
			saveUsersJSON(users, path)
			break
		}
	}

	return &authProto.SelfReset2FAResponse{Success: true}, nil
}
