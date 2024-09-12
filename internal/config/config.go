package config

import (
	"fmt"
	"os"
)

type Config struct {
	DebugMode bool
	RedisCfg  RedisConfig
}

type RedisConfig struct {
	Password string
	Address  string
	DB       int
}

func New() Config {
	return Config{
		DebugMode: os.Getenv("DEBUG_MODE") == "true",
		RedisCfg: RedisConfig{
			Password: os.Getenv("REDIS_PASSWORD"),
			Address:  fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			DB:       0,
		},
	}
}
