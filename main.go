package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/handler"
	"github.com/WenChunTech/OpenAICompatible/src/manager"
	"github.com/WenChunTech/OpenAICompatible/src/provider/codegeex"
	"github.com/WenChunTech/OpenAICompatible/src/provider/qwen"
)

func main() {
	context := context.Background()
	manager := manager.NewHandlerManager()
	manager.RegisterHandler(context, "codegeex", codegeex.Provider)
	manager.RegisterHandler(context, "qwen", qwen.Provider)

	http.HandleFunc("/", handler.NotFoundHandler)

	// handler := handler.CodeGeexHandler
	// handler := handler.QwenHandler

	http.HandleFunc("/v1/models", manager.HandleListModel)
	http.HandleFunc("/v1/chat/completions", manager.HandleChatComplete)

	slog.Info("Server starting...", "host", config.Config.Host, "port", config.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port), nil); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
