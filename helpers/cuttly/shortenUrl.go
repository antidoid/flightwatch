package cuttly

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func GetShortUrl(url string) (string, error) {
	reqUrl := "https://cutt.ly/api/api.php?key=" + os.Getenv("CUTTLY_API_KEY") + "&short=" + url
	res, err := http.Get(reqUrl)
	if err != nil {
		return "", err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	type Response struct {
		Url map[string]string `json:"url"`
	}
	var resBody Response
	err = json.Unmarshal(resData, &resBody)

	return resBody.Url["shortLink"], nil

}
