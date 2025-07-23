package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/handler"
	"github.com/WenChunTech/OpenAICompatible/src/model"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
)

func main() {
	provider := provider.Provider[*model.CodeGeexChatRequest, model.CodeGeexSSEData, *model.OpenAPIChatCompletionStreamResponse]{
		Provider: &model.CodeGeexChatRequest{},
		Url:      constant.CodeGeexChatURL,
		Method:   http.MethodPost,
		Headers: map[string]string{
			constant.Accept:    constant.ContentTypeEventStream,
			constant.CodeToken: config.Config.Token,
			constant.UserAgent: constant.DefaultUserAgent,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "path", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/v1/chat/completions", provider.Handle)
	http.HandleFunc("/v1/models", handler.ChatProxyModelHandler)

	slog.Info("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
