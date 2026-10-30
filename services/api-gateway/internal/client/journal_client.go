package client

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Daily journal entries of interns: several activities per day, each with an evidence file and its own
// approval decision. Stored in the attendance database; evidence bytes live in the row so no shared
// volume is needed between containers.

const MaxJournalEvidenceBytes = 3 << 20

var allowedEvidenceTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

type JournalEntry struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	UserNIP       string `json:"user_nip"`
	UserName      string `json:"user_name"`
	ActivityDate  string `json:"activity_date"`
	ActivityTime  string `json:"activity_time"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	EvidenceName  string `json:"evidence_name"`
	EvidenceMime  string `json:"evidence_mime"`
	HasEvidence   bool   `json:"has_evidence"`
	Status        string `json:"status"`
	ReviewerNote  string `json:"reviewer_note"`
	ReviewedBy    string `json:"reviewed_by"`
	ReviewedAt    string `json:"reviewed_at"`
	CreatedAt     string `json:"created_at"`
	EvidenceBytes []byte `json:"-"`
}

type JournalFilter struct {
	UserID int // 0 = all users
	Date   string
	Status string
}

func ensureJournalTable(ctx context.Context, db *sql.DB) {
	_, _ = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS intern_journal_entries (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			user_nip VARCHAR(50) NOT NULL DEFAULT '',
			user_name VARCHAR(255) NOT NULL DEFAULT '',
			activity_date DATE NOT NULL,
			activity_time TIME NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			evidence_name VARCHAR(255) NOT NULL DEFAULT '',
			evidence_mime VARCHAR(100) NOT NULL DEFAULT '',
			evidence_data BYTEA,
			status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
			reviewer_note TEXT NOT NULL DEFAULT '',
			reviewed_by VARCHAR(255) NOT NULL DEFAULT '',
			reviewed_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_journal_user_date ON intern_journal_entries (user_id, activity_date);
	`)
}

func openJournalDB(ctx context.Context) (*sql.DB, context.Context, context.CancelFunc, error) {
	db, err := openInternshipDB()
	if err != nil {
		return nil, nil, nil, err
	}
	qctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	ensureJournalTable(qctx, db)
	return db, qctx, cancel, nil
}

// ValidateJournalEntry normalises the submitted fields and checks the evidence file.
func ValidateJournalEntry(e *JournalEntry) error {
	e.Title = strings.TrimSpace(e.Title)
	e.Description = strings.TrimSpace(e.Description)
	if e.Title == "" {
		return fmt.Errorf("nama kegiatan wajib diisi")
	}
	if len(e.Title) > 255 {
		return fmt.Errorf("nama kegiatan maksimal 255 karakter")
	}
	d, err := time.Parse("2006-01-02", strings.TrimSpace(e.ActivityDate))
	if err != nil {
		return fmt.Errorf("tanggal kegiatan tidak valid")
	}
	if d.After(time.Now().In(witaLoc()).AddDate(0, 0, 1)) {
		return fmt.Errorf("tanggal kegiatan tidak boleh di masa depan")
	}
	if _, err := time.Parse("15:04", strings.TrimSpace(e.ActivityTime)); err != nil {
		return fmt.Errorf("jam kegiatan tidak valid (format JJ:MM)")
	}
	if len(e.EvidenceBytes) == 0 {
		return fmt.Errorf("bukti dukung wajib diunggah")
	}
	if len(e.EvidenceBytes) > MaxJournalEvidenceBytes {
		return fmt.Errorf("ukuran bukti dukung maksimal 3 MB")
	}
	if !allowedEvidenceTypes[e.EvidenceMime] {
		return fmt.Errorf("bukti dukung harus berupa gambar (JPG/PNG/WEBP) atau PDF")
	}
	return nil
}

func witaLoc() *time.Location {
	if loc, err := time.LoadLocation("Asia/Makassar"); err == nil {
		return loc
	}
	return time.FixedZone("WITA", 8*3600)
}

const journalSelectColumns = `
	id, user_id, user_nip, user_name, to_char(activity_date, 'YYYY-MM-DD'), to_char(activity_time, 'HH24:MI'),
	title, description, evidence_name, evidence_mime, evidence_data IS NOT NULL, status, reviewer_note, reviewed_by,
	COALESCE(to_char(reviewed_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD HH24:MI'), ''),
	to_char(created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD HH24:MI')`

