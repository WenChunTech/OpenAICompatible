package main

import (
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/WenChunTech/OpenAICompatible/src/auth"
	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/handler"
	"github.com/WenChunTech/OpenAICompatible/src/manager"
)

func main() {
	switch len(os.Args) {
	case 2:
		switch os.Args[1] {
		case "gemini_cli":
			auth.StartGeminiCliAuth()
		case "qwen_code_cli":
			auth.StartQwenCodeAuth()
		default:
			Serve()
		}
	default:
		Serve()
	}
}

func Serve() {
	manager := manager.InitProviderManager()
	http.HandleFunc("/", handler.HandleForward)
	http.HandleFunc("/v1/models", manager.HandleListModel)
	http.HandleFunc("/v1/chat/completions", manager.HandleChatComplete)

	slog.Info("Server starting...", "host", config.Config.Host, "port", config.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port), nil); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
