package client

import (
	"context"

	repProto "proto/reporting"
)

type ReportingClientDirectStub struct{}

var memReports = []*repProto.AttendanceRecapRecord{
	{Id: 1, UserId: 101, UserName: "Budi Santoso", UserNip: "IN-301", Department: "Engineering", Date: "24 Okt 2023", ClockIn: "08:50 AM", ClockOut: "05:15 PM", TotalHours: "8h 25m", Status: "HADIR"},
	{Id: 2, UserId: 102, UserName: "Siti Aminah", UserNip: "IN-302", Department: "Design", Date: "24 Okt 2023", ClockIn: "09:15 AM", ClockOut: "06:00 PM", TotalHours: "8h 45m", Status: "TERLAMBAT"},
	{Id: 3, UserId: 103, UserName: "Andi Wijaya", UserNip: "IN-303", Department: "Marketing", Date: "24 Okt 2023", ClockIn: "--:--", ClockOut: "--:--", TotalHours: "0h 0m", Status: "ABSEN"},
	{Id: 4, UserId: 104, UserName: "Rina Permata", UserNip: "IN-304", Department: "Engineering", Date: "24 Okt 2023", ClockIn: "--:--", ClockOut: "--:--", TotalHours: "0h 0m", Status: "CUTI"},
}

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
			TotalHadir:     4821,
			PctTepatWaktu:  92,
			TotalTerlambat: 312,
			TotalAbsen:     84,
		},
		Records: filtered,
	}, nil
}
