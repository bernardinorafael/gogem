package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/bernardinorafael/gogem/pkg/fault"
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

type SetParams struct {
	Client *Client
	Key    string
	TTL    time.Duration
}

type Client struct {
	redis *redis.Client
	log   *log.Logger
}

func New(redis *redis.Client, log *log.Logger) *Client {
	return &Client{
		redis: redis,
		log:   log,
	}
}

func GetOrSet[T any](ctx context.Context, params SetParams, callback func() (T, error)) (T, error) {
	var zero T

	cached, err := get[T](ctx, params.Client, params.Key)
	if err == nil {
		return cached, nil
	}

	value, err := callback()
	if err != nil {
		return zero, err
	}

	if err = set(ctx, params.Client, params.Key, value, params.TTL); err != nil {
		params.Client.log.Warn("failed to save to cache", "key", params.Key, "error", err)
	}

	return value, nil
}

func Delete(ctx context.Context, c *Client, keys ...string) error {
	if err := c.redis.Del(ctx, keys...).Err(); err != nil {
		return fault.New("failed to delete from cache", fault.WithTag(fault.DB), fault.WithErr(err))
	}

	return nil
}

func get[T any](ctx context.Context, c *Client, key string) (T, error) {
	var zero T

	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return zero, fault.New("key not found", fault.WithTag(fault.NotFound))
		}
		return zero, fault.New("failed to get from cache", fault.WithTag(fault.DB), fault.WithErr(err))
	}

	var value T
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return zero, fault.New("failed to deserialize cached value", fault.WithTag(fault.DB), fault.WithErr(err))
	}

	return value, nil
}

func set[T any](ctx context.Context, c *Client, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fault.New("failed to serialize value", fault.WithTag(fault.DB), fault.WithErr(err))
	}

	if err := c.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return fault.New("failed to set cache", fault.WithTag(fault.DB), fault.WithErr(err))
	}

	return nil
}
