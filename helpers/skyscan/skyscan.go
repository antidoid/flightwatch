package skyscan

import (
	"fmt"

	"github.com/antidoid/flightwatch/helpers/notify"
	"github.com/antidoid/flightwatch/initializers"
	"github.com/antidoid/flightwatch/models"

	"github.com/fxtlabs/date"
)

func getDate(d string) (date.Date) {
    res, _ := date.Parse("2006-01-02", d)
    return res
}

func scanTrack(track *models.Track) {
    // Loop from start date to end date
    startDate := getDate(track.StartAt)
    endDate := getDate(track.EndAt)

    for d := startDate; d.Sub(endDate) <= 0; d = d.Add(1) {
        // Check if price has reached threshold
        flightNo, price, link := getCheapestFlight(track.Origin, track.Destination, d)
        if (hasHitThreshold(price, track.Threshold)) {
            message := fmt.Sprintf("\nGreeting from FlightWatch\nYour tracked flight (%s) is currently priced at Rs%s\n Book now at: %s\nHave a nice day",
                flightNo, price, link)
            notify.SendSMS(track.Contact, message)
        }
    }
}

func scanAllTracks() error {
    var tracks []models.Track
    // Get the database
    res := initializers.DB.Find(&tracks)
    if res.Error != nil {
        return res.Error
    }

    // query over each row
    for _, track := range tracks {
        scanTrack(&track)
    }
    return nil
}

// return Flight number, price and booking link
func getCheapestFlight (ogn string, dsn string, date date.Date) (string, string, string)

// If current price has hit threshold -> notify user
func hasHitThreshold(price string, threshold string) bool

