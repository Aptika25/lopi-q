package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api-gateway/internal/client"
	"api-gateway/internal/middleware"
)

// InternshipHTTPHandler serves the admin "Proyek" and "Laporan Akhir" pages.
type InternshipHTTPHandler struct{}

func NewInternshipHTTPHandler() *InternshipHTTPHandler {
	return &InternshipHTTPHandler{}
}

// pathID returns the trailing numeric segment of /api/admin/<resource>/{id}, or 0 for the collection URL.
func pathID(r *http.Request, prefix string) (int, error) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")
	if rest == "" {
		return 0, nil
	}
	id, err := strconv.Atoi(rest)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("ID tidak valid")
	}
	return id, nil
}

func (h *InternshipHTTPHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "/api/admin/projects")
	if err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	switch {
	case r.Method == http.MethodGet && id == 0:
		projects, err := client.ListProjects(r.Context())
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal memuat data proyek: " + err.Error()})
			return
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "projects": projects})

	case (r.Method == http.MethodPost && id == 0) || (r.Method == http.MethodPut && id > 0):
		var p client.InternProject
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Payload JSON tidak valid."})
			return
		}
		p.ID = id
		if err := client.NormalizeProject(&p); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		savedID, err := client.SaveProject(r.Context(), p)
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal menyimpan proyek: " + err.Error()})
			return
		}
		logAdminAction(r, "SAVE_INTERN_PROJECT", fmt.Sprintf("Admin menyimpan proyek magang #%d: %s", savedID, p.Title))
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "id": savedID, "message": "Proyek berhasil disimpan."})

	case r.Method == http.MethodDelete && id > 0:
		if err := client.DeleteProject(r.Context(), id); err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal menghapus proyek: " + err.Error()})
			return
		}
		logAdminAction(r, "DELETE_INTERN_PROJECT", fmt.Sprintf("Admin menghapus proyek magang #%d", id))
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Proyek berhasil dihapus."})

	default:
		middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
	}
}

func (h *InternshipHTTPHandler) HandleFinalReports(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "/api/admin/final-reports")
	if err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	switch {
	case r.Method == http.MethodGet && id == 0:
		reports, err := client.ListFinalReports(r.Context())
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal memuat laporan akhir: " + err.Error()})
			return
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "reports": reports})

	case (r.Method == http.MethodPost && id == 0) || (r.Method == http.MethodPut && id > 0):
		var rep client.FinalReport
		if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Payload JSON tidak valid."})
			return
		}
		rep.ID = id
		if err := client.NormalizeFinalReport(&rep); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		savedID, err := client.SaveFinalReport(r.Context(), rep)
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal menyimpan laporan akhir: " + err.Error()})
			return
		}
		logAdminAction(r, "SAVE_FINAL_REPORT", fmt.Sprintf("Admin menyimpan laporan akhir #%d (%s) - status %s", savedID, rep.InternName, rep.Status))
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "id": savedID, "message": "Laporan akhir berhasil disimpan."})

	case r.Method == http.MethodDelete && id > 0:
		if err := client.DeleteFinalReport(r.Context(), id); err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal menghapus laporan akhir: " + err.Error()})
			return
		}
		logAdminAction(r, "DELETE_FINAL_REPORT", fmt.Sprintf("Admin menghapus laporan akhir #%d", id))
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Laporan akhir berhasil dihapus."})

	default:
		middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
	}
}
