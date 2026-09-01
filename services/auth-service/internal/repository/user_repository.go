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

	"auth-service/internal/model"
)

type UserRepository struct {
	mu        sync.Mutex
	dbPath    string
	sqlDB     *sql.DB
	sqlUserDB *sql.DB
	users     []model.User
	nextID    int
}

func NewUserRepository(dbPath string, sqlDB *sql.DB, sqlUserDB *sql.DB) *UserRepository {
	dbDir := filepath.Dir(dbPath)
	_ = os.MkdirAll(dbDir, 0755)

	repo := &UserRepository{
		dbPath:    dbPath,
		sqlDB:     sqlDB,
		sqlUserDB: sqlUserDB,
		nextID:    1,
	}
	repo.LoadDB()
	return repo
}

func (r *UserRepository) ResetAll2FA() {
	if r.sqlDB != nil {
		_, err := r.sqlDB.Exec("UPDATE auth_users SET totp_enabled = false, totp_secret = NULL;")
		if err == nil {
			log.Println("[Auth-Service] Successfully reset 2FA for all users in PostgreSQL database!")
		}
	}
	if r.sqlUserDB != nil {
		_, _ = r.sqlUserDB.Exec("UPDATE users SET totp_enabled = false, totp_secret = NULL;")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.users {
		r.users[i].TotpEnabled = false
		r.users[i].TotpSecret = ""
		r.users[i].BackupCodes = nil
	}
	r.saveDBLocked()
}

func (r *UserRepository) LoadDB() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure 2FA columns exist in PostgreSQL if database connection is available
	if r.sqlDB != nil {
		_, _ = r.sqlDB.Exec("ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS totp_enabled BOOLEAN DEFAULT false;")
		_, _ = r.sqlDB.Exec("ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS totp_secret TEXT;")
	}
	if r.sqlUserDB != nil {
		_, _ = r.sqlUserDB.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS totp_enabled BOOLEAN DEFAULT false;")
		_, _ = r.sqlUserDB.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS totp_secret TEXT;")
	}

	// 1. Try loading from PostgreSQL if database connection is available
	if r.sqlDB != nil {
		rows, err := r.sqlDB.Query("SELECT id, COALESCE(nip, ''), email, name, role, COALESCE(jabatan, ''), COALESCE(unit_kerja, ''), password, COALESCE(totp_secret, ''), COALESCE(totp_enabled, false), COALESCE(is_active, true), created_at FROM auth_users;")
		if err == nil {
			defer rows.Close()
			var loaded []model.User
			for rows.Next() {
				var u model.User
				if err := rows.Scan(&u.ID, &u.NIP, &u.Email, &u.Name, &u.Role, &u.Jabatan, &u.UnitKerja, &u.PasswordHash, &u.TotpSecret, &u.TotpEnabled, &u.IsActive, &u.CreatedAt); err == nil {
					// Preserve active 2FA state if in-memory has it enabled
					for _, existing := range r.users {
						if (existing.ID == u.ID || (existing.Email != "" && strings.EqualFold(existing.Email, u.Email))) && existing.TotpEnabled {
							u.TotpEnabled = true
							if u.TotpSecret == "" {
								u.TotpSecret = existing.TotpSecret
							}
							u.BackupCodes = existing.BackupCodes
							break
						}
					}
					loaded = append(loaded, u)
					if u.ID >= r.nextID {
						r.nextID = u.ID + 1
					}
				}
			}
			if len(loaded) > 0 {
				r.users = loaded
				log.Printf("[Auth-Service] Loaded %d users directly from PostgreSQL database (auth_users)", len(r.users))
				return
			}
		}
	}

	// 2. Fallback to users.json file
	if _, err := os.Stat(r.dbPath); err == nil {
		raw, err := os.ReadFile(r.dbPath)
		if err == nil {
			_ = json.Unmarshal(raw, &r.users)
			for _, u := range r.users {
				if u.ID >= r.nextID {
					r.nextID = u.ID + 1
				}
			}
			log.Printf("[Auth-Service] Loaded %d users from JSON file repository", len(r.users))
			return
		}
	}

	log.Println("[Auth-Service] Seeding default auth users...")
	r.seedUsers()
	r.saveDBLocked()
}

func (r *UserRepository) saveDBLocked() {
	raw, err := json.MarshalIndent(r.users, "", "  ")
	if err == nil {
		_ = os.WriteFile(r.dbPath, raw, 0644)
	}
}

