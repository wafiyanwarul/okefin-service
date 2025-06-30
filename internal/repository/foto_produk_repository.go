package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"

)

type FotoProdukRepository interface {
	CreateFotoProduk(foto *models.FotoProduk) error
	GetFotoProdukByProdukID(id uint) ([]models.FotoProduk, error)
	DeleteFotoProdukByProdukID(id uint) error
}

type fotoProdukRepository struct {
	db *gorm.DB
}

func NewFotoProdukRepository(db *gorm.DB) FotoProdukRepository {
	return &fotoProdukRepository{db: db}
}

func (r *fotoProdukRepository) CreateFotoProduk(foto *models.FotoProduk) error {
	return r.db.Create(foto).Error
}

func (r *fotoProdukRepository) GetFotoProdukByProdukID(id uint) ([]models.FotoProduk, error) {
	var fotos []models.FotoProduk
	err := r.db.Where("id_produk = ?", id).Find(&fotos).Error
	return fotos, err
}

func (r *fotoProdukRepository) DeleteFotoProdukByProdukID(id uint) error {
	return r.db.Where("id_produk = ?", id).Delete(&models.FotoProduk{}).Error
}