package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/middleware"
	"github.com/wafiydev/okefin-service/internal/service"
)

type ProdukHandler struct {
	produkService service.ProdukService
	validator     *validator.Validate
}

func NewProdukHandler(produkService service.ProdukService) *ProdukHandler {
	return &ProdukHandler{
		produkService: produkService,
		validator:     validator.New(),
	}
}

func (h *ProdukHandler) CreateProduk(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	var req dto.CreateProdukRequest
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

	produk, err := h.produkService.CreateProduk(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to create produk",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:  true,
		Message: "Produk created successfully",
		Errors:  nil,
		Data:    produk,
	})
}

func (h *ProdukHandler) GetAllProdukByTokoID(c *fiber.Ctx) error {
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

	produkList, totalPages, totalItems, err := h.produkService.GetAllProdukByTokoID(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to get produk list",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"produk": produkList,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  totalItems,
			"limit":        limit,
		},
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Produk list retrieved successfully",
		Errors:  nil,
		Data:    response,
	})
}

func (h *ProdukHandler) GetProdukByID(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	produkID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid produk ID",
			Errors:  []string{"Produk ID must be a valid number"},
			Data:    nil,
		})
	}

	produk, err := h.produkService.GetProdukByID(userID, uint(produkID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Produk not found",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Produk retrieved successfully",
		Errors:  nil,
		Data:    produk,
	})
}

func (h *ProdukHandler) UpdateProduk(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	produkID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid produk ID",
			Errors:  []string{"Produk ID must be a valid number"},
			Data:    nil,
		})
	}

	var req dto.UpdateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	produk, err := h.produkService.UpdateProduk(userID, uint(produkID), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to update produk",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Produk updated successfully",
		Errors:  nil,
		Data:    produk,
	})
}

func (h *ProdukHandler) DeleteProduk(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	produkID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid produk ID",
			Errors:  []string{"Produk ID must be a valid number"},
			Data:    nil,
		})
	}

	err = h.produkService.DeleteProduk(userID, uint(produkID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to delete produk",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Produk deleted successfully",
		Errors:  nil,
		Data:    nil,
	})
}
