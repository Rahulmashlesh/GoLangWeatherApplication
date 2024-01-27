package main

import (
	"GoWeaterAPI/config"
	_ "GoWeaterAPI/config"
	"GoWeaterAPI/internal/client"
	"GoWeaterAPI/internal/poller"
	"GoWeaterAPI/internal/weather"
	"GoWeaterAPI/metrics"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/spf13/viper"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	//TODO:  Create a new poller with a poll period of 1 seconds

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	p := poller.NewPoller(1 * time.Second)

	fmt.Printf("apiKey: ", config.AppConfig.Apikey)

	client := client.NewOpenWeatherMapClient(config.AppConfig.Apikey)
	metric := metrics.NewMetrics()

	p.Add(weather.NewCurrentWeather(client, logger, "95134", metric))
	p.Add(weather.NewCurrentWeather(client, logger, "78759", metric))
	p.Add(weather.NewCurrentWeather(client, logger, "11213", metric))

	// StartPollingWeatherAPI the poller
	go p.StartPollingWeatherAPI()

	// sleep just to see the code in action,
	// Should this be indefinite sleep?
	time.Sleep(5 * time.Second)

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	//http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))

}
