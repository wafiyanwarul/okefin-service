package dto

type UpdateUserRequest struct {
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type UserResponse struct {
	ID           uint              `json:"id"`
	Nama         string            `json:"nama"`
	NoTelp       string            `json:"no_telp"`
	TanggalLahir string            `json:"tanggal_lahir"`
	JenisKelamin string            `json:"jenis_kelamin"`
	Tentang      string            `json:"tentang"`
	Pekerjaan    string            `json:"pekerjaan"`
	Email        string            `json:"email"`
	IDProvinsi   *ProvinsiResponse `json:"id_provinsi"`
	IDKota       *KotaResponse     `json:"id_kota"`
	IsAdmin      bool              `json:"is_admin"`
}

type UploadResponse struct {
	URL string `json:"url"`
}