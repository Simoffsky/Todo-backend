package config

import (
	"time"
)

type Config struct {
	LoggerLevel string
	HTTPTimeout time.Duration
	ServerPort  string

	AuthAddr  string
	JwtSecret string

	DbConn string
}
