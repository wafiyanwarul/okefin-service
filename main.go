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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDatabase()

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

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	alamatRepo := repository.NewAlamatRepository(db)
	tokoRepo := repository.NewTokoRepository(db)
	produkRepo := repository.NewProdukRepository(db)
	fotoRepo := repository.NewFotoProdukRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	authService := service.NewAuthService(authRepo, service.NewTokoService(tokoRepo))
	userService := service.NewUserService(userRepo, authService)
	alamatService := service.NewAlamatService(alamatRepo)
	tokoService := service.NewTokoService(tokoRepo)
	produkService := service.NewProdukService(produkRepo, tokoRepo, fotoRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	alamatHandler := handler.NewAlamatHandler(alamatService)
	tokoHandler := handler.NewTokoHandler(tokoService)
	produkHandler := handler.NewProdukHandler(produkService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	middleware.SetAuthRepo(authRepo)

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

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Okefin-Service!")
	})

	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	user := app.Group("/user")
	user.Get("/my", middleware.JWTProtected(), userHandler.GetMyProfile)
	user.Put("/my", middleware.JWTProtected(), userHandler.UpdateMyProfile)
	user.Get("/", middleware.JWTProtected(), userHandler.GetAllUsers)
	user.Get("/:id", middleware.JWTProtected(), userHandler.GetUserByID)

	alamat := app.Group("/alamat")
	alamat.Post("/", middleware.JWTProtected(), alamatHandler.CreateAlamat)
	alamat.Get("/my", middleware.JWTProtected(), alamatHandler.GetMyAlamat)
	alamat.Get("/:id", middleware.JWTProtected(), alamatHandler.GetAlamatByID)
	alamat.Put("/:id", middleware.JWTProtected(), alamatHandler.UpdateAlamat)
	alamat.Delete("/:id", middleware.JWTProtected(), alamatHandler.DeleteAlamat)

	toko := app.Group("/toko")
	toko.Post("/", middleware.JWTProtected(), tokoHandler.CreateToko)
	toko.Get("/my", middleware.JWTProtected(), tokoHandler.GetTokoByUserID)
	toko.Get("/:id", middleware.JWTProtected(), tokoHandler.GetTokoByID)
	toko.Put("/:id", middleware.JWTProtected(), tokoHandler.UpdateToko)
	toko.Delete("/:id", middleware.JWTProtected(), tokoHandler.DeleteToko)

	produk := app.Group("/produk")
	produk.Post("/", middleware.JWTProtected(), produkHandler.CreateProduk)
	produk.Get("/toko", middleware.JWTProtected(), produkHandler.GetAllProdukByTokoID)
	produk.Get("/:id", middleware.JWTProtected(), produkHandler.GetProdukByID)
	produk.Put("/:id", middleware.JWTProtected(), produkHandler.UpdateProduk)
	produk.Delete("/:id", middleware.JWTProtected(), produkHandler.DeleteProduk)

	category := app.Group("/category")
	category.Post("/", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.CreateCategory)
	category.Get("/", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.GetAllCategories)
	category.Get("/:id", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.GetCategoryByID)
	category.Put("/:id", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.UpdateCategory)
	category.Delete("/:id", middleware.JWTProtected(), middleware.AdminOnly(), categoryHandler.DeleteCategory)

	app.Post("/upload", middleware.JWTProtected(), userHandler.UploadFile)
	app.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
