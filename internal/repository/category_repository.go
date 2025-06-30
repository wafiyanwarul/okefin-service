package repository

import (
	"github.com/wafiydev/okefin-service/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories(limit, offset int) ([]models.Category, int64, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
	CheckCategoryExists(id uint) bool
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAllCategories(limit, offset int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Get total count
	r.db.Model(&models.Category{}).Count(&total)

	// Get categories with pagination
	err := r.db.Limit(limit).Offset(offset).Find(&categories).Error
	return categories, total, err
}

func (r *categoryRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) CheckCategoryExists(id uint) bool {
	var count int64
	r.db.Model(&models.Category{}).Where("id = ?", id).Count(&count)
	return count > 0
}
