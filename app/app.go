package app

import (
	"GoWeatherAPI/config"
	"GoWeatherAPI/internal/client"
	handlers "GoWeatherAPI/internal/handler"
	"GoWeatherAPI/internal/metrics"
	"GoWeatherAPI/internal/models"
	"GoWeatherAPI/internal/my_logger"
	"GoWeatherAPI/internal/poller"
	"GoWeatherAPI/internal/pubsub"
	"GoWeatherAPI/internal/service"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "k8s.io/client-go/rest"
)

type App struct {
	config *config.Config
}

func New() *App {
	c := &config.Config{}
	if err := c.Parse(); err != nil {
		panic(err)
	}
	return &App{config: c}
}

var zipcodes = []string{"95134"}

func (a *App) Run() int {

	unit := config.AppConfig.Unit
	ctx, cancel := context.WithCancel(context.Background())
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	logger := my_logger.SetLogLevel()
	redisPubSub := pubsub.NewRedisPubSub(redisClient, logger)
	weatherMapClient := client.NewOpenWeatherMapClient(config.AppConfig.Apikey)
	registry := prometheus.NewRegistry()
	metric := metrics.NewMetrics(registry)
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis", err)
		os.Exit(1)
	} else {
		logger.Info("Connected to Redis successfully")

	}
	redisStore := models.NewRedisStore(redisClient, logger)
	locationService := service.NewLocationService(weatherMapClient, redisStore, logger, metric, redisPubSub)
	locationService.Start(ctx)
	expireTimeDuration := config.AppConfig.PollInterval*time.Second
	pLocker := poller.NewRedisLock(redisClient, expireTimeDuration , logger, "pollerLock")
	qlocker := poller.NewRedisLock(redisClient, expireTimeDuration , logger, "queueLock")
	p := service.NewPoller(expireTimeDuration, logger, redisPubSub, qlocker, pLocker)
	p.Start(ctx)

	done := a.handleExit(logger)
	go func() {
		defer cancel()
		<-done
		logger.Info("executing cancel()")
	}()

	for _, zipcode := range a.config.Locations {
		w := models.NewLocation(zipcode)

		if err := redisStore.Create(ctx, w); err != nil {
			logger.Error("error creating location", "location", zipcode)
			return 1
		}
	}

	_ = handlers.NewLocation(redisStore, logger)
	// register new "GET /hello" route
	e := echo.New()
	e.GET("/metrics", echo.WrapHandler(promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})))

	go func() {
		err := e.Start(":8080")
		if err != nil {
			logger.Error("http server hit error", "error", err)
		}
	}()

	p.StartPollingWeatherAPI(ctx, logger)
	logger.Info("ending app")

	return 0
}

func (a *App) handleExit(logger *slog.Logger) <-chan struct{} {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)
	done := make(chan struct{})
	go func() {
		<-sig
		logger.Info("Handling exit signal")
		done <- struct{}{}
		//close(sig)
		//close(done)
	}()
	return done
}

// TODO: UI :gorilla mux web socket for UI dynamic update.

// IMPLEMENT THE GET UPDATE FUNC, CHECK ONLY FOR THE TEMP CHANGE.
