package grpc

import (
	"context"
	"fmt"
	"time"

	attProto "proto/attendance"
	"attendance-service/internal/model"
	"attendance-service/internal/repository"
)

type AttendanceHandler struct {
	attProto.UnimplementedAttendanceServiceServer
	repo repository.AttendanceRepository
}

func NewAttendanceHandler(repo repository.AttendanceRepository) *AttendanceHandler {
	return &AttendanceHandler{repo: repo}
}

func (h *AttendanceHandler) GetPoskoQR(ctx context.Context, req *attProto.GetPoskoQRRequest) (*attProto.GetPoskoQRResponse, error) {
	name, address, lat, lng, radius := h.repo.GetPoskoInfo()
	qrToken := fmt.Sprintf("POSKO-NTPD112-%d", time.Now().Unix())

	return &attProto.GetPoskoQRResponse{
		Success:         true,
		PoskoName:       name,
		Address:         address,
		QrToken:         qrToken,
		QrImage:         "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='100' height='100'><rect width='100' height='100' fill='pink'/></svg>",
		Latitude:        lat,
		Longitude:       lng,
		MaxRadiusMeters: radius,
	}, nil
}

func (h *AttendanceHandler) ClockIn(ctx context.Context, req *attProto.ClockInRequest) (*attProto.AttendanceResponse, error) {
	if req.UserId <= 0 {
		return &attProto.AttendanceResponse{
			Success: false,
			Error:   "User ID tidak valid",
		}, nil
	}

	record := &model.Attendance{
		UserID:    req.UserId,
		UserNIP:   req.UserNip,
		UserName:  req.UserName,
		Type:      "MASUK",
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	saved, err := h.repo.CreateRecord(ctx, record)
	if err != nil {
		return &attProto.AttendanceResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	msg := fmt.Sprintf("Presensi MASUK berhasil dicatat (Jarak: %.1fm)", saved.DistanceMeters)
	if !saved.WithinRadius {
		msg += " - Catatan: Di luar radius Posko (100m)"
	}

	return &attProto.AttendanceResponse{
		Success:        true,
		Message:        msg,
		DistanceMeters: saved.DistanceMeters,
		Record: &attProto.AttendanceRecord{
			Id:             saved.ID,
			UserId:         saved.UserID,
			UserNip:        saved.UserNIP,
			UserName:       saved.UserName,
			Type:           saved.Type,
			Timestamp:      saved.Timestamp.Format(time.RFC3339),
			Latitude:       saved.Latitude,
			Longitude:      saved.Longitude,
			DistanceMeters: saved.DistanceMeters,
			WithinRadius:   saved.WithinRadius,
		},
	}, nil
}

func (h *AttendanceHandler) ClockOut(ctx context.Context, req *attProto.ClockOutRequest) (*attProto.AttendanceResponse, error) {
	if req.UserId <= 0 {
		return &attProto.AttendanceResponse{
			Success: false,
			Error:   "User ID tidak valid",
		}, nil
	}

	record := &model.Attendance{
		UserID:    req.UserId,
		UserNIP:   req.UserNip,
		UserName:  req.UserName,
		Type:      "PULANG",
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	saved, err := h.repo.CreateRecord(ctx, record)
	if err != nil {
		return &attProto.AttendanceResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	msg := fmt.Sprintf("Presensi PULANG berhasil dicatat (Jarak: %.1fm)", saved.DistanceMeters)
	if !saved.WithinRadius {
		msg += " - Catatan: Di luar radius Posko (100m)"
	}

	return &attProto.AttendanceResponse{
		Success:        true,
		Message:        msg,
		DistanceMeters: saved.DistanceMeters,
		Record: &attProto.AttendanceRecord{
			Id:             saved.ID,
			UserId:         saved.UserID,
			UserNip:        saved.UserNIP,
			UserName:       saved.UserName,
			Type:           saved.Type,
			Timestamp:      saved.Timestamp.Format(time.RFC3339),
			Latitude:       saved.Latitude,
			Longitude:      saved.Longitude,
			DistanceMeters: saved.DistanceMeters,
			WithinRadius:   saved.WithinRadius,
		},
	}, nil
}

func (h *AttendanceHandler) GetHistory(ctx context.Context, req *attProto.GetHistoryRequest) (*attProto.GetHistoryResponse, error) {
	list, err := h.repo.GetHistory(ctx, req.UserId, req.Limit)
	if err != nil {
		return &attProto.GetHistoryResponse{Success: false}, nil
	}

	var protoList []*attProto.AttendanceRecord
	for _, item := range list {
		protoList = append(protoList, &attProto.AttendanceRecord{
			Id:             item.ID,
			UserId:         item.UserID,
			UserNip:        item.UserNIP,
			UserName:       item.UserName,
			Type:           item.Type,
			Timestamp:      item.Timestamp.Format(time.RFC3339),
			Latitude:       item.Latitude,
			Longitude:      item.Longitude,
			DistanceMeters: item.DistanceMeters,
			WithinRadius:   item.WithinRadius,
		})
	}

	return &attProto.GetHistoryResponse{
		Success: true,
		History: protoList,
	}, nil
}

func (h *AttendanceHandler) GetTodayStatus(ctx context.Context, req *attProto.GetTodayStatusRequest) (*attProto.GetTodayStatusResponse, error) {
	masuk, pulang, err := h.repo.GetTodayStatus(ctx, req.UserId)
	if err != nil {
		return &attProto.GetTodayStatusResponse{Success: false, Error: err.Error()}, nil
	}

	res := &attProto.GetTodayStatusResponse{
		Success:  true,
		IsMasuk:  masuk != nil,
		IsPulang: pulang != nil,
	}

	if masuk != nil {
		res.Masuk = &attProto.AttendanceRecord{
			Id:             masuk.ID,
			UserId:         masuk.UserID,
			UserNip:        masuk.UserNIP,
			UserName:       masuk.UserName,
			Type:           masuk.Type,
			Timestamp:      masuk.Timestamp.Format(time.RFC3339),
			Latitude:       masuk.Latitude,
			Longitude:      masuk.Longitude,
			DistanceMeters: masuk.DistanceMeters,
			WithinRadius:   masuk.WithinRadius,
		}
	}

	if pulang != nil {
		res.Pulang = &attProto.AttendanceRecord{
			Id:             pulang.ID,
			UserId:         pulang.UserID,
			UserNip:        pulang.UserNIP,
			UserName:       pulang.UserName,
			Type:           pulang.Type,
			Timestamp:      pulang.Timestamp.Format(time.RFC3339),
			Latitude:       pulang.Latitude,
			Longitude:      pulang.Longitude,
			DistanceMeters: pulang.DistanceMeters,
			WithinRadius:   pulang.WithinRadius,
		}
	}

	return res, nil
}
