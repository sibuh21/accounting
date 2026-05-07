package config

import (
	"time"
)

type Config struct {
	Server struct {
		Port         string        `json:"port"`
		ReadTimeout  time.Duration `json:"read_timeout"`
		WriteTimeout time.Duration `json:"write_timeout"`
		IdleTimeout  time.Duration `json:"idle_timeout"`
	} `json:"server"`
}

func LoadConfig() *Config {
	// Hardcoded defaults for simplicity
	cfg := &Config{}
	cfg.Server.Port = "8080"
	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = 10 * time.Second
	cfg.Server.IdleTimeout = 120 * time.Second
	return cfg
}
