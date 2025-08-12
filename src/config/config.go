package config

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"golang.org/x/oauth2"
)

const (
	AppConfigFile      = "config.json"
	APPLocalConfigFile = "config.local.json"
)

var RollMap = map[string]int{
	constant.CodeGeexPrefix:  0,
	constant.QwenPrefix:      0,
	constant.GeminiCliPrefix: 0,
}

type AppConfig struct {
	Host      string             `json:"host,omitempty"`
	Port      int                `json:"port"`
	CodeGeex  []*CodeGeexConfig  `json:"codegeex,omitempty"`
	Qwen      []*QwenConfig      `json:"qwen,omitempty"`
	GeminiCli []*GeminiCliConfig `json:"gemini_cli,omitempty"`
}

type GeminiCliConfig struct {
	Prefix    string        `json:"prefix,omitempty"`
	ProjectID string        `json:"project_id,omitempty"`
	Token     *oauth2.Token `json:"token,omitempty"`
}

type CodeGeexConfig struct {
	UserID        string `json:"user_id,omitempty"`
	UserRole      int    `json:"user_role,omitempty"`
	IDE           string `json:"ide,omitempty"`
	IDEVersion    string `json:"ide_version,omitempty"`
	PluginVersion string `json:"plugin_version,omitempty"`
	MachineID     string `json:"machine_id,omitempty"`
	TalkID        string `json:"talk_id,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Token         string `json:"token"`
	Prefix        string `json:"prefix,omitempty"`
}

type QwenConfig struct {
	Prefix string `json:"prefix,omitempty"`
	Token  string `json:"token"`
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

	if len(Config.CodeGeex) != 0 && len(Config.CodeGeex[0].Prefix) != 0 {
		Config.CodeGeex[0].Prefix = constant.CodeGeexPrefix
	}
	if len(Config.Qwen) != 0 && len(Config.Qwen[0].Prefix) != 0 {
		Config.Qwen[0].Prefix = constant.QwenPrefix
	}
	if Config.GeminiCli != nil && len(Config.GeminiCli[0].Prefix) != 0 {
		Config.GeminiCli[0].Prefix = constant.GeminiCliPrefix
	}
}

func GetCodeGeexConfig() *CodeGeexConfig {
	if len(Config.CodeGeex) == 0 {
		return nil
	}

	return Config.CodeGeex[RollMap[constant.CodeGeexPrefix]]
}

func NextCodeGeexConfig() *CodeGeexConfig {
	config := GetCodeGeexConfig()
	if config == nil {
		return nil
	}

	RollMap[constant.CodeGeexPrefix] = (RollMap[constant.CodeGeexPrefix] + 1) % len(Config.CodeGeex)
	return config
}

func GetQwenConfig() *QwenConfig {
	if len(Config.Qwen) == 0 {
		return nil
	}

	return Config.Qwen[RollMap[constant.QwenPrefix]]
}

func NextQwenConfig() *QwenConfig {
	config := GetQwenConfig()
	if config == nil {
		return nil
	}

	RollMap[constant.QwenPrefix] = (RollMap[constant.QwenPrefix] + 1) % len(Config.Qwen)
	return config
}

func GetGeminiCliConfig() *GeminiCliConfig {
	if len(Config.GeminiCli) == 0 {
		return nil
	}

	return Config.GeminiCli[RollMap[constant.GeminiCliPrefix]]
}

func NextGeminiCliConfig() *GeminiCliConfig {
	config := GetGeminiCliConfig()
	if config == nil {
		return nil
	}

	RollMap[constant.GeminiCliPrefix] = (RollMap[constant.GeminiCliPrefix] + 1) % len(Config.GeminiCli)
	return config
}
