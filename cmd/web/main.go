package main

import (
	"net/http"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fiberServer()
	muxServer()
}

func muxServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHanlder)
	http.ListenAndServe(":8080", mux)
}

func fiberServer() {
	app := fiber.New()

	app.Get("/toma", func(c *fiber.Ctx) error {
		return c.SendString("Hello, toma")
	})
	app.Listen(":8090")
}

func helloHanlder(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("ui/html/pages/home.html"))
	templ.Execute(w, nil)
}
