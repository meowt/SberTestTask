package cMap

import (
	"errors"
	"fmt"

	"github.com/meowt/SberTestTask/internal/config"
	"github.com/meowt/SberTestTask/pkg/cMap"
)

type Storage struct {
	CMap *cMap.CMap
}

func Setup(cfg *config.CMapConfig) (s *Storage, err error) {
	s = &Storage{}
	s.CMap = cMap.New(cfg.DeleteTick)
	s.CMap, err = s.CMap.LoadFromFile(cfg.StorageFile)
	if err != nil && !errors.Is(err, cMap.ErrEmptyLoadFile) {
		err = fmt.Errorf("cMap setup: %w", err)
		return
	}
	err = nil

	s.CMap.Mu.Lock()
	s.CMap.Mu.Unlock()
	err = s.CMap.StoreToFile(cfg.StorageFile, cfg.StoreTick)
	if err != nil {
		err = fmt.Errorf("cMap setup: %w", err)
		return
	}

	return
}
