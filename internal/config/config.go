package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type CMapConfig struct {
	StoragePath string `yaml:"storagePath"`
	StorageFile *os.File
	DeleteTick  time.Duration `yaml:"deleteTick"`
	StoreTick   time.Duration `yaml:"storeTick"`
}

type StorageConfig struct {
	//Redis RedisConfig `yaml:"redis"`
	CMap CMapConfig `yaml:"cMap"`
}

type RouterConfig struct {
	Address string
}

type Config struct {
	Storage StorageConfig `yaml:"storage"`
	Router  RouterConfig  `yaml:"router"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	cfgFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config file reading error: %w", err)
	}

	err = yaml.Unmarshal(cfgFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("config file unmarshalling error: %w", err)
	}
	cfg.Storage.CMap.StorageFile, err = os.OpenFile(cfg.Storage.CMap.StoragePath, os.O_RDWR|os.O_CREATE, 0755)
	return cfg, err
}
