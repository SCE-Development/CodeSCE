package redisclient

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"os"
)

var client *redis.Client

func Init() error {
	client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis connect is fail: %w", err)
	}

	log.Printf("connected to redis at %s", os.Getenv("REDIS_URL"))
	return nil
}

func GetClient() *redis.Client {
	return client
}

func Close() {
	if client != nil {
		client.Close()
	}
}
