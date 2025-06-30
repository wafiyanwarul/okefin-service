package dto

type CreateProdukRequest struct {
	NamaProduk    string    `json:"nama_produk" validate:"required"`
	Harga         float64   `json:"harga" validate:"required,gt=0"`
	Stok          int       `json:"stok" validate:"required,gt=0"`
	IDCategory    uint      `json:"id_category" validate:"required"`
	UrlFotos      []string  `json:"url_fotos" validate:"required,min=1"` // Multiple URLs
	Deskripsi     string    `json:"deskripsi" validate:"required"`
}

type UpdateProdukRequest struct {
	NamaProduk    string    `json:"nama_produk"`
	Harga         float64   `json:"harga"`
	Stok          int       `json:"stok"`
	IDCategory    uint      `json:"id_category"`
	UrlFotos      []string  `json:"url_fotos"` // Optional for updates
	Deskripsi     string    `json:"deskripsi"`
}

type ProdukResponse struct {
	ID           uint            `json:"id"`
	NamaProduk   string          `json:"nama_produk"`
	Harga        float64         `json:"harga"`
	Stok         int             `json:"stok"`
	IDCategory   uint            `json:"id_category"`
	UrlFotos     []string        `json:"url_fotos"`
	Deskripsi    string          `json:"deskripsi"`
	TokoID       uint            `json:"toko_id"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}