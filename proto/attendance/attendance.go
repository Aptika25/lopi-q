package attendance

import "context"

type GetPoskoQRRequest struct{}

type GetPoskoQRResponse struct {
	Success         bool    `json:"success"`
	PoskoName       string  `json:"posko_name"`
	Address         string  `json:"address"`
	QrToken         string  `json:"qr_token"`
	QrImage         string  `json:"qr_image"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	MaxRadiusMeters float64 `json:"max_radius_meters"`
	Error           string  `json:"error"`
}

type ClockInRequest struct {
	UserId    int32   `json:"user_id"`
	UserNip   string  `json:"user_nip"`
	UserName  string  `json:"user_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	QrToken   string  `json:"qr_token"`
}

type ClockOutRequest struct {
	UserId    int32   `json:"user_id"`
	UserNip   string  `json:"user_nip"`
	UserName  string  `json:"user_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	QrToken   string  `json:"qr_token"`
}

type AttendanceRecord struct {
	Id             int32   `json:"id"`
	UserId         int32   `json:"user_id"`
	UserNip        string  `json:"user_nip"`
	UserName       string  `json:"user_name"`
	Type           string  `json:"type"`
	Timestamp      string  `json:"timestamp"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	DistanceMeters float64 `json:"distance_meters"`
	WithinRadius   bool    `json:"within_radius"`
}

type AttendanceResponse struct {
	Success        bool              `json:"success"`
	Message        string            `json:"message"`
	DistanceMeters float64           `json:"distance_meters"`
	Record         *AttendanceRecord `json:"record"`
	Error          string            `json:"error"`
}

type GetHistoryRequest struct {
	Limit   int32  `json:"limit"`
	UserId  int32  `json:"user_id"`
	UserNip string `json:"user_nip"`
}

type GetHistoryResponse struct {
	Success bool                `json:"success"`
	History []*AttendanceRecord `json:"history"`
}

type GetTodayStatusRequest struct {
	UserId  int32  `json:"user_id"`
	UserNip string `json:"user_nip"`
}

type GetTodayStatusResponse struct {
	Success  bool              `json:"success"`
	Masuk    *AttendanceRecord `json:"masuk"`
	Pulang   *AttendanceRecord `json:"pulang"`
	IsMasuk  bool              `json:"is_masuk"`
	IsPulang bool              `json:"is_pulang"`
	Error    string            `json:"error"`
}

type UnimplementedAttendanceServiceServer struct{}

type AttendanceServiceServer interface {
	GetPoskoQR(context.Context, *GetPoskoQRRequest) (*GetPoskoQRResponse, error)
	ClockIn(context.Context, *ClockInRequest) (*AttendanceResponse, error)
	ClockOut(context.Context, *ClockOutRequest) (*AttendanceResponse, error)
	GetHistory(context.Context, *GetHistoryRequest) (*GetHistoryResponse, error)
	GetTodayStatus(context.Context, *GetTodayStatusRequest) (*GetTodayStatusResponse, error)
}

func RegisterAttendanceServiceServer(s interface{}, srv AttendanceServiceServer) {
	_ = s
	_ = srv
}
