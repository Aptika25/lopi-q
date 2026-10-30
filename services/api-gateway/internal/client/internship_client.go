package client

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Projects and final reports live in the attendance database next to presensi_records.

type ProjectMember struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type InternProject struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Mentor      string          `json:"mentor"`
	Members     []ProjectMember `json:"members"`
	Status      string          `json:"status"`
	Progress    int             `json:"progress"`
	StartDate   string          `json:"start_date"`
	DueDate     string          `json:"due_date"`
	UpdatedAt   string          `json:"updated_at"`
}

type FinalReport struct {
	ID            int      `json:"id"`
	InternID      int      `json:"intern_id"`
	InternName    string   `json:"intern_name"`
	Institution   string   `json:"institution"`
	Title         string   `json:"title"`
	DocumentURL   string   `json:"document_url"`
	SubmittedAt   string   `json:"submitted_at"`
	Status        string   `json:"status"`
	Score         *float64 `json:"score"`
	ReviewerNotes string   `json:"reviewer_notes"`
	UpdatedAt     string   `json:"updated_at"`
}

var validProjectStatuses = map[string]bool{"PERENCANAAN": true, "BERJALAN": true, "SELESAI": true, "DITUNDA": true}
var validReportStatuses = map[string]bool{"BELUM_DIKUMPULKAN": true, "MENUNGGU_REVIEW": true, "REVISI": true, "DISETUJUI": true}

func openInternshipDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}
	db, err := getPostgresDB(dbHost)
	if err != nil {
		return nil, fmt.Errorf("koneksi database gagal: %v", err)
	}
	return db, nil
}

