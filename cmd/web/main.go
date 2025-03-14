package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rfidsystem/internal/config"
	"rfidsystem/internal/handlers"
	"rfidsystem/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Load db config
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// init db toma
	db, err := repositories.NewDatabaseClient(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the database connection
	if err := testDBConnection(db.DB); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}
	log.Println("Database connection test successful.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Shutting down...")
		db.Close()
		os.Exit(0)
	}()
	// defer db.Close()

	viewsEngine := html.New("./ui/html/pages", ".html")

	app := fiber.New(fiber.Config{
		Views:                 viewsEngine,
		DisableStartupMessage: false,
		IdleTimeout:           time.Second * 60,
		ReadTimeout:           time.Second * 60,
		WriteTimeout:          time.Second * 60,
	})

	// Configure middleware with enhanced CORS settings for SSE
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080", // Specify exact origins instead of wildcard
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true, // Keep credentials enabled for cookies/auth
		ExposeHeaders:    "Content-Type, Content-Length, Content-Disposition",
	}))

	// Serve static files
	app.Static("/ui/static", "./ui/static")

	// Create static folder for profile images if it doesn't exist
	if _, err := os.Stat("./ui/static/images"); os.IsNotExist(err) {
		os.MkdirAll("./ui/static/images", 0755)
	}

	// Pass the database client to the handler
	appHandler := handlers.NewHandler(db)

	// Routes
	app.Get("/", appHandler.HandleGetIndex)
	app.Get("/grades", appHandler.HandleGrades)

	// SSE endpoint - crucial for real-time updates
	app.Get("/stream", appHandler.HandleSSE)

	// Card scan endpoint - receives data from Arduino
	app.Post("/card-scan", appHandler.HandleCardScan)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Fiber Web Server is running")
	})

	log.Fatal(app.Listen(":8080"))
}

func testDBConnection(db *sql.DB) error {
	// Ping the database
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}

	// Try a simple query
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Students").Scan(&count)
	if err != nil {
		return fmt.Errorf("test query failed: %v", err)
	}
	fmt.Printf("Found %d students in database\n", count)

	return nil
}
