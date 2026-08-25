package dto

type ClockInOutRequestDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	QRToken   string  `json:"qr_token"`
}
