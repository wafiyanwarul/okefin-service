package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"Invalid JSON format: " + err.Error()},
			Data:    nil,
		})
	}

	// Validate request using struct tags
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				validationErrors = append(validationErrors, err.Field()+" is required")
			case "email":
				validationErrors = append(validationErrors, err.Field()+" must be a valid email")
			case "min":
				validationErrors = append(validationErrors, err.Field()+" must be at least "+err.Param()+" characters")
			default:
				validationErrors = append(validationErrors, err.Field()+" is invalid")
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  validationErrors,
			Data:    nil,
		})
	}

	// Additional custom validation
	errors := make([]string, 0)
	if req.Nama == "" {
		errors = append(errors, "nama is required")
	}
	if req.KataSandi == "" {
		errors = append(errors, "kata_sandi is required")
	} else if len(req.KataSandi) < 6 {
		errors = append(errors, "kata_sandi must be at least 6 characters")
	}
	if req.NoTelp == "" {
		errors = append(errors, "no_telp is required")
	}
	if req.TanggalLahir == "" {
		errors = append(errors, "tanggal_lahir is required")
	}
	if req.Pekerjaan == "" {
		errors = append(errors, "pekerjaan is required")
	}
	if req.Email == "" {
		errors = append(errors, "email is required")
	}
	if req.IDProvinsi == "" {
		errors = append(errors, "id_provinsi is required")
	}
	if req.IDKota == "" {
		errors = append(errors, "id_kota is required")
	}

	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  errors,
			Data:    nil,
		})
	}

	response, err := h.authService.Register(&req)
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
		Data:    response,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"Invalid JSON format: " + err.Error()},
			Data:    nil,
		})
	}

	// Detailed validation
	errors := make([]string, 0)
	if req.NoTelp == "" {
		errors = append(errors, "no_telp is required")
	}
	if req.KataSandi == "" {
		errors = append(errors, "kata_sandi is required")
	}

	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  errors,
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
