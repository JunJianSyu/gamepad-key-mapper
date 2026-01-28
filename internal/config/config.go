package config

import (
	"gamepad-key-mapper/internal/mapper"
)

// Config 应用配置
type Config struct {
	Rules          []*mapper.MappingRule `json:"rules"`
	MinimizeToTray bool                  `json:"minimize_to_tray"`
	StartMinimized bool                  `json:"start_minimized"`
}

// NewDefault 创建默认配置
func NewDefault() *Config {
	return &Config{
		Rules:          make([]*mapper.MappingRule, 0),
		MinimizeToTray: true,
		StartMinimized: false,
	}
}
