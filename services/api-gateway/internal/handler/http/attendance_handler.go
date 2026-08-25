package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"api-gateway/internal/client"
	"api-gateway/internal/middleware"
	attendanceProto "proto/attendance"
	userProto "proto/user"
)

type AttendanceHTTPHandler struct {
	attendanceSvc attendanceProto.AttendanceServiceServer
	userSvc       userProto.UserServiceServer
}

func NewAttendanceHTTPHandler(attendanceSvc attendanceProto.AttendanceServiceServer, userSvc userProto.UserServiceServer) *AttendanceHTTPHandler {
	return &AttendanceHTTPHandler{attendanceSvc: attendanceSvc, userSvc: userSvc}
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

func (h *AttendanceHTTPHandler) HandleGetPoskoQR(w http.ResponseWriter, r *http.Request) {
	res, _ := h.attendanceSvc.GetPoskoQR(r.Context(), &attendanceProto.GetPoskoQRRequest{})
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    res.Success,
		"posko_name": res.PoskoName,
		"address":    res.Address,
		"qr_token":   res.QrToken,
		"qr_image":   res.QrImage,
		"coordinates": map[string]interface{}{
			"latitude":          res.Latitude,
			"longitude":         res.Longitude,
			"max_radius_meters": res.MaxRadiusMeters,
		},
	})
}

