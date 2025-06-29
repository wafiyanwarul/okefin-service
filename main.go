package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/wafiydev/okefin-service/config"
	"github.com/wafiydev/okefin-service/internal/handler"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"
	"github.com/wafiydev/okefin-service/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	config.ConnectDatabase()

	// Auto migrate database
	db := config.GetDB()
	err := db.AutoMigrate(
		&models.User{},
		&models.Alamat{},
		&models.Toko{},
		&models.Category{},
		&models.Produk{},
		&models.FotoProduk{},
		&models.Trx{},
		&models.DetailTrx{},
		&models.LogProduk{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)

	// Initialize services
	authService := service.NewAuthService(authRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  false,
				"message": "Internal Server Error",
				"errors":  []string{err.Error()},
				"data":    nil,
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Okefin-Service!")
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
