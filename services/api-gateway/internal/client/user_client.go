package client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	userProto "proto/user"
)

type UserClient struct {
	server userProto.UserServiceServer
}

func NewUserClient(server userProto.UserServiceServer) *UserClient {
	return &UserClient{server: server}
}

type UserClientDirectStub struct {
	mu sync.Mutex
}

func (s *UserClientDirectStub) toProtoUser(u *UserDataJSON) *userProto.User {
	perms := u.Permissions
	if perms == nil {
		perms = []string{}
	}
	return &userProto.User{
		Id:             int32(u.ID),
		Nip:            u.NIP,
		Email:          u.Email,
		Name:           u.Name,
		Role:           u.Role,
		Jabatan:        u.Jabatan,
		UnitKerja:      u.UnitKerja,
		Permissions:    perms,
		TotpEnabled:    u.TotpEnabled,
		IsActive:       u.IsActive,
		PeriodeMulai:   u.PeriodeMulai,
		PeriodeSelesai: u.PeriodeSelesai,
	}
}

func syncPostgresCreateUser(req *userProto.CreateUserRequest, passwordHash string) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))
	cleanNip := strings.TrimSpace(req.Nip)

	// 1. Insert/Update into auth_users in db_lopiq_auth
	for _, conn := range getAuthConnStrings(dbHost) {
		if dbAuth, err := sql.Open("postgres", conn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, _ = dbAuth.ExecContext(ctx, `ALTER TABLE auth_users ALTER COLUMN nip DROP NOT NULL;`)

			var exists bool
			_ = dbAuth.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth_users WHERE LOWER(email)=LOWER($1) OR ($2 <> '' AND REPLACE(nip,' ','')=REPLACE($2,' ','')));`, cleanEmail, cleanNip).Scan(&exists)

			if exists {
				_, errUpd := dbAuth.ExecContext(ctx,
					`UPDATE auth_users SET nip=NULLIF($1, ''), name=$3, role=$4, jabatan=$5, unit_kerja=$6, password=$7, is_active=true
					 WHERE LOWER(email)=$2 OR ($1 <> '' AND REPLACE(nip,' ','')=REPLACE($1,' ',''));`,
					cleanNip, cleanEmail, req.Name, req.Role, req.Jabatan, req.UnitKerja, passwordHash,
				)
				if errUpd != nil {
					log.Printf("[syncPostgresCreateUser UPDATE Auth Error] %v", errUpd)
				} else {
					log.Printf("[syncPostgresCreateUser Auth Success] Updated existing user %s", cleanEmail)
				}
			} else {
				_, errIns := dbAuth.ExecContext(ctx,
					`INSERT INTO auth_users (nip, email, name, role, jabatan, unit_kerja, password, is_active)
					 VALUES (NULLIF($1, ''), $2, $3, $4, $5, $6, $7, true);`,
					cleanNip, cleanEmail, req.Name, req.Role, req.Jabatan, req.UnitKerja, passwordHash,
				)
				if errIns != nil {
					log.Printf("[syncPostgresCreateUser INSERT Auth Error] %v", errIns)
				} else {
					log.Printf("[syncPostgresCreateUser Auth Success] Inserted new user %s", cleanEmail)
				}
			}
			cancel()
			dbAuth.Close()
		}
	}

	// 2. Insert/Update into users in db_lopiq_user
	for _, conn := range getUserConnStrings(dbHost) {
		if dbUser, err := sql.Open("postgres", conn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, _ = dbUser.ExecContext(ctx, `ALTER TABLE users ALTER COLUMN nip DROP NOT NULL;`)

			var exists bool
			_ = dbUser.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email)=LOWER($1) OR ($2 <> '' AND REPLACE(nip,' ','')=REPLACE($2,' ','')));`, cleanEmail, cleanNip).Scan(&exists)

			if exists {
				_, errUpd := dbUser.ExecContext(ctx,
					`UPDATE users SET nip=NULLIF($1, ''), name=$3, role=$4, jabatan=$5, unit_kerja=$6, password_hash=$7, is_active=true
					 WHERE LOWER(email)=$2 OR ($1 <> '' AND REPLACE(nip,' ','')=REPLACE($1,' ',''));`,
					cleanNip, cleanEmail, req.Name, req.Role, req.Jabatan, req.UnitKerja, passwordHash,
				)
				if errUpd != nil {
					log.Printf("[syncPostgresCreateUser UPDATE User Error] %v", errUpd)
				} else {
					log.Printf("[syncPostgresCreateUser User Success] Updated existing user %s", cleanEmail)
				}
			} else {
				_, errIns := dbUser.ExecContext(ctx,
					`INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
					 VALUES (NULLIF($1, ''), $2, $3, $4, $5, $6, $7, true);`,
					cleanNip, cleanEmail, req.Name, req.Role, req.Jabatan, req.UnitKerja, passwordHash,
				)
				if errIns != nil {
					log.Printf("[syncPostgresCreateUser INSERT User Error] %v", errIns)
				} else {
					log.Printf("[syncPostgresCreateUser User Success] Inserted new user %s", cleanEmail)
				}
			}
			cancel()
			dbUser.Close()
		}
	}
}

