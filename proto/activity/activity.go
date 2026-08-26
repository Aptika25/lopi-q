package activity

import "context"

type DailyActivityRecord struct {
	Id           int32  `json:"id"`
	UserId       int32  `json:"user_id"`
	UserNip      string `json:"user_nip"`
	UserName     string `json:"user_name"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ActivityDate string `json:"activity_date"`
	PhotoUrl     string `json:"photo_url"`
	Status       string `json:"status"` // PENDING, APPROVED, REJECTED
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateActivityRequest struct {
	UserId       int32  `json:"user_id"`
	UserNip      string `json:"user_nip"`
	UserName     string `json:"user_name"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ActivityDate string `json:"activity_date"`
	PhotoUrl     string `json:"photo_url"`
}

type CreateActivityResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Record  *DailyActivityRecord `json:"record"`
	Error   string               `json:"error"`
}

type GetActivitiesRequest struct {
	UserId  int32  `json:"user_id"`
	UserNip string `json:"user_nip"`
	Status  string `json:"status"`
	Limit   int32  `json:"limit"`
}

type GetActivitiesResponse struct {
	Success    bool                   `json:"success"`
	Activities []*DailyActivityRecord `json:"activities"`
	Error      string                 `json:"error"`
}

type GetActivityByIDRequest struct {
	Id int32 `json:"id"`
}

type GetActivityByIDResponse struct {
	Success bool                 `json:"success"`
	Record  *DailyActivityRecord `json:"record"`
	Error   string               `json:"error"`
}

type UpdateActivityStatusRequest struct {
	Id     int32  `json:"id"`
	Status string `json:"status"` // APPROVED, REJECTED
}

type UpdateActivityStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type DeleteActivityRequest struct {
	Id     int32 `json:"id"`
	UserId int32 `json:"user_id"`
}

type DeleteActivityResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type UnimplementedActivityServiceServer struct{}

type ActivityServiceServer interface {
	CreateActivity(context.Context, *CreateActivityRequest) (*CreateActivityResponse, error)
	GetActivities(context.Context, *GetActivitiesRequest) (*GetActivitiesResponse, error)
	GetActivityByID(context.Context, *GetActivityByIDRequest) (*GetActivityByIDResponse, error)
	UpdateActivityStatus(context.Context, *UpdateActivityStatusRequest) (*UpdateActivityStatusResponse, error)
	DeleteActivity(context.Context, *DeleteActivityRequest) (*DeleteActivityResponse, error)
}

func RegisterActivityServiceServer(s interface{}, srv ActivityServiceServer) {
	_ = s
	_ = srv
}
