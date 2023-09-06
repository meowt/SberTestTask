package cMap

import (
	"sync"
	"time"
)

// CMap - custom map with concurrency support & TTL mechanism
type CMap struct {
	m  map[string]Value
	mu sync.Mutex
}

type Value struct {
	v         interface{}
	expiresAt time.Time
}

func New() *CMap {
	cm := &CMap{
		m:  map[string]Value{},
		mu: sync.Mutex{},
	}
	go func() {
		for now := range time.Tick(time.Second) {
			cm.mu.Lock()
			for i, v := range cm.m {
				if !v.expiresAt.IsZero() && v.expiresAt.After(now) {
					delete(cm.m, i)
				}
			}
			cm.mu.Unlock()
		}
	}()
	return cm
}

func (cm *CMap) Put(key string, v interface{}, ttl time.Duration) {
	value := Value{
		v: v,
	}
	if ttl != 0 {
		value.expiresAt = time.Now().Add(ttl)
	}

	cm.mu.Lock()
	cm.m[key] = value
	cm.mu.Unlock()
}

func (cm *CMap) Get(key string) (value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.m[key]
}