func scanJournalEntry(scan func(dest ...interface{}) error) (JournalEntry, error) {
	var e JournalEntry
	err := scan(&e.ID, &e.UserID, &e.UserNIP, &e.UserName, &e.ActivityDate, &e.ActivityTime,
		&e.Title, &e.Description, &e.EvidenceName, &e.EvidenceMime, &e.HasEvidence, &e.Status, &e.ReviewerNote, &e.ReviewedBy,
		&e.ReviewedAt, &e.CreatedAt)
	return e, err
}

func ListJournalEntries(ctx context.Context, f JournalFilter) ([]JournalEntry, error) {
	db, qctx, cancel, err := openJournalDB(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer cancel()

	rows, err := db.QueryContext(qctx, `
		SELECT `+journalSelectColumns+`
		FROM intern_journal_entries
		WHERE ($1 = 0 OR user_id = $1)
		  AND ($2 = '' OR activity_date = NULLIF($2, '')::date)
		  AND ($3 = '' OR status = $3)
		ORDER BY activity_date DESC, user_name ASC, activity_time ASC, id ASC
		LIMIT 2000;
	`, f.UserID, f.Date, f.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []JournalEntry{}
	for rows.Next() {
		e, err := scanJournalEntry(rows.Scan)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func CreateJournalEntry(ctx context.Context, e JournalEntry) (JournalEntry, error) {
	db, qctx, cancel, err := openJournalDB(ctx)
	if err != nil {
		return e, err
	}
	defer db.Close()
	defer cancel()

	row := db.QueryRowContext(qctx, `
		INSERT INTO intern_journal_entries (user_id, user_nip, user_name, activity_date, activity_time, title, description, evidence_name, evidence_mime, evidence_data)
		VALUES ($1, $2, $3, $4::date, $5::time, $6, $7, $8, $9, $10)
		RETURNING `+journalSelectColumns+`;
	`, e.UserID, e.UserNIP, e.UserName, e.ActivityDate, e.ActivityTime, e.Title, e.Description, e.EvidenceName, e.EvidenceMime, e.EvidenceBytes)
	saved, err := scanJournalEntry(row.Scan)
	if err != nil {
		return e, fmt.Errorf("gagal menyimpan kegiatan: %v", err)
	}
	return saved, nil
}

// GetJournalEvidence returns the entry (with evidence bytes) so the caller can check ownership.
func GetJournalEvidence(ctx context.Context, id int) (JournalEntry, error) {
	db, qctx, cancel, err := openJournalDB(ctx)
	if err != nil {
		return JournalEntry{}, err
	}
	defer db.Close()
	defer cancel()

	var e JournalEntry
	err = db.QueryRowContext(qctx, `SELECT id, user_id, evidence_name, evidence_mime, evidence_data FROM intern_journal_entries WHERE id = $1;`, id).
		Scan(&e.ID, &e.UserID, &e.EvidenceName, &e.EvidenceMime, &e.EvidenceBytes)
	if err == sql.ErrNoRows {
		return e, fmt.Errorf("kegiatan tidak ditemukan")
	}
	return e, err
}

// DeleteOwnPendingJournalEntry lets an intern withdraw an entry that has not been reviewed yet.
func DeleteOwnPendingJournalEntry(ctx context.Context, id, userID int) error {
	db, qctx, cancel, err := openJournalDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()

	res, err := db.ExecContext(qctx, `DELETE FROM intern_journal_entries WHERE id = $1 AND user_id = $2 AND status = 'PENDING';`, id, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("kegiatan tidak ditemukan atau sudah direview sehingga tidak dapat dihapus")
	}
	return nil
}

func ReviewJournalEntry(ctx context.Context, id int, status, note, reviewer string) (JournalEntry, error) {
	if status != "APPROVED" && status != "REJECTED" && status != "PENDING" {
		return JournalEntry{}, fmt.Errorf("status review tidak valid")
	}
	note = strings.TrimSpace(note)
	if status == "REJECTED" && note == "" {
		return JournalEntry{}, fmt.Errorf("catatan wajib diisi saat kegiatan tidak disetujui")
	}

	db, qctx, cancel, err := openJournalDB(ctx)
	if err != nil {
		return JournalEntry{}, err
	}
	defer db.Close()
	defer cancel()

	row := db.QueryRowContext(qctx, `
		UPDATE intern_journal_entries
		SET status = $2, reviewer_note = $3, reviewed_by = $4,
		    reviewed_at = CASE WHEN $2 = 'PENDING' THEN NULL ELSE NOW() END
		WHERE id = $1
		RETURNING `+journalSelectColumns+`;
	`, id, status, note, reviewer)
	saved, err := scanJournalEntry(row.Scan)
	if err == sql.ErrNoRows {
		return saved, fmt.Errorf("kegiatan tidak ditemukan")
	}
	return saved, err
}
