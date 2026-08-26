package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-gateway/internal/client"
	actProto "proto/activity"
)

type ActivityHTTPHandler struct {
	activityClient *client.ActivityClientDirectStub
}

func NewActivityHTTPHandler(actClient *client.ActivityClientDirectStub) *ActivityHTTPHandler {
	return &ActivityHTTPHandler{
		activityClient: actClient,
	}
}

func (h *ActivityHTTPHandler) HandleActivities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		userIDStr := r.URL.Query().Get("user_id")
		userNIP := r.URL.Query().Get("user_nip")
		status := r.URL.Query().Get("status")

		var userID int32
		if userIDStr != "" {
			if id, err := strconv.Atoi(userIDStr); err == nil {
				userID = int32(id)
			}
		}

		res, err := h.activityClient.GetActivities(r.Context(), &actProto.GetActivitiesRequest{
			UserId:  userID,
			UserNip: userNIP,
			Status:  status,
		})

		if err != nil || !res.Success {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Gagal mengambil data kegiatan harian",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    true,
			"activities": res.Activities,
		})

	case http.MethodPost:
		var req actProto.CreateActivityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Payload JSON tidak valid",
			})
			return
		}

		res, err := h.activityClient.CreateActivity(r.Context(), &req)
		if err != nil || !res.Success {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   res.Error,
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)

	case http.MethodPut:
		var req actProto.UpdateActivityStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Payload JSON tidak valid",
			})
			return
		}

		res, err := h.activityClient.UpdateActivityStatus(r.Context(), &req)
		if err != nil || !res.Success {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   res.Error,
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
