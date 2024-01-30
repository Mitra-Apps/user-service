package redis

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func Connection() *redis.Client {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	// Initialize Redis connection
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddress, // Your Redis server address
		Password: "",           // No password
		DB:       0,            // Default DB
	})
	return redisClient
}
