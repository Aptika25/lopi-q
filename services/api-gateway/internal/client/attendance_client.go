package client

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	attendanceProto "proto/attendance"
)


type LocationConfig struct {
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters float64 `json:"radius_meters"`
	QRToken      string  `json:"qr_token"`
}

type AttendanceClient struct {
	server attendanceProto.AttendanceServiceServer
}

func NewAttendanceClient(server attendanceProto.AttendanceServiceServer) *AttendanceClient {
	return &AttendanceClient{server: server}
}

type AttendanceClientDirectStub struct {
	mu sync.Mutex
}

func getPostgresDB(dbHost string) (*sql.DB, error) {
	// 1. Try db_lopi-q_presensi with super_admin_apps
	connSuper := fmt.Sprintf("host=%s port=5432 user=super_admin_apps password=superAdminAppsPassword@2k26# dbname=db_lopi-q_presensi sslmode=disable", dbHost)
	if db, err := sql.Open("postgres", connSuper); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err == nil {
			return db, nil
		}
		db.Close()
	}

	// 2. Try db_lopi-q_presensi with user_lopi-q_presensi
	connUser := fmt.Sprintf("host=%s port=5432 user=user_lopi-q_presensi password=lopi-qpresensiPassword@2k26# dbname=db_lopi-q_presensi sslmode=disable", dbHost)
	if db, err := sql.Open("postgres", connUser); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err == nil {
			return db, nil
		}
		db.Close()
	}

	// 3. Fallback to db_lopiq_auth with super_admin_apps
	connAuth := fmt.Sprintf("host=%s port=5432 user=super_admin_apps password=superAdminAppsPassword@2k26# dbname=db_lopiq_auth sslmode=disable", dbHost)
	return sql.Open("postgres", connAuth)
}