func syncPostgresUpdateUser(u *UserDataJSON) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	for _, authConn := range getAuthConnStrings(dbHost) {
		if dbAuth, err := sql.Open("postgres", authConn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			if u.PasswordHash != "" {
				_, _ = dbAuth.ExecContext(ctx,
					`UPDATE auth_users SET nip=$1, email=$2, name=$3, role=$4, jabatan=$5, unit_kerja=$6, password=$7 WHERE id=$8 OR (nip <> '' AND REPLACE(nip,' ','')=REPLACE($1,' ','')) OR email=$2;`,
					u.NIP, u.Email, u.Name, u.Role, u.Jabatan, u.UnitKerja, u.PasswordHash, u.ID,
				)
			} else {
				_, _ = dbAuth.ExecContext(ctx,
					`UPDATE auth_users SET nip=$1, email=$2, name=$3, role=$4, jabatan=$5, unit_kerja=$6 WHERE id=$7 OR (nip <> '' AND REPLACE(nip,' ','')=REPLACE($1,' ','')) OR email=$2;`,
					u.NIP, u.Email, u.Name, u.Role, u.Jabatan, u.UnitKerja, u.ID,
				)
			}
			cancel()
			dbAuth.Close()
		}
	}

	for _, userConn := range getUserConnStrings(dbHost) {
		if dbUser, err := sql.Open("postgres", userConn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			if u.PasswordHash != "" {
				_, _ = dbUser.ExecContext(ctx,
					`UPDATE users SET nip=$1, email=$2, name=$3, role=$4, jabatan=$5, unit_kerja=$6, password_hash=$7 WHERE id=$8 OR (nip <> '' AND REPLACE(nip,' ','')=REPLACE($1,' ','')) OR email=$2;`,
					u.NIP, u.Email, u.Name, u.Role, u.Jabatan, u.UnitKerja, u.PasswordHash, u.ID,
				)
			} else {
				_, _ = dbUser.ExecContext(ctx,
					`UPDATE users SET nip=$1, email=$2, name=$3, role=$4, jabatan=$5, unit_kerja=$6 WHERE id=$7 OR (nip <> '' AND REPLACE(nip,' ','')=REPLACE($1,' ','')) OR email=$2;`,
					u.NIP, u.Email, u.Name, u.Role, u.Jabatan, u.UnitKerja, u.ID,
				)
			}

			// Also update user_name in presensi_records so past/existing presensi records reflect the updated name!
			_, _ = dbUser.ExecContext(ctx,
				`UPDATE presensi_records SET user_name=$1 WHERE user_id=$2 OR (user_nip <> '' AND REPLACE(user_nip,' ','')=REPLACE($3,' ',''));`,
				u.Name, u.ID, u.NIP,
			)
			cancel()
			dbUser.Close()
		}
	}
}

