package main

import (
	"os"

	"github.com/antidoid/flightwatch/controllers"
	"github.com/antidoid/flightwatch/initializers"

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html"
)

func init() {
    initializers.LoadEnvVars()
    initializers.ConnectToDB()
}

func main() {
    // Load templates
    engine := html.New("./views", ".html")

    // Setup App
	app := fiber.New(fiber.Config{
        Views: engine,
    })

    // Backend Routes
    app.Post("/api/searchs", controllers.CreateSearch)
    app.Put("/api/searchs/:id", controllers.UpdateSearch)

	app.Get("/api/searchs", controllers.GetSearchs)
    app.Get("/api/searchs/:id", controllers.GetSearch)

    app.Delete("/api/searchs/:id", controllers.DeleteSearch)

    // Frontend Routes
    app.Get("/", controllers.Home)

    // Start the app
	app.Listen(":" + os.Getenv("PORT"))
}

