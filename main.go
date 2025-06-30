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
	"github.com/wafiydev/okefin-service/internal/middleware"
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
	userRepo := repository.NewUserRepository(db)
	alamatRepo := repository.NewAlamatRepository(db)
	tokoRepo := repository.NewTokoRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	authService := service.NewAuthService(authRepo, service.NewTokoService(tokoRepo))
	userService := service.NewUserService(userRepo, authService)
	alamatService := service.NewAlamatService(alamatRepo)
	tokoService := service.NewTokoService(tokoRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	alamatHandler := handler.NewAlamatHandler(alamatService)
	tokoHandler := handler.NewTokoHandler(tokoService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Set authRepo for middleware
	middleware.SetAuthRepo(authRepo) // Inject repository into middleware

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

	// User routes
	user := app.Group("/user")
	user.Get("/my", middleware.JWTProtected(), userHandler.GetMyProfile)
	user.Put("/my", middleware.JWTProtected(), userHandler.UpdateMyProfile)
	user.Get("/", middleware.JWTProtected(), userHandler.GetAllUsers)
	user.Get("/:id", middleware.JWTProtected(), userHandler.GetUserByID)

	// Alamat routes
	alamat := app.Group("/alamat")
	alamat.Post("/", middleware.JWTProtected(), alamatHandler.CreateAlamat)
	alamat.Get("/my", middleware.JWTProtected(), alamatHandler.GetMyAlamat)
	alamat.Get("/:id", middleware.JWTProtected(), alamatHandler.GetAlamatByID)
	alamat.Put("/:id", middleware.JWTProtected(), alamatHandler.UpdateAlamat)
	alamat.Delete("/:id", middleware.JWTProtected(), alamatHandler.DeleteAlamat)

	// Toko routes
	toko := app.Group("/toko")
	toko.Post("/", middleware.JWTProtected(), tokoHandler.CreateToko)
	toko.Get("/my", middleware.JWTProtected(), tokoHandler.GetTokoByUserID)
	toko.Get("/:id", middleware.JWTProtected(), tokoHandler.GetTokoByID)
	toko.Put("/:id", middleware.JWTProtected(), tokoHandler.UpdateToko)
	toko.Delete("/:id", middleware.JWTProtected(), tokoHandler.DeleteToko)

	// Category routes (admin-only)
	category := app.Group("/category")
	category.Post("/", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.CreateCategory)
	category.Get("/", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.GetAllCategories)
	category.Get("/:id", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.GetCategoryByID)
	category.Put("/:id", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.UpdateCategory)
	category.Delete("/:id", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.DeleteCategory)

	// File upload route
	app.Post("/upload", middleware.JWTProtected(), userHandler.UploadFile)

	// Serve static files
	app.Static("/uploads", "./uploads")

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}