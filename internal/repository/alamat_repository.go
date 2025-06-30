package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"
)

type AlamatRepository interface {
	CreateAlamat(alamat *models.Alamat) error
	GetAlamatByID(id uint) (*models.Alamat, error)
	GetAlamatByUserID(userID uint, limit, offset int) ([]models.Alamat, int64, error)
	UpdateAlamat(alamat *models.Alamat) error
	DeleteAlamat(id uint) error
	CheckAlamatOwnership(alamatID, userID uint) bool
}

type alamatRepository struct {
	db *gorm.DB
}

func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db: db}
}

func (r *alamatRepository) CreateAlamat(alamat *models.Alamat) error {
	return r.db.Create(alamat).Error
}

func (r *alamatRepository) GetAlamatByID(id uint) (*models.Alamat, error) {
	var alamat models.Alamat
	err := r.db.Where("id = ?", id).First(&alamat).Error
	if err != nil {
		return nil, err
	}
	return &alamat, nil
}

func (r *alamatRepository) GetAlamatByUserID(userID uint, limit, offset int) ([]models.Alamat, int64, error) {
	var alamat []models.Alamat
	var total int64

	// Get total count
	r.db.Model(&models.Alamat{}).Where("id_user = ?", userID).Count(&total)

	// Get alamat with pagination
	err := r.db.Where("id_user = ?", userID).Limit(limit).Offset(offset).Find(&alamat).Error
	return alamat, total, err
}

func (r *alamatRepository) UpdateAlamat(alamat *models.Alamat) error {
	return r.db.Save(alamat).Error
}

func (r *alamatRepository) DeleteAlamat(id uint) error {
	return r.db.Delete(&models.Alamat{}, id).Error
}

func (r *alamatRepository) CheckAlamatOwnership(alamatID, userID uint) bool {
	var count int64
	r.db.Model(&models.Alamat{}).Where("id = ? AND id_user = ?", alamatID, userID).Count(&count)
	return count > 0
}
