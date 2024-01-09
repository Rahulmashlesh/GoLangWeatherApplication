package main

import (
	"GoWeaterAPI/config"
	_ "GoWeaterAPI/config"
	"GoWeaterAPI/internal/client"
	"GoWeaterAPI/internal/poller"
	"GoWeaterAPI/internal/weather"
	"fmt"
	_ "github.com/spf13/viper"
	"log/slog"
	"os"
	"time"
)

func main() {
	//TODO:  Create a new poller with a poll period of 1 seconds

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	p := poller.NewPoller(1 * time.Second)

	fmt.Printf("apiKey: ", config.AppConfig.Apikey)

	client := client.NewOpenWeatherMapClient(config.AppConfig.Apikey)

	p.Add(weather.NewCurrentWeather(client, logger, "95134"))
	p.Add(weather.NewCurrentWeather(client, logger, "78759"))
	p.Add(weather.NewCurrentWeather(client, logger, "11213"))

	// StartPollingWeatherAPI the poller
	go p.StartPollingWeatherAPI()

	// sleep just to see the code in action,
	// Should this be indefinite sleep?
	time.Sleep(5 * time.Second)
}
