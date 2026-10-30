package http

import (
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
		// Moving the geofence or rotating the QR token affects every intern, so only the Super Admin may do it.
		if !isSuperAdmin(r) {
			middleware.RespondJSON(w, http.StatusForbidden, map[string]interface{}{"success": false, "error": "Hanya Super Admin yang dapat mengubah lokasi Posko dan QR Code presensi."})
			return
		}
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
			QRToken      string  `json:"qr_token"`
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
			if req.QRToken != "" {
				cfg.QRToken = req.QRToken
			}
			updatedToken = cfg.QRToken
			_ = stub.SaveLocationConfig(cfg)
		}
		logAdminAction(r, "UPDATE_LOCATION_CONFIG", fmt.Sprintf("Admin memperbarui kordinat / radius Geofence Posko 112 [Lat: %.6f, Lon: %.6f, Radius: %.1fm]", req.Latitude, req.Longitude, req.RadiusMeters))
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"message":  "Pengaturan Lokasi Posko & Radius Geofence berhasil disimpan.",
			"qr_token": updatedToken,
		})
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

type presensiScanRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	QRToken   string  `json:"qr_token"`
}

// validatePresensiScan checks the scanned posko QR and the reported GPS position. On rejection it writes
// the response itself and returns ok=false, so nothing is recorded.
func (h *AttendanceHTTPHandler) validatePresensiScan(w http.ResponseWriter, r *http.Request, label string) (body presensiScanRequest, dist, maxRad float64, ok bool) {
	if r.Method != http.MethodPost {
		middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Payload presensi tidak valid."})
		return
	}
	if body.Latitude == 0 && body.Longitude == 0 {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Lokasi GPS tidak terdeteksi. Aktifkan GPS dan izin lokasi, lalu coba lagi."})
		return
	}

	posko, _ := h.attendanceSvc.GetPoskoQR(r.Context(), &attendanceProto.GetPoskoQRRequest{})
	poskoLat := -5.5578602
	poskoLon := 120.1937020
	maxRad = 25.0
	expectedToken := ""
	if posko != nil {
		expectedToken = strings.TrimSpace(posko.QrToken)
		if posko.Latitude != 0 {
			poskoLat = posko.Latitude
			poskoLon = posko.Longitude
			maxRad = posko.MaxRadiusMeters
		}
	}

	if expectedToken == "" || strings.TrimSpace(body.QRToken) != expectedToken {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "QR Code tidak valid. Scan QR Code presensi resmi yang terpasang di Posko."})
		return
	}

	dist = calculateHaversine(poskoLat, poskoLon, body.Latitude, body.Longitude)
	if dist > maxRad {
		middleware.RespondJSON(w, http.StatusForbidden, map[string]interface{}{
			"success":           false,
			"error":             fmt.Sprintf("%s Ditolak! Lokasi Anda berjarak %.2fm dari Posko 112 (Batas Radius <= %.1fm).", label, dist, maxRad),
			"distance_meters":   dist,
			"max_radius_meters": maxRad,
		})
		return
	}
	ok = true
	return
}

func (h *AttendanceHTTPHandler) presensiUserName(r *http.Request, userID int) string {
	userRes, _ := h.userSvc.GetProfile(r.Context(), &userProto.GetProfileRequest{UserId: int32(userID)})
	if userRes != nil && userRes.Success && userRes.User != nil && userRes.User.Name != "" {
		return userRes.User.Name
	}
	return "Peserta Magang"
}

// respondPresensiResult maps the attendance service result to HTTP and logs successful scans.
func respondPresensiResult(w http.ResponseWriter, r *http.Request, res *attendanceProto.AttendanceResponse, err error, userID int, nip, name string, dist, maxRad float64) {
	if err != nil || res == nil {
		errMsg := "Gagal menyimpan presensi ke database."
		if err != nil {
			errMsg = err.Error()
		}
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": errMsg})
		return
	}
	if !res.Success {
		middleware.RespondJSON(w, http.StatusConflict, map[string]interface{}{"success": false, "error": res.Message})
		return
	}

	recordType := "MASUK"
	if res.Record != nil && res.Record.Type != "" {
		recordType = res.Record.Type
	}
	client.RecordActivityLog(userID, nip, name, "PRESENSI_"+recordType+"_SCAN", fmt.Sprintf("Presensi %s via Scan QR Geofence Posko (Jarak: %.2fm)", recordType, dist), client.GetClientIP(r), r.UserAgent())

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"type":    recordType,
		"message": fmt.Sprintf("%s (Batas <= %.1fm).", strings.TrimSuffix(res.Message, "."), maxRad),
		"record":  res.Record,
	})
}

func (h *AttendanceHTTPHandler) HandleClockIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	nip, _ := r.Context().Value("nip").(string)

	body, dist, maxRad, ok := h.validatePresensiScan(w, r, "Presensi")
	if !ok {
		return
	}
	name := h.presensiUserName(r, userID)

	res, err := h.attendanceSvc.ClockIn(r.Context(), &attendanceProto.ClockInRequest{
		UserId:    int32(userID),
		UserNip:   nip,
		UserName:  name,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
		QrToken:   body.QRToken,
	})
	respondPresensiResult(w, r, res, err, userID, nip, name, dist, maxRad)
}

