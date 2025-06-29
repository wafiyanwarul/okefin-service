package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"Invalid request format"},
			Data:    nil,
		})
	}

	// Basic validation
	if req.Email == "" || req.NoTelp == "" || req.KataSandi == "" || req.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"Required fields are missing"},
			Data:    nil,
		})
	}

	err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Errors:  nil,
		Data:    "Register Succeed",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"Invalid request format"},
			Data:    nil,
		})
	}

	// Basic validation
	if req.NoTelp == "" || req.KataSandi == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"No Telp and kata_sandi are required"},
			Data:    nil,
		})
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Errors:  nil,
		Data:    response,
	})
}
