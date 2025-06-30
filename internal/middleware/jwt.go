package middleware

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/repository"

)

var authRepo repository.AuthRepository // Declare as a package-level variable

// SetAuthRepo sets the AuthRepository for middleware
func SetAuthRepo(repo repository.AuthRepository) {
	authRepo = repo
}

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
		if authRepo == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status:  false,
				Message: "Middleware not initialized",
				Errors:  []string{"Auth repository not set"},
				Data:    nil,
			})
		}

		// Get userID from context (set by JWTProtected)
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:  false,
				Message: "Unauthorized",
				Errors:  []string{"Invalid user ID"},
				Data:    nil,
			})
		}

		// Fetch user from repository to check IsAdmin
		user, err := authRepo.GetUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status:  false,
				Message: "Failed to fetch user",
				Errors:  []string{err.Error()},
				Data:    nil,
			})
		}

		if !user.IsAdmin {
			return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
				Status:  false,
				Message: "Admin access required",
				Errors:  []string{"Only admins can perform this action"},
				Data:    nil,
			})
		}

		return c.Next()
	}
}

// Helper function to get user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	raw := c.Locals("userID")

	switch val := raw.(type) {
	case float64:
		return uint(val), nil
	case int:
		return uint(val), nil
	case string:
		userID, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID string")
		}
		return uint(userID), nil
	default:
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID type")
	}
}