package initializers

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
    var err error
    DB, err = gorm.Open(sqlite.Open(os.Getenv("DB")), &gorm.Config{})

    if err != nil {
        log.Fatal("Error connecting to the database")
    }
}

