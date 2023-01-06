package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
    var err error
    DB, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})

    if err != nil {
        log.Fatal("Error connecting to the database")
    }
}

