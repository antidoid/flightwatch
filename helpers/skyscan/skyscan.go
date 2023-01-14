package skyscan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/antidoid/flightwatch/helpers/notify"
	"github.com/antidoid/flightwatch/initializers"
	"github.com/antidoid/flightwatch/models"

	"github.com/fxtlabs/date"
)

func getDate(d string) (date.Date) {
    res, _ := date.Parse("2006-01-02", d)
    return res
}

func scanTrack(track *models.Track) error {
    // Loop from start date to end date
    startDate := getDate(track.StartAt)
    endDate := getDate(track.EndAt)

    for d := startDate; d.Sub(endDate) <= 0; d = d.Add(1) {
        // Check if price has reached threshold
        price, link, err := getCheapestFlight(track.Origin, track.Destination, d)
        if err != nil {
            return err
        }

        if (hasHitThreshold(price, track.Threshold)) {
            message := fmt.Sprintf("\nGreeting from FlightWatch\nYour tracked flight is currently priced at Rs%s\n Book now at: %s\nHave a nice day",
                price, link)
            notify.SendSMS(track.Contact, message)
        }
    }
    return nil
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

// TODO: Add Date and origin and destination support
// return Flight number, price and booking link
func getCheapestFlight(ogn string, dsn string, date date.Date) (string, string, error) {
    payload := map[string]map[string]interface{}{
        "query": {
            "market": "IN",
            "locale": "en-GB",
            "currency": "INR",
            "cabinClass": "CABIN_CLASS_ECONOMY",
            "adults": 1,
            "queryLegs": []map[string]interface{}{{
                "originPlaceId": map[string]string{"iata": ogn,},
                "destinationPlaceId": map[string]string{"iata": dsn,},
                "date": map[string]int{
                    "year": date.Year(),
                    "month": int(date.Month()),
                    "day": date.Day(),
                },
            }},
        },
    }

    postBody, err := json.Marshal(payload)
    if err != nil {
        return "", "", err
    }

    url := "https://partners.api.skyscanner.net/apiservices/v3/flights/live/search/create"
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-api-key", "prtl6749387986743898559646983194")
    if err != nil {
        return "", "", err
    }

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", "", err
    }
    defer res.Body.Close()

    responseData, _ := ioutil.ReadAll(res.Body)

    type Content struct {
        Results map[string]map[string]map[string][]map[string][]map[string]interface{} `json:"results"`
        SortingOptions map[string][]map[string]string `json:"sortingOptions"`
    }

    type Response struct {
        SessionToken string `json:"sessionToken"`
        Status string `json:"status"`
        Content `json:"content"`
    }

    var respBody Response
    json.Unmarshal(responseData, &respBody)

    itenaryId := respBody.Content.SortingOptions["cheapest"][0]["itineraryId"]
    cheapestFlight := respBody.Content.Results["itineraries"][itenaryId]["pricingOptions"][0]["items"][0]

    var link string
    var price string

    chvalue := reflect.ValueOf(cheapestFlight)
    for _, e := range chvalue.MapKeys() {
        key := e.Interface().(string)
        if key == "deepLink" {
            link = chvalue.MapIndex(e).Interface().(string)
        } else if key == "price" {
            temp := chvalue.MapIndex(e).Interface().(map[string]interface{})["amount"]
            price = fmt.Sprintf("%v", temp)
        }
    }

    return price, link, nil
}

// If current price has hit threshold -> notify user
func hasHitThreshold(price string, threshold string) bool 

