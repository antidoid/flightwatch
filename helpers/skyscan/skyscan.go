package skyscan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"

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

    // If today > endDate => notify the user that flight never hit threshold and delte from db
    if (date.TodayUTC().Sub(endDate) > 0) {
        message := fmt.Sprintf("\nGreeting from FlightWatch\n, This is to inform you that your tracked flight from %s to %s never went below %s",
            track.Origin, track.Destination, track.Threshold)
        notify.SendSMS(track.Contact, message)
        tx := initializers.DB.Unscoped().Delete(&track)
        return tx.Error
    }

    for d := startDate; d.Sub(endDate) <= 0; d = d.Add(1) {
        // Check if price has reached threshold
        price, link, err := getCheapestFlight(track.Origin, track.Destination, d)
        if err != nil {
            return err
        }

        if (hasHitThreshold(price, track.Threshold)) {
            message := fmt.Sprintf("\nGreeting from FlightWatch\nYour tracked flight from %s to %s on %s is currently priced at Rs%s\n Book now at: %s\nHave a nice day",
                track.Origin, track.Destination, d, price, link)
            notify.SendSMS(track.Contact, message)
            tx := initializers.DB.Unscoped().Delete(&track)
            return tx.Error
        }
    }
    return nil
}

func ScanAllTracks() error {
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

    createUrl := "https://partners.api.skyscanner.net/apiservices/v3/flights/live/search/create"
    createReq, err := http.NewRequest("POST", createUrl, bytes.NewBuffer(postBody))
    createReq.Header.Add("Content-Type", "application/json")
    createReq.Header.Add("x-api-key", os.Getenv("SKYSCANNER_API_KEY"))
    if err != nil {
        return "", "", err
    }

    res, err := http.DefaultClient.Do(createReq)
    if err != nil {
        return "", "", err
    }
    defer res.Body.Close()

    createResponseData, _ := ioutil.ReadAll(res.Body)

    type CreateRespone struct {
        SessionToken string `json:"sessionToken"`
        Status string `json:"status"`
    }

    var createRespBody CreateRespone
    err = json.Unmarshal(createResponseData, &createRespBody)
    if err != nil {
        return "", "", err
    }

    pollUrl :=  "https://partners.api.skyscanner.net/apiservices/v3/flights/live/search/poll/" + createRespBody.SessionToken
    pollReq, err := http.NewRequest("POST", pollUrl, bytes.NewBuffer(postBody))
    pollReq.Header.Add("Content-Type", "application/json")
    pollReq.Header.Add("x-api-key", os.Getenv("SKYSCANNER_API_KEY"))
    if err != nil {
        return "", "", err
    }

    pollRes, err := http.DefaultClient.Do(pollReq)
    if err != nil {
        return "", "", err
    }
    defer res.Body.Close()

    pollRespData, _ := ioutil.ReadAll(pollRes.Body)

    type Content struct {
        Results map[string]map[string]map[string][]map[string][]map[string]interface{} `json:"results"`
        SortingOptions map[string][]map[string]string `json:"sortingOptions"`
    }

    type Response struct {
        SessionToken string `json:"sessionToken"`
        Status string `json:"status"`
        Content `json:"content"`
    }

    var pollRespBody Response
    err = json.Unmarshal(pollRespData, &pollRespBody)
    if err != nil {
        return "", "", err
    }

    itenaryId := pollRespBody.Content.SortingOptions["cheapest"][0]["itineraryId"]
    cheapestFlight := pollRespBody.Content.Results["itineraries"][itenaryId]["pricingOptions"][0]["items"][0]

    var link string
    var price string

    cheapestFlightValue := reflect.ValueOf(cheapestFlight)
    for _, e := range cheapestFlightValue.MapKeys() {
        key := e.Interface().(string)
        if key == "deepLink" {
            link = cheapestFlightValue.MapIndex(e).Interface().(string)
        } else if key == "price" {
            temp := cheapestFlightValue.MapIndex(e).Interface().(map[string]interface{})["amount"]
            price = fmt.Sprintf("%v", temp)
        }
    }

    return price, link, nil
}

// Unit of price is in MILLI so in DB threshold should also be in MILLI
func hasHitThreshold(price string, threshold string) bool {
    priceVal, _ := strconv.ParseInt(price[:len(price) - 3], 10, 64)
    thresholdVal, _ := strconv.ParseInt(threshold, 10, 64)

    return priceVal <= thresholdVal
}

