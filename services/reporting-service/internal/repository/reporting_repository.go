package repository

import (
	"database/sql"
	"fmt"
	"log"

	repProto "proto/reporting"
	"reporting-service/internal/model"
)

type ReportingRepository interface {
	GetAttendanceRecap(month, department, search string, limit int32) (*repProto.ReportingStats, []*model.AttendanceReport, error)
}

type reportingRepository struct {
	db       *sql.DB
	memStore []*model.AttendanceReport
}

func NewReportingRepository(db *sql.DB) ReportingRepository {
	repo := &reportingRepository{
		db: db,
	}

	if db == nil {
		repo.initSeedData()
	}

	return repo
}

func (r *reportingRepository) initSeedData() {
	// No dummy attendance reports
	r.memStore = []*model.AttendanceReport{}
}

func (r *reportingRepository) GetAttendanceRecap(month, department, search string, limit int32) (*repProto.ReportingStats, []*model.AttendanceReport, error) {
	stats := &repProto.ReportingStats{
		TotalHadir:     0,
		PctTepatWaktu:  0,
		TotalTerlambat: 0,
		TotalAbsen:     0,
	}

	if r.db != nil {
		query := `SELECT id, user_id, user_name, user_nip, department, date, clock_in, clock_out, total_hours, status FROM tb_lopiq_attendance_reports WHERE 1=1`
		args := []interface{}{}
		argIdx := 1

		if department != "" && department != "Semua Departemen" {
			query += fmt.Sprintf(" AND department = $%d", argIdx)
			args = append(args, department)
			argIdx++
		}
		if search != "" {
			query += fmt.Sprintf(" AND (user_name ILIKE $%d OR user_nip ILIKE $%d)", argIdx, argIdx)
			args = append(args, "%"+search+"%")
			argIdx++
		}

		query += " ORDER BY id DESC LIMIT 50"

		rows, err := r.db.Query(query, args...)
		if err != nil {
			log.Printf("Error querying reporting db: %v", err)
			return stats, r.memStore, nil
		}
		defer rows.Close()

		var results []*model.AttendanceReport
		for rows.Next() {
			item := &model.AttendanceReport{}
			if err := rows.Scan(
				&item.ID, &item.UserID, &item.UserName, &item.UserNIP,
				&item.Department, &item.Date, &item.ClockIn, &item.ClockOut,
				&item.TotalHours, &item.Status,
			); err == nil {
				results = append(results, item)
			}
		}
		if len(results) > 0 {
			return stats, results, nil
		}
	}

	// Memory fallback
	var filtered []*model.AttendanceReport
	for _, item := range r.memStore {
		if department != "" && department != "Semua Departemen" && item.Department != department {
			continue
		}
		filtered = append(filtered, item)
	}

	return stats, filtered, nil
}
