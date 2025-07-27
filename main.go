package main

import (
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/handler"
)

func main() {
	http.HandleFunc("/", handler.NotFoundHandler)

	handler := handler.CodeGeexHandler
	// handler := handler.QwenHandler
	http.HandleFunc("/v1/models", handler.HandleListModel)
	http.HandleFunc("/v1/chat/completions", handler.HandleChatComplete)

	slog.Info("Server starting...", "host", config.Config.Host, "port", config.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port), nil); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
