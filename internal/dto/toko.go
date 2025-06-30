package dto

type CreateTokoRequest struct {
	NamaToko string `json:"nama_toko" validate:"required"`
	UrlFoto  string `json:"url_foto"`
}

type UpdateTokoRequest struct {
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
}

type TokoWithUserResponse struct {
	ID       uint         `json:"id"`
	NamaToko string       `json:"nama_toko"`
	UrlFoto  string       `json:"url_foto"`
	User     UserResponse `json:"user"`
}
