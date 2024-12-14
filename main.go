package main

import (
	_ "GoWeatherAPI/app"
	app2 "GoWeatherAPI/app"
	_ "GoWeatherAPI/config"
	_ "github.com/spf13/viper"
)

func main() {

	app := app2.New()
	app.Run()
}

/*client := client.NewOpenWeatherMapClient(config.AppConfig.Apikey)
p := poller.NewPoller(config.AppConfig.PollInterval * time.Second)
metric := metrics.NewMetrics()
unit := config.AppConfig.Unit*/

// use a new hook, on before serve.
/*fmt.Println("Starting UI ...")
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	zipCodeHandler(w, r, p, client, logger, unit, metric)
})

addToPoller(client, logger, "78759", unit, metric, p)

// use a new hook, on before bootstrap.
go p.StartPollingWeatherAPI()*/

// TODO need to move prometheus stuff to hooks.
/*reg := prometheus.NewRegistry()
reg.MustRegister(
	collectors.NewGoCollector(),
	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
)*/

//http.Handle("/metrics", promhttp.Handler())
//log.Fatal(http.ListenAndServe(":8081", nil))
