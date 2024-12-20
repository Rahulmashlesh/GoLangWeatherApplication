package models

import (
	"GoWeatherAPI/internal/pubsub"
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
	rs.logger.Debug("RedisStore:", "1 Create:", location.Zipcode)
	
	//if zipcode was already in set, 
	//	counterAdded = 0 
	//		else 
	//  counterAdded = 1
	countAdded, err := rs.addToSet(ctx, location.Zipcode)
	if err != nil {
		return err
	}
	rs.logger.Debug("RedisStore:", "2 Create:", location.Zipcode, "countAdded:", countAdded)
		
	if countAdded > 0 {
		if err := rs.publisher.Publish(ctx, "create", location.Zipcode); err != nil {
			return err
		}
			
		rs.logger.Debug("RedisStore: inside if countAdded", "X2.5 Create:", location.Zipcode, "countAdded:", countAdded)
			
		data, err := json.Marshal(location)
		if err != nil {
			rs.logger.Error("Redis Store", "Error json marshal", err)
			return err
		}
		rs.logger.Debug("RedisStore:", "X3 Create:", location.Zipcode)
		if _, err = rs.rdb.Set(ctx, location.Zipcode, data, 0).Result(); err != nil {
			rs.logger.Error("Redis Store", "Error rdn.set location", err)
			return err
		}
		rs.logger.Debug("RedisStore:", "X4 Create: trying to publish", location.Zipcode)
		if err := rs.publisher.Publish(ctx, "update", location.Zipcode); err != nil {
			return err
		
		}	
	} else {
		rs.logger.Debug("RedisStore: inside if countAdded", "2.5 Create:", location.Zipcode, "countAdded:", countAdded)
		
		data, err := json.Marshal(location)
		if err != nil {
			rs.logger.Error("Redis Store", "Error json marshal", err)
			return err
		}
		rs.logger.Debug("RedisStore:", "3 Create:", location.Zipcode)
		if _, err = rs.rdb.Set(ctx, location.Zipcode, data, 0).Result(); err != nil {
			rs.logger.Error("Redis Store", "Error rdn.set location", err)
			return err
		}
		rs.logger.Debug("RedisStore:", "4 Create: trying to publish", location.Zipcode)
		if err := rs.publisher.Publish(ctx, "update", location.Zipcode); err != nil {
			return err
		}
	}
// 
	return nil
}

func (rs *RedisStore) Update(ctx context.Context, location *Location) error {
	rs.logger.Info("RedisStore:", "Update location:" , location.Zipcode)
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
		rs.logger.Debug("RedisStore: List Func", "Appending location:", zipCode)
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
	countAdded, err := rs.rdb.SAdd (ctx, LOCATION_SET, zipcode).Result()
	if err != nil {
		return 0, err
	}
	rs.logger.Debug("RedisStore:", "addToSet:", zipcode, "countAdded:", countAdded)
	return countAdded, nil
}

func (rs *RedisStore) Notify(w *CurrentWeather) {
	err := rs.Create(context.Background(), &Location{Zipcode: w.Zipcode, Name: w.Name, Temperature: w.Main.Temp})
	if err != nil {
		rs.logger.Error("Cant insert to Redis DB","error", err)
	}
}
