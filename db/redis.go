package db

import (
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

// NewSessionRedis provide a connection to redis as a session database
func NewSessionRedis() *redis.Client {
	options := &redis.Options{
		Addr:        os.Getenv("REDIS_URL"),
		Password:    "",
		DB:          0,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
	}

	return redis.NewClient(options)
}
