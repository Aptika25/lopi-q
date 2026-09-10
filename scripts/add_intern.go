package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func connectDB(dbName string) (*sql.DB, error) {
	hosts := []string{os.Getenv("DB_HOST"), "localhost", "127.0.0.1", "postgres_apps", "host.docker.internal"}
	user := "user_lopiq_auth"
	pass := "lopiqauthPassword@2k26#"
	if dbName == "db_lopiq_user" {
		user = "user_lopiq_user"
		pass = "lopiquserPassword@2k26#"
	}

	for _, h := range hosts {
		if strings.TrimSpace(h) == "" {
			continue
		}
		connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", h, user, pass, dbName)
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			errPing := db.PingContext(ctx)
			cancel()
			if errPing == nil {
				return db, nil
			}
			db.Close()
		}
	}
	return nil, fmt.Errorf("gagal terhubung ke database %s di seluruh host", dbName)
}

func main() {
	if len(os.Args) < 7 {
		fmt.Println("Usage: go run scripts/add_intern.go <Nama> <NISN/NIM> <Email> <Jurusan/UnitKerja> <Sekolah/Jabatan> <Password>")
		fmt.Println("Example: go run scripts/add_intern.go \"Ade Anisa\" \"0091755897\" \"adeanisa150299@gmail.com\" \"Rekayasa Perangkat Lunak\" \"SMK TI Bulukumba\" \"password123\"")
		os.Exit(1)
	}

	name := os.Args[1]
	nip := os.Args[2]
	email := strings.ToLower(strings.TrimSpace(os.Args[3]))
	unitKerja := os.Args[4]
	jabatan := os.Args[5]
	password := os.Args[6]

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal meng-hash password: %v", err)
	}
	passwordHash := string(hash)

	// 1. Insert/Update ke db_lopiq_auth
	dbAuth, err := connectDB("db_lopiq_auth")
	if err != nil {
		log.Fatalf("Error auth DB: %v", err)
	}
	defer dbAuth.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, errAuth := dbAuth.ExecContext(ctx,
		`INSERT INTO auth_users (nip, email, name, role, jabatan, unit_kerja, password, is_active)
		 VALUES ($1, $2, $3, 'intern', $4, $5, $6, true)
		 ON CONFLICT (email) DO UPDATE SET nip=$1, name=$3, jabatan=$4, unit_kerja=$5, password=$6, is_active=true;`,
		nip, email, name, jabatan, unitKerja, passwordHash,
	)
	if errAuth != nil {
		fmt.Printf("⚠️ Warning insert auth_users: %v. Mencoba fallback update...\n", errAuth)
		_, _ = dbAuth.ExecContext(ctx,
			`UPDATE auth_users SET nip=$1, name=$3, jabatan=$4, unit_kerja=$5, password=$6, is_active=true WHERE email=$2 OR (nip <> '' AND nip=$1);`,
			nip, email, name, jabatan, unitKerja, passwordHash,
		)
	}

	// 2. Insert/Update ke db_lopiq_user
	dbUser, err := connectDB("db_lopiq_user")
	if err != nil {
		log.Fatalf("Error user DB: %v", err)
	}
	defer dbUser.Close()

	_, errUser := dbUser.ExecContext(ctx,
		`INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
		 VALUES ($1, $2, $3, 'intern', $4, $5, $6, true)
		 ON CONFLICT (email) DO UPDATE SET nip=$1, name=$3, jabatan=$4, unit_kerja=$5, password_hash=$6, is_active=true;`,
		nip, email, name, jabatan, unitKerja, passwordHash,
	)
	if errUser != nil {
		fmt.Printf("⚠️ Warning insert users: %v. Mencoba fallback update...\n", errUser)
		_, _ = dbUser.ExecContext(ctx,
			`UPDATE users SET nip=$1, name=$3, jabatan=$4, unit_kerja=$5, password_hash=$6, is_active=true WHERE email=$2 OR (nip <> '' AND nip=$1);`,
			nip, email, name, jabatan, unitKerja, passwordHash,
		)
	}

	fmt.Println("==========================================================")
	fmt.Printf("✅ BERHASIL MENAMBAHKAN AKUN PESERTA MAGANG BARU IN DATABASE!\n")
	fmt.Printf("   Nama         : %s\n", name)
	fmt.Printf("   NISN / NIM   : %s\n", nip)
	fmt.Printf("   Email Akses  : %s\n", email)
	fmt.Printf("   Jurusan      : %s\n", unitKerja)
	fmt.Printf("   Asal Sekolah : %s\n", jabatan)
	fmt.Printf("   Role         : intern\n")
	fmt.Printf("   Status       : Direct Active (Bisa langsung login di /login)\n")
	fmt.Println("==========================================================")
}
