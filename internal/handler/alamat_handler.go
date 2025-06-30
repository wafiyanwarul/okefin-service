package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/middleware"
	"github.com/wafiydev/okefin-service/internal/service"
)

type AlamatHandler struct {
	alamatService service.AlamatService
	validator     *validator.Validate
}

func NewAlamatHandler(alamatService service.AlamatService) *AlamatHandler {
	return &AlamatHandler{
		alamatService: alamatService,
		validator:     validator.New(),
	}
}

func (h *AlamatHandler) CreateAlamat(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	var req dto.CreateAlamatRequest
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

	alamat, err := h.alamatService.CreateAlamat(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to create alamat",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:  true,
		Message: "Alamat created successfully",
		Errors:  nil,
		Data:    alamat,
	})
}

func (h *AlamatHandler) GetMyAlamat(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	alamats, totalPages, totalItems, err := h.alamatService.GetMyAlamat(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to get alamat",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"alamats": alamats,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  totalItems,
			"limit":        limit,
		},
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Alamat retrieved successfully",
		Errors:  nil,
		Data:    response,
	})
}

func (h *AlamatHandler) GetAlamatByID(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	alamatID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid alamat ID",
			Errors:  []string{"Alamat ID must be a valid number"},
			Data:    nil,
		})
	}

	alamat, err := h.alamatService.GetAlamatByID(uint(alamatID), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Alamat not found",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Alamat retrieved successfully",
		Errors:  nil,
		Data:    alamat,
	})
}

func (h *AlamatHandler) UpdateAlamat(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	alamatID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid alamat ID",
			Errors:  []string{"Alamat ID must be a valid number"},
			Data:    nil,
		})
	}

	var req dto.UpdateAlamatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	alamat, err := h.alamatService.UpdateAlamat(uint(alamatID), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to update alamat",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Alamat updated successfully",
		Errors:  nil,
		Data:    alamat,
	})
}

func (h *AlamatHandler) DeleteAlamat(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	alamatID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid alamat ID",
			Errors:  []string{"Alamat ID must be a valid number"},
			Data:    nil,
		})
	}

	err = h.alamatService.DeleteAlamat(uint(alamatID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to delete alamat",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Alamat deleted successfully",
		Errors:  nil,
		Data:    nil,
	})
}
