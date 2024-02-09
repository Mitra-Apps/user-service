package redis

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// var redisClient *redis.Client

type redisClient struct {
	client *redis.Client
}

//go:generate mockgen -source=redis.go -destination=mock/redis.go -package=mock
type RedisInterface interface {
	GetContext() context.Context
	GetStringKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

func Connection() *redisClient {
	redisServer := os.Getenv("REDIS_SERVER")
	// Initialize Redis connection
	client := redis.NewClient(&redis.Options{
		Addr:     redisServer, // Your Redis server address
		Password: "",          // No password
		DB:       0,           // Default DB
	})
	client.Context()
	return &redisClient{
		client: client,
	}
}

func (r *redisClient) GetStringKey(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisClient) GetContext() context.Context {
	return r.client.Context()
}
