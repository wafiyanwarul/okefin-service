package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/service"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	validator       *validator.Validate
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validator:       validator.New(),
	}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" is required")
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  errors,
			Data:    nil,
		})
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to create category",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:  true,
		Message: "Category created successfully",
		Errors:  nil,
		Data:    category,
	})
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	categories, totalPages, totalItems, err := h.categoryService.GetAllCategories(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to get categories",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"categories": categories,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  totalItems,
			"limit":        limit,
		},
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Categories retrieved successfully",
		Errors:  nil,
		Data:    response,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid category ID",
			Errors:  []string{"Category ID must be a valid number"},
			Data:    nil,
		})
	}

	category, err := h.categoryService.GetCategoryByID(uint(categoryID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Category not found",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Category retrieved successfully",
		Errors:  nil,
		Data:    category,
	})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid category ID",
			Errors:  []string{"Category ID must be a valid number"},
			Data:    nil,
		})
	}

	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	category, err := h.categoryService.UpdateCategory(uint(categoryID), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to update category",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Category updated successfully",
		Errors:  nil,
		Data:    category,
	})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid category ID",
			Errors:  []string{"Category ID must be a valid number"},
			Data:    nil,
		})
	}

	err = h.categoryService.DeleteCategory(uint(categoryID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to delete category",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Category deleted successfully",
		Errors:  nil,
		Data:    nil,
	})
}
