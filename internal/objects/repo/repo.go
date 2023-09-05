package repo

import (
	"context"
	"time"

	"github.com/meowt/SberTestTask/internal/storage"
	models "github.com/meowt/SberTestTask/pkg/models/objects"
)

type RepoImpl struct {
	Storage *storage.Storage
}

func SetupObjectsRepo(s *storage.Storage) *RepoImpl {
	return &RepoImpl{
		Storage: s,
	}
}

func (r *RepoImpl) PutObject(key, value string, ttl time.Duration) (res string, err error) {
	status := r.Storage.Redis.Set(context.Background(), key, value, ttl)
	return status.Result()
}

func (r *RepoImpl) GetObject(key string) (obj models.Object, err error) {
	res := r.Storage.Redis.Get(context.Background(), key)
	obj.Data, err = res.Result()
	return
}