func (h *AttendanceHTTPHandler) HandleAdminLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		res, _ := h.attendanceSvc.GetPoskoQR(r.Context(), &attendanceProto.GetPoskoQRRequest{})
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success":       true,
			"name":          res.PoskoName,
			"address":       res.Address,
			"latitude":      res.Latitude,
			"longitude":     res.Longitude,
			"radius_meters": res.MaxRadiusMeters,
			"qr_token":      res.QrToken,
		})
		return
	}

	if r.Method == http.MethodPut || r.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil || len(bodyBytes) == 0 {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Payload JSON kosong atau tidak valid."})
			return
		}
		var req struct {
			Name         string  `json:"name"`
			Address      string  `json:"address"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			RadiusMeters float64 `json:"radius_meters"`
		}
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Payload JSON tidak valid: " + err.Error()})
			return
		}

		var updatedToken string
		if stub, ok := h.attendanceSvc.(*client.AttendanceClientDirectStub); ok {
			cfg := stub.GetLocationConfig()
			if req.Name != "" {
				cfg.Name = req.Name
			}
			if req.Address != "" {
				cfg.Address = req.Address
			}
			if req.Latitude != 0 {
				cfg.Latitude = req.Latitude
			}
			if req.Longitude != 0 {
				cfg.Longitude = req.Longitude
			}
			if req.RadiusMeters > 0 {
				cfg.RadiusMeters = req.RadiusMeters
			}

			// Generate a new unique QR Code token automatically whenever location is updated
			cfg.QRToken = fmt.Sprintf("LOPI-Q-POSKO-BULUKUMBA-%X", time.Now().UnixNano()%1000000)
			updatedToken = cfg.QRToken
			_ = stub.SaveLocationConfig(cfg)
		}
		client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "UPDATE_LOCATION_CONFIG", fmt.Sprintf("Admin memperbarui kordinat / radius Geofence Posko 112 [Lat: %.6f, Lon: %.6f, Radius: %.1fm]", req.Latitude, req.Longitude, req.RadiusMeters), client.GetClientIP(r), r.UserAgent())
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"message":  "Pengaturan Lokasi Posko & Radius Geofence berhasil disimpan.",
			"qr_token": updatedToken,
		})
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (h *AttendanceHTTPHandler) HandleClockIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	nip := r.Context().Value("nip").(string)

	userRes, _ := h.userSvc.GetProfile(r.Context(), &userProto.GetProfileRequest{UserId: int32(userID)})

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		QRToken   string  `json:"qr_token"`
	}
	_ = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&body)

	posko, _ := h.attendanceSvc.GetPoskoQR(r.Context(), &attendanceProto.GetPoskoQRRequest{})
	poskoLat := -5.5645
	poskoLon := 120.1945
	maxRad := 2.0
	if posko != nil && posko.Latitude != 0 {
		poskoLat = posko.Latitude
		poskoLon = posko.Longitude
		maxRad = posko.MaxRadiusMeters
	}

	dist := calculateHaversine(poskoLat, poskoLon, body.Latitude, body.Longitude)
	if dist > maxRad {
		middleware.RespondJSON(w, http.StatusForbidden, map[string]interface{}{
			"success":           false,
			"error":             fmt.Sprintf("Presensi Ditolak! Lokasi Anda berjarak %.2fm dari Posko 112 (Batas Radius <= %.1fm).", dist, maxRad),
			"distance_meters":   dist,
			"max_radius_meters": maxRad,
		})
		return
	}

	name := "Call Taker"
	if userRes != nil && userRes.User != nil {
		name = userRes.User.Name
	}

	res, err := h.attendanceSvc.ClockIn(r.Context(), &attendanceProto.ClockInRequest{
		UserId:    int32(userID),
		UserNip:   nip,
		UserName:  name,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
		QrToken:   body.QRToken,
	})
	if err != nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	client.RecordActivityLog(userID, nip, name, "PRESENSI_MASUK_SCAN", fmt.Sprintf("Presensi MASUK Siaga Call Taker via Kamera Geofence Posko (Jarak: %.2fm)", dist), client.GetClientIP(r), r.UserAgent())

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Presensi MASUK Berhasil! Petugas: %s. Jarak: %.2fm (Batas <= %.1fm).", name, dist, maxRad),
		"record":  res.Record,
	})
}

func (h *AttendanceHTTPHandler) HandleClockOut(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	nip := r.Context().Value("nip").(string)

	userRes, _ := h.userSvc.GetProfile(r.Context(), &userProto.GetProfileRequest{UserId: int32(userID)})
	name := "Call Taker"
	if userRes != nil && userRes.User != nil {
		name = userRes.User.Name
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	_ = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&body)

	posko, _ := h.attendanceSvc.GetPoskoQR(r.Context(), &attendanceProto.GetPoskoQRRequest{})
	poskoLat := -5.5645
	poskoLon := 120.1945
	maxRad := 2.0
	if posko != nil && posko.Latitude != 0 {
		poskoLat = posko.Latitude
		poskoLon = posko.Longitude
		maxRad = posko.MaxRadiusMeters
	}

	dist := calculateHaversine(poskoLat, poskoLon, body.Latitude, body.Longitude)
	if dist > maxRad {
		middleware.RespondJSON(w, http.StatusForbidden, map[string]interface{}{
			"success":           false,
			"error":             fmt.Sprintf("Presensi PULANG Ditolak! Lokasi Anda berjarak %.2fm dari Posko 112 (Batas Radius <= %.1fm).", dist, maxRad),
			"distance_meters":   dist,
			"max_radius_meters": maxRad,
		})
		return
	}

	res, err := h.attendanceSvc.ClockOut(r.Context(), &attendanceProto.ClockOutRequest{
		UserId:    int32(userID),
		UserNip:   nip,
		UserName:  name,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	})
	if err != nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	client.RecordActivityLog(userID, nip, name, "PRESENSI_PULANG_SCAN", fmt.Sprintf("Presensi PULANG Siaga Call Taker via Kamera Geofence Posko (Jarak: %.2fm)", dist), client.GetClientIP(r), r.UserAgent())

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Presensi PULANG Berhasil! Jarak: %.2fm (Batas <= %.1fm).", dist, maxRad),
		"record":  res.Record,
	})
}

func (h *AttendanceHTTPHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	nip := ""
	if v := r.Context().Value("nip"); v != nil {
		nip = v.(string)
	}
	role := ""
	if v := r.Context().Value("role"); v != nil {
		role = v.(string)
	}

	req := &attendanceProto.GetHistoryRequest{Limit: 10000}
	if role != "admin" && role != "superadmin" {
		req.UserId = int32(userID)
		req.UserNip = nip
	}

	res, _ := h.attendanceSvc.GetHistory(r.Context(), req)
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "history": res.History})
}


func (h *AttendanceHTTPHandler) HandleGetTodayStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	userNIP := ""
	if nipVal := r.Context().Value("nip"); nipVal != nil {
		userNIP = nipVal.(string)
	}
	if userNIP == "" {
		if nipVal := r.Context().Value("user_nip"); nipVal != nil {
			userNIP = nipVal.(string)
		}
	}

	res, _ := h.attendanceSvc.GetTodayStatus(r.Context(), &attendanceProto.GetTodayStatusRequest{
		UserId:  int32(userID),
		UserNip: userNIP,
	})

	clockInTime := ""
	clockOutTime := ""
	var inDist, outDist float64
	if res.Masuk != nil && res.Masuk.Timestamp != "" {
		clockInTime = res.Masuk.Timestamp + " WITA"
		inDist = res.Masuk.DistanceMeters
	}
	if res.Pulang != nil && res.Pulang.Timestamp != "" {
		clockOutTime = res.Pulang.Timestamp + " WITA"
		outDist = res.Pulang.DistanceMeters
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":         true,
		"is_masuk":        res.IsMasuk,
		"is_pulang":       res.IsPulang,
		"has_clocked_in":  res.IsMasuk,
		"has_clocked_out": res.IsPulang,
		"clock_in_time":   clockInTime,
		"clock_out_time":  clockOutTime,
		"masuk": map[string]interface{}{
			"timestamp":       clockInTime,
			"time":            clockInTime,
			"distance_meters": inDist,
		},
		"pulang": map[string]interface{}{
			"timestamp":       clockOutTime,
			"time":            clockOutTime,
			"distance_meters": outDist,
		},
	})
}

func (h *AttendanceHTTPHandler) HandleAdminSchedules(w http.ResponseWriter, r *http.Request) {
	stub := &client.AttendanceClientDirectStub{}
	if r.Method == http.MethodGet {
		cfg := stub.GetSchedulesConfig()
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "schedules": cfg})
		return
	}
	if r.Method == http.MethodPut || r.Method == http.MethodPost {
		var raw interface{}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
			return
		}
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "JSON format error: " + err.Error()})
			return
		}
		if err := stub.SaveSchedulesConfig(raw); err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "UPDATE_SCHEDULE", "Super Admin menyimpan perubahan pembagian roster shift bulan ini", client.GetClientIP(r), r.UserAgent())
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Jadwal shift berhasil disimpan."})
		return
	}
	middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
}

func (h *AttendanceHTTPHandler) HandleLeaveRequests(w http.ResponseWriter, r *http.Request) {
	stub := &client.AttendanceClientDirectStub{}
	if r.Method == http.MethodGet {
		reqs := stub.GetLeaveRequests()
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "requests": reqs})
		return
	}
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		var raw interface{}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
			return
		}
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "JSON format error: " + err.Error()})
			return
		}
		if err := stub.SaveLeaveRequests(raw); err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Pengajuan berhasil disimpan."})
		return
	}
	middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
}

func (h *AttendanceHTTPHandler) HandleCleanupDuplicates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	deletedCount, err := client.CleanupDuplicatePresensiRecords()
	if err != nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("Berhasil membersihkan %d data presensi ganda dari database PostgreSQL.", deletedCount),
		"deleted_count": deletedCount,
	})
}

func (h *AttendanceHTTPHandler) HandleDeleteRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "ID record presensi tidak valid."})
		return
	}
	idStr := parts[len(parts)-1]
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id <= 0 {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "ID record presensi tidak valid."})
		return
	}
	if err := client.DeletePresensiRecord(id); err != nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "DELETE_ATTENDANCE_RECORD", fmt.Sprintf("Admin menghapus record presensi ID #%d dari database", id), client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Record presensi ID #%d berhasil dihapus dari database PostgreSQL.", id),
	})
}

func (h *AttendanceHTTPHandler) HandleAdminManualEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Request body tidak valid."})
		return
	}
	var req client.AdminManualPresensiReq
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "JSON format error: " + err.Error()})
		return
	}
	if err := client.SaveAdminManualPresensi(req); err != nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	client.RecordActivityLog(1, "199501012020011000", "Muhammad Aswan", "MANUAL_ATTENDANCE_ENTRY", fmt.Sprintf("Admin menyimpan perubahan presensi manual petugas %s (NIP: %s) tanggal %s [Status: %s, Masuk: %s, Pulang: %s]", req.Name, req.NIP, req.Date, req.Status, req.TimeMasuk, req.TimePulang), client.GetClientIP(r), r.UserAgent())
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Presensi %s pada tanggal %s berhasil di-update manual oleh Admin.", req.Name, req.Date),
	})
}

