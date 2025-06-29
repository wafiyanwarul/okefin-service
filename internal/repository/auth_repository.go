package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByNoTelp(noTelp string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateToko(toko *models.Toko) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByNoTelp(noTelp string) (*models.User, error) {
	var user models.User
	err := r.db.Where("no_telp = ?", noTelp).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateToko(toko *models.Toko) error {
	return r.db.Create(toko).Error
}
