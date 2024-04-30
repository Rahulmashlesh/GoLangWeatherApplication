package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

var AppConfig Config

type Config struct {
	Lang         string
	Unit         string
	Url          string
	Apikey       string
	Loglevel     string
	PollInterval time.Duration
}

func init() {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/GoWeatherAPI/config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Can add defaults if needed.
	//viper.SetDefault("lang", "EN")
	viper.SetDefault("unit", "imperial")
	viper.SetDefault("pollInterval", "20")

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

}
