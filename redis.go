package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func startRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}

func storeMessage(roomID, message string) error {
	return rdb.LPush(ctx, roomID, message).Err()
}

func getRecentMessage(roomID string, n int64) ([]string, error) {
	return rdb.LRange(ctx, roomID, 0, n-1).Result()
}

func publishMessage(roomID, message string) error {
	return rdb.Publish(ctx, roomID, message).Err()
}

func subscribeRoom(roomID string) *redis.PubSub {
	return rdb.Subscribe(ctx, roomID)
}
