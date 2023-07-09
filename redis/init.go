package redis

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/108356037/torn-root-fetcher/config"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient redis.Client
)

func Init() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_URL,
		Username: "redisgo",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Cannot connect to redis", "error", err.Error())
	}

	_, err = client.FlushAll(context.Background()).Result()
	if err != nil {
		log.Fatal("Cannot flush redis on init", "error", err.Error())
	}

	RedisClient = *client
	log.Info("Successfully connected to redis, ", RedisClient.ClientID(context.Background()))
}

func Close() {
	err := RedisClient.Close()
	if err != nil {
		log.Error(err.Error())
	}
}
