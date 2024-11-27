package service

import (
	"GoWeatherAPI/internal/client"
	"GoWeatherAPI/internal/metrics"
	"GoWeatherAPI/internal/models"
	"GoWeatherAPI/internal/pubsub"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type LocationService struct {
	client     client.HttpGetter
	model      models.Model
	logger     *slog.Logger
	metrics    metrics.Metrics
	subscriber pubsub.Subscriber[string]
}

func (ls *LocationService) Start(ctx context.Context) {
	// subscribe to the channel
	dataChan := ls.subscriber.Subscribe(ctx, "call")

	go func() {
		for {
			// we want to select the data channel

			select {
			case <-ctx.Done():
				ls.logger.Info("Received cancellation signal")
				return
			case data := <-dataChan:
				ls.logger.Info("Received message", "data", data.Data())
				location, err := ls.model.Get(ctx, data.Data())

				if err != nil {
					ls.logger.Error("Error getting location", "error", err)
					return
				}
				changed, err := ls.getWeather(ctx, location)
				if err != nil {
					ls.logger.Error("Error getting weather", "error", err)
				}
				if changed {
					if err := ls.model.Update(ctx, location); err != nil {
						ls.logger.Error("Error updating location", "error", err)
					}
				}
			}
		}
	}()

}

func (ls *LocationService) getWeather(ctx context.Context, location *models.Location) (bool, error) {
	rsp, err := ls.client.Get(location.Zipcode)
	changed := false
	if err != nil {
		ls.logger.Error("Error during HTTP GET req:", err)
		return changed, errors.New("Http Client Error")
	}
	ls.logger.Debug("Processing", " Zipcode:", location.Zipcode)

	if rsp.StatusCode != http.StatusOK {
		// Handle non-OK status codes
		ls.logger.Error("Non-OK HTTP status code", "StatusCode", rsp.StatusCode)
		return changed, errors.New("Http Client Error") // Change this line to return "Http Client Error"

	}

	newCurrentWeather := &models.CurrentWeather{}
	err = json.NewDecoder(rsp.Body).Decode(newCurrentWeather)
	if err != nil {
		ls.logger.Error("Error Decoding Json", err)
		return changed, errors.New("json Decoding Error")
	}
	if newCurrentWeather.Main.Temp != location.Temperature {
		// Update only when there is a change in Temp.
		changed = true
		location.Name = newCurrentWeather.Name
		location.Temperature = newCurrentWeather.Main.Temp
	}
	ls.logger.Debug("API Response:", "Temperature", location.Temperature, "City", location.Name)
	ls.metrics.TempGage.WithLabelValues(location.Name, location.Zipcode).Set(location.Temperature)

	return changed, nil
}

func NewLocationService(client client.HttpGetter, model models.Model, logger *slog.Logger, metrics metrics.Metrics, subscriber pubsub.Subscriber[string]) *LocationService {

	return &LocationService{
		subscriber: subscriber,
		client:     client,
		model:      model,
		logger:     logger,
		metrics:    metrics,
	}
}