func fetchPostgresLocationConfig() (LocationConfig, bool) {
	var cfg LocationConfig
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	db, err := getPostgresDB(dbHost)
	if err != nil {
		return cfg, false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _ = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posko_locations (
			id INT PRIMARY KEY DEFAULT 1,
			name VARCHAR(255) NOT NULL,
			address TEXT NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			radius_meters DOUBLE PRECISION NOT NULL,
			qr_token VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	row := db.QueryRowContext(ctx, "SELECT name, address, latitude, longitude, radius_meters, qr_token FROM posko_locations WHERE id = 1;")
	if err := row.Scan(&cfg.Name, &cfg.Address, &cfg.Latitude, &cfg.Longitude, &cfg.RadiusMeters, &cfg.QRToken); err == nil && cfg.Latitude != 0 {
		return cfg, true
	}

	return cfg, false
}

func savePostgresLocationConfig(cfg LocationConfig) error {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	db, err := getPostgresDB(dbHost)
	if err != nil {
		log.Printf("[Location DB Error] Failed to open PostgreSQL connection: %v", err)
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posko_locations (
			id INT PRIMARY KEY DEFAULT 1,
			name VARCHAR(255) NOT NULL,
			address TEXT NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			radius_meters DOUBLE PRECISION NOT NULL,
			qr_token VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Printf("[Location DB Warning] Table creation check: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO posko_locations (id, name, address, latitude, longitude, radius_meters, qr_token, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			address = EXCLUDED.address,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			radius_meters = EXCLUDED.radius_meters,
			qr_token = EXCLUDED.qr_token,
			updated_at = NOW();
	`, cfg.Name, cfg.Address, cfg.Latitude, cfg.Longitude, cfg.RadiusMeters, cfg.QRToken)

	if err != nil {
		log.Printf("[Location DB Error] Failed to upsert posko_locations: %v", err)
		return err
	}

	log.Printf("[Location DB Success] Location and QR token successfully saved to PostgreSQL database posko_locations (ID=1)")
	return nil
}

func (s *AttendanceClientDirectStub) getLocationFilePath() string {
	candidates := []string{
		"data/location.json",
		"../data/location.json",
		"../../data/location.json",
		"/app/data/location.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/location.json"
}

func (s *AttendanceClientDirectStub) GetLocationConfig() LocationConfig {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pgCfg, ok := fetchPostgresLocationConfig(); ok {
		return pgCfg
	}

	p := s.getLocationFilePath()
	data, err := os.ReadFile(p)
	if err == nil {
		var cfg LocationConfig
		if err := json.Unmarshal(data, &cfg); err == nil && cfg.Latitude != 0 {
			return cfg
		}
	}

	return LocationConfig{
		Name:         "Posko Siaga NTPD 112 Kabupaten Bulukumba",
		Address:      "Jl. Jend. Sudirman No. 1, Caile, Kec. Ujung Bulu, Kabupaten Bulukumba, Sulawesi Selatan",
		Latitude:     -5.5645,
		Longitude:    120.1945,
		RadiusMeters: 2.0,
		QRToken:      "LOPI-Q-POSKO-BULUKUMBA-2026-NTPD112",
	}
}

func (s *AttendanceClientDirectStub) SaveLocationConfig(cfg LocationConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := s.getLocationFilePath()
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err == nil {
		if data, err := json.MarshalIndent(cfg, "", "  "); err == nil {
			_ = os.WriteFile(p, data, 0644)
		}
	}

	return savePostgresLocationConfig(cfg)
}

func (s *AttendanceClientDirectStub) GetPoskoQR(ctx context.Context, req *attendanceProto.GetPoskoQRRequest) (*attendanceProto.GetPoskoQRResponse, error) {
	cfg := s.GetLocationConfig()
	return &attendanceProto.GetPoskoQRResponse{
		Success:         true,
		PoskoName:       cfg.Name,
		Address:         cfg.Address,
		QrToken:         cfg.QRToken,
		Latitude:        cfg.Latitude,
		Longitude:       cfg.Longitude,
		MaxRadiusMeters: cfg.RadiusMeters,
	}, nil
}

func calculateHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return math.Round((R*c)*100) / 100
}

func (s *AttendanceClientDirectStub) ClockIn(ctx context.Context, req *attendanceProto.ClockInRequest) (*attendanceProto.AttendanceResponse, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	cfg := s.GetLocationConfig()
	dist := calculateHaversine(cfg.Latitude, cfg.Longitude, req.Latitude, req.Longitude)

	if db, err := getPostgresDB(dbHost); err == nil {
		defer db.Close()
		queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		_, _ = db.ExecContext(queryCtx, `
			CREATE TABLE IF NOT EXISTS presensi_records (
				id SERIAL PRIMARY KEY,
				user_id INT NOT NULL,
				user_nip VARCHAR(50) NOT NULL,
				user_name VARCHAR(255) NOT NULL,
				type VARCHAR(20) NOT NULL,
				timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				latitude DOUBLE PRECISION NOT NULL,
				longitude DOUBLE PRECISION NOT NULL,
				distance_meters DOUBLE PRECISION NOT NULL,
				within_radius BOOLEAN DEFAULT TRUE
			);
		`)

		// AUTO-REPAIR FIX: Revert corrupted evening records (>= 18:00 WITA) back to MASUK if there was no prior MASUK on that same date
		_, _ = db.ExecContext(queryCtx, `
			UPDATE presensi_records p1
			SET type = 'MASUK'
			WHERE p1.type = 'PULANG'
			  AND EXTRACT(HOUR FROM (p1.timestamp AT TIME ZONE 'Asia/Makassar')) >= 18
			  AND NOT EXISTS (
			      SELECT 1 FROM presensi_records p2
			      WHERE (p2.user_id = p1.user_id OR p2.user_nip = p1.user_nip)
			        AND p2.type = 'MASUK'
			        AND (p2.timestamp AT TIME ZONE 'Asia/Makassar')::date = (p1.timestamp AT TIME ZONE 'Asia/Makassar')::date
			        AND p2.id < p1.id
			  );
		`)

		// Check if user already has an active MASUK record within the last 30 hours without a PULANG record
		var hasActiveMasuk bool
		_ = db.QueryRowContext(queryCtx, `
			SELECT EXISTS (
				SELECT 1 FROM presensi_records p1
				WHERE (($2 <> '' AND REPLACE(user_nip, ' ', '') = REPLACE($2, ' ', '')) OR ($2 = '' AND user_id = $1))
				  AND type = 'MASUK'
				  AND (timestamp AT TIME ZONE 'Asia/Makassar') >= (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Makassar') - INTERVAL '30 hours'
				  AND NOT EXISTS (
				      SELECT 1 FROM presensi_records p2
				      WHERE ((p2.user_nip <> '' AND REPLACE(p2.user_nip, ' ', '') = REPLACE(p1.user_nip, ' ', '')) OR (p2.user_nip = '' AND p2.user_id = p1.user_id))
				        AND p2.type = 'PULANG'
				        AND p2.id > p1.id
				  )
			);
		`, req.UserId, req.UserNip).Scan(&hasActiveMasuk)

		recordType := "MASUK"
		msgType := "MASUK"
		if hasActiveMasuk {
			recordType = "PULANG"
			msgType = "PULANG (Otomatis dari Presensi Masuk Shift Malam Sebelumnya)"
		}

		_, err = db.ExecContext(queryCtx, `
			INSERT INTO presensi_records (user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius)
			VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, $8)
		`, req.UserId, req.UserNip, req.UserName, recordType, req.Latitude, req.Longitude, dist, dist <= cfg.RadiusMeters)

		if err != nil {
			log.Printf("[ClockIn DB Error] Failed to insert record: %v", err)
		} else {
			log.Printf("[ClockIn DB Success] Presensi %s for user %s (NIP %s) saved to PostgreSQL presensi_records", recordType, req.UserName, req.UserNip)
		}

		return &attendanceProto.AttendanceResponse{
			Success: true,
			Message: fmt.Sprintf("Presensi %s Berhasil! Petugas: %s. Jarak: %.2fm dari Posko.", msgType, req.UserName, dist),
		}, nil
	}

	return &attendanceProto.AttendanceResponse{
		Success: true,
		Message: fmt.Sprintf("Presensi MASUK Berhasil! Petugas: %s. Jarak: %.2fm dari Posko.", req.UserName, dist),
	}, nil
}

func (s *AttendanceClientDirectStub) ClockOut(ctx context.Context, req *attendanceProto.ClockOutRequest) (*attendanceProto.AttendanceResponse, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	cfg := s.GetLocationConfig()
	dist := calculateHaversine(cfg.Latitude, cfg.Longitude, req.Latitude, req.Longitude)

	if db, err := getPostgresDB(dbHost); err == nil {
		defer db.Close()
		queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		_, _ = db.ExecContext(queryCtx, `
			CREATE TABLE IF NOT EXISTS presensi_records (
				id SERIAL PRIMARY KEY,
				user_id INT NOT NULL,
				user_nip VARCHAR(50) NOT NULL,
				user_name VARCHAR(255) NOT NULL,
				type VARCHAR(20) NOT NULL,
				timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				latitude DOUBLE PRECISION NOT NULL,
				longitude DOUBLE PRECISION NOT NULL,
				distance_meters DOUBLE PRECISION NOT NULL,
				within_radius BOOLEAN DEFAULT TRUE
			);
		`)

		_, err = db.ExecContext(queryCtx, `
			INSERT INTO presensi_records (user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius)
			VALUES ($1, $2, $3, 'PULANG', NOW(), $4, $5, $6, $7)
		`, req.UserId, req.UserNip, req.UserName, req.Latitude, req.Longitude, dist, dist <= cfg.RadiusMeters)

		if err != nil {
			log.Printf("[ClockOut DB Error] Failed to insert record: %v", err)
		} else {
			log.Printf("[ClockOut DB Success] Presensi PULANG for user %s (NIP %s) saved to PostgreSQL presensi_records", req.UserName, req.UserNip)
		}
	}

	return &attendanceProto.AttendanceResponse{
		Success: true,
		Message: fmt.Sprintf("Presensi PULANG Berhasil! Petugas: %s. Jarak: %.2fm dari Posko.", req.UserName, dist),
	}, nil
}

func (s *AttendanceClientDirectStub) GetHistory(ctx context.Context, req *attendanceProto.GetHistoryRequest) (*attendanceProto.GetHistoryResponse, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	resp := &attendanceProto.GetHistoryResponse{Success: true, History: []*attendanceProto.AttendanceRecord{}}

	if db, err := getPostgresDB(dbHost); err == nil {
		defer db.Close()
		queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		cleanNIP := strings.ReplaceAll(req.UserNip, " ", "")
		limit := req.Limit
		if limit <= 0 {
			limit = 500
		}

		rows, err := db.QueryContext(queryCtx, `
			SELECT id, user_id, COALESCE(user_nip, ''), COALESCE(user_name, ''), type, to_char(timestamp AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD HH24:MI:SS'), distance_meters 
			FROM presensi_records 
			WHERE ($1 = 0 AND $2 = '') 
			   OR ($2 <> '' AND REPLACE(user_nip, ' ', '') = $2)
			   OR ($2 = '' AND user_id = $1)
			ORDER BY id ASC LIMIT $3;
		`, req.UserId, cleanNIP, limit)

		if err == nil {
			for rows.Next() {
				var rId, rUserId int32
				var rNip, rName, pType, pTime string
				var dist float64
				if err := rows.Scan(&rId, &rUserId, &rNip, &rName, &pType, &pTime, &dist); err == nil {
					resp.History = append(resp.History, &attendanceProto.AttendanceRecord{
						Id:             rId,
						UserId:         rUserId,
						UserNip:        rNip,
						UserName:       rName,
						Type:           pType,
						Timestamp:      pTime,
						DistanceMeters: dist,
					})
				}
			}
			rows.Close()
		}
	}

	return resp, nil
}


func (s *AttendanceClientDirectStub) GetTodayStatus(ctx context.Context, req *attendanceProto.GetTodayStatusRequest) (*attendanceProto.GetTodayStatusResponse, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	resp := &attendanceProto.GetTodayStatusResponse{Success: true}

	if db, err := getPostgresDB(dbHost); err == nil {
		defer db.Close()
		queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		cleanNIP := strings.ReplaceAll(req.UserNip, " ", "")

		// 1. Fetch records for TODAY (Asia/Makassar date = TODAY)
		rows, err := db.QueryContext(queryCtx, `
			SELECT type, to_char(timestamp AT TIME ZONE 'Asia/Makassar', 'HH24:MI:SS'), distance_meters 
			FROM presensi_records 
			WHERE (($2 <> '' AND REPLACE(user_nip, ' ', '') = $2) OR ($2 = '' AND user_id = $1))
			  AND (timestamp AT TIME ZONE 'Asia/Makassar')::date = (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Makassar')::date 
			ORDER BY id ASC;
		`, req.UserId, cleanNIP)

		if err == nil {
			for rows.Next() {
				var pType, pTime string
				var dist float64
				if err := rows.Scan(&pType, &pTime, &dist); err == nil {
					if pType == "MASUK" {
						resp.IsMasuk = true
						resp.Masuk = &attendanceProto.AttendanceRecord{
							Timestamp:      pTime,
							DistanceMeters: dist,
						}
					} else if pType == "PULANG" {
						resp.IsPulang = true
						resp.Pulang = &attendanceProto.AttendanceRecord{
							Timestamp:      pTime,
							DistanceMeters: dist,
						}
					}
				}
			}
			rows.Close()
		}

		// 2. OVERNIGHT SHIFT CHECK:
		// If NO 'MASUK' found today, check if user has an unclosed 'MASUK' from YESTERDAY (within last 30 hours)
		// which has no subsequent 'PULANG' record yet.
		if !resp.IsMasuk {
			var pTime string
			var dist float64
			err := db.QueryRowContext(queryCtx, `
				SELECT to_char(timestamp AT TIME ZONE 'Asia/Makassar', 'HH24:MI:SS'), distance_meters
				FROM presensi_records p1
				WHERE (($2 <> '' AND REPLACE(user_nip, ' ', '') = $2) OR ($2 = '' AND user_id = $1))
				  AND type = 'MASUK'
				  AND (timestamp AT TIME ZONE 'Asia/Makassar') >= (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Makassar') - INTERVAL '30 hours'
				  AND NOT EXISTS (
				      SELECT 1 FROM presensi_records p2 
				      WHERE ((p2.user_nip <> '' AND REPLACE(p2.user_nip, ' ', '') = REPLACE(p1.user_nip, ' ', '')) OR (p2.user_nip = '' AND p2.user_id = p1.user_id))
				        AND p2.type = 'PULANG' 
				        AND p2.id > p1.id
				  )
				ORDER BY id DESC LIMIT 1;
			`, req.UserId, cleanNIP).Scan(&pTime, &dist)

			if err == nil && pTime != "" {
				resp.IsMasuk = true
				resp.Masuk = &attendanceProto.AttendanceRecord{
					Timestamp:      pTime + " (Shift Kemarin)",
					DistanceMeters: dist,
				}
			}
		}
	}

	return resp, nil
}


func fetchPostgresSchedulesConfig() (interface{}, bool) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	db, err := getPostgresDB(dbHost)
	if err != nil {
		return nil, false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _ = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posko_schedules (
			id INT PRIMARY KEY DEFAULT 1,
			schedule_data JSONB NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	var jsonRaw []byte
	row := db.QueryRowContext(ctx, "SELECT schedule_data FROM posko_schedules WHERE id = 1;")
	if err := row.Scan(&jsonRaw); err == nil && len(jsonRaw) > 0 {
		var result interface{}
		if err := json.Unmarshal(jsonRaw, &result); err == nil {
			return result, true
		}
	}
	return nil, false
}

func savePostgresSchedulesConfig(raw interface{}) error {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	db, err := getPostgresDB(dbHost)
	if err != nil {
		log.Printf("[Schedules DB Error] Connection failed: %v", err)
		return err
	}
	defer db.Close()

	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posko_schedules (
			id INT PRIMARY KEY DEFAULT 1,
			schedule_data JSONB NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Printf("[Schedules DB Warning] Table creation check: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO posko_schedules (id, schedule_data, updated_at)
		VALUES (1, $1, NOW())
		ON CONFLICT (id) DO UPDATE SET
			schedule_data = EXCLUDED.schedule_data,
			updated_at = NOW();
	`, string(jsonBytes))

	if err != nil {
		log.Printf("[Schedules DB Error] Failed to upsert posko_schedules: %v", err)
		return err
	}

	log.Printf("[Schedules DB Success] Schedules config saved to PostgreSQL database posko_schedules (ID=1)")
	return nil
}

func fetchPostgresLeaveRequests() (interface{}, bool) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	db, err := getPostgresDB(dbHost)
	if err != nil {
		return nil, false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _ = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posko_leave_requests (
			id INT PRIMARY KEY DEFAULT 1,
			requests_data JSONB NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	var jsonRaw []byte
	row := db.QueryRowContext(ctx, "SELECT requests_data FROM posko_leave_requests WHERE id = 1;")
	if err := row.Scan(&jsonRaw); err == nil && len(jsonRaw) > 0 {
		var result interface{}
		if err := json.Unmarshal(jsonRaw, &result); err == nil {
			return result, true
		}
	}
	return nil, false
}

func savePostgresLeaveRequests(raw interface{}) error {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}

	db, err := getPostgresDB(dbHost)
	if err != nil {
		log.Printf("[LeaveRequests DB Error] Connection failed: %v", err)
		return err
	}
	defer db.Close()

	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posko_leave_requests (
			id INT PRIMARY KEY DEFAULT 1,
			requests_data JSONB NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Printf("[LeaveRequests DB Warning] Table creation check: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO posko_leave_requests (id, requests_data, updated_at)
		VALUES (1, $1, NOW())
		ON CONFLICT (id) DO UPDATE SET
			requests_data = EXCLUDED.requests_data,
			updated_at = NOW();
	`, string(jsonBytes))

	if err != nil {
		log.Printf("[LeaveRequests DB Error] Failed to upsert posko_leave_requests: %v", err)
		return err
	}

	log.Printf("[LeaveRequests DB Success] Leave requests saved to PostgreSQL database posko_leave_requests (ID=1)")
	return nil
}

func (s *AttendanceClientDirectStub) getSchedulesFilePath() string {
	candidates := []string{
		"data/schedules.json",
		"../data/schedules.json",
		"../../data/schedules.json",
		"/app/data/schedules.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/schedules.json"
}

func sanitizePastSchedules(raw interface{}) interface{} {
	if raw == nil {
		return raw
	}
	mMap, ok := raw.(map[string]interface{})
	if !ok {
		return raw
	}
	days, okDays := mMap["daysInMonth"].([]interface{})
	if !okDays {
		return raw
	}

	cutoffStr := "2026-08-13"

	for _, day := range days {
		dMap, ok := day.(map[string]interface{})
		if !ok {
			continue
		}
		dStr, okStr := dMap["dateStr"].(string)
		if !okStr {
			continue
		}

		if dStr < cutoffStr {
			snaps, okSnaps := dMap["teamsSnapshot"].([]interface{})
			if okSnaps && len(snaps) > 0 {
				for _, t := range snaps {
					if tMap, ok := t.(map[string]interface{}); ok {
						if code, ok := tMap["code"].(string); ok && code == "B" {
							if members, ok := tMap["members"].([]interface{}); ok && len(members) > 0 {
								if m0, ok := members[0].(map[string]interface{}); ok {
									nipStr, _ := m0["nip"].(string)
									nameStr, _ := m0["name"].(string)
									cleanNip := strings.ReplaceAll(nipStr, " ", "")
									if cleanNip == "199107052025211081" || strings.Contains(nameStr, "Aswar") {
										m0["name"] = "Riswandi Risman"
										m0["nip"] = "20000206 202521 1 166"
										m0["jabatan"] = "OPERATOR LAYANAN OPERASIONAL"
										m0["unit"] = "Dinas Kesehatan"
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return raw
}

func (s *AttendanceClientDirectStub) GetSchedulesConfig() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pgResult, ok := fetchPostgresSchedulesConfig(); ok {
		return sanitizePastSchedules(pgResult)
	}

	p := s.getSchedulesFilePath()
	data, err := os.ReadFile(p)
	if err == nil {
		var result interface{}
		if err := json.Unmarshal(data, &result); err == nil {
			return sanitizePastSchedules(result)
		}
	}
	return nil
}

func preservePastSchedules(newRaw interface{}, existingRaw interface{}) interface{} {
	if newRaw == nil || existingRaw == nil {
		return newRaw
	}
	loc, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		loc = time.FixedZone("WITA", 8*3600)
	}
	todayStr := time.Now().In(loc).Format("2006-01-02")

	newMap, okNew := newRaw.(map[string]interface{})
	existingMap, okExisting := existingRaw.(map[string]interface{})
	if !okNew || !okExisting {
		return newRaw
	}

	newDays, okNewDays := newMap["daysInMonth"].([]interface{})
	existingDays, okExistDays := existingMap["daysInMonth"].([]interface{})
	if !okNewDays || !okExistDays {
		return newRaw
	}

	existPastMap := make(map[string]interface{})
	for _, day := range existingDays {
		if dMap, ok := day.(map[string]interface{}); ok {
			if dStr, ok := dMap["dateStr"].(string); ok && dStr < todayStr {
				existPastMap[dStr] = dMap
			}
		}
	}

	for i, day := range newDays {
		if dMap, ok := day.(map[string]interface{}); ok {
			if dStr, ok := dMap["dateStr"].(string); ok && dStr < todayStr {
				if oldDay, found := existPastMap[dStr]; found {
					newDays[i] = oldDay
				}
			}
		}
	}

	newMap["daysInMonth"] = newDays
	return newMap
}

func (s *AttendanceClientDirectStub) SaveSchedulesConfig(raw interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Preserve historic schedules for past dates
	if pgResult, ok := fetchPostgresSchedulesConfig(); ok && pgResult != nil {
		raw = preservePastSchedules(raw, pgResult)
	}
	raw = sanitizePastSchedules(raw)

	p := s.getSchedulesFilePath()
	dir := filepath.Dir(p)
	_ = os.MkdirAll(dir, 0755)

	if data, err := json.MarshalIndent(raw, "", "  "); err == nil {
		_ = os.WriteFile(p, data, 0644)
	}

	return savePostgresSchedulesConfig(raw)
}

func (s *AttendanceClientDirectStub) getLeaveRequestsFilePath() string {
	candidates := []string{
		"data/leave_requests.json",
		"../data/leave_requests.json",
		"../../data/leave_requests.json",
		"/app/data/leave_requests.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/leave_requests.json"
}

func (s *AttendanceClientDirectStub) GetLeaveRequests() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pgResult, ok := fetchPostgresLeaveRequests(); ok {
		return pgResult
	}

	p := s.getLeaveRequestsFilePath()
	data, err := os.ReadFile(p)
	if err == nil {
		var result interface{}
		if err := json.Unmarshal(data, &result); err == nil {
			return result
		}
	}
	return []interface{}{}
}

func (s *AttendanceClientDirectStub) SaveLeaveRequests(raw interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := s.getLeaveRequestsFilePath()
	dir := filepath.Dir(p)
	_ = os.MkdirAll(dir, 0755)

	return savePostgresLeaveRequests(raw)
}

func DeletePresensiRecord(id int) error {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}
	db, err := getPostgresDB(dbHost)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, "DELETE FROM presensi_records WHERE id = $1;", id)
	return err
}

func CleanupDuplicatePresensiRecords() (int64, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}
	db, err := getPostgresDB(dbHost)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, `
		DELETE FROM presensi_records p1
		WHERE EXISTS (
			SELECT 1 FROM presensi_records p2
			WHERE (p1.user_id = p2.user_id OR p1.user_nip = p2.user_nip)
			  AND p1.type = p2.type
			  AND (p1.timestamp AT TIME ZONE 'Asia/Makassar')::date = (p2.timestamp AT TIME ZONE 'Asia/Makassar')::date
			  AND p1.id > p2.id
		);
	`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type AdminManualPresensiReq struct {
	NIP        string `json:"nip"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	TimeMasuk  string `json:"time_masuk"`
	TimePulang string `json:"time_pulang"`
	Status     string `json:"status"`
	Note       string `json:"note"`
}

func SaveAdminManualPresensi(req AdminManualPresensiReq) error {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}
	db, err := getPostgresDB(dbHost)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dateISO := req.Date
	if strings.Contains(req.Date, "-") && len(strings.Split(req.Date, "-")[0]) == 2 {
		parts := strings.Split(req.Date, "-")
		dateISO = fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])
	}

	// Calculate nextDayISO for night shift / cross-midnight clock-out
	parsedDate, _ := time.Parse("2006-01-02", dateISO)
	nextDayISO := parsedDate.AddDate(0, 0, 1).Format("2006-01-02")

	cleanNip := strings.ReplaceAll(req.NIP, " ", "")

	// 1. Delete previous manual/auto records for this shift date
	// Delete MASUK on dateISO and PULANG on dateISO or nextDayISO (morning < 11:00)
	_, err = db.ExecContext(ctx, `
		DELETE FROM presensi_records 
		WHERE REPLACE(user_nip, ' ', '') = $1 
		  AND (
		    ((timestamp AT TIME ZONE 'Asia/Makassar')::date = $2::date)
		    OR ((timestamp AT TIME ZONE 'Asia/Makassar')::date = $3::date AND type = 'PULANG' AND EXTRACT(HOUR FROM (timestamp AT TIME ZONE 'Asia/Makassar')) < 11)
		  );
	`, cleanNip, dateISO, nextDayISO)
	if err != nil {
		return fmt.Errorf("gagal menghapus data lama: %v", err)
	}

	// Helper to format HH:MM into HH:MM:SS
	formatTimeStr := func(tStr string) string {
		tStr = strings.TrimSpace(tStr)
		if tStr == "" {
			return ""
		}
		parts := strings.Split(tStr, ":")
		if len(parts) == 2 {
			return fmt.Sprintf("%s:%s:00", parts[0], parts[1])
		}
		return tStr
	}

	timeMasukClean := formatTimeStr(req.TimeMasuk)
	timePulangClean := formatTimeStr(req.TimePulang)

	// 2. Insert MASUK if provided
	if timeMasukClean != "" && req.Status != "Tidak Hadir (Alpha)" && req.Status != "Sakit / Izin (Resmi)" {
		tsMasuk := fmt.Sprintf("%s %s+08", dateISO, timeMasukClean)
		_, err = db.ExecContext(ctx, `
			INSERT INTO presensi_records (user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius)
			VALUES (
				COALESCE((SELECT user_id FROM presensi_records WHERE REPLACE(user_nip, ' ', '') = $1 LIMIT 1), 1),
				$1, $2, 'MASUK', $3::timestamp with time zone, -5.5645, 120.1945, 0, true
			);
		`, cleanNip, req.Name, tsMasuk)
		if err != nil {
			return fmt.Errorf("gagal menyimpan jam masuk: %v", err)
		}
	}

	// 3. Insert PULANG if provided
	if timePulangClean != "" && req.Status != "Tidak Hadir (Alpha)" && req.Status != "Sakit / Izin (Resmi)" {
		pulangDate := dateISO

		masukHour := 0
		pulangHour := 0
		if timeMasukClean != "" {
			fmt.Sscanf(timeMasukClean, "%d", &masukHour)
		}
		if timePulangClean != "" {
			fmt.Sscanf(timePulangClean, "%d", &pulangHour)
		}

		if (masukHour >= 18 && pulangHour < 12) ||
			(timeMasukClean != "" && timePulangClean < timeMasukClean && pulangHour < 12) ||
			(timeMasukClean == "" && pulangHour < 11) {
			pulangDate = nextDayISO
		}

		tsPulang := fmt.Sprintf("%s %s+08", pulangDate, timePulangClean)
		_, err = db.ExecContext(ctx, `
			INSERT INTO presensi_records (user_id, user_nip, user_name, type, timestamp, latitude, longitude, distance_meters, within_radius)
			VALUES (
				COALESCE((SELECT user_id FROM presensi_records WHERE REPLACE(user_nip, ' ', '') = $1 LIMIT 1), 1),
				$1, $2, 'PULANG', $3::timestamp with time zone, -5.5645, 120.1945, 0, true
			);
		`, cleanNip, req.Name, tsPulang)
		if err != nil {
			return fmt.Errorf("gagal menyimpan jam pulang: %v", err)
		}
	}

	return nil
}

