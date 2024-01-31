package lib

import (
	"log"
	"os"
)

// GetEnv ...
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("Env %s value not exist \n", key)
	return ""
}
