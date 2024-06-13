package main

import (
	"os"

	"github.com/antidoid/flightwatch/controllers"
	"github.com/antidoid/flightwatch/helpers/skyscan"
	"github.com/antidoid/flightwatch/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	go skyscan.ScanAllTracks()
}

func main() {
	// Setup App
	app := fiber.New()

	// Setup CORS
	app.Use(cors.New())

	// Backend Routes
	app.Post("/api/tracks", controllers.CreateTrack)
	app.Put("/api/tracks/:id", controllers.UpdateTrack)

	app.Get("/api/tracks", controllers.GetTracks)
	app.Get("/api/tracks/:id", controllers.GetTrack)

	app.Delete("/api/tracks/:id", controllers.DeleteTrack)

	// Start the app
	app.Listen(":" + os.Getenv("PORT"))
}
