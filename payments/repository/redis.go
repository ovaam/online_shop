package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(addr string) *RedisRepository {
	return &RedisRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
		ctx: context.Background(),
	}
}

func (r *RedisRepository) SubscribeToChannel(channel string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channel)
}

func (r *RedisRepository) AddToInbox(queue string, message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.client.LPush(r.ctx, queue, msg).Err()
}

func (r *RedisRepository) ProcessInbox(queue string, processFunc func(string) error) {
	for {
		msg, err := r.client.RPop(r.ctx, queue).Result()
		if err == nil {
			if err := processFunc(msg); err != nil {
				r.client.RPush(r.ctx, queue, msg)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *RedisRepository) PublishResult(channel string, result interface{}) error {
	msg, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return r.client.Publish(r.ctx, channel, msg).Err()
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}
