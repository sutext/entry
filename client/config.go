package client

import (
	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/types"
)

type Config struct {
	Host         string         `json:"host"`
	Port         string         `json:"port"`
	Platform     types.Platform `json:"platform"`
	KeepAlive    int64          `json:"keepalive"`
	PingTimeout  int64          `json:"ping_timeout"`
	LoggerLevel  logger.Level   `json:"logger_level"`
	LoggerFormat logger.Format  `json:"logger_format"`
}

func NewConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         "8080",
		Platform:     types.PlatformMobile,
		KeepAlive:    60,
		PingTimeout:  5,
		LoggerLevel:  logger.LevelInfo,
		LoggerFormat: logger.FormatJSON,
	}
}
