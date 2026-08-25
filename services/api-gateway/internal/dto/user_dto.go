package dto

type CreateUserDTO struct {
	NIP       string `json:"nip"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Jabatan   string `json:"jabatan"`
	UnitKerja string `json:"unit_kerja"`
	Password  string `json:"password"`
}

type UpdateUserDTO struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Jabatan   string `json:"jabatan"`
	UnitKerja string `json:"unit_kerja"`
	Password  string `json:"password"`
}

type ToggleActiveDTO struct {
	IsActive bool `json:"is_active"`
}
