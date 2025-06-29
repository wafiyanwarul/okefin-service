package middleware

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wafiydev/okefin-service/internal/dto"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:  false,
				Message: "Missing authorization header",
				Errors:  []string{"Authorization header is required"},
				Data:    nil,
			})
		}

		// Check if token starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:  false,
				Message: "Invalid token format",
				Errors:  []string{"Token must start with 'Bearer '"},
				Data:    nil,
			})
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:  false,
				Message: "Invalid or expired token",
				Errors:  []string{"Token validation failed"},
				Data:    nil,
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:  false,
				Message: "Invalid token claims",
				Errors:  []string{"Unable to extract token claims"},
				Data:    nil,
			})
		}

		// Set user information in context
		c.Locals("userID", claims["id"])
		c.Locals("userEmail", claims["email"])

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// This middleware should be used after JWTProtected
		// You can implement admin check logic here
		// For now, we'll just pass through
		return c.Next()
	}
}

// Helper function to get user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}

	return uint(userID), nil
}
