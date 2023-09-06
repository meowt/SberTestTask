package repo

import (
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

func (r *RepoImpl) PutObject(key, value string, ttl time.Duration) {
	r.Storage.CM.CMap.Put(key, value, ttl)
	return
}

func (r *RepoImpl) GetObject(key string) (obj models.Object) {
	value := r.Storage.CM.CMap.Get(key)
	if value != nil {
		obj.Data = value.(string)
	}
	return
}
