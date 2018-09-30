package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := initConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/v1/current/", currentWeatherHandler)
	serverMux.HandleFunc("/v1/forecast/", forecastWeatherHandler)

	err = http.ListenAndServe(":"+os.Getenv("LISTEN_PORT"), serverMux)
	fmt.Println(err.Error())
}
