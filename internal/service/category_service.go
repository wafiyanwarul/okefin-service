package service

import (
	"errors"
	"time"

	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"
)

type CategoryService interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetAllCategories(page, limit int) ([]dto.CategoryResponse, int, int64, error)
	GetCategoryByID(id uint) (*dto.CategoryResponse, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &models.Category{
		NamaCategory: req.NamaCategory,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.categoryRepo.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	response := &dto.CategoryResponse{
		ID:           category.ID,
		NamaCategory: category.NamaCategory,
	}

	return response, nil
}

func (s *categoryService) GetAllCategories(page, limit int) ([]dto.CategoryResponse, int, int64, error) {
	offset := (page - 1) * limit
	categories, total, err := s.categoryRepo.GetAllCategories(limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []dto.CategoryResponse
	for _, c := range categories {
		response := dto.CategoryResponse{
			ID:           c.ID,
			NamaCategory: c.NamaCategory,
		}
		responses = append(responses, response)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return responses, totalPages, total, nil
}

func (s *categoryService) GetCategoryByID(id uint) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	response := &dto.CategoryResponse{
		ID:           category.ID,
		NamaCategory: category.NamaCategory,
	}

	return response, nil
}

func (s *categoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Update fields
	if req.NamaCategory != "" {
		category.NamaCategory = req.NamaCategory
	}

	category.UpdatedAt = time.Now()

	err = s.categoryRepo.UpdateCategory(category)
	if err != nil {
		return nil, err
	}

	response := &dto.CategoryResponse{
		ID:           category.ID,
		NamaCategory: category.NamaCategory,
	}

	return response, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	if !s.categoryRepo.CheckCategoryExists(id) {
		return errors.New("category not found")
	}

	return s.categoryRepo.DeleteCategory(id)
}
