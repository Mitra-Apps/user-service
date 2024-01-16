package redis

import (
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func Connection() *redis.Client {
	// Initialize Redis connection
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Your Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})
	return redisClient
}
