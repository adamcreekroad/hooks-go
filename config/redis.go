package config

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisConn *redis.Client

func configureRedis() {
	options, err := redis.ParseURL(url())

	if err != nil {
		panic(err)
	}

	RedisConn = redis.NewClient(options)
}

func url() string {
	url, present := os.LookupEnv("REDIS_URL")

	if !present {
		url = "redis://localhost:6379/0"
	}

	return url
}
