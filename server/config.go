package server

import (
	"time"

	"sutext.github.io/entry/logger"
)

type Config struct {
	Port         string        `json:"port"`
	KeepAlive    time.Duration `json:"keepalive"`
	PingTimeout  time.Duration `json:"ping_timeout"`
	LoggerLevel  logger.Level  `json:"logger_level"`
	LoggerFormat logger.Format `json:"logger_format"`
}

func NewConfig() *Config {
	return &Config{
		Port:         "8080",
		KeepAlive:    60,
		PingTimeout:  5,
		LoggerLevel:  logger.LevelInfo,
		LoggerFormat: logger.FormatJSON,
	}
}
