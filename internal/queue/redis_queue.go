package queue

import (
	"context"
	"log/slog"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisQueue struct {
	client     *redis.Client
	list       string
	set        string
	identifier string
	logger     *slog.Logger
}

func NewRedisQueue(client *redis.Client, logger *slog.Logger, set string, list string) *RedisQueue {
	id := uuid.New()

	return &RedisQueue{
		client:     client,
		list:       list,
		set:        set,
		identifier: id.String(),
		logger:     logger,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, s string) error {
	q.logger.Info("Adding to the list", "item", s, " in list ", q.list)
	_, err := q.client.LPush(ctx, q.list, s).Result()
	if err != nil {
		q.logger.Error("Failed to add to the list", "item", s, " in list ", q.list)
		return err
	}
	return nil
}

func (q *RedisQueue) Dequeue(ctx context.Context) (string, error) {
	return q.Next(ctx) // Because we really dont dequeue, we just move to the next item
}

func (q *RedisQueue) Next(ctx context.Context) (string, error) {
	return q.client.LMove(ctx, q.list, q.list, "RIGHT", "LEFT").Result()
}

func (q *RedisQueue) Delete(ctx context.Context, zip string) error {
	if	_, err := q.client.LRem(ctx, q.list, -1, zip).Result()  ; err != nil {
		return err
	}
	return nil
}
