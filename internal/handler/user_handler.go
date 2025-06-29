package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/middleware"
	"github.com/wafiydev/okefin-service/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMyProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	response, err := h.userService.GetUserProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    response,
	})
}

func (h *UserHandler) UpdateMyProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{"Invalid request format"},
			Data:    nil,
		})
	}

	response, err := h.userService.UpdateUserProfile(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Succeed to PUT data",
		Errors:  nil,
		Data:    response,
	})
}

func (h *UserHandler) UploadFile(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to upload file",
			Errors:  []string{"No file provided"},
			Data:    nil,
		})
	}

	response, err := h.userService.UploadFile(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to upload file",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "File uploaded successfully",
		Errors:  nil,
		Data:    response,
	})
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Get pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, totalPages, totalData, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	response := fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_data":   totalData,
			"per_page":     limit,
		},
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    response,
	})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{"Invalid user ID"},
			Data:    nil,
		})
	}

	response, err := h.userService.GetUserProfile(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    response,
	})
}
