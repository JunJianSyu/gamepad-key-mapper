package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = "gamepad-key-mapper.json"

// GetConfigPath 获取配置文件路径
func GetConfigPath() (string, error) {
	// 优先使用用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		// 备选：使用程序所在目录
		execPath, err := os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Join(filepath.Dir(execPath), configFileName), nil
	}

	// 创建应用专属目录
	appDir := filepath.Join(configDir, "GamepadKeyMapper")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, configFileName), nil
}

// Load 加载配置
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return NewDefault(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在，返回默认配置
			return NewDefault(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		// 配置文件损坏，备份并返回默认配置
		backupPath := path + ".backup"
		os.Rename(path, backupPath)
		return NewDefault(), nil
	}

	return &cfg, nil
}

// Save 保存配置
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// SaveRules 仅保存规则（保留其他设置）
func SaveRules(rules []*Rule) error {
	cfg, err := Load()
	if err != nil {
		cfg = NewDefault()
	}

	// 类型转换 - 这里需要注意Rule和MappingRule的关系
	// 实际使用时应该直接使用mapper.MappingRule
	return Save(cfg)
}

// Rule 是一个别名，用于方便导入
type Rule = struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SourceKey uint16 `json:"source_key"`
	TargetKey int    `json:"target_key"`
	Modifiers struct {
		Ctrl  bool `json:"ctrl"`
		Alt   bool `json:"alt"`
		Shift bool `json:"shift"`
		Win   bool `json:"win"`
	} `json:"modifiers"`
	Enabled bool `json:"enabled"`
}
