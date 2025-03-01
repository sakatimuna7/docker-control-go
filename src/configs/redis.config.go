package configs

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB")) // Konversi string ke int
	if err != nil {
		panic("Invalid REDIS_DB value: must be an integer")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	err = RedisClient.Ping(context.Background()).Err()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	fmt.Println("✅ Connected to Redis")
}
