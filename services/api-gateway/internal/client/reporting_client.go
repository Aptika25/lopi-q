package client

import (
	"context"

	repProto "proto/reporting"
)

type ReportingClientDirectStub struct{}

var memReports = []*repProto.AttendanceRecapRecord{}

func (c *ReportingClientDirectStub) GetAttendanceRecap(ctx context.Context, req *repProto.GetAttendanceRecapRequest) (*repProto.GetAttendanceRecapResponse, error) {
	var filtered []*repProto.AttendanceRecapRecord
	for _, item := range memReports {
		if req.Department != "" && req.Department != "Semua Departemen" && item.Department != req.Department {
			continue
		}
		filtered = append(filtered, item)
	}

	return &repProto.GetAttendanceRecapResponse{
		Success: true,
		Stats: &repProto.ReportingStats{
			TotalHadir:     0,
			PctTepatWaktu:  0,
			TotalTerlambat: 0,
			TotalAbsen:     0,
		},
		Records: filtered,
	}, nil
}
