package dto

type LoginRequestDTO struct {
	Identifier string `json:"identifier"`
	NIP        string `json:"nip"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginResponseDTO struct {
	Success     bool        `json:"success"`
	OTPRequired bool        `json:"otp_required,omitempty"`
	TempToken   string      `json:"temp_token,omitempty"`
	Token       string      `json:"token,omitempty"`
	User        interface{} `json:"user,omitempty"`
	Error       string      `json:"error,omitempty"`
}

type Verify2FARequestDTO struct {
	TempToken string `json:"temp_token"`
	Code      string `json:"code"`
}

type Enable2FARequestDTO struct {
	Code string `json:"code"`
}
