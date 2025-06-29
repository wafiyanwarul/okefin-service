package dto

type RegisterRequest struct {
	Nama         string `json:"nama" validate:"required"`
	KataSandi    string `json:"kata_sandi" validate:"required,min=6"`
	NoTelp       string `json:"no_telp" validate:"required"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	Pekerjaan    string `json:"pekerjaan" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	IDProvinsi   string `json:"id_provinsi" validate:"required"`
	IDKota       string `json:"id_kota" validate:"required"`
}

type LoginRequest struct {
	NoTelp    string `json:"no_telp" validate:"required"`
	KataSandi string `json:"kata_sandi" validate:"required"`
}

type ProvinsiResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type KotaResponse struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

type LoginResponse struct {
	Nama         string            `json:"nama"`
	NoTelp       string            `json:"no_telp"`
	TanggalLahir string            `json:"tanggal_Lahir"`
	Tentang      string            `json:"tentang"`
	Pekerjaan    string            `json:"pekerjaan"`
	Email        string            `json:"email"`
	IDProvinsi   *ProvinsiResponse `json:"id_provinsi"`
	IDKota       *KotaResponse     `json:"id_kota"`
	Token        string            `json:"token"`
}

type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}
