package main

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/Qodarrz/fiber-app/middleware"
	"github.com/Qodarrz/fiber-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

var (
	db   *sql.DB
	once sync.Once
)

// initEnv loads environment variables
func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}
}

// initDB initializes database only once
func initDB() *sql.DB {
	once.Do(func() {
		var err error
		// pakai DATABASE_URL biar simple (contoh: postgres://user:pass@host:port/dbname)
		db, err = sql.Open("pgx", os.Getenv("DATABASE_URL")+"?prefer_simple_protocol=true")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
	})
	return db
}

func main() {
	initEnv()
	db := initDB()
	defer db.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Fiber Auth App",
	})

	// Serve static files
	app.Static("/api/uploads", "./uploads")

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://reizentech-innoventure.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Middleware & routes
	mw := middleware.InitMiddlewares(db)
	routes.Setup(app, db, mw)

	// Port default :8080 (bisa override via .env PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
