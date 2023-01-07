package models

import "gorm.io/gorm"

type Search struct {
    gorm.Model
    Origin string 
    Destination string
    StartAt string
    EndAt string
    Contact string
    WayToContact string
}