func (h *AttendanceHTTPHandler) HandleClockOut(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	nip, _ := r.Context().Value("nip").(string)

	body, dist, maxRad, ok := h.validatePresensiScan(w, r, "Presensi PULANG")
	if !ok {
		return
	}
	name := h.presensiUserName(r, userID)

	res, err := h.attendanceSvc.ClockOut(r.Context(), &attendanceProto.ClockOutRequest{
		UserId:    int32(userID),
		UserNip:   nip,
		UserName:  name,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	})
	respondPresensiResult(w, r, res, err, userID, nip, name, dist, maxRad)
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
		logAdminAction(r, "UPDATE_SCHEDULE", "Super Admin menyimpan perubahan pembagian roster shift bulan ini")
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Jadwal shift berhasil disimpan."})
		return
	}
	middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
}

// leaveRequestList normalises the stored JSONB document into a list of request objects.
func leaveRequestList(raw interface{}) []map[string]interface{} {
	list := []map[string]interface{}{}
	items, _ := raw.([]interface{})
	for _, it := range items {
		if m, ok := it.(map[string]interface{}); ok {
			list = append(list, m)
		}
	}
	return list
}

// leaveRequestOwnedBy matches on user_id (set by the server since this fix) and falls back to NIP for older entries.
func leaveRequestOwnedBy(item map[string]interface{}, userID int, nip string) bool {
	if id, ok := item["user_id"].(float64); ok && int(id) == userID && userID > 0 {
		return true
	}
	cleanNIP := strings.ReplaceAll(nip, " ", "")
	itemNIP, _ := item["user_nip"].(string)
	return cleanNIP != "" && strings.ReplaceAll(itemNIP, " ", "") == cleanNIP
}

func (h *AttendanceHTTPHandler) HandleLeaveRequests(w http.ResponseWriter, r *http.Request) {
	stub := &client.AttendanceClientDirectStub{}
	role := requestRole(r)
	isAdmin := role == "admin" || role == "superadmin"
	userID, _ := r.Context().Value("user_id").(int)
	nip, _ := r.Context().Value("nip").(string)

	if r.Method == http.MethodGet {
		reqs := stub.GetLeaveRequests()
		if !isAdmin {
			// Interns only see their own submissions
			own := []map[string]interface{}{}
			for _, item := range leaveRequestList(reqs) {
				if leaveRequestOwnedBy(item, userID, nip) {
					own = append(own, item)
				}
			}
			reqs = own
		}
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
		if _, ok := raw.([]interface{}); !ok {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Format daftar pengajuan tidak valid."})
			return
		}

		if !isAdmin {
			// An intern may only add new requests of their own. Everything already stored (other interns'
			// requests and admin decisions) is kept exactly as is.
			existing := leaveRequestList(stub.GetLeaveRequests())
			knownIDs := map[string]bool{}
			for _, item := range existing {
				knownIDs[fmt.Sprint(item["id"])] = true
			}
			name := h.presensiUserName(r, userID)
			added := []interface{}{}
			for i, item := range leaveRequestList(raw) {
				if knownIDs[fmt.Sprint(item["id"])] {
					continue
				}
				category, _ := item["category"].(string)
				if category != "IZIN" && category != "SAKIT" {
					continue
				}
				shiftDate, _ := item["shift_date"].(string)
				reason, _ := item["reason"].(string)
				if strings.TrimSpace(shiftDate) == "" || strings.TrimSpace(reason) == "" {
					continue
				}
				now := time.Now()
				displayNIP := nip
				if displayNIP == "" {
					displayNIP = "-"
				}
				added = append(added, map[string]interface{}{
					"id":         now.UnixMilli() + int64(i),
					"created_at": now.In(witaLocation()).Format("2006-01-02 15:04"),
					"user_id":    userID,
					"user_name":  name,
					"user_nip":   displayNIP,
					"category":   category,
					"shift_date": shiftDate,
					"reason":     strings.TrimSpace(reason),
					"status":     "PENDING",
				})
			}
			if len(added) == 0 {
				middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Tidak ada pengajuan baru yang valid untuk disimpan."})
				return
			}
			merged := added
			for _, item := range existing {
				merged = append(merged, item)
			}
			raw = merged
		}

		if err := stub.SaveLeaveRequests(raw); err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		if isAdmin {
			logAdminAction(r, "UPDATE_LEAVE_REQUESTS", "Admin memperbarui status pengajuan izin / sakit peserta magang")
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Pengajuan berhasil disimpan."})
		return
	}
	middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
}

func witaLocation() *time.Location {
	if loc, err := time.LoadLocation("Asia/Makassar"); err == nil {
		return loc
	}
	return time.FixedZone("WITA", 8*3600)
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
	logAdminAction(r, "DELETE_ATTENDANCE_RECORD", fmt.Sprintf("Admin menghapus record presensi ID #%d dari database", id))
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
	logAdminAction(r, "MANUAL_ATTENDANCE_ENTRY", fmt.Sprintf("Admin menyimpan perubahan presensi manual petugas %s (NIP: %s) tanggal %s [Status: %s, Masuk: %s, Pulang: %s]", req.Name, req.NIP, req.Date, req.Status, req.TimeMasuk, req.TimePulang))
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Presensi %s pada tanggal %s berhasil di-update manual oleh Admin.", req.Name, req.Date),
	})
}

