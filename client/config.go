package client

import (
	"time"

	"sutext.github.io/entry/code"
	"sutext.github.io/entry/logger"
)

type Config struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	Platform     code.Platform `json:"platform"`
	KeepAlive    time.Duration `json:"keepalive"`
	PingTimeout  time.Duration `json:"ping_timeout"`
	LoggerLevel  logger.Level  `json:"logger_level"`
	LoggerFormat logger.Format `json:"logger_format"`
}

func NewConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         "8080",
		Platform:     code.PlatformMobile,
		KeepAlive:    60,
		PingTimeout:  5,
		LoggerLevel:  logger.LevelInfo,
		LoggerFormat: logger.FormatJSON,
	}
}
