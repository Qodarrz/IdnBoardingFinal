package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Qodarrz/fiber-app/middleware"
	"github.com/Qodarrz/fiber-app/routes"
	"github.com/gofiber/adaptor/v2"
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

// initDB initializes database only once (efficient for serverless)
func initDB() *sql.DB {
	once.Do(func() {
		var err error
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

// Handler is the function called by Vercel for every request
func Handler(w http.ResponseWriter, r *http.Request) {
	initEnv()

	db := initDB()

	

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Fiber Auth App",
	})

	app.Static("/api/uploads", "./uploads")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://reizentech-innoventure.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
	}))

	// Initialize middleware dan routes
	mw := middleware.InitMiddlewares(db)
	routes.Setup(app, db, mw)

	// Serve Fiber app as http.Handler for Vercel
	adaptor.FiberApp(app).ServeHTTP(w, r)
}