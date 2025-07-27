package config

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
)

const (
	AppConfigFile      = "config.json"
	APPLocalConfigFile = "config.local.json"
)

type AppConfig struct {
	Host     string          `json:"host,omitempty"`
	Port     int             `json:"port"`
	CodeGeex *CodeGeexConfig `json:"codegeex,omitempty"`
	Qwen     *QwenConfig     `json:"qwen,omitempty"`
}

type CodeGeexConfig struct {
	UserID        string `json:"user_id,omitempty"`        // 用户ID
	UserRole      int    `json:"user_role,omitempty"`      // 用户角色
	IDE           string `json:"ide,omitempty"`            // IDE类型
	IDEVersion    string `json:"ide_version,omitempty"`    // IDE版本
	PluginVersion string `json:"plugin_version,omitempty"` // 插件版本
	MachineID     string `json:"machine_id,omitempty"`     // 机器ID
	TalkID        string `json:"talk_id,omitempty"`        // 对话ID
	Locale        string `json:"locale,omitempty"`         // 语言
	Token         string `json:"token"`                    // 访问令牌
}

type QwenConfig struct {
	Token string `json:"token"`
}

var Config = &AppConfig{}

func init() {
	if _, err := os.Stat(APPLocalConfigFile); err == nil {
		buffer, err := os.ReadFile(APPLocalConfigFile)
		if err != nil {
			slog.Error("Failed to read app local config file", "error", err)
		}
		err = json.NewDecoder(bytes.NewReader(buffer)).Decode(Config)
		if err != nil {
			slog.Error("Failed to parse app config file", "error", err)
		}

	} else {
		slog.Error("App config file not found", "file", AppConfigFile)
		buffer, err := os.ReadFile(AppConfigFile)
		if err != nil {
			slog.Error("Failed to read app config file", "error", err)
		}
		err = json.NewDecoder(bytes.NewReader(buffer)).Decode(Config)
		if err != nil {
			slog.Error("Failed to parse app local config file", "error", err)
			panic("App local config file not found")
		}
	}
}
