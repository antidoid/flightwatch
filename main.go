package main

import (
	"os"

	"github.com/antidoid/fligthwatch/controllers"
	"github.com/antidoid/fligthwatch/initializers"

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html"
)

func init() {
    intializers.LoadEnvVars()
}

func main() {
    // Load templates
    engine := html.New("./views", ".html")

    // Setup App
	app := fiber.New(fiber.Config{
        Views: engine,
    })

    // Backend Routes
	app.Get("/api/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Universe!")
	})

    // Frontend Routes
    app.Get("/", controllers.Home)

    // Start the app
	app.Listen(":" + os.Getenv("PORT"))
}

