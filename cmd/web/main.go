package main

import (
	"rfidsystem/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	viewsEngine := html.New("./ui/html/pages", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	// app.Static("/ui/static", "./static")
	appHandler := handlers.NewHandler()

	app.Get("/", appHandler.HandleGetIndex)

	app.Get("/toma", func(c *fiber.Ctx) error {
		return c.SendString("Hello, toma")
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Fiber Web Server is running")
	})

	app.Listen(":8080")
}
