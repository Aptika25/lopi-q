package reporting

import "context"

type AttendanceRecapRecord struct {
	Id         int32  `json:"id"`
	UserId     int32  `json:"user_id"`
	UserName   string `json:"user_name"`
	UserNip    string `json:"user_nip"`
	Department string `json:"department"`
	Date       string `json:"date"`
	ClockIn    string `json:"clock_in"`
	ClockOut   string `json:"clock_out"`
	TotalHours string `json:"total_hours"`
	Status     string `json:"status"` // HADIR, TERLAMBAT, ABSEN, CUTI
}

type ReportingStats struct {
	TotalHadir      int32 `json:"total_hadir"`
	PctTepatWaktu   int32 `json:"pct_tepat_waktu"`
	TotalTerlambat  int32 `json:"total_terlambat"`
	TotalAbsen      int32 `json:"total_absen"`
}

type GetAttendanceRecapRequest struct {
	Month      string `json:"month"`
	Department string `json:"department"`
	Search     string `json:"search"`
	Limit      int32  `json:"limit"`
}

type GetAttendanceRecapResponse struct {
	Success bool                     `json:"success"`
	Stats   *ReportingStats          `json:"stats"`
	Records []*AttendanceRecapRecord `json:"records"`
	Error   string                   `json:"error"`
}

type UnimplementedReportingServiceServer struct{}

type ReportingServiceServer interface {
	GetAttendanceRecap(context.Context, *GetAttendanceRecapRequest) (*GetAttendanceRecapResponse, error)
}

func RegisterReportingServiceServer(s interface{}, srv ReportingServiceServer) {
	_ = s
	_ = srv
}
