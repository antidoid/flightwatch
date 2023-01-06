package models

import (
	"gorm.io/gorm"

    "time"
)

type Search struct {
    gorm.Model
    Origin string 
    Destination string
    StartAt time.Time
    EndAt time.Time
    Contact string
    WayToContact string
}
