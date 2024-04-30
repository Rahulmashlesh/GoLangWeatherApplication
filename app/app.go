package app

import (
	"GoWeatherAPI/config"
	"GoWeatherAPI/internal/client"
	"GoWeatherAPI/internal/poller"
	"GoWeatherAPI/internal/weather"
	"GoWeatherAPI/metrics"
	"GoWeatherAPI/model"
	"fmt"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"time"
)

type App struct {
}

var zipcodes = []string{"95131"}

func (a *App) Run() {

	pocketBaseApp := pocketbase.New()
	weatherMapClient := client.NewOpenWeatherMapClient(config.AppConfig.Apikey)
	p := poller.NewPoller(config.AppConfig.PollInterval * time.Second)
	metric := metrics.NewMetrics()
	unit := config.AppConfig.Unit

	// Starting PocketBase
	pocketBaseApp.OnAfterBootstrap().Add(
		func(_ *core.BootstrapEvent) error {
			for _, zip := range zipcodes {
				dbLocations := model.NewLocations(pocketBaseApp.Dao())
				currentWeather := weather.NewCurrentWeather(weatherMapClient, pocketBaseApp.Logger(), zip, unit, metric, dbLocations)
				p.Add(currentWeather)
				fmt.Println("temp:", currentWeather.Main.Temp)
				fmt.Println("zip: ", zip)
				err := dbLocations.Create(&model.Location{Zipcode: zip})
				if err != nil {
					return err
				}

			}
			fmt.Println("staring ticking !!!")
			go p.StartPollingWeatherAPI()
			return nil
		})
	fmt.Println("starting DB...")
	if err := pocketBaseApp.Start(); err != nil {
		log.Fatal(err)
	}

}
