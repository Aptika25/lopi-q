package model

import "time"

type Attendance struct {
	ID             int32     `json:"id"`
	UserID         int32     `json:"user_id"`
	UserNIP        string    `json:"user_nip"`
	UserName       string    `json:"user_name"`
	Type           string    `json:"type"` // 'MASUK' or 'PULANG'
	Timestamp      time.Time `json:"timestamp"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	DistanceMeters float64   `json:"distance_meters"`
	WithinRadius   bool      `json:"within_radius"`
	CreatedAt      time.Time `json:"created_at"`
}
