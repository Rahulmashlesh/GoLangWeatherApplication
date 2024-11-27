package poller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/go-redis/redis/v8"
)

// RedisLock structure
type RedisLock struct {
	client    *redis.Client
	expiry    time.Duration // Expiration time for the lock
	id        uuid.UUID
	logger    *slog.Logger
	lockerKey string
}

// NewRedisLock creates a new RedisLock
func NewRedisLock(client *redis.Client, expiry time.Duration, logger *slog.Logger, lockerKey string) *RedisLock {
	return &RedisLock{
		client:    client,
		expiry:    expiry,
		id:        uuid.New(),
		logger:    logger.With("context", "redis_lock"),
		lockerKey: lockerKey,
	}
}

// Lock tries to acquire a lock
func (l *RedisLock) Lock(ctx context.Context) bool {
	// Use SETNX to attempt to acquire the lock with an expiry time

	lock, err := l.client.Get(ctx, l.lockerKey).Result()
	if errors.Is(err, redis.Nil) {
		success, err := l.client.SetNX(ctx, l.lockerKey, l.id, l.expiry).Result()
		if err != nil {
			l.logger.Error("Error while trying to acquire lock", "error", err)
			return false
		}

		// If success is true, the lock was acquired
		if success {
			l.logger.Info("Lock acquired with key", "id", l.id)
			return true
		}
	} else if err != nil {
		l.logger.Error("Error while trying to acquire lock", "error", err)
		return false
	}
	if lock == l.id.String() {
		_, err = l.client.Expire(ctx, l.lockerKey, l.expiry).Result()
		if err != nil {
			l.logger.Error("Error while trying to set expiration time:", "error", err)
			return false
		}
		l.logger.Info("Lock acquired with key", "id", l.id)
		return true
	}

	l.logger.Info("Failed to acquire lock, key already held", "id", lock)
	return false
}

// dont need this
//func (l *RedisLock) Unlock(ctx context.Context) error {
//	// Delete the lock key to release the lock
//	_, err := l.client.Del(ctx, l.lockerKey).Result()
//	if err != nil {
//		l.logger.Error("Error while trying to release lock", "error", err)
//		return err
//	}
//	return nil
//}
