package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

var AppConfig Config

type Config struct {
	Lang         string        `yaml:"language"`
	Unit         string        `yaml:"units"`
	Url          string        `yaml:"url"`
	Apikey       string        `yaml:"apikey"`
	Loglevel     string        `yaml:"loglevel"`
	PollInterval time.Duration `yaml:"poll_interval"`
	Name         string        `yaml:"name"`
	Secret       string        `yaml:"secret"`
	Locations    []string      `yaml:"locations"`
	Redis        Redis         `yaml:"redis"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
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

	viper.AutomaticEnv()
	//redisPaswrd := viper.GetString("REDIS_PASWRD")

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

}

func (c *Config) Parse() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
