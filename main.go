package main

import (
	"os"

	"github.com/antidoid/fligthwatch/initializers"

	"github.com/gofiber/fiber/v2"
)

func init() {
    intializers.LoadEnvVars()
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Universe!")
	})

	app.Listen(":" + os.Getenv("PORT"))
}

