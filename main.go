package main

import (
	"os"

	"github.com/antidoid/flightwatch/controllers"
	"github.com/antidoid/flightwatch/initializers"
	"github.com/antidoid/flightwatch/helpers/skyscan"

	"github.com/gofiber/fiber/v2"
)

func init() {
    initializers.LoadEnvVars()
    initializers.ConnectToDB()
}

func main() {
    // Backend Tasks
    go skyscan.ScanAllTracks()

    // Setup App
	app := fiber.New()

    // Backend Routes
    app.Post("/api/tracks", controllers.CreateTrack)
    app.Put("/api/tracks/:id", controllers.UpdateTrack)

    app.Get("/api/tracks", controllers.GetTracks)
    app.Get("/api/tracks/:id", controllers.GetTrack)

    app.Delete("/api/tracks/:id", controllers.DeleteTrack)

    // Frontend Routes
    app.Get("/", controllers.Home)

    // Start the app
    app.Listen(":" + os.Getenv("PORT"))
}

