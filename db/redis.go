package db

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

// NewSessionRedis provide a connection to redis as a session database
func NewSessionRedis() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("UPSTASH_URL"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Redis ready!")

	return rdb
}
