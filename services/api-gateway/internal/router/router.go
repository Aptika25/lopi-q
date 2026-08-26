package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	httpHandler "api-gateway/internal/handler/http"
	"api-gateway/internal/middleware"
	"api-gateway/internal/utils"
)

func NewRouter(authHandler *httpHandler.AuthHTTPHandler, userHandler *httpHandler.UserHTTPHandler, attendanceHandler *httpHandler.AttendanceHTTPHandler, activityHandler *httpHandler.ActivityHTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()
	cors := middleware.CORSHandler
	authMw := middleware.AuthenticateMiddleware

	// AUTH ENDPOINTS
	mux.HandleFunc("/api/auth/login", cors(authHandler.HandleLogin))
	mux.HandleFunc("/api/auth/verify-2fa", cors(authHandler.HandleVerify2FA))
	mux.HandleFunc("/api/auth/2fa/setup", cors(authHandler.HandleSetup2FA))
	mux.HandleFunc("/api/auth/2fa/enable", cors(authHandler.HandleEnable2FA))
	mux.HandleFunc("/api/auth/2fa/disable", cors(authMw(authHandler.HandleDisable2FA)))
	mux.HandleFunc("/api/auth/2fa/self-reset", cors(authHandler.HandleSelfReset2FA))

	// USER & ADMIN ENDPOINTS
	mux.HandleFunc("/api/users/me", cors(authMw(userHandler.HandleGetProfile)))
	mux.HandleFunc("/api/admin/users", cors(authMw(userHandler.HandleAdminUsers)))
	mux.HandleFunc("/api/admin/users/", cors(authMw(userHandler.HandleAdminUserDetail)))
	mux.HandleFunc("/api/admin/location", cors(authMw(attendanceHandler.HandleAdminLocation)))
	mux.HandleFunc("/api/admin/schedules", cors(authMw(attendanceHandler.HandleAdminSchedules)))
	mux.HandleFunc("/api/admin/leave-requests", cors(authMw(attendanceHandler.HandleLeaveRequests)))
	mux.HandleFunc("/api/admin/presensi/cleanup-duplicates", cors(authMw(attendanceHandler.HandleCleanupDuplicates)))
	mux.HandleFunc("/api/admin/presensi/manual-entry", cors(authMw(attendanceHandler.HandleAdminManualEntry)))
	mux.HandleFunc("/api/admin/presensi/records/", cors(authMw(attendanceHandler.HandleDeleteRecord)))
	mux.HandleFunc("/api/admin/activity-logs", cors(authMw(userHandler.HandleAdminActivityLogs)))
	mux.HandleFunc("/api/admin/activities", cors(authMw(activityHandler.HandleActivities)))

	// DAILY ACTIVITY LOGBOOK ENDPOINTS
	mux.HandleFunc("/api/activities", cors(authMw(activityHandler.HandleActivities)))

	// ATTENDANCE ENDPOINTS
	mux.HandleFunc("/api/presensi/posko-qr", cors(attendanceHandler.HandleGetPoskoQR))
	mux.HandleFunc("/api/presensi/leave-requests", cors(authMw(attendanceHandler.HandleLeaveRequests)))
	mux.HandleFunc("/api/presensi/today", cors(authMw(attendanceHandler.HandleGetTodayStatus)))
	mux.HandleFunc("/api/presensi/history", cors(authMw(attendanceHandler.HandleGetHistory)))
	mux.HandleFunc("/api/presensi/clock-in", cors(authMw(attendanceHandler.HandleClockIn)))
	mux.HandleFunc("/api/presensi/clock-out", cors(authMw(attendanceHandler.HandleClockOut)))

	// FRONTEND SPA STATIC SERVING
	webDist := filepath.Join("..", "..", "apps", "web", "dist")
	if _, err := os.Stat(webDist); err != nil {
		webDist = filepath.Join("frontend", "dist")
	}

	fileServer := http.FileServer(http.Dir(webDist))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Endpoint API tidak ditemukan."})
			return
		}
		path := filepath.Join(webDist, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(webDist, "index.html"))
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	return mux
}