func syncPostgresToggleActive(email, nip string, isActive bool) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	for _, authConn := range getAuthConnStrings(dbHost) {
		if dbAuth, err := sql.Open("postgres", authConn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, _ = dbAuth.ExecContext(ctx,
				`UPDATE auth_users SET is_active=$1 WHERE email=$2 OR (nip <> '' AND REPLACE(nip,' ','')=REPLACE($3,' ',''));`,
				isActive, email, nip,
			)
			cancel()
			dbAuth.Close()
		}
	}

	for _, userConn := range getUserConnStrings(dbHost) {
		if dbUser, err := sql.Open("postgres", userConn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, _ = dbUser.ExecContext(ctx,
				`UPDATE users SET is_active=$1 WHERE email=$2 OR (nip <> '' AND REPLACE(nip,' ','')=REPLACE($3,' ',''));`,
				isActive, email, nip,
			)
			cancel()
			dbUser.Close()
		}
	}
}

func syncPostgresReset2FA(email, nip string) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	cleanNIP := strings.ReplaceAll(nip, " ", "")

	for _, conn := range getAuthConnStrings(dbHost) {
		if dbAuth, err := sql.Open("postgres", conn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, _ = dbAuth.ExecContext(ctx, `UPDATE auth_users SET totp_enabled=false, totp_secret='' WHERE LOWER(email)=LOWER($1) OR (REPLACE(nip, ' ', '')=$2 AND $2 != '');`, cleanEmail, cleanNIP)
			cancel()
			dbAuth.Close()
		}
	}

	for _, conn := range getUserConnStrings(dbHost) {
		if dbUser, err := sql.Open("postgres", conn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, _ = dbUser.ExecContext(ctx, `UPDATE users SET totp_enabled=false, totp_secret='' WHERE LOWER(email)=LOWER($1) OR (REPLACE(nip, ' ', '')=$2 AND $2 != '');`, cleanEmail, cleanNIP)
			cancel()
			dbUser.Close()
		}
	}
}

func (s *UserClientDirectStub) GetProfile(ctx context.Context, req *userProto.GetProfileRequest) (*userProto.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, _ := loadUsersJSON()
	for i := range users {
		if users[i].ID == int(req.UserId) {
			return &userProto.UserResponse{
				Success: true,
				User:    s.toProtoUser(&users[i]),
			}, nil
		}
	}

	// Never fall back to another account's profile: the caller would inherit that user's role.
	return &userProto.UserResponse{Success: false, Error: "Pengguna tidak ditemukan."}, nil
}

func (s *UserClientDirectStub) ListUsers(ctx context.Context, req *userProto.ListUsersRequest) (*userProto.ListUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, _ := loadUsersJSON()

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	existingEmails := make(map[string]bool)
	for _, u := range users {
		existingEmails[strings.ToLower(u.Email)] = true
	}

	for _, conn := range getUserConnStrings(dbHost) {
		if dbUser, err := sql.Open("postgres", conn); err == nil {
			ctxQ, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			rows, errQ := dbUser.QueryContext(ctxQ, `SELECT id, nip, email, name, role, jabatan, unit_kerja, is_active FROM users;`)
			if errQ == nil {
				for rows.Next() {
					var id int32
					var nip, email, name, role, jabatan, unitKerja string
					var isActive bool
					if errS := rows.Scan(&id, &nip, &email, &name, &role, &jabatan, &unitKerja, &isActive); errS == nil {
						cleanE := strings.ToLower(email)
						if !existingEmails[cleanE] {
							existingEmails[cleanE] = true
							users = append(users, UserDataJSON{
								ID:          int(id),
								NIP:         nip,
								Email:       email,
								Name:        name,
								Role:        role,
								Jabatan:     jabatan,
								UnitKerja:   unitKerja,
								Permissions: []string{"submit_attendance"},
								IsActive:    isActive,
							})
						}
					}
				}
				rows.Close()
			}
			cancel()
			dbUser.Close()
		}
	}

	result := make([]*userProto.User, 0, len(users))
	for i := range users {
		result = append(result, s.toProtoUser(&users[i]))
	}

	return &userProto.ListUsersResponse{
		Success: true,
		Users:   result,
	}, nil
}

