package repository

import (
	"context"
	"database/sql"
	"math"
	"sync"
	"time"

	"attendance-service/internal/model"
)

// Posko NTPD 112 Bulukumba Center Coordinates
const (
	PoskoLat        = -5.548981
	PoskoLng        = 120.197943
	PoskoMaxRadiusM = 100.0
	PoskoName       = "Posko NTPD 112 Bulukumba"
	PoskoAddress    = "Jl. Jend. Sudirman No. 1, Caile, Ujung Bulu, Kabupaten Bulukumba"
)

type AttendanceRepository interface {
	GetPoskoInfo() (string, string, float64, float64, float64)
	CreateRecord(ctx context.Context, rec *model.Attendance) (*model.Attendance, error)
	GetHistory(ctx context.Context, userID int32, limit int32) ([]*model.Attendance, error)
	GetTodayStatus(ctx context.Context, userID int32) (*model.Attendance, *model.Attendance, error)
}

type attendanceRepository struct {
	db        *sql.DB
	mu        sync.RWMutex
	memRecords []*model.Attendance
}

func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{
		db:         db,
		memRecords: make([]*model.Attendance, 0),
	}
}

// CalculateDistanceMeters calculates Haversine distance between 2 coordinates in meters
func CalculateDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (r *attendanceRepository) GetPoskoInfo() (string, string, float64, float64, float64) {
	return PoskoName, PoskoAddress, PoskoLat, PoskoLng, PoskoMaxRadiusM
}

func (r *attendanceRepository) CreateRecord(ctx context.Context, rec *model.Attendance) (*model.Attendance, error) {
	rec.DistanceMeters = CalculateDistanceMeters(rec.Latitude, rec.Longitude, PoskoLat, PoskoLng)
	rec.WithinRadius = rec.DistanceMeters <= PoskoMaxRadiusM
	rec.Timestamp = time.Now()
	rec.CreatedAt = time.Now()

	if r.db != nil {
		query := `
			INSERT INTO tb_lopiq_attendances 
			(user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`
		var id int32
		err := r.db.QueryRowContext(ctx, query,
			rec.UserID, rec.UserNIP, rec.UserName, rec.Type, rec.Timestamp,
			rec.Latitude, rec.Longitude, rec.DistanceMeters, rec.WithinRadius, rec.CreatedAt,
		).Scan(&id)

		if err == nil {
			rec.ID = id
			return rec, nil
		}
	}

	// Fallback to In-Memory
	r.mu.Lock()
	defer r.mu.Unlock()
	rec.ID = int32(len(r.memRecords) + 1)
	r.memRecords = append(r.memRecords, rec)

	return rec, nil
}

func (r *attendanceRepository) GetHistory(ctx context.Context, userID int32, limit int32) ([]*model.Attendance, error) {
	if limit <= 0 {
		limit = 20
	}

	if r.db != nil {
		query := `
			SELECT id, user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius, created_at
			FROM tb_lopiq_attendances
			WHERE user_id = $1
			ORDER BY timestamp DESC
			LIMIT $2
		`
		rows, err := r.db.QueryContext(ctx, query, userID, limit)
		if err == nil {
			defer rows.Close()
			var res []*model.Attendance
			for rows.Next() {
				var item model.Attendance
				if err := rows.Scan(&item.ID, &item.UserID, &item.UserNIP, &item.UserName, &item.Type, &item.Timestamp, &item.Latitude, &item.Longitude, &item.DistanceMeters, &item.WithinRadius, &item.CreatedAt); err == nil {
					res = append(res, &item)
				}
			}
			return res, nil
		}
	}

	// Fallback to In-Memory
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*model.Attendance
	for i := len(r.memRecords) - 1; i >= 0; i-- {
		if r.memRecords[i].UserID == userID {
			res = append(res, r.memRecords[i])
			if int32(len(res)) >= limit {
				break
			}
		}
	}
	return res, nil
}

func (r *attendanceRepository) GetTodayStatus(ctx context.Context, userID int32) (*model.Attendance, *model.Attendance, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if r.db != nil {
		query := `
			SELECT id, user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius, created_at
			FROM tb_lopiq_attendances
			WHERE user_id = $1 AND timestamp >= $2
			ORDER BY timestamp ASC
		`
		rows, err := r.db.QueryContext(ctx, query, userID, startOfDay)
		if err == nil {
			defer rows.Close()
			var masuk *model.Attendance
			var pulang *model.Attendance

			for rows.Next() {
				var item model.Attendance
				if err := rows.Scan(&item.ID, &item.UserID, &item.UserNIP, &item.UserName, &item.Type, &item.Timestamp, &item.Latitude, &item.Longitude, &item.DistanceMeters, &item.WithinRadius, &item.CreatedAt); err == nil {
					if item.Type == "MASUK" && masuk == nil {
						masuk = &item
					} else if item.Type == "PULANG" {
						pulang = &item
					}
				}
			}
			return masuk, pulang, nil
		}
	}

	// Fallback to In-Memory
	r.mu.RLock()
	defer r.mu.RUnlock()
	var masuk *model.Attendance
	var pulang *model.Attendance

	for _, item := range r.memRecords {
		if item.UserID == userID && item.Timestamp.After(startOfDay) {
			if item.Type == "MASUK" && masuk == nil {
				masuk = item
			} else if item.Type == "PULANG" {
				pulang = item
			}
		}
	}

	return masuk, pulang, nil
}
