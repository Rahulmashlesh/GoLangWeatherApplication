package config

import (
	"github.com/spf13/viper"
	"log"
)

var AppConfig Config

type Config struct {
	Lang   string
	Unit   string
	Url    string
	Apikey string
}

func init() {
	viper.SetConfigFile("/Users/rahul.ashlesh/GolandProjects/GoWeaterAPI/config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Can add defaults if needed.
	//viper.SetDefault("lang", "EN")
	//viper.SetDefault("unit", "F")

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

}
