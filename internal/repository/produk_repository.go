package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"

)

type ProdukRepository interface {
	CreateProduk(produk *models.Produk) error
	GetAllProdukByTokoID(id uint, limit, offset int) ([]models.Produk, int64, error)
	GetProdukByID(id uint, userID uint) (*models.Produk, error)
	UpdateProduk(produk *models.Produk) error
	DeleteProduk(id uint) error
	CheckProdukOwnership(id uint, userID uint) bool
	CreateLogProduk(log *models.LogProduk) error
	GetLogProdukByID(id uint) (*models.LogProduk, error) // Tambahkan ini
}

type produkRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db: db}
}

func (r *produkRepository) CreateProduk(produk *models.Produk) error {
	return r.db.Create(produk).Error
}

func (r *produkRepository) GetProdukByID(id uint, userID uint) (*models.Produk, error) {
	var produk models.Produk
	err := r.db.Where("id = ? AND toko_id IN (SELECT id FROM toko WHERE id_user = ?)", id, userID).First(&produk).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *produkRepository) GetAllProdukByTokoID(tokoID uint, limit, offset int) ([]models.Produk, int64, error) {
	var produk []models.Produk
	var total int64

	r.db.Model(&models.Produk{}).Where("toko_id = ?", tokoID).Count(&total)
	err := r.db.Where("toko_id = ?", tokoID).Limit(limit).Offset(offset).Find(&produk).Error
	return produk, total, err
}

func (r *produkRepository) UpdateProduk(produk *models.Produk) error {
	return r.db.Save(produk).Error
}

func (r *produkRepository) DeleteProduk(id uint) error {
	return r.db.Delete(&models.Produk{}, id).Error
}

func (r *produkRepository) CheckProdukOwnership(id uint, userID uint) bool {
	var count int64
	r.db.Model(&models.Produk{}).Joins("JOIN toko ON toko.id = produk.toko_id").
		Where("produk.id = ? AND toko.id_user = ?", id, userID).Count(&count)
	return count > 0
}

func (r *produkRepository) CreateLogProduk(log *models.LogProduk) error {
	return r.db.Create(log).Error
}

func (r *produkRepository) GetLogProdukByID(id uint) (*models.LogProduk, error) {
	var log models.LogProduk
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
