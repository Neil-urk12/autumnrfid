// Package main implements the web server for the RFID system.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"rfidsystem/internal/config"
	"rfidsystem/internal/handlers"
	"rfidsystem/internal/model"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

// main initializes and runs the RFID web server.
func main() {
	// Load db config
	dbClient, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbClient.Close()

	// Handle graceful shutdown on interrupt signal
	go handleShutdown(dbClient)

	// Initialize view engine and Fiber app
	engine := initViewEngine()
	app := configureApp(engine)

	// Create handler with repository
	rfidRepo := repositories.NewRFIDRepository(dbClient)
	handler := handlers.NewHandler(dbClient, rfidRepo)

	// Register all routes
	registerRoutes(app, handler)

	// Start server
	log.Fatal(app.Listen(":8080"))
}

// initDatabase loads configuration, connects to the database, and verifies connectivity.
func initDatabase() (*repositories.DatabaseClient, error) {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("load database config: %v", err)
	}

	dbClient, err := repositories.NewDatabaseClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %v", err)
	}

	if err := testDBConnection(dbClient.DB); err != nil {
		return nil, fmt.Errorf("test database connection: %v", err)
	}

	log.Println("Database connection test successful.")
	return dbClient, nil
}

// initViewEngine sets up the HTML template engine with development options and custom functions.
func initViewEngine() *html.Engine {
	engine := html.New("./ui/html", ".html")
	engine.Reload(true)
	engine.Debug(true)

	engine.AddFunc("lower", strings.ToLower)
	engine.AddFunc("feesByCategory", func(fees []model.FeeBreakdown, category string) []model.FeeBreakdown {
		var filtered []model.FeeBreakdown
		for _, fee := range fees {
			if fee.Category == category {
				filtered = append(filtered, fee)
			}
		}
		return filtered
	})
	// Format time as "YYYY-MM-DD hh:mm am/pm" without seconds
	engine.AddFunc("formatTime", func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return strings.ToLower(t.Format("2006-01-02 03:04 PM"))
	})
	return engine
}

// configureApp creates a Fiber app, applies middleware, and ensures static assets are served.
func configureApp(engine *html.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: false,
		IdleTimeout:           time.Hour * 24,
		ReadTimeout:           time.Second * 60,
		WriteTimeout:          0,
		ColorScheme: fiber.Colors{
			Black:   "\u001b[93m",
			Red:     "\u001b[91m",
			Green:   "\u001b[92m",
			Yellow:  "\u001b[93m",
			Blue:    "\u001b[94m",
			Magenta: "\u001b[95m",
			Cyan:    "\u001b[96m",
			White:   "\u001b[97m",
			Reset:   "\u001b[0m",
		},
	})

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080",
		AllowHeaders:     "Origin, Content-Type, Accept, Cache-Control",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Type, Content-Length, Content-Disposition",
	}))

	// Logger middleware (basic HTTP logging)
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output: os.Stdout,
	}))

	// Serve static assets and ensure images directory exists
	app.Static("/ui/static", "./ui/static")
	if _, err := os.Stat("./ui/static/images"); os.IsNotExist(err) {
		os.MkdirAll("./ui/static/images", 0755)
	}

	return app
}

// registerRoutes maps URL paths to their handler functions.
func registerRoutes(app *fiber.App, h *handlers.AppHandler) {
	app.Get("/", h.HandleGetIndex)
	app.Get("/docs", h.HandleDocs)
	app.Get("/grades", h.HandleGrades)
	app.Get("/grades/semester/:studentId", h.HandleSemesterGrades)
	app.Get("/error", h.HandleError)
	app.Get("/student-partial/:rfid", h.HandleStudentInfo)
	app.Get("/students/v1", h.RetrieveStudentsHandler)
	app.Get("/students/:id", h.GetStudentById)
	app.Get("/stream", h.HandleSSE)
	app.Get("/log", h.HandleLog)
	app.Get("/logs", h.HandleLog)
	// HTMX polling endpoint for log container
	app.Get("/log/partial", h.HandleLogPartial)
	// HTMX polling endpoint for stats cards
	app.Get("/stats/partial", h.HandleStatsPartial)
	app.Post("/card-scan", h.HandleCardScan)
	app.Get("/card-scan-ws", websocket.New(h.HandleCardScanWS))
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("Fiber Web Server is running") })
	app.Get("/bills", h.HandleBills)
	// support HTMX POST navigation with hidden RFID
	app.Post("/student-partial", h.HandleStudentInfo)
	app.Post("/grades", h.HandleGrades)
	app.Post("/bills", h.HandleBills)
}

// handleShutdown listens for interrupt signals to gracefully close resources.
func handleShutdown(dbClient *repositories.DatabaseClient) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down...")
	dbClient.Close()
	os.Exit(0)
}

// testDBConnection pings the database and runs a simple query to verify connectivity.
func testDBConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM Students").Scan(&count); err != nil {
		return fmt.Errorf("test query failed: %v", err)
	}
	fmt.Printf("Found %d students in database\n", count)
	return nil
}
