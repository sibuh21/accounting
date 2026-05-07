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
	Database struct {
		URL string `json:"url"`
	} `json:"database"`
}

func LoadConfig() *Config {
	// ... defaults ...
	cfg := &Config{}
	cfg.Server.Port = "8080"
	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = 10 * time.Second
	cfg.Server.IdleTimeout = 120 * time.Second
	cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/accounting?sslmode=disable"
	return cfg
}
