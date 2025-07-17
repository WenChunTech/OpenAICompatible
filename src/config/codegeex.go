package config

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	ConfigFile      = ".env"
	LocalConfigFile = ".env.local"
)

const (
	UserIDField        = "UserID"
	UserRoleField      = "UserRole"
	IDEField           = "IDE"
	IDEVersionField    = "IDEVersion"
	PluginVersionField = "PluginVersion"
	MachineIDField     = "MachineID"
	TalkIDField        = "TalkID"
	LocaleField        = "Locale"
	TokenField         = "Token"
)

type CodeGeexConfig struct {
	UserID        string
	UserRole      int
	IDE           string
	IDEVersion    string
	PluginVersion string
	MachineID     string
	TalkID        string
	Locale        string
	Token         string
}

var Config = &CodeGeexConfig{}

func init() {
	if _, err := os.Stat(LocalConfigFile); err == nil {
		buffer, err := os.ReadFile(LocalConfigFile)
		if err != nil {
			slog.Error("Failed to read local config file", "error", err)
		}
		parseConfig(buffer)
	} else {
		buffer, err := os.ReadFile(ConfigFile)
		if err != nil {
			slog.Error("Failed to read config file", "error", err)
			panic("Config file not found")
		}
		parseConfig(buffer)
	}
}

func parseConfig(buffer []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		parts := bytes.SplitN([]byte(line), []byte("="), 2)
		if len(parts) < 2 {
			slog.Error("Invalid config line, skipping", "line", line)
			continue
		}

		key := strings.TrimSpace(string(parts[0]))
		value := strings.Trim(strings.TrimSpace(string(parts[1])), `"`)

		switch key {
		case UserIDField:
			Config.UserID = value
		case UserRoleField:
			Config.UserRole, _ = strconv.Atoi(value)
		case IDEField:
			Config.IDE = value
		case IDEVersionField:
			Config.IDEVersion = value
		case PluginVersionField:
			Config.PluginVersion = value
		case MachineIDField:
			Config.MachineID = value
		case TalkIDField:
			Config.TalkID = value
		case LocaleField:
			Config.Locale = value
		case TokenField:
			Config.Token = value
		default:
			slog.Warn("Unknown config key, skipping", "key", key)
		}
	}
}
