package repository

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"user-service/internal/model"
)

type UserRepository struct {
	mu        sync.Mutex
	dbPath    string
	sqlDB     *sql.DB
	sqlAuthDB *sql.DB
	users     []model.User
	nextID    int
}

func NewUserRepository(dbPath string, sqlDB *sql.DB, sqlAuthDB *sql.DB) *UserRepository {
	dbDir := filepath.Dir(dbPath)
	_ = os.MkdirAll(dbDir, 0755)

	r := &UserRepository{
		dbPath:    dbPath,
		sqlDB:     sqlDB,
		sqlAuthDB: sqlAuthDB,
		nextID:    1,
	}
	r.LoadDB()
	return r
}

func (r *UserRepository) LoadDB() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Build real-time 2FA status map from auth_users in db_lopiq_auth
	totpMap := make(map[string]bool)
	if r.sqlAuthDB != nil {
		authRows, err := r.sqlAuthDB.Query("SELECT COALESCE(email, ''), COALESCE(nip, ''), COALESCE(totp_enabled, false) FROM auth_users;")
		if err == nil {
			defer authRows.Close()
			for authRows.Next() {
				var email, nip string
				var enabled bool
				if err := authRows.Scan(&email, &nip, &enabled); err == nil {
					if email != "" {
						totpMap[strings.ToLower(email)] = enabled
					}
					if nip != "" {
						totpMap[strings.ReplaceAll(nip, " ", "")] = enabled
					}
				}
			}
		}
	}

	// 1. Try loading from PostgreSQL if database connection is available
	if r.sqlDB != nil {
		rows, err := r.sqlDB.Query("SELECT id, COALESCE(nip, ''), email, name, role, COALESCE(jabatan, ''), COALESCE(unit_kerja, ''), password_hash, COALESCE(totp_secret, ''), COALESCE(totp_enabled, false), COALESCE(is_active, true), created_at FROM users ORDER BY id;")
		if err == nil {
			defer rows.Close()
			var loaded []model.User
			for rows.Next() {
				var u model.User
				if err := rows.Scan(&u.ID, &u.NIP, &u.Email, &u.Name, &u.Role, &u.Jabatan, &u.UnitKerja, &u.PasswordHash, &u.TotpSecret, &u.TotpEnabled, &u.IsActive, &u.CreatedAt); err == nil {
					// Check if auth_users override exists for 2FA status
					if enabled, ok := totpMap[strings.ToLower(u.Email)]; ok {
						u.TotpEnabled = enabled
					} else if enabled, ok := totpMap[strings.ReplaceAll(u.NIP, " ", "")]; ok {
						u.TotpEnabled = enabled
					}

					// Set default permissions according to role
					if u.Role == "superadmin" {
						u.Permissions = []string{"manage_users", "manage_attendance", "manage_locations", "view_reports"}
					} else if u.Role == "admin" {
						u.Permissions = []string{"manage_attendance", "manage_locations", "view_reports"}
					} else {
						u.Permissions = []string{"submit_attendance"}
					}
					loaded = append(loaded, u)
					if u.ID >= r.nextID {
						r.nextID = u.ID + 1
					}
				} else {
					log.Printf("[User-Service Warning] Scan row failed: %v", err)
				}
			}
			if len(loaded) > 0 {
				// Ensure superadmin is present in the list
				hasSuperAdmin := false
				for _, u := range loaded {
					if u.Role == "superadmin" || u.Email == "aswan@bulukumbakab.go.id" {
						hasSuperAdmin = true
						break
					}
				}
				if !hasSuperAdmin {
					aswanHash, _ := bcrypt.GenerateFromPassword([]byte("Asw&a198"), bcrypt.DefaultCost)
					totpVal := false
					if enabled, ok := totpMap["aswan@bulukumbakab.go.id"]; ok {
						totpVal = enabled
					}
					superAdmin := model.User{
						NIP:          "199708192025061003",
						Email:        "aswan@bulukumbakab.go.id",
						Name:         "Muhammad Aswan, S.T.",
						Role:         "superadmin",
						Jabatan:      "JF Pranata Komputer Ahli Pertama",
						UnitKerja:    "Diskominfo Kab. Bulukumba",
						PasswordHash: string(aswanHash),
						Permissions:  []string{"manage_users", "manage_attendance", "manage_locations", "view_reports"},
						TotpEnabled:  totpVal,
						IsActive:     true,
						CreatedAt:    time.Now(),
					}
					var id int
					err := r.sqlDB.QueryRow(
						`INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, totp_enabled, is_active)
						 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (email) DO UPDATE SET nip='199708192025061003', name='Muhammad Aswan, S.T.', role='superadmin' RETURNING id`,
						superAdmin.NIP, superAdmin.Email, superAdmin.Name, superAdmin.Role, superAdmin.Jabatan, superAdmin.UnitKerja, superAdmin.PasswordHash, superAdmin.TotpEnabled, superAdmin.IsActive,
					).Scan(&id)
					if err == nil {
						superAdmin.ID = id
					} else {
						superAdmin.ID = r.nextID
						r.nextID++
					}
					loaded = append([]model.User{superAdmin}, loaded...)
				}
				r.users = loaded
				log.Printf("[User-Service] Loaded %d users directly from PostgreSQL database (users)", len(r.users))
				return
			}
		} else {
			log.Printf("[User-Service Warning] Error querying PostgreSQL users table: %v", err)
		}
	}

	// 2. Fallback to users.json file
	if _, err := os.Stat(r.dbPath); err == nil {
		raw, err := os.ReadFile(r.dbPath)
		if err == nil {
			_ = json.Unmarshal(raw, &r.users)
			for i := range r.users {
				u := &r.users[i]
				if enabled, ok := totpMap[strings.ToLower(u.Email)]; ok {
					u.TotpEnabled = enabled
				} else if enabled, ok := totpMap[strings.ReplaceAll(u.NIP, " ", "")]; ok {
					u.TotpEnabled = enabled
				}
				if u.ID >= r.nextID {
					r.nextID = u.ID + 1
				}
			}
			if len(r.users) > 0 {
				log.Printf("[User-Service] Loaded %d users from JSON file repository", len(r.users))
				return
			}
		}
	}

	// 3. Fallback: Seed users if database and JSON file are empty
	log.Println("[User-Service] Seeding default single superadmin user...")
	r.seedUsers()
	r.SaveDBLocked()
}

