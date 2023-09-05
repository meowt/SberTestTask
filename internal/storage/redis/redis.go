package redis

import (
	"github.com/meowt/SberTestTask/internal/config"
	"github.com/redis/go-redis/v9"
)

func Setup(cfg *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
	})
}
