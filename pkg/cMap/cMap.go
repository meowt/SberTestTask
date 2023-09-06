package cMap

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// CMap - custom map with concurrency support & TTL mechanism
type CMap struct {
	m  map[string]Value
	mu sync.Mutex
}

type Value struct {
	V         interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func New(d time.Duration) *CMap {
	cm := &CMap{
		m:  map[string]Value{},
		mu: sync.Mutex{},
	}
	go func() {
		for now := range time.Tick(d) {
			cm.mu.Lock()
			for i, v := range cm.m {
				if !v.ExpiresAt.IsZero() && now.After(v.ExpiresAt) {
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
		V: v,
	}
	if ttl != 0 {
		value.ExpiresAt = time.Now().Add(ttl)
	}

	cm.mu.Lock()
	cm.m[key] = value
	cm.mu.Unlock()
}

func (cm *CMap) Get(key string) (value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.m[key].V
}

func (cm *CMap) StoreToFile(f *os.File, interval time.Duration) (err error) {
	var mapData []byte
	go func() {
		for range time.Tick(interval) {
			cm.mu.Lock()
			if mapData, err = json.Marshal(cm.m); err != nil {
				log.Println(err)
				return
			}
			cm.mu.Unlock()

			//Check of file's openness
			_, err = f.Read([]byte{})
			if err != nil && err != io.EOF {
				break
			}

			//Truncate old data
			if err = f.Truncate(0); err != nil {
				log.Println(err)
				return
			}

			//Move file's cursor to 0 byte
			if _, err = f.Seek(0, 0); err != nil {
				log.Println(err)
				return
			}

			//Write map's data to file
			if _, err = f.Write(mapData); err != nil {
				log.Println(err)
				return
			}

			//Guarantee that data is on disk
			if err = f.Sync(); err != nil {
				log.Println(err)
				return
			}
		}
	}()
	return
}

func (cm *CMap) LoadFromFile(f *os.File) (*CMap, error) {
	jsonData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	//Empty file check
	if string(jsonData) == "" {
		return nil, errors.New("empty load-file")
	}

	//Unmarshalling json into map
	cm.mu.Lock()
	if err = json.Unmarshal(jsonData, &cm.m); err != nil {
		return nil, err
	}
	cm.mu.Unlock()

	return cm, nil
}
