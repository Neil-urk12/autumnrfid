package main

import (
	"log"
	"rfidsystem/internal/config"
	"rfidsystem/internal/handlers"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
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
	defer db.Close()

	viewsEngine := html.New("./ui/html/pages", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	app.Static("/ui/static", "./ui/static")

	// Pass the database client to the handler
	appHandler := handlers.NewHandler(db)
	app.Get("/", appHandler.HandleGetIndex)
	app.Get("/stream", appHandler.HandleSSE)

	app.Post("/card-scan", appHandler.HandleCardScan)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Fiber Web Server is running")
	})

	log.Fatal(app.Listen(":8080"))
}
