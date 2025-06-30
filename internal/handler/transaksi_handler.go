package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/middleware"
	"github.com/wafiydev/okefin-service/internal/service"
)

type TransaksiHandler struct {
	transaksiService service.TransaksiService
	validator        *validator.Validate
}

func NewTransaksiHandler(transaksiService service.TransaksiService) *TransaksiHandler {
	return &TransaksiHandler{
		transaksiService: transaksiService,
		validator:        validator.New(),
	}
}

func (h *TransaksiHandler) CreateTransaksi(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	var req dto.CreateTransaksiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	var details []dto.DetailTransaksiRequest
	if err := c.BodyParser(&details); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid details",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	if err := h.validator.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" is invalid")
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  errors,
			Data:    nil,
		})
	}

	for _, d := range details {
		if err := h.validator.Struct(d); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid")
			}
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status:  false,
				Message: "Validation failed for details",
				Errors:  errors,
				Data:    nil,
			})
		}
	}

	transaksi, err := h.transaksiService.CreateTransaksi(userID, &req, details)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to create transaksi",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:  true,
		Message: "Transaksi created successfully",
		Errors:  nil,
		Data:    transaksi,
	})
}

func (h *TransaksiHandler) GetAllTransaksiByUserID(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transaksiList, totalPages, totalItems, err := h.transaksiService.GetAllTransaksiByUserID(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to get transaksi list",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"transaksi": transaksiList,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  totalItems,
			"limit":        limit,
		},
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Transaksi list retrieved successfully",
		Errors:  nil,
		Data:    response,
	})
}

func (h *TransaksiHandler) GetTransaksiByID(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	transaksiID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid transaksi ID",
			Errors:  []string{"Transaksi ID must be a valid number"},
			Data:    nil,
		})
	}

	transaksi, err := h.transaksiService.GetTransaksiByID(userID, uint(transaksiID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Transaksi not found",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Transaksi retrieved successfully",
		Errors:  nil,
		Data:    transaksi,
	})
}

func (h *TransaksiHandler) UpdateTransaksi(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	transaksiID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid transaksi ID",
			Errors:  []string{"Transaksi ID must be a valid number"},
			Data:    nil,
		})
	}

	var req dto.UpdateTransaksiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	if err := h.validator.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" is invalid")
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  errors,
			Data:    nil,
		})
	}

	transaksi, err := h.transaksiService.UpdateTransaksi(userID, uint(transaksiID), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to update transaksi",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Transaksi updated successfully",
		Errors:  nil,
		Data:    transaksi,
	})
}
