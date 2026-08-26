package model

import "time"

type DailyActivity struct {
	ID           int32     `json:"id"`
	UserID       int32     `json:"user_id"`
	UserNIP      string    `json:"user_nip"`
	UserName     string    `json:"user_name"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ActivityDate string    `json:"activity_date"`
	PhotoURL     string    `json:"photo_url"`
	Status       string    `json:"status"` // PENDING, APPROVED, REJECTED
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
