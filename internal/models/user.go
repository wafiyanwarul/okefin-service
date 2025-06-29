package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama         string    `json:"nama" gorm:"type:varchar(255);not null"`
	KataSandi    string    `json:"kata_sandi" gorm:"type:varchar(255);not null"`
	NoTelp       string    `json:"no_telp" gorm:"type:varchar(255);uniqueIndex;not null"`
	TanggalLahir string    `json:"tanggal_lahir" gorm:"type:date"`
	JenisKelamin string    `json:"jenis_kelamin" gorm:"type:varchar(255)"`
	Tentang      string    `json:"tentang" gorm:"type:text"`
	Pekerjaan    string    `json:"pekerjaan" gorm:"type:varchar(255)"`
	Email        string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	IDProvinsi   string    `json:"id_provinsi" gorm:"type:varchar(255)"`
	IDKota       string    `json:"id_kota" gorm:"type:varchar(255)"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	Alamat []Alamat `json:"alamat,omitempty" gorm:"foreignKey:IDUser"`
	Toko   *Toko    `json:"toko,omitempty" gorm:"foreignKey:IDUser"`
	Trx    []Trx    `json:"trx,omitempty" gorm:"foreignKey:IDUser"`
}

type Alamat struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser       uint      `json:"id_user" gorm:"not null"`
	JudulAlamat  string    `json:"judul_alamat" gorm:"type:varchar(255)"`
	NamaPenerima string    `json:"nama_penerima" gorm:"type:varchar(255)"`
	NoTelp       string    `json:"no_telp" gorm:"type:varchar(255)"`
	DetailAlamat string    `json:"detail_alamat" gorm:"type:varchar(255)"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:IDUser"`
}

type Toko struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser    uint      `json:"id_user" gorm:"not null"`
	NamaToko  string    `json:"nama_toko" gorm:"type:varchar(255)"`
	UrlFoto   string    `json:"url_foto" gorm:"type:varchar(255)"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User   User     `json:"user,omitempty" gorm:"foreignKey:IDUser"`
	Produk []Produk `json:"produk,omitempty" gorm:"foreignKey:IDToko"`
}

type Category struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaCategory string    `json:"nama_category" gorm:"type:varchar(255)"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Produk []Produk `json:"produk,omitempty" gorm:"foreignKey:IDCategory"`
}

type Produk struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaProduk    string    `json:"nama_produk" gorm:"type:varchar(255)"`
	Slug          string    `json:"slug" gorm:"type:varchar(255)"`
	HargaReseller string    `json:"harga_reseller" gorm:"type:varchar(255)"`
	HargaKonsumen string    `json:"harga_konsumen" gorm:"type:varchar(255)"`
	Stok          int       `json:"stok"`
	Deskripsi     string    `json:"deskripsi" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IDToko        uint      `json:"id_toko" gorm:"not null"`
	IDCategory    uint      `json:"id_category" gorm:"not null"`

	// Relations
	Toko       Toko         `json:"toko,omitempty" gorm:"foreignKey:IDToko"`
	Category   Category     `json:"category,omitempty" gorm:"foreignKey:IDCategory"`
	FotoProduk []FotoProduk `json:"foto_produk,omitempty" gorm:"foreignKey:IDProduk"`
	LogProduk  []LogProduk  `json:"log_produk,omitempty" gorm:"foreignKey:IDProduk"`
}

type FotoProduk struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDProduk  uint      `json:"id_produk" gorm:"not null"`
	Url       string    `json:"url" gorm:"type:varchar(255)"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Produk Produk `json:"produk,omitempty" gorm:"foreignKey:IDProduk"`
}

type Trx struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser           uint      `json:"id_user" gorm:"not null"`
	AlamatPengiriman uint      `json:"alamat_pengiriman" gorm:"not null"`
	HargaTotal       int       `json:"harga_total"`
	KodeInvoice      string    `json:"kode_invoice" gorm:"type:varchar(255)"`
	MethodBayar      string    `json:"method_bayar" gorm:"type:varchar(255)"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`

	// Relations
	User      User        `json:"user,omitempty" gorm:"foreignKey:IDUser"`
	Alamat    Alamat      `json:"alamat,omitempty" gorm:"foreignKey:AlamatPengiriman"`
	DetailTrx []DetailTrx `json:"detail_trx,omitempty" gorm:"foreignKey:IDTrx"`
}

type DetailTrx struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDTrx       uint      `json:"id_trx" gorm:"not null"`
	IDLogProduk uint      `json:"id_log_produk" gorm:"not null"`
	IDToko      uint      `json:"id_toko" gorm:"not null"`
	Kuantitas   int       `json:"kuantitas"`
	HargaTotal  int       `json:"harga_total"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Trx       Trx       `json:"trx,omitempty" gorm:"foreignKey:IDTrx"`
	LogProduk LogProduk `json:"log_produk,omitempty" gorm:"foreignKey:IDLogProduk"`
}

type LogProduk struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDProduk      uint      `json:"id_produk" gorm:"not null"`
	NamaProduk    string    `json:"nama_produk" gorm:"type:varchar(255)"`
	Slug          string    `json:"slug" gorm:"type:varchar(255)"`
	HargaReseller string    `json:"harga_reseller" gorm:"type:varchar(255)"`
	HargaKonsumen string    `json:"harga_konsumen" gorm:"type:varchar(255)"`
	Deskripsi     string    `json:"deskripsi" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IDToko        uint      `json:"id_toko" gorm:"not null"`
	IDCategory    uint      `json:"id_category" gorm:"not null"`

	// Relations
	Produk    Produk      `json:"produk,omitempty" gorm:"foreignKey:IDProduk"`
	DetailTrx []DetailTrx `json:"detail_trx,omitempty" gorm:"foreignKey:IDLogProduk"`
}