func (s *UserClientDirectStub) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, path := loadUsersJSON()

	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))
	cleanNIP := strings.ReplaceAll(req.Nip, " ", "")

	var existingIdx = -1
	for i := range users {
		if strings.ToLower(users[i].Email) == cleanEmail || (cleanNIP != "" && strings.ReplaceAll(users[i].NIP, " ", "") == cleanNIP) {
			existingIdx = i
			break
		}
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	perms := req.Permissions
	if perms == nil {
		perms = []string{"submit_attendance"}
	}

	if existingIdx >= 0 {
		users[existingIdx].Name = req.Name
		users[existingIdx].NIP = req.Nip
		users[existingIdx].Email = req.Email
		users[existingIdx].Role = req.Role
		users[existingIdx].Jabatan = req.Jabatan
		users[existingIdx].UnitKerja = req.UnitKerja
		users[existingIdx].PasswordHash = string(hash)
		users[existingIdx].PeriodeMulai = req.PeriodeMulai
		users[existingIdx].PeriodeSelesai = req.PeriodeSelesai
		users[existingIdx].IsActive = true
		saveUsersJSON(users, path)

		syncPostgresCreateUser(req, string(hash))
		_ = syncPostgresInternPeriod(req.Email, req.Nip, req.PeriodeMulai, req.PeriodeSelesai)

		return &userProto.UserResponse{
			Success: true,
			User:    s.toProtoUser(&users[existingIdx]),
			Message: "Data akun peserta magang berhasil diperbarui.",
		}, nil
	}

	nextID := 1
	for _, u := range users {
		if u.ID >= nextID {
			nextID = u.ID + 1
		}
	}

	newUser := UserDataJSON{
		ID:             nextID,
		NIP:            req.Nip,
		Email:          req.Email,
		Name:           req.Name,
		Role:           req.Role,
		Jabatan:        req.Jabatan,
		UnitKerja:      req.UnitKerja,
		Permissions:    perms,
		PasswordHash:   string(hash),
		IsActive:       true,
		PeriodeMulai:   req.PeriodeMulai,
		PeriodeSelesai: req.PeriodeSelesai,
	}

	users = append(users, newUser)
	saveUsersJSON(users, path)

	syncPostgresCreateUser(req, string(hash))
	_ = syncPostgresInternPeriod(req.Email, req.Nip, req.PeriodeMulai, req.PeriodeSelesai)

	return &userProto.UserResponse{
		Success: true,
		User:    s.toProtoUser(&newUser),
		Message: "Peserta magang berhasil ditambahkan.",
	}, nil
}

func (s *UserClientDirectStub) UpdateUser(ctx context.Context, req *userProto.UpdateUserRequest) (*userProto.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, path := loadUsersJSON()
	for i := range users {
		if users[i].ID == int(req.Id) {
			if req.Name != "" {
				users[i].Name = req.Name
			}
			if req.Nip != "" {
				users[i].NIP = req.Nip
			}
			if req.Email != "" {
				users[i].Email = req.Email
			}
			if req.Role != "" {
				users[i].Role = req.Role
			}
			if req.Jabatan != "" {
				users[i].Jabatan = req.Jabatan
			}
			if req.UnitKerja != "" {
				users[i].UnitKerja = req.UnitKerja
			}
			if req.Permissions != nil {
				users[i].Permissions = req.Permissions
			}
			if req.Password != "" {
				hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
				users[i].PasswordHash = string(hash)
			}

			// The edit form always sends both dates, so empty means the period was cleared
			users[i].PeriodeMulai = req.PeriodeMulai
			users[i].PeriodeSelesai = req.PeriodeSelesai

			saveUsersJSON(users, path)
			go syncPostgresUpdateUser(&users[i])
			go syncPostgresInternPeriod(users[i].Email, users[i].NIP, users[i].PeriodeMulai, users[i].PeriodeSelesai)

			return &userProto.UserResponse{
				Success: true,
				User:    s.toProtoUser(&users[i]),
				Message: "User berhasil diperbarui.",
			}, nil
		}
	}

	return &userProto.UserResponse{Success: false, Error: "User tidak ditemukan."}, nil
}

