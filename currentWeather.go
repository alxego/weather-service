package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jmoiron/jsonq"
)

func getCurrentTemperature(place string) (temp int, err error) {
	var client = &http.Client{Timeout: 30 * time.Second}
	q := config.currentWeatherURL.Query()
	q.Set("q", place)
	config.currentWeatherURL.RawQuery = q.Encode()

	resp, err := client.Get(config.currentWeatherURL.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data := map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&data)
	if err != nil {
		return
	}
	jq := jsonq.NewQuery(data)

	cod, err := jq.Int("cod")
	if err != nil {
		return
	} else if cod != 200 {
		err = errors.New("error while working with external API")
		return
	}

	tempF, err := jq.Float("main", "temp")
	if err != nil {
		return
	}
	temp = int(tempF - celsiusAbsoluteZero)
	return
}

func currentWeatherHandler(w http.ResponseWriter, req *http.Request) {
	city := req.URL.Query().Get("city")

	if city == "" {
		setHTTPJsonError(w, req, "Bad request", http.StatusBadRequest)
		return
	}

	temp, err := getCurrentTemperature(city)
	if err != nil {
		setHTTPJsonError(w, req, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := newCelsiusWeather(city, temp)
	respByteJSON, err := json.Marshal(resp)
	if err != nil {
		setHTTPJsonError(w, req, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(respByteJSON)
}
