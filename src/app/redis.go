package app

import (
	"github.com/gofiber/storage/redis"
)

var RedisStorage *redis.Storage

func (a *App) RedisStorage() *redis.Storage {
	if RedisStorage != nil {
		return RedisStorage
	}

	// setup redis based storage for cache, rate limiter
	RedisStorage = redis.New(redis.Config{
		Host:     "redis-18660.c1.asia-northeast1-1.gce.cloud.redislabs.com",
		Port:     18660,
		Password: "lS6uZuYFeDfMzIHKpEWb6bhHrReutUTX",
		Database: 0,
		Reset:    false,
	})

	return RedisStorage
}
