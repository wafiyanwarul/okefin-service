package dto

type CreateTransaksiRequest struct {
	AlamatID   uint    `json:"alamat_id" validate:"required"`
	TotalHarga float64 `json:"total_harga" validate:"required,gt=0"`
}

type UpdateTransaksiRequest struct {
	Status string `json:"status" validate:"required,oneof=pending completed cancelled"`
}

type DetailTransaksiRequest struct {
	ProdukID uint    `json:"produk_id" validate:"required"`
	Jumlah   int     `json:"jumlah" validate:"required,gt=0"`
	Harga    float64 `json:"harga" validate:"required,gt=0"`
}

type TransaksiResponse struct {
	ID         uint                  `json:"id"`
	AlamatID   uint                  `json:"alamat_id"`
	UserID     uint                  `json:"user_id"`
	TotalHarga float64               `json:"total_harga"`
	Status     string                `json:"status"`
	CreatedAt  string                `json:"created_at"`
	UpdatedAt  string                `json:"updated_at"`
	Details    []DetailTransaksiResp `json:"details"`
}

type DetailTransaksiResp struct {
	ID       uint    `json:"id"`
	ProdukID uint    `json:"produk_id"`
	Jumlah   int     `json:"jumlah"`
	Harga    float64 `json:"harga"`
}
