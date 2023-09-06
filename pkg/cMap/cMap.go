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

var (
	ErrEmptyLoadFile = errors.New("empty load-file")
)

// CMap - custom map with concurrency support & TTL mechanism
type CMap struct {
	M  map[string]Value
	Mu sync.Mutex
}

type Value struct {
	V         interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func New(deleteTick time.Duration) *CMap {
	cm := &CMap{
		M:  map[string]Value{},
		Mu: sync.Mutex{},
	}
	go func() {
		for now := range time.Tick(deleteTick) {
			cm.Mu.Lock()
			for i, v := range cm.M {
				if !v.ExpiresAt.IsZero() && now.After(v.ExpiresAt) {
					delete(cm.M, i)
				}
			}
			cm.Mu.Unlock()
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

	cm.Mu.Lock()
	cm.M[key] = value
	cm.Mu.Unlock()
}

func (cm *CMap) Get(key string) (value interface{}) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	return cm.M[key].V
}

func (cm *CMap) StoreToFile(f *os.File, interval time.Duration) (err error) {
	var mapData []byte
	go func() {
		for range time.Tick(interval) {
			cm.Mu.Lock()
			if mapData, err = json.Marshal(cm.M); err != nil {
				log.Println(err)
				return
			}
			cm.Mu.Unlock()

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
		return nil, ErrEmptyLoadFile
	}

	//Unmarshalling json into map
	cm.Mu.Lock()
	if err = json.Unmarshal(jsonData, &cm.M); err != nil {
		return nil, err
	}
	cm.Mu.Unlock()

	return cm, nil
}
