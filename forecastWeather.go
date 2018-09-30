package main

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/jsonq"
)

func getForecastTemperature(place string, dt uint64) (temp int, err error) {
	now := time.Now().Unix()
	diff := int64(dt) - now

	const threeHoursInSeconds = 3.0 * 60 * 60
	if diff < threeHoursInSeconds*0.5 || diff > threeHoursInSeconds*40.5 {
		temp, err = getCurrentTemperature(place)
		return
	}
	var client = &http.Client{Timeout: 30 * time.Second}
	q := config.forecastWeatherURL.Query()
	q.Set("q", place)
	config.forecastWeatherURL.RawQuery = q.Encode()

	resp, err := client.Get(config.forecastWeatherURL.String())
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
		err = errors.New("Error while working with external API")
		return
	}

	cnt, err := jq.Int("cnt")
	if err != nil {
		return
	}
	var dtTemp int
	var tempF float64
	for i := 0; i < cnt; i++ {
		dtTemp, err = jq.Int("list", strconv.Itoa(i), "dt")
		if err != nil {
			return
		}
		if math.Abs(float64(dtTemp-int(dt))) <= threeHoursInSeconds/2 {
			tempF, err = jq.Float("list", strconv.Itoa(i), "main", "temp")
			if err != nil {
				return
			}
			temp = int(tempF - celsiusAbsoluteZero)
			return
		}
	}

	err = errors.New("Something strange happened (API die)")
	return
}

func forecastWeatherHandler(w http.ResponseWriter, req *http.Request) {
	city := req.URL.Query().Get("city")
	dtStr := req.URL.Query().Get("dt")
	dt, err := strconv.ParseUint(dtStr, 10, 64)

	if city == "" || err != nil {
		setHTTPJsonError(w, req, "Bad request", http.StatusBadRequest)
		return
	}

	temp, err := getForecastTemperature(city, dt)
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
