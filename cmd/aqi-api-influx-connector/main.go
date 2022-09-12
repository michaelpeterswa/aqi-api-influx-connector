package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaelpeterswa/aqi-api-influx-connector/internal/handlers"
	"github.com/michaelpeterswa/aqi-api-influx-connector/internal/influx"
	"github.com/michaelpeterswa/aqi-api-influx-connector/internal/logging"
	"github.com/michaelpeterswa/aqi-api-influx-connector/internal/requests"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.InitZap()
	if err != nil {
		log.Panicf("could not acquire zap logger: %s", err.Error())
	}
	logger.Info("aqi-api-influx-connector init...")

	influxEndpoint := os.Getenv("INFLUX_ENDPOINT")
	if influxEndpoint == "" {
		logger.Fatal("INFLUX_ENDPOINT not set")
	}

	influxToken := os.Getenv("INFLUX_TOKEN")
	if influxToken == "" {
		logger.Fatal("INFLUX_TOKEN not set")
	}

	apiEndpoint := os.Getenv("API_ENDPOINT")
	if apiEndpoint == "" {
		logger.Fatal("API_ENDPOINT not set")
	}

	requestsClient := requests.NewRequestClient()

	influxClient := influx.InitInflux(influxEndpoint, influxToken)
	defer influxClient.Close()

	requestTicker := time.NewTicker(time.Second * 30)
	go func() {
		for range requestTicker.C {
			aqi, err := requestsClient.GetCurrentAQI(apiEndpoint)
			if err != nil {
				logger.Error("could not get current AQI", zap.Error(err))
				continue
			}
			if aqi.PrimaryPollutant == "" {
				logger.Error("primary pollutant is empty", zap.Any("aqi", aqi))
				continue
			}
			influxClient.WriteAQIPoint(aqi)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", handlers.HealthcheckHandler)
	r.Handle("/metrics", promhttp.Handler())
	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal("could not start http server", zap.Error(err))
	}
}
