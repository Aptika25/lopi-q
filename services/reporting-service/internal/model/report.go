package model

type AttendanceReport struct {
	ID         int32  `json:"id"`
	UserID     int32  `json:"user_id"`
	UserName   string `json:"user_name"`
	UserNIP    string `json:"user_nip"`
	Department string `json:"department"`
	Date       string `json:"date"`
	ClockIn    string `json:"clock_in"`
	ClockOut   string `json:"clock_out"`
	TotalHours string `json:"total_hours"`
	Status     string `json:"status"` // HADIR, TERLAMBAT, ABSEN, CUTI
}
