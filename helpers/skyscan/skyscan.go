package skyscan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

type Culture struct {
    Market map[string]string `json:"market"`
    Locale map[string]string `json:"locale"`
}

func getNearestCulture(ip string) (Culture, error) {
    var culture Culture
    url := "https://partners.api.skyscanner.net/apiservices/v3/culture/nearestculture?ipAddress=" + ip
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return culture, err
    }
    req.Header.Add("x-api-key", os.Getenv("SKYSCANNER_API_KEY"))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return culture, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return culture, err
    }

    err = json.Unmarshal(body, &culture)
    if err != nil {
        return culture, err
    }

    return culture, nil
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
        price, link, err := getCheapestFlight(track.Origin, track.Destination, track.UserIp, d)
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

// return price and booking link
func getCheapestFlight(ogn string, dsn string, ip string, date date.Date) (string, string, error) {
    cl, err := getNearestCulture(ip)
    if err != nil {
        return "", "", err 
    }

    payload := map[string]map[string]interface{}{
        "query": {
            "market": cl.Market["code"],
            "locale": cl.Locale["code"],
            "currency": cl.Market["currency"],
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

    createRespData, _ := ioutil.ReadAll(res.Body)

    type CreateRespone struct {
        SessionToken string `json:"sessionToken"`
        Status string `json:"status"`
    }

    var createRespBody CreateRespone
    err = json.Unmarshal(createRespData, &createRespBody)
    if err != nil {
        return "", "", err
    }

    pollUrl :=  "https://partners.api.skyscanner.net/apiservices/v3/flights/live/search/poll/" + createRespBody.SessionToken
    pollReq, err := http.NewRequest("POST", pollUrl, nil)
    pollReq.Header.Add("Content-Type", "application/json")
    pollReq.Header.Add("x-api-key", os.Getenv("SKYSCANNER_API_KEY"))
    if err != nil {
        return "", "", err
    }

    pollRes, err := http.DefaultClient.Do(pollReq)
    if err != nil {
        return "", "", err
    }
    defer pollRes.Body.Close()

    pollRespData, _ := ioutil.ReadAll(pollRes.Body)

    type Iternary struct {
        PricingOptions []map[string]interface{} `json:"pricingOptions"`
    }

    type Results struct {
        Itenararies map[string]Iternary `json:"itineraries"`
    }

    type Content struct {
        SortingOptions map[string][]map[string]interface{} `json:"sortingOptions"`
        Results `json:"results"`
    }

    type Response struct {
        SessionToken string `json:"sessionToken"`
        Content `json:"content"`
    }

    var pollRespBody Response
    err = json.Unmarshal(pollRespData, &pollRespBody)
    if err != nil {
        return "", "", err
    }

    itineraryId := fmt.Sprintf("%v", pollRespBody.Content.SortingOptions["cheapest"][0]["itineraryId"])
    price := fmt.Sprintf("%v", pollRespBody.Content.Results.Itenararies[itineraryId].PricingOptions[0]["price"].(map[string]interface{})["amount"])
    link := fmt.Sprintf("%v", pollRespBody.Content.Results.Itenararies[itineraryId].PricingOptions[0]["items"].([]interface{})[0].(map[string]interface{})["deepLink"])

    return price, link, nil

}

func hasHitThreshold(price string, threshold string) bool {
    priceVal, _ := strconv.ParseInt(price[:len(price) - 3], 10, 64) // Price unit in mill
    thresholdVal, _ := strconv.ParseInt(threshold, 10, 64)

    return priceVal <= thresholdVal
}

