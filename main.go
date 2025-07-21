package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/handler"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "path", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/v1/chat/completions", handler.ChatProxyChatHandler)
	http.HandleFunc("/v1/models", handler.ChatProxyModelHandler)

	slog.Info("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
