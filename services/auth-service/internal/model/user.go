package model

import "time"

type User struct {
	ID           int       `json:"id"`
	NIP          string    `json:"nip"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Jabatan      string    `json:"jabatan"`
	UnitKerja    string    `json:"unit_kerja"`
	PasswordHash string    `json:"password_hash"`
	TotpSecret   string    `json:"totp_secret,omitempty"`
	TotpEnabled  bool      `json:"totp_enabled"`
	BackupCodes  []string  `json:"backup_codes,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
