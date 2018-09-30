package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

const celsiusAbsoluteZero = 273.15

type configURLs struct {
	currentWeatherURL  *url.URL
	forecastWeatherURL *url.URL
}

var config configURLs

func initConfig() (err error) {

	file, err := os.Open("config.json")
	if err != nil {
		return
	}

	configMap := map[string]interface{}{}
	dec := json.NewDecoder(file)
	err = dec.Decode(&configMap)
	if err != nil {
		return
	}

	currentStr, currentOk := configMap["current"].(string)
	forecastStr, forecastOk := configMap["forecast"].(string)
	if !currentOk || !forecastOk {
		err = errors.New("config file error")
	}

	config.currentWeatherURL, err = url.Parse(currentStr)
	if err != nil {
		return
	}

	config.forecastWeatherURL, err = url.Parse(forecastStr)
	if err != nil {
		return
	}

	return
}

// Weather response JSON structure
type weather struct {
	City        string `json:"city,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Temperature int    `json:"temperature,omitempty"`
}

// Weather response constructor for temperature in celsius
func newCelsiusWeather(city string, temperature int) (w weather) {
	return weather{City: city, Unit: "celsius", Temperature: temperature}
}

// Set HTTP status code and write JSON with error field
func setHTTPJsonError(w http.ResponseWriter, req *http.Request, errorMsg string, errorHTTPStatusCode int) {
	w.WriteHeader(errorHTTPStatusCode)

	errorMap := make(map[string]interface{})
	errorMap["error"] = errorMsg
	respByteJSON, err := json.Marshal(errorMap)
	if err != nil {
		return
	}
	w.Write(respByteJSON)
}
