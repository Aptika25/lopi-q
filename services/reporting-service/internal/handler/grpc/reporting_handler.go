package grpc

import (
	"context"

	repProto "proto/reporting"
	"reporting-service/internal/model"
	"reporting-service/internal/repository"
)

type ReportingHandler struct {
	repProto.UnimplementedReportingServiceServer
	repo repository.ReportingRepository
}

func NewReportingHandler(repo repository.ReportingRepository) *ReportingHandler {
	return &ReportingHandler{
		repo: repo,
	}
}

func (h *ReportingHandler) GetAttendanceRecap(ctx context.Context, req *repProto.GetAttendanceRecapRequest) (*repProto.GetAttendanceRecapResponse, error) {
	stats, items, err := h.repo.GetAttendanceRecap(req.Month, req.Department, req.Search, req.Limit)
	if err != nil {
		return &repProto.GetAttendanceRecapResponse{
			Success: false,
			Error:   "Gagal mengambil rekapan kehadiran",
		}, nil
	}

	var protoItems []*repProto.AttendanceRecapRecord
	for _, item := range items {
		protoItems = append(protoItems, toProtoRecord(item))
	}

	return &repProto.GetAttendanceRecapResponse{
		Success: true,
		Stats:   stats,
		Records: protoItems,
	}, nil
}

func toProtoRecord(m *model.AttendanceReport) *repProto.AttendanceRecapRecord {
	if m == nil {
		return nil
	}
	return &repProto.AttendanceRecapRecord{
		Id:         m.ID,
		UserId:     m.UserID,
		UserName:   m.UserName,
		UserNip:    m.UserNIP,
		Department: m.Department,
		Date:       m.Date,
		ClockIn:    m.ClockIn,
		ClockOut:   m.ClockOut,
		TotalHours: m.TotalHours,
		Status:     m.Status,
	}
}
