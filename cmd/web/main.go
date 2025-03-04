package main

import (
	"log"
	"rfidsystem/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	viewsEngine := html.New("./ui/html/pages", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	app.Static("/ui/static", "./ui/static")

	appHandler := handlers.NewHandler()
	app.Get("/", appHandler.HandleGetIndex)
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Fiber Web Server is running")
	})

	log.Fatal(app.Listen(":8080"), nil)
	// app.Listen(":8080")
}