func (s *UserClientDirectStub) ToggleActive(ctx context.Context, req *userProto.ToggleActiveRequest) (*userProto.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, path := loadUsersJSON()
	for i := range users {
		if users[i].ID == int(req.Id) {
			users[i].IsActive = req.IsActive
			saveUsersJSON(users, path)
			go syncPostgresToggleActive(users[i].Email, users[i].NIP, req.IsActive)

			return &userProto.UserResponse{
				Success: true,
				User:    s.toProtoUser(&users[i]),
			}, nil
		}
	}

	return &userProto.UserResponse{Success: false, Error: "User tidak ditemukan."}, nil
}

func (s *UserClientDirectStub) Reset2FA(ctx context.Context, req *userProto.Reset2FARequest) (*userProto.Reset2FAResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, path := loadUsersJSON()
	var targetEmail, targetNIP string
	for i := range users {
		if users[i].ID == int(req.Id) {
			users[i].TotpEnabled = false
			users[i].TotpSecret = ""
			users[i].BackupCodes = nil
			targetEmail = users[i].Email
			targetNIP = users[i].NIP
			break
		}
	}
	saveUsersJSON(users, path)

	// Always sync reset to PostgreSQL database
	syncPostgresReset2FA(targetEmail, targetNIP)

	return &userProto.Reset2FAResponse{
		Success: true,
		Message: "Keamanan 2FA pengguna berhasil direset.",
	}, nil
}

type ActivityLogItem struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	UserName  string `json:"user_name"`
	UserNIP   string `json:"user_nip"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

func GetClientIP(r *http.Request) string {
	if r == nil {
		return "-"
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// The first entry is the original client; later ones are proxies
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}
	return "-"
}

const createActivityLogsTableSQL = `
	CREATE TABLE IF NOT EXISTS activity_logs (
		id SERIAL PRIMARY KEY,
		user_id INT,
		user_nip VARCHAR(50),
		user_name VARCHAR(255),
		action VARCHAR(100) NOT NULL,
		details TEXT NOT NULL,
		ip_address VARCHAR(100) NOT NULL,
		user_agent TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
`

// RecordActivityLog writes the audit entry asynchronously so a slow database never delays the request.
func RecordActivityLog(userID int, userNIP, userName, action, details, ip, userAgent string) {
	if ip == "" {
		ip = "-"
	}
	if userAgent == "" {
		userAgent = "-"
	}

	go func() {
		db, err := openAuthDB()
		if err != nil {
			log.Printf("[ActivityLog Error] %s tidak tercatat: %v", action, err)
			return
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := db.ExecContext(ctx, createActivityLogsTableSQL); err != nil {
			log.Printf("[ActivityLog Warning] Table creation check: %v", err)
		}
		if _, err := db.ExecContext(ctx, `
			INSERT INTO activity_logs (user_id, user_nip, user_name, action, details, ip_address, user_agent, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW());
		`, userID, userNIP, userName, action, details, ip, userAgent); err != nil {
			log.Printf("[ActivityLog Error] Gagal menyimpan %s: %v", action, err)
		}
	}()
}

func (s *UserClientDirectStub) GetActivityLogs(ctx context.Context) ([]ActivityLogItem, error) {
	db, err := openAuthDB()
	if err != nil {
		return nil, fmt.Errorf("koneksi database log audit gagal: %v", err)
	}
	defer db.Close()

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.ExecContext(queryCtx, createActivityLogsTableSQL); err != nil {
		log.Printf("[ActivityLog Warning] Table creation check: %v", err)
	}

	rows, err := db.QueryContext(queryCtx, `
		SELECT id, COALESCE(NULLIF(user_name, ''), '-'), COALESCE(NULLIF(user_nip, ''), '-'), action, details, ip_address, user_agent,
		       to_char(created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD HH24:MI:SS')
		FROM activity_logs
		ORDER BY created_at DESC
		LIMIT 1000;
	`)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca log audit: %v", err)
	}
	defer rows.Close()

	list := make([]ActivityLogItem, 0)
	for rows.Next() {
		var item ActivityLogItem
		if err := rows.Scan(&item.ID, &item.UserName, &item.UserNIP, &item.Action, &item.Details, &item.IPAddress, &item.UserAgent, &item.Timestamp); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