func (r *UserRepository) seedUsers() {
	aswanHash, _ := bcrypt.GenerateFromPassword([]byte("Asw&a198"), bcrypt.DefaultCost)

	// Super Admin Aswan
	superAdmin := model.User{
		ID:           r.nextID,
		NIP:          "199708192025061003",
		Email:        "aswan@bulukumbakab.go.id",
		Name:         "Muhammad Aswan, S.T.",
		Role:         "superadmin",
		Jabatan:      "JF Pranata Komputer Ahli Pertama",
		UnitKerja:    "Diskominfo Kab. Bulukumba",
		PasswordHash: string(aswanHash),
		Permissions:  []string{"manage_users", "manage_attendance", "manage_locations", "view_reports"},
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	r.users = append(r.users, superAdmin)
	r.nextID++

	// Also insert into PostgreSQL if DB connection exists
	if r.sqlDB != nil {
		for _, u := range r.users {
			_, _ = r.sqlDB.Exec(
				`INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (email) DO UPDATE SET nip=$1, name=$3;`,
				u.NIP, u.Email, u.Name, u.Role, u.Jabatan, u.UnitKerja, u.PasswordHash, u.IsActive,
			)
		}
	}
}

func (r *UserRepository) SaveDBLocked() {
	raw, err := json.MarshalIndent(r.users, "", "  ")
	if err == nil {
		_ = os.WriteFile(r.dbPath, raw, 0644)
	}
}

func (r *UserRepository) FindByID(id int) *model.User {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == id {
			return &r.users[i]
		}
	}
	return nil
}

func (r *UserRepository) ListAll() []model.User {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]model.User{}, r.users...)
}

