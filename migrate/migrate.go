package main

import (
	"github.com/antidoid/flightwatch/initializers"
	"github.com/antidoid/flightwatch/models"
)

func init() {
    initializers.LoadEnvVars()
    initializers.ConnectToDB()
}

func main() {
    initializers.DB.AutoMigrate(&models.Track{})
}
