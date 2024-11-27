package my_logger

import (
	"GoWeatherAPI/config"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func SetLogLevel() *slog.Logger {
	logLevel := new(slog.LevelVar)
	switch strings.ToLower(config.AppConfig.Loglevel) {
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "error":
		logLevel.Set(slog.LevelError)
	case "warning":
		logLevel.Set(slog.LevelWarn)
	default:
		fmt.Println("Error: Invalid Log level: ", config.AppConfig.Loglevel)
		os.Exit(1)
	}
	fmt.Println("Logger level: ", config.AppConfig.Loglevel)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	return logger
}
