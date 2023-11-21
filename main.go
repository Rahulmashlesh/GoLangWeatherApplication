package main

import (
	"GoWeaterAPI/internal/poller"
	"GoWeaterAPI/internal/weather"
	"log/slog"
	"os"
	"time"
)

func main() {
	//TODO:  Create a new poller with a poll period of 1 seconds

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	p := poller.NewPoller(1 * time.Second)

	p.Add(weather.NewCurrentWeather(logger, "95134"))
	p.Add(weather.NewCurrentWeather(logger, "78759"))
	p.Add(weather.NewCurrentWeather(logger, "11213"))

	// StartPollingWeatherAPI the poller
	go p.StartPollingWeatherAPI()

	// sleep just to see the code in action,
	// Should this be indefinite sleep?
	time.Sleep(5 * time.Second)
}
