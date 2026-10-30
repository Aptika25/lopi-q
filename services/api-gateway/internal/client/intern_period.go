package client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// The period columns come from migration 004_add_internship_fields. They are read in a separate query so
// that a database without them never breaks the main user load (and therefore login).

func mergeInternPeriods(ctx context.Context, db *sql.DB, users []UserDataJSON) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, COALESCE(to_char(periode_mulai, 'YYYY-MM-DD'), ''), COALESCE(to_char(periode_selesai, 'YYYY-MM-DD'), '')
		FROM auth_users;
	`)
	if err != nil {
		log.Printf("[InternPeriod Warning] Kolom periode magang tidak terbaca: %v", err)
		return
	}
	defer rows.Close()

	periods := map[int][2]string{}
	for rows.Next() {
		var id int
		var start, end string
		if rows.Scan(&id, &start, &end) == nil {
			periods[id] = [2]string{start, end}
		}
	}
	for i := range users {
		if p, ok := periods[users[i].ID]; ok {
			users[i].PeriodeMulai, users[i].PeriodeSelesai = p[0], p[1]
		}
	}
}

// ValidateInternPeriod accepts empty values (period not set) or two dates where the end is not before the start.
func ValidateInternPeriod(start, end string) error {
	start, end = strings.TrimSpace(start), strings.TrimSpace(end)
	var startT, endT time.Time
	var err error
	if start != "" {
		if startT, err = time.Parse("2006-01-02", start); err != nil {
			return fmt.Errorf("format tanggal mulai magang tidak valid")
		}
	}
	if end != "" {
		if endT, err = time.Parse("2006-01-02", end); err != nil {
			return fmt.Errorf("format tanggal berakhir magang tidak valid")
		}
	}
	if start != "" && end != "" && endT.Before(startT) {
		return fmt.Errorf("tanggal berakhir magang tidak boleh sebelum tanggal mulai")
	}
	return nil
}

// syncPostgresInternPeriod stores the period on both the auth and the profile table, matched by email / NIP.
func syncPostgresInternPeriod(email, nip, start, end string) error {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	cleanNIP := strings.ReplaceAll(nip, " ", "")
	targets := []struct {
		creds dbCreds
		table string
	}{
		{authDBCreds(), "auth_users"},
		{userDBCreds(), "users"},
	}

	var firstErr error
	for _, t := range targets {
		db, err := openFirstReachable(t.creds.name, buildConnStrings(t.creds))
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, _ = db.ExecContext(ctx, `ALTER TABLE `+t.table+` ADD COLUMN IF NOT EXISTS periode_mulai DATE, ADD COLUMN IF NOT EXISTS periode_selesai DATE;`)
		_, err = db.ExecContext(ctx, `
			UPDATE `+t.table+` SET periode_mulai = NULLIF($1, '')::date, periode_selesai = NULLIF($2, '')::date
			WHERE ($3 <> '' AND LOWER(email) = $3) OR ($4 <> '' AND REPLACE(nip, ' ', '') = $4);
		`, strings.TrimSpace(start), strings.TrimSpace(end), cleanEmail, cleanNIP)
		cancel()
		db.Close()
		if err != nil {
			log.Printf("[InternPeriod Error] Gagal menyimpan periode magang di %s: %v", t.table, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}
