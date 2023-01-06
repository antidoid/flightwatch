package models

import (
	"gorm.io/gorm"

    "time"
)

type Search struct {
    gorm.Model
    ORGIN string 
    DESTINATION string
    START_DATE time.Time
    END_DATE time.Time
    CONTACT string
    MEDIUM_OF_CONTACT string
}
