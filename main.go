package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func servePrometheusMetrics() {
	serverMux := http.NewServeMux()

	serverMux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", serverMux)
}

func main() {
	err := initConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/v1/current/", currentWeatherHandler)
	serverMux.HandleFunc("/v1/forecast/", forecastWeatherHandler)

	go servePrometheusMetrics()

	err = http.ListenAndServe(":"+os.Getenv("LISTEN_PORT"), serverMux)
	fmt.Println(err.Error())
}
