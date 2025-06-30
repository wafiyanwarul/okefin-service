package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/middleware"
	"github.com/wafiydev/okefin-service/internal/service"
)

type TokoHandler struct {
	tokoService service.TokoService
	validator   *validator.Validate
}

func NewTokoHandler(tokoService service.TokoService) *TokoHandler {
	return &TokoHandler{
		tokoService: tokoService,
		validator:   validator.New(),
	}
}

func (h *TokoHandler) CreateToko(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	var req dto.CreateTokoRequest
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

	toko, err := h.tokoService.CreateToko(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to create toko",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:  true,
		Message: "Toko created successfully",
		Errors:  nil,
		Data:    toko,
	})
}

func (h *TokoHandler) GetTokoByUserID(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	toko, err := h.tokoService.GetTokoByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Toko not found",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Toko retrieved successfully",
		Errors:  nil,
		Data:    toko,
	})
}

func (h *TokoHandler) GetTokoByID(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	tokoID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid toko ID",
			Errors:  []string{"Toko ID must be a valid number"},
			Data:    nil,
		})
	}

	toko, err := h.tokoService.GetTokoByID(uint(tokoID), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Toko not found",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Toko retrieved successfully",
		Errors:  nil,
		Data:    toko,
	})
}

func (h *TokoHandler) UpdateToko(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	tokoID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid toko ID",
			Errors:  []string{"Toko ID must be a valid number"},
			Data:    nil,
		})
	}

	var req dto.UpdateTokoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	toko, err := h.tokoService.UpdateToko(uint(tokoID), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to update toko",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Toko updated successfully",
		Errors:  nil,
		Data:    toko,
	})
}

func (h *TokoHandler) DeleteToko(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	tokoID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid toko ID",
			Errors:  []string{"Toko ID must be a valid number"},
			Data:    nil,
		})
	}

	err = h.tokoService.DeleteToko(uint(tokoID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to delete toko",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Toko deleted successfully",
		Errors:  nil,
		Data:    nil,
	})
}
