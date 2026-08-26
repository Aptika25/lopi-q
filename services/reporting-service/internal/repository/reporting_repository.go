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
	r.memStore = []*model.AttendanceReport{
		{ID: 1, UserID: 101, UserName: "Budi Santoso", UserNIP: "IN-301", Department: "Engineering", Date: "24 Okt 2023", ClockIn: "08:50 AM", ClockOut: "05:15 PM", TotalHours: "8h 25m", Status: "HADIR"},
		{ID: 2, UserID: 102, UserName: "Siti Aminah", UserNIP: "IN-302", Department: "Design", Date: "24 Okt 2023", ClockIn: "09:15 AM", ClockOut: "06:00 PM", TotalHours: "8h 45m", Status: "TERLAMBAT"},
		{ID: 3, UserID: 103, UserName: "Andi Wijaya", UserNIP: "IN-303", Department: "Marketing", Date: "24 Okt 2023", ClockIn: "--:--", ClockOut: "--:--", TotalHours: "0h 0m", Status: "ABSEN"},
		{ID: 4, UserID: 104, UserName: "Rina Permata", UserNIP: "IN-304", Department: "Engineering", Date: "24 Okt 2023", ClockIn: "--:--", ClockOut: "--:--", TotalHours: "0h 0m", Status: "CUTI"},
	}
}

func (r *reportingRepository) GetAttendanceRecap(month, department, search string, limit int32) (*repProto.ReportingStats, []*model.AttendanceReport, error) {
	stats := &repProto.ReportingStats{
		TotalHadir:     4821,
		PctTepatWaktu:  92,
		TotalTerlambat: 312,
		TotalAbsen:     84,
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
