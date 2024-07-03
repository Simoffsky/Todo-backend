package config

import (
	"fmt"
	"time"
)

type Config struct {
	LoggerLevel string
	HTTPTimeout time.Duration
	ServerPort  string
}

func (c Config) String() string {
	return fmt.Sprintf("LoggerLevel: %s, HTTPTimeout: %s, ServerPort: %s", c.LoggerLevel, c.HTTPTimeout, c.ServerPort)
}
