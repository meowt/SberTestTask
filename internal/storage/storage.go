package storage

import (
	"context"
	"time"

	"github.com/meowt/SberTestTask/internal/config"
	r "github.com/meowt/SberTestTask/internal/storage/redis"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Redis *redis.Client
}

func New(cfg *config.StorageConfig) (*Storage, error) {
	redisClient := r.Setup(&cfg.Redis)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := redisClient.Ping(ctx)
	return &Storage{Redis: redisClient}, status.Err()
}
