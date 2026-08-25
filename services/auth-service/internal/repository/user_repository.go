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

	// 1. Try loading from PostgreSQL if database connection is available
	if r.sqlDB != nil {
		rows, err := r.sqlDB.Query("SELECT id, COALESCE(nip, ''), email, name, role, COALESCE(jabatan, ''), COALESCE(unit_kerja, ''), password, COALESCE(totp_secret, ''), COALESCE(totp_enabled, false), COALESCE(is_active, true), created_at FROM auth_users;")
		if err == nil {
			defer rows.Close()
			var loaded []model.User
			for rows.Next() {
				var u model.User
				if err := rows.Scan(&u.ID, &u.NIP, &u.Email, &u.Name, &u.Role, &u.Jabatan, &u.UnitKerja, &u.PasswordHash, &u.TotpSecret, &u.TotpEnabled, &u.IsActive, &u.CreatedAt); err == nil {
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
	aswanHash, _ := bcrypt.GenerateFromPassword([]byte("Asw&a198"), bcrypt.DefaultCost)
	callTakerHash, _ := bcrypt.GenerateFromPassword([]byte("Asw&a198"), bcrypt.DefaultCost)

	// Super Admin Aswan
	r.users = append(r.users, model.User{
		ID:           r.nextID,
		NIP:          "199501012020011000",
		Email:        "aswan@bulukumbakab.go.id",
		Name:         "Muhammad Aswan",
		Role:         "superadmin",
		Jabatan:      "HEAD OF DISKOMINFO",
		UnitKerja:    "Diskominfo Kab. Bulukumba",
		PasswordHash: string(aswanHash),
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
		_, _ = r.sqlDB.Exec("UPDATE auth_users SET totp_enabled = $1, totp_secret = $2 WHERE id = $3 OR email = $4 OR nip = $5;", u.TotpEnabled, u.TotpSecret, u.ID, u.Email, u.NIP)
	}
	if r.sqlUserDB != nil {
		_, _ = r.sqlUserDB.Exec("UPDATE users SET totp_enabled = $1, totp_secret = $2 WHERE id = $3 OR email = $4 OR nip = $5;", u.TotpEnabled, u.TotpSecret, u.ID, u.Email, u.NIP)
	}

	for i := range r.users {
		if r.users[i].ID == u.ID {
			r.users[i] = *u
			r.saveDBLocked()
			return
		}
	}
}
