package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"
)

type TokoRepository interface {
	CreateToko(toko *models.Toko) error
	GetTokoByID(id uint) (*models.Toko, error)
	GetTokoByUserID(userID uint) (*models.Toko, error)
	GetAllTokos(limit, offset int) ([]models.Toko, int64, error)
	UpdateToko(toko *models.Toko) error
	DeleteToko(id uint) error
	CheckTokoOwnership(tokoID, userID uint) bool
}

type tokoRepository struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepository{db: db}
}

func (r *tokoRepository) CreateToko(toko *models.Toko) error {
	return r.db.Create(toko).Error
}

func (r *tokoRepository) GetTokoByID(id uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.Preload("User").Where("id = ?", id).First(&toko).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

func (r *tokoRepository) GetTokoByUserID(userID uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.Where("id_user = ?", userID).First(&toko).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

func (r *tokoRepository) GetAllTokos(limit, offset int) ([]models.Toko, int64, error) {
	var tokos []models.Toko
	var total int64

	// Get total count
	r.db.Model(&models.Toko{}).Count(&total)

	// Get tokos with pagination and preload user
	err := r.db.Preload("User").Limit(limit).Offset(offset).Find(&tokos).Error
	return tokos, total, err
}

func (r *tokoRepository) UpdateToko(toko *models.Toko) error {
	return r.db.Save(toko).Error
}

func (r *tokoRepository) DeleteToko(id uint) error {
	return r.db.Delete(&models.Toko{}, id).Error
}

func (r *tokoRepository) CheckTokoOwnership(tokoID, userID uint) bool {
	var count int64
	r.db.Model(&models.Toko{}).Where("id = ? AND id_user = ?", tokoID, userID).Count(&count)
	return count > 0
}
