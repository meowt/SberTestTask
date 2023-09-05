package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type StorageConfig struct {
	Redis RedisConfig `yaml:"redis"`
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
	return cfg, nil
}
