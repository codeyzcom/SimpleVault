package config

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Port       int
	Host       string
	SessionTTL time.Duration
	DataStore  string
}

func LoadConfig() *Config {
	cfg := Config{}

	flag.StringVar(&cfg.Host, "host", "localhost", "HTTP server host")
	flag.IntVar(&cfg.Port, "port", 7879, "HTTP server port")
	flag.DurationVar(&cfg.SessionTTL, "session_ttl", 15*time.Minute, "session TTL  (e.g. 10s, 5m, 1h)")
	flag.StringVar(&cfg.DataStore, "data_store", "data/", "path to store vaults")

	flag.Parse()

	return &cfg
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
