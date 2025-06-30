package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"
)

type TransaksiRepository interface {
	CreateTransaksi(transaksi *models.Trx) error
	GetTransaksiByID(id uint, userID uint) (*models.Trx, error)
	GetAllTransaksiByUserID(userID uint, limit, offset int) ([]models.Trx, int64, error)
	UpdateTransaksi(transaksi *models.Trx) error
	CreateDetailTransaksi(detail *models.DetailTrx) error
	GetDetailsByTransaksiID(transaksiID uint) ([]models.DetailTrx, error)
}

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepository(db *gorm.DB) TransaksiRepository {
	return &transaksiRepository{db: db}
}

func (r *transaksiRepository) CreateTransaksi(transaksi *models.Trx) error {
	return r.db.Create(transaksi).Error
}

func (r *transaksiRepository) GetTransaksiByID(id uint, userID uint) (*models.Trx, error) {
	var transaksi models.Trx
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&transaksi).Error
	if err != nil {
		return nil, err
	}
	return &transaksi, nil
}

func (r *transaksiRepository) GetAllTransaksiByUserID(userID uint, limit, offset int) ([]models.Trx, int64, error) {
	var transaksi []models.Trx
	var total int64

	r.db.Model(&models.Trx{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&transaksi).Error
	return transaksi, total, err
}

func (r *transaksiRepository) UpdateTransaksi(transaksi *models.Trx) error {
	return r.db.Save(transaksi).Error
}

func (r *transaksiRepository) CreateDetailTransaksi(detail *models.DetailTrx) error {
	return r.db.Create(detail).Error
}

func (r *transaksiRepository) GetDetailsByTransaksiID(transaksiID uint) ([]models.DetailTrx, error) {
	var details []models.DetailTrx
	err := r.db.Where("trx_id = ?", transaksiID).Find(&details).Error
	return details, err
}
