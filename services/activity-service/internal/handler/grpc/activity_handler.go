package grpc

import (
	"context"
	"fmt"
	"time"

	actProto "proto/activity"
	"activity-service/internal/model"
	"activity-service/internal/repository"
)

type ActivityHandler struct {
	actProto.UnimplementedActivityServiceServer
	repo repository.ActivityRepository
}

func NewActivityHandler(repo repository.ActivityRepository) *ActivityHandler {
	return &ActivityHandler{
		repo: repo,
	}
}

func (h *ActivityHandler) CreateActivity(ctx context.Context, req *actProto.CreateActivityRequest) (*actProto.CreateActivityResponse, error) {
	if req.Title == "" || req.Description == "" {
		return &actProto.CreateActivityResponse{
			Success: false,
			Error:   "Judul dan Deskripsi kegiatan harian wajib diisi",
		}, nil
	}

	actDate := req.ActivityDate
	if actDate == "" {
		actDate = time.Now().Format("2006-01-02")
	}

	m := &model.DailyActivity{
		UserID:       req.UserId,
		UserNIP:      req.UserNip,
		UserName:     req.UserName,
		Title:        req.Title,
		Description:  req.Description,
		ActivityDate: actDate,
		PhotoURL:     req.PhotoUrl,
		Status:       "PENDING",
	}

	created, err := h.repo.Create(m)
	if err != nil {
		return &actProto.CreateActivityResponse{
			Success: false,
			Error:   fmt.Sprintf("Gagal menyimpan kegiatan harian: %v", err),
		}, nil
	}

	return &actProto.CreateActivityResponse{
		Success: true,
		Message: "Kegiatan harian berhasil dikirim untuk peninjauan admin",
		Record:  toProtoRecord(created),
	}, nil
}

func (h *ActivityHandler) GetActivities(ctx context.Context, req *actProto.GetActivitiesRequest) (*actProto.GetActivitiesResponse, error) {
	items, err := h.repo.GetByFilter(req.UserId, req.UserNip, req.Status, req.Limit)
	if err != nil {
		return &actProto.GetActivitiesResponse{
			Success: false,
			Error:   fmt.Sprintf("Gagal mengambil data kegiatan: %v", err),
		}, nil
	}

	var protoItems []*actProto.DailyActivityRecord
	for _, item := range items {
		protoItems = append(protoItems, toProtoRecord(item))
	}

	return &actProto.GetActivitiesResponse{
		Success:    true,
		Activities: protoItems,
	}, nil
}

func (h *ActivityHandler) GetActivityByID(ctx context.Context, req *actProto.GetActivityByIDRequest) (*actProto.GetActivityByIDResponse, error) {
	item, err := h.repo.GetByID(req.Id)
	if err != nil {
		return &actProto.GetActivityByIDResponse{
			Success: false,
			Error:   fmt.Sprintf("Kegiatan dengan ID %d tidak ditemukan", req.Id),
		}, nil
	}

	return &actProto.GetActivityByIDResponse{
		Success: true,
		Record:  toProtoRecord(item),
	}, nil
}

func (h *ActivityHandler) UpdateActivityStatus(ctx context.Context, req *actProto.UpdateActivityStatusRequest) (*actProto.UpdateActivityStatusResponse, error) {
	if req.Status != "APPROVED" && req.Status != "REJECTED" {
		return &actProto.UpdateActivityStatusResponse{
			Success: false,
			Error:   "Status harus APPROVED atau REJECTED",
		}, nil
	}

	err := h.repo.UpdateStatus(req.Id, req.Status)
	if err != nil {
		return &actProto.UpdateActivityStatusResponse{
			Success: false,
			Error:   fmt.Sprintf("Gagal memperbarui status kegiatan: %v", err),
		}, nil
	}

	return &actProto.UpdateActivityStatusResponse{
		Success: true,
		Message: fmt.Sprintf("Status kegiatan #%d berhasil diubah menjadi %s", req.Id, req.Status),
	}, nil
}

func (h *ActivityHandler) DeleteActivity(ctx context.Context, req *actProto.DeleteActivityRequest) (*actProto.DeleteActivityResponse, error) {
	err := h.repo.Delete(req.Id, req.UserId)
	if err != nil {
		return &actProto.DeleteActivityResponse{
			Success: false,
			Error:   fmt.Sprintf("Gagal menghapus kegiatan: %v", err),
		}, nil
	}

	return &actProto.DeleteActivityResponse{
		Success: true,
		Message: fmt.Sprintf("Kegiatan #%d berhasil dihapus", req.Id),
	}, nil
}

func toProtoRecord(m *model.DailyActivity) *actProto.DailyActivityRecord {
	if m == nil {
		return nil
	}
	return &actProto.DailyActivityRecord{
		Id:           m.ID,
		UserId:       m.UserID,
		UserNip:      m.UserNIP,
		UserName:     m.UserName,
		Title:        m.Title,
		Description:  m.Description,
		ActivityDate: m.ActivityDate,
		PhotoUrl:     m.PhotoURL,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}
