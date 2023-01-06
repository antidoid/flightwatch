package main

import (
	"github.com/antidoid/fligthwatch/initializers"
	"github.com/antidoid/fligthwatch/models"
)

func init() {
    initializers.LoadEnvVars()
    initializers.ConnectToDB()
}

func main() {
    initializers.DB.AutoMigrate(&models.Search{})
}