func (r *UserRepository) seedUsers() {
	superAdminHash := "$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO"
	callTakerHash := superAdminHash

	// Super Admin Aswan (Strictly from 002_seed_super_admin.up.sql)
	r.users = append(r.users, model.User{
		ID:           r.nextID,
		NIP:          "199708192025061003",
		Email:        "aswan@bulukumbakab.go.id",
		Name:         "Muhammad Aswan, S.T.",
		Role:         "superadmin",
		Jabatan:      "JF Pranata Komputer Ahli Pertama",
		UnitKerja:    "Diskominfo Kab. Bulukumba",
		PasswordHash: superAdminHash,
		IsActive:     true,
		CreatedAt:    time.Now(),
	})
	r.nextID++

	seeds := []struct {
		nip, email, name, jabatan, unit string
	}{
		{"19940503202521 1 138", "amappalua@bulukumbakab.go.id", "A.Mappalua, S.Pd", "PENATA LAYANAN OPERASIONAL", "Dinas Sosial"},
		{"19870304202521 1 061", "suherman@bulukumbakab.go.id", "Suherman, S.Pd", "PENATA LAYANAN OPERASIONAL", "Badan Penanggulangan Bencana Daerah"},
		{"20000206202521 1 166", "riswandirisman@bulukumbakab.go.id", "Riswandi Risman", "OPERATOR LAYANAN OPERASIONAL", "Dinas Kesehatan"},
		{"19900215202521 1 114", "abilkizri@bulukumbakab.go.id", "Abil Kizri", "OPERATOR LAYANAN OPERASIONAL", "Dinas Perhubungan"},
		{"19911005202521 1 087", "imamardiyansah@bulukumbakab.go.id", "Imam Ardiyansah", "OPERATOR LAYANAN OPERASIONAL", "Satpol, Pemadam Kebakaran dan Penyelamatan"},
		{"19861130202521 1 101", "abdrahim@bulukumbakab.go.id", "Abd.Rahim", "OPERATOR LAYANAN OPERASIONAL", "Dinas Sosial"},
		{"19860304202521 1 147", "munawir@bulukumbakab.go.id", "Munawir Syadzali", "PENATA LAYANAN OPERASIONAL", "Badan Penanggulangan Bencana Daerah"},
		{"19760802200604 1 017", "abdullah@bulukumbakab.go.id", "Abdullah, S.Kep., Ns", "PERENCANA", "Dinas Kesehatan"},
		{"19860712202521 1 089", "ismail@bulukumbakab.go.id", "Ismail, S.Sos", "PENATA LAYANAN OPERASIONAL", "Dinas Perhubungan"},
		{"19960328202521 1 050", "aldiafdal@bulukumbakab.go.id", "Aldi Afdali Saputra", "OPERATOR LAYANAN OPERASIONAL", "Satpol, Pemadam Kebakaran dan Penyelamatan"},
	}

	for _, ct := range seeds {
		r.users = append(r.users, model.User{
			ID:           r.nextID,
			NIP:          ct.nip,
			Email:        ct.email,
			Name:         ct.name,
			Role:         "call_taker",
			Jabatan:      ct.jabatan,
			UnitKerja:    ct.unit,
			PasswordHash: string(callTakerHash),
			IsActive:     true,
			CreatedAt:    time.Now(),
		})
		r.nextID++
	}
}

func (r *UserRepository) FindByIdentifier(identifier string) *model.User {
	r.LoadDB()
	r.mu.Lock()
	defer r.mu.Unlock()

	clean := strings.TrimSpace(identifier)
	cleanNIP := strings.ReplaceAll(clean, " ", "")
	cleanEmail := strings.ToLower(clean)

	for i := range r.users {
		u := &r.users[i]
		if strings.ReplaceAll(u.NIP, " ", "") == cleanNIP || strings.ToLower(u.Email) == cleanEmail {
			return u
		}
	}
	return nil
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

func (r *UserRepository) Save(u *model.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.sqlDB != nil {
		res, err := r.sqlDB.Exec("UPDATE auth_users SET totp_enabled = $1, totp_secret = $2 WHERE id = $3 OR LOWER(email) = LOWER($4) OR REPLACE(nip, ' ', '') = REPLACE($5, ' ', '');", u.TotpEnabled, u.TotpSecret, u.ID, u.Email, u.NIP)
		if err != nil {
			log.Printf("[Auth-Service Error] Failed to update auth_users DB: %v", err)
		} else if n, _ := res.RowsAffected(); n > 0 {
			log.Printf("[Auth-Service] Successfully persisted 2FA state (enabled=%v) in auth_users DB for user ID %d (%s)", u.TotpEnabled, u.ID, u.Email)
		}
	}
	if r.sqlUserDB != nil {
		_, err := r.sqlUserDB.Exec("UPDATE users SET totp_enabled = $1, totp_secret = $2 WHERE id = $3 OR LOWER(email) = LOWER($4) OR REPLACE(nip, ' ', '') = REPLACE($5, ' ', '');", u.TotpEnabled, u.TotpSecret, u.ID, u.Email, u.NIP)
		if err != nil {
			log.Printf("[Auth-Service Error] Failed to update users DB: %v", err)
		}
	}

	for i := range r.users {
		if r.users[i].ID == u.ID || (u.Email != "" && strings.EqualFold(r.users[i].Email, u.Email)) {
			r.users[i] = *u
			r.saveDBLocked()
			log.Printf("[Auth-Service] Saved user 2FA state to JSON file repository (enabled=%v) for %s", u.TotpEnabled, u.Email)
			return
		}
	}
}
