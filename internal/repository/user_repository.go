package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	GetAllUsers(limit, offset int) ([]models.User, int64, error)
	CheckEmailExists(email string, userID uint) bool
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) GetAllUsers(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	r.db.Model(&models.User{}).Count(&total)

	// Get users with pagination
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) CheckEmailExists(email string, userID uint) bool {
	var count int64
	r.db.Model(&models.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
	return count > 0
}
