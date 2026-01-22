package main

import (
	"log"
	"os"
	"qr-shorten-go/handlers"
	"qr-shorten-go/middleware"
	"qr-shorten-go/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Koneksi Postgres
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// 3. Auto-Migrate (Otomatis bikin tabel di Postgres)
	db.AutoMigrate(&models.User{}, &models.Link{})

	

	// 4. Setup Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000", // Ganti dengan URL frontend kamu nanti
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
        AllowMethods: "GET, POST, PUT, DELETE",
    }))
	// Inisialisasi Handler
    linkHandler := handlers.LinkHandler{DB: db}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend Go Jalan!")
	})
	app.Get("/api/auth/google", handlers.LoginGoogle)
	app.Get("/api/auth/google/callback", func(c *fiber.Ctx) error {
    return handlers.CallbackGoogle(c, db)
		})
	
	app.Get("/api/profile", middleware.AuthGuard, linkHandler.GetProfile)
	// LINK & QR ROUTES
    app.Post("/api/shorten", middleware.AuthGuard, linkHandler.CreateShorten)
    apiGroup := app.Group("/api", middleware.AuthGuard) // Tips: Bisa dikelompokkan pakai Group
    apiGroup.Get("/my-links", linkHandler.GetUserLinks)
    apiGroup.Delete("/links/:id", linkHandler.DeleteLink)
	app.Get("/api/stats", middleware.AuthGuard, linkHandler.GetStats)
    
    // REDIRECT ROUTE (Harus di paling bawah agar tidak bentrok)
    app.Get("/:code", linkHandler.Resolve)
	// Jalankan Server
	log.Fatal(app.Listen(":8080"))
}