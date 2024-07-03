package config

import "time"

type Config struct {
	LoggerLevel string
	HTTPTimeout time.Duration
	ServerPort  string
}
