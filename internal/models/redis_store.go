package models

import (
	"GoWeatherAPI/internal/pubsub"
	"GoWeatherAPI/internal/weather"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	rdb       *redis.Client
	logger    *slog.Logger
	publisher pubsub.Publisher[string]
}

func NewRedisStore(redisClient *redis.Client, logger *slog.Logger, publisher pubsub.Publisher[string]) *RedisStore {
	return &RedisStore{
		rdb:       redisClient,
		logger:    logger.With("context", "redis_store"),
		publisher: publisher,
	}
}

var LOCATION_SET = "location_set"

func (rs *RedisStore) Get(ctx context.Context, zipcode string) (*Location, error) {
	result, err := rs.rdb.Get(ctx, zipcode).Result()
	if err == redis.Nil {
		return nil, err // Return nil error on key not found
	} else if err != nil {
		return nil, err // Return error if any other error occurred
	}
	loc := &Location{}
	err = json.Unmarshal([]byte(result), loc)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

func (rs *RedisStore) Create(ctx context.Context, location *Location) error {
	// Create a hash with location and temperature
	countAdded, err := rs.addToSet(ctx, location.Zipcode)
	if err != nil {
		return err
	}

	if countAdded > 0 {
		data, err := json.Marshal(location)
		if err != nil {
			return err
		}
		if _, err = rs.rdb.Set(ctx, location.Zipcode, data, 0).Result(); err != nil {
			return err
		}

		if err := rs.publisher.Publish(ctx, "create", location.Zipcode); err != nil {
			return err
		}
	}

	return nil
}

func (rs *RedisStore) Update(ctx context.Context, location *Location) error {
	return rs.Create(ctx, location)
}

func (rs *RedisStore) List(ctx context.Context) ([]Location, error) {
	zipCodes, err := rs.rdb.SMembers(ctx, LOCATION_SET).Result()
	var locations []Location
	if err != nil {
		return nil, err
	}
	for _, zipCode := range zipCodes {
		loc, err := rs.Get(ctx, zipCode)
		if err != nil {
			return nil, err
		}
		locations = append(locations, *loc)
	}
	return locations, nil
}

func (rs *RedisStore) Delete(ctx context.Context, zipcode string) error {
	result, err := rs.rdb.SRem(ctx, LOCATION_SET, zipcode).Result()
	if err != nil {
		return err
	}
	if result > 0 {
		_, err := rs.rdb.Del(ctx, zipcode).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs *RedisStore) addToSet(ctx context.Context, zipcode string) (int64, error) {
	countAdded, err := rs.rdb.SAdd(ctx, LOCATION_SET, zipcode).Result()
	if err != nil {
		return 0, err
	}
	return countAdded, nil
}

func (rs *RedisStore) Notify(w *weather.CurrentWeather) {
	err := rs.Create(context.Background(), &Location{Zipcode: w.Zipcode, Name: w.Name, Temperature: w.Main.Temp})
	if err != nil {
		rs.logger.Error("Cant insert to Redis DB", err)
	}
}

func (rs *RedisStore) Call(ctx context.Context, zip string, weatherClient weather.CurrentWeather) {
	// Retrieve location data from Redis
	locationData, err := rs.Get(ctx, zip)
	if err != nil {
		if err == redis.Nil {
			rs.logger.Error("No data found in Redis for zipcode", "zipcode", zip)
		} else {
			rs.logger.Error("Error retrieving location data from Redis", "error", err)
		}
		return
	}

	// Translate location data to CurrentWeather
	currentWeather := &weather.CurrentWeather{
		Zipcode: zip,
		Name:    locationData.Name,
		Main: weather.Main{
			Temp: locationData.Temperature,
		},
		Client:  weatherClient.Client,
		Logger:  rs.logger,
		Metrics: weatherClient.Metrics,
	}

	// Fetch fresh weather data from the weather API
	hasChanged, err := currentWeather.GetWeather(ctx)
	if err != nil {
		rs.logger.Error("Error fetching weather data from API", "error", err)
		return
	}

	// Update Redis if temperature or other data has changed
	if hasChanged {
		rs.logger.Info("Weather data has changed, updating Redis", "zipcode", zip)
		err := rs.Update(ctx, &Location{
			Zipcode:     zip,
			Name:        currentWeather.Name,
			Temperature: currentWeather.Main.Temp,
		})
		if err != nil {
			rs.logger.Error("Error updating Redis with new weather data", "error", err)
		}
	} else {
		rs.logger.Debug("No changes in weather data for zipcode", "zipcode", zip)
	}

}
