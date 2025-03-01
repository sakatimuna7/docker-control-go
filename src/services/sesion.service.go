package services

import (
	"docker-control-go/src/configs"
	"time"

	"context"
)

var (
	ctx = context.Background()
)

// Store session in Redis with expiration
func StoreSessionInRedis(userID string, token string, duration time.Duration) error {
	expiration := time.Now().Add(duration)
	return configs.RedisClient.Set(ctx, "session:"+token, userID, time.Until(expiration)).Err()
}

// Delete session from Redis
func DeleteSessionFromRedis(token string) error {
	key := "session:" + token
	return configs.RedisClient.Del(ctx, key).Err()
}
