package storage

import (
	"github.com/meowt/SberTestTask/internal/config"
	"github.com/meowt/SberTestTask/internal/storage/cMap"
)

type Storage struct {
	//Redis *redis.Client
	CM *cMap.Storage
}

func New(cfg *config.StorageConfig) (*Storage, error) {
	//redisClient := r.Setup(&cfg.Redis)
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//
	//status := redisClient.Ping(ctx)
	//return &Storage{Redis: redisClient}, status.Err()

	cm, err := cMap.Setup(&cfg.CMap)
	return &Storage{CM: cm}, err
}
