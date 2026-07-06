package config

import (
	"fmt"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Koanf *koanf.Koanf
}

func NewConfig(path string) *Config {
	knf := koanf.New(".")

	if err := knf.Load(file.Provider(path), yaml.Parser()); err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	return &Config{
		Koanf: knf,
	}
}

func (c *Config) PublicHTTPPort() string {
	return c.Koanf.String("port.http.public")
}

func (c *Config) PublicHTTPTimeout() string {
	return c.Koanf.String("port.http.timeout")
}

func (c *Config) StorageType() string {
	return c.Koanf.String("storage.type")
}

func (c *Config) StorageConnectionString(storageType string) string {
	return c.Koanf.String(fmt.Sprintf("%s.connection_string", storageType))
}

func (c *Config) CoingeckoAPIKey() string {
	return c.Koanf.String("providers.coingecko.api_key")
}

func (c *Config) ActualizeInterval() time.Duration {
	t := c.Koanf.Duration("actualizerates_interval")
	if t > 0 {
		return c.Koanf.Duration("actualizerates_interval")
	}
	return time.Minute * 5
}