func (r *UserRepository) Create(u model.User) model.User {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	u.CreatedAt = time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	u.PasswordHash = string(hash)
	if u.Permissions == nil {
		u.Permissions = []string{}
	}

	if r.sqlDB != nil {
		var id int
		err := r.sqlDB.QueryRow(
			`INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			u.NIP, u.Email, u.Name, u.Role, u.Jabatan, u.UnitKerja, u.PasswordHash, u.IsActive,
		).Scan(&id)
		if err == nil {
			u.ID = id
		} else {
			log.Printf("[User-Service Warning] PostgreSQL Create failed: %v", err)
			u.ID = r.nextID
			r.nextID++
		}
	} else {
		u.ID = r.nextID
		r.nextID++
	}

	r.users = append(r.users, u)
	r.SaveDBLocked()
	return u
}

func (r *UserRepository) Update(u *model.User) bool {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == u.ID {
			if u.Name != "" {
				r.users[i].Name = u.Name
			}
			if u.Email != "" {
				r.users[i].Email = u.Email
			}
			if u.NIP != "" {
				r.users[i].NIP = u.NIP
			}
			if u.Jabatan != "" {
				r.users[i].Jabatan = u.Jabatan
			}
			if u.UnitKerja != "" {
				r.users[i].UnitKerja = u.UnitKerja
			}
			if u.Role != "" {
				r.users[i].Role = u.Role
			}
			if u.Permissions != nil {
				r.users[i].Permissions = u.Permissions
			}
			if u.PasswordHash != "" {
				hash, _ := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
				r.users[i].PasswordHash = string(hash)
			}

			if r.sqlDB != nil {
				if u.PasswordHash != "" {
					_, err := r.sqlDB.Exec(
						`UPDATE users SET nip=$1, email=$2, name=$3, role=$4, jabatan=$5, unit_kerja=$6, password_hash=$7 WHERE id=$8`,
						r.users[i].NIP, r.users[i].Email, r.users[i].Name, r.users[i].Role, r.users[i].Jabatan, r.users[i].UnitKerja, r.users[i].PasswordHash, u.ID,
					)
					if err != nil {
						log.Printf("[User-Service Warning] PostgreSQL Update failed: %v", err)
					}
				} else {
					_, err := r.sqlDB.Exec(
						`UPDATE users SET nip=$1, email=$2, name=$3, role=$4, jabatan=$5, unit_kerja=$6 WHERE id=$7`,
						r.users[i].NIP, r.users[i].Email, r.users[i].Name, r.users[i].Role, r.users[i].Jabatan, r.users[i].UnitKerja, u.ID,
					)
					if err != nil {
						log.Printf("[User-Service Warning] PostgreSQL Update failed: %v", err)
					}
				}
			}

			r.SaveDBLocked()
			*u = r.users[i]
			return true
		}
	}
	return false
}

func (r *UserRepository) ToggleActive(id int, isActive bool) *model.User {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == id {
			r.users[i].IsActive = isActive
			if r.sqlDB != nil {
				_, err := r.sqlDB.Exec("UPDATE users SET is_active = $1 WHERE id = $2", isActive, id)
				if err != nil {
					log.Printf("[User-Service Warning] PostgreSQL ToggleActive failed: %v", err)
				}
			}
			r.SaveDBLocked()
			return &r.users[i]
		}
	}
	return nil
}

func (r *UserRepository) Reset2FA(id int) bool {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == id {
			r.users[i].TotpEnabled = false
			r.users[i].TotpSecret = ""
			if r.sqlDB != nil {
				_, err := r.sqlDB.Exec("UPDATE users SET totp_enabled = false, totp_secret = NULL WHERE id = $1", id)
				if err != nil {
					log.Printf("[User-Service Warning] PostgreSQL Reset2FA failed: %v", err)
				}
			}
			if r.sqlAuthDB != nil {
				_, err := r.sqlAuthDB.Exec("UPDATE auth_users SET totp_enabled = false, totp_secret = NULL WHERE email = $1 OR nip = $2", r.users[i].Email, r.users[i].NIP)
				if err != nil {
					log.Printf("[User-Service Warning] PostgreSQL auth_users Reset2FA failed: %v", err)
				}
			}
			r.SaveDBLocked()
			return true
		}
	}
	return false
}
