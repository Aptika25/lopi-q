package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"api-gateway/internal/client"
	"api-gateway/internal/middleware"
)

// JournalHTTPHandler serves the interns' daily activity journal and the admin review of each entry.
type JournalHTTPHandler struct{}

func NewJournalHTTPHandler() *JournalHTTPHandler {
	return &JournalHTTPHandler{}
}

func isAdminRole(role string) bool {
	return role == "admin" || role == "superadmin"
}

// HandleJournal: GET lists entries (interns see only their own), POST adds an entry with evidence (multipart).
func (h *JournalHTTPHandler) HandleJournal(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)
	nip, _ := r.Context().Value("nip").(string)
	admin := isAdminRole(requestRole(r))

	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		filter := client.JournalFilter{Date: strings.TrimSpace(q.Get("date")), Status: strings.TrimSpace(q.Get("status"))}
		if admin {
			filter.UserID, _ = strconv.Atoi(q.Get("user_id"))
		} else {
			filter.UserID = userID
		}
		entries, err := client.ListJournalEntries(r.Context(), filter)
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Gagal memuat jurnal kegiatan: " + err.Error()})
			return
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "entries": entries})

	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, client.MaxJournalEvidenceBytes+(1<<20))
		if err := r.ParseMultipartForm(client.MaxJournalEvidenceBytes + (1 << 20)); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Form tidak valid atau ukuran bukti dukung melebihi 3 MB."})
			return
		}
		entry := client.JournalEntry{
			UserID:       userID,
			UserNIP:      nip,
			ActivityDate: r.FormValue("activity_date"),
			ActivityTime: r.FormValue("activity_time"),
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
		}
		if file, header, err := r.FormFile("evidence"); err == nil {
			data, readErr := io.ReadAll(io.LimitReader(file, client.MaxJournalEvidenceBytes+1))
			file.Close()
			if readErr != nil {
				middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Gagal membaca file bukti dukung."})
				return
			}
			entry.EvidenceBytes = data
			entry.EvidenceName = filepath.Base(header.Filename)
			// Trust the content, not the client-declared type
			entry.EvidenceMime = strings.Split(http.DetectContentType(data), ";")[0]
		}
		if err := client.ValidateJournalEntry(&entry); err != nil {
			middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		entry.UserName = (&AttendanceHTTPHandler{userSvc: profileLookup}).presensiUserName(r, userID)

		saved, err := client.CreateJournalEntry(r.Context(), entry)
		if err != nil {
			middleware.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		client.RecordActivityLog(userID, nip, saved.UserName, "JOURNAL_ENTRY_CREATE", fmt.Sprintf("Peserta menambahkan kegiatan jurnal %s %s: %s", saved.ActivityDate, saved.ActivityTime, saved.Title), client.GetClientIP(r), r.UserAgent())
		middleware.RespondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "entry": saved, "message": "Kegiatan berhasil dikirim dan menunggu persetujuan pembimbing."})

	default:
		middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
	}
}

// HandleJournalItem: GET /api/journal/{id}/evidence streams the evidence file, DELETE /api/journal/{id} withdraws a pending entry.
func (h *JournalHTTPHandler) HandleJournalItem(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/journal/"), "/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "ID kegiatan tidak valid."})
		return
	}
	userID, _ := r.Context().Value("user_id").(int)

	switch {
	case r.Method == http.MethodGet && len(parts) == 2 && parts[1] == "evidence":
		entry, err := client.GetJournalEvidence(r.Context(), id)
		if err != nil {
			middleware.RespondJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		if entry.UserID != userID && !isAdminRole(requestRole(r)) {
			middleware.RespondJSON(w, http.StatusForbidden, map[string]interface{}{"success": false, "error": "Akses ditolak."})
			return
		}
		if len(entry.EvidenceBytes) == 0 {
			middleware.RespondJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Bukti dukung tidak tersedia."})
			return
		}
		w.Header().Set("Content-Type", entry.EvidenceMime)
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", entry.EvidenceName))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "private, max-age=300")
		_, _ = w.Write(entry.EvidenceBytes)

	case r.Method == http.MethodDelete && len(parts) == 1:
		if err := client.DeleteOwnPendingJournalEntry(r.Context(), id, userID); err != nil {
			middleware.RespondJSON(w, http.StatusConflict, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Kegiatan dihapus."})

	default:
		middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
	}
}

// HandleJournalReview: PUT /api/admin/journal/{id}/review {status, note} by an admin / pembimbing.
func (h *JournalHTTPHandler) HandleJournalReview(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/journal/"), "/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 || len(parts) != 2 || parts[1] != "review" {
		middleware.RespondJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Endpoint tidak ditemukan."})
		return
	}
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		middleware.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"success": false, "error": "Method not allowed"})
		return
	}
	var body struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Payload JSON tidak valid."})
		return
	}

	reviewerID, _ := r.Context().Value("user_id").(int)
	reviewer := (&AttendanceHTTPHandler{userSvc: profileLookup}).presensiUserName(r, reviewerID)
	saved, err := client.ReviewJournalEntry(r.Context(), id, strings.ToUpper(strings.TrimSpace(body.Status)), body.Note, reviewer)
	if err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	logAdminAction(r, "JOURNAL_REVIEW", fmt.Sprintf("Review kegiatan #%d (%s, %s %s) -> %s", saved.ID, saved.UserName, saved.ActivityDate, saved.ActivityTime, saved.Status))
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "entry": saved})
}