func ensureInternshipTables(ctx context.Context, db *sql.DB) {
	_, _ = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS intern_projects (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			mentor VARCHAR(255) NOT NULL DEFAULT '',
			members JSONB NOT NULL DEFAULT '[]',
			status VARCHAR(30) NOT NULL DEFAULT 'PERENCANAAN',
			progress INT NOT NULL DEFAULT 0,
			start_date DATE,
			due_date DATE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS intern_final_reports (
			id SERIAL PRIMARY KEY,
			intern_id INT NOT NULL DEFAULT 0,
			intern_name VARCHAR(255) NOT NULL,
			institution VARCHAR(255) NOT NULL DEFAULT '',
			title VARCHAR(255) NOT NULL,
			document_url TEXT NOT NULL DEFAULT '',
			submitted_at DATE,
			status VARCHAR(30) NOT NULL DEFAULT 'BELUM_DIKUMPULKAN',
			score DOUBLE PRECISION,
			reviewer_notes TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
}

// nullableDate turns "" into NULL so optional DATE columns accept empty form fields.
func nullableDate(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return s
}

func NormalizeProject(p *InternProject) error {
	p.Title = strings.TrimSpace(p.Title)
	if p.Title == "" {
		return fmt.Errorf("judul proyek wajib diisi")
	}
	if p.Status == "" {
		p.Status = "PERENCANAAN"
	}
	if !validProjectStatuses[p.Status] {
		return fmt.Errorf("status proyek tidak valid")
	}
	if p.Progress < 0 {
		p.Progress = 0
	}
	if p.Progress > 100 {
		p.Progress = 100
	}
	if p.Members == nil {
		p.Members = []ProjectMember{}
	}
	return nil
}

func NormalizeFinalReport(r *FinalReport) error {
	r.InternName = strings.TrimSpace(r.InternName)
	r.Title = strings.TrimSpace(r.Title)
	if r.InternName == "" || r.Title == "" {
		return fmt.Errorf("nama peserta dan judul laporan wajib diisi")
	}
	if r.Status == "" {
		r.Status = "BELUM_DIKUMPULKAN"
	}
	if !validReportStatuses[r.Status] {
		return fmt.Errorf("status laporan tidak valid")
	}
	if r.Score != nil && (*r.Score < 0 || *r.Score > 100) {
		return fmt.Errorf("nilai harus di antara 0 - 100")
	}
	r.DocumentURL = strings.TrimSpace(r.DocumentURL)
	if r.DocumentURL != "" && !strings.HasPrefix(r.DocumentURL, "https://") && !strings.HasPrefix(r.DocumentURL, "http://") {
		return fmt.Errorf("link dokumen harus diawali http:// atau https://")
	}
	return nil
}

func ListProjects(ctx context.Context) ([]InternProject, error) {
	db, err := openInternshipDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ensureInternshipTables(ctx, db)

	rows, err := db.QueryContext(ctx, `
		SELECT id, title, description, mentor, members, status, progress,
		       COALESCE(to_char(start_date, 'YYYY-MM-DD'), ''), COALESCE(to_char(due_date, 'YYYY-MM-DD'), ''),
		       to_char(updated_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD HH24:MI')
		FROM intern_projects ORDER BY id DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []InternProject{}
	for rows.Next() {
		var p InternProject
		var membersRaw []byte
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Mentor, &membersRaw, &p.Status, &p.Progress, &p.StartDate, &p.DueDate, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(membersRaw, &p.Members); err != nil || p.Members == nil {
			p.Members = []ProjectMember{}
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// SaveProject inserts when p.ID is 0, otherwise updates; it returns the stored ID.
func SaveProject(ctx context.Context, p InternProject) (int, error) {
	db, err := openInternshipDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ensureInternshipTables(ctx, db)

	members, _ := json.Marshal(p.Members)
	if p.ID == 0 {
		err = db.QueryRowContext(ctx, `
			INSERT INTO intern_projects (title, description, mentor, members, status, progress, start_date, due_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;
		`, p.Title, p.Description, p.Mentor, string(members), p.Status, p.Progress, nullableDate(p.StartDate), nullableDate(p.DueDate)).Scan(&p.ID)
		return p.ID, err
	}

	res, err := db.ExecContext(ctx, `
		UPDATE intern_projects SET title=$1, description=$2, mentor=$3, members=$4, status=$5, progress=$6,
		       start_date=$7, due_date=$8, updated_at=NOW()
		WHERE id=$9;
	`, p.Title, p.Description, p.Mentor, string(members), p.Status, p.Progress, nullableDate(p.StartDate), nullableDate(p.DueDate), p.ID)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, fmt.Errorf("proyek #%d tidak ditemukan", p.ID)
	}
	return p.ID, nil
}

func DeleteProject(ctx context.Context, id int) error {
	return deleteInternshipRow(ctx, "intern_projects", id)
}

func ListFinalReports(ctx context.Context) ([]FinalReport, error) {
	db, err := openInternshipDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ensureInternshipTables(ctx, db)

	rows, err := db.QueryContext(ctx, `
		SELECT id, intern_id, intern_name, institution, title, document_url,
		       COALESCE(to_char(submitted_at, 'YYYY-MM-DD'), ''), status, score, reviewer_notes,
		       to_char(updated_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD HH24:MI')
		FROM intern_final_reports ORDER BY id DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []FinalReport{}
	for rows.Next() {
		var r FinalReport
		var score sql.NullFloat64
		if err := rows.Scan(&r.ID, &r.InternID, &r.InternName, &r.Institution, &r.Title, &r.DocumentURL, &r.SubmittedAt, &r.Status, &score, &r.ReviewerNotes, &r.UpdatedAt); err != nil {
			return nil, err
		}
		if score.Valid {
			r.Score = &score.Float64
		}
		reports = append(reports, r)
	}
	return reports, rows.Err()
}

// SaveFinalReport inserts when r.ID is 0, otherwise updates; it returns the stored ID.
func SaveFinalReport(ctx context.Context, r FinalReport) (int, error) {
	db, err := openInternshipDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ensureInternshipTables(ctx, db)

	if r.ID == 0 {
		err = db.QueryRowContext(ctx, `
			INSERT INTO intern_final_reports (intern_id, intern_name, institution, title, document_url, submitted_at, status, score, reviewer_notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;
		`, r.InternID, r.InternName, r.Institution, r.Title, r.DocumentURL, nullableDate(r.SubmittedAt), r.Status, r.Score, r.ReviewerNotes).Scan(&r.ID)
		return r.ID, err
	}

	res, err := db.ExecContext(ctx, `
		UPDATE intern_final_reports SET intern_id=$1, intern_name=$2, institution=$3, title=$4, document_url=$5,
		       submitted_at=$6, status=$7, score=$8, reviewer_notes=$9, updated_at=NOW()
		WHERE id=$10;
	`, r.InternID, r.InternName, r.Institution, r.Title, r.DocumentURL, nullableDate(r.SubmittedAt), r.Status, r.Score, r.ReviewerNotes, r.ID)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, fmt.Errorf("laporan akhir #%d tidak ditemukan", r.ID)
	}
	return r.ID, nil
}

func DeleteFinalReport(ctx context.Context, id int) error {
	return deleteInternshipRow(ctx, "intern_final_reports", id)
}

// deleteInternshipRow only ever receives one of the two constant table names above.
func deleteInternshipRow(ctx context.Context, table string, id int) error {
	db, err := openInternshipDB()
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ensureInternshipTables(ctx, db)

	res, err := db.ExecContext(ctx, "DELETE FROM "+table+" WHERE id = $1;", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("data #%d tidak ditemukan", id)
	}
	return nil
}
