package models

import (
	"github.com/antidoid/flightwatch/initializers"
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	StartAt     string `json:"startat"`
	EndAt       string `json:"endat"`
	Contact     string `json:"contact"`
	Threshold   string `json:"threshold"`
	UserIp      string `json:"userip"`
}

func CreateTrack(track Track) error {
	tx := initializers.DB.Create(&track)
	return tx.Error
}

func GetTracks() ([]Track, error) {
	var tracks []Track
	tx := initializers.DB.Order("ID asc").Find(&tracks)
	if tx.Error != nil {
		return []Track{}, tx.Error
	}

	return tracks, nil
}

func GetTrack(id string) (Track, error) {
	var track Track
	tx := initializers.DB.First(&track, id)
	if tx.Error != nil {
		return Track{}, tx.Error
	}

	return track, nil
}

func UpdateTrack(track Track, newTrack Track) error {
	tx := initializers.DB.Model(&track).Updates(Track{
		Origin:      newTrack.Origin,
		Destination: newTrack.Destination,
		StartAt:     newTrack.StartAt,
		EndAt:       newTrack.EndAt,
		Contact:     newTrack.Contact,
		Threshold:   newTrack.Threshold,
		UserIp:      newTrack.UserIp,
	})

	return tx.Error
}

func DeleteTrack(track Track) error {
	tx := initializers.DB.Unscoped().Delete(&track)
	return tx.Error
}
