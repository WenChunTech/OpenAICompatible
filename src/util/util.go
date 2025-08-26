package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/config"
)

const Nonce = "Qibednl6_AZRzQLle-gdA"

func Sign(timestamp int64, payload string) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d%s%s", timestamp, payload, Nonce)))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func GenerateUUID() string {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		slog.Error("Failed to generate UUID")
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:])
}

func SaveConfig(config *config.AppConfig, path string) {
	data, err := json.Marshal(config)
	if err != nil {
		slog.Error("Failed to marshal config", "error", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		slog.Error("Failed to save config", "error", err)
	}
}
