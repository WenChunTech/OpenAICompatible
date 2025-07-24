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
	provider := provider.BaseProvider[*model.CodeGeexChatRequest, *model.CodeGeexModelOptions, model.CodeGeexEventSourceData]{
		ChatCompleteResp:   &model.CodeGeexChatRequest{},
		ModelResp:          &model.CodeGeexModelOptions{},
		ChatCompleteUrl:    constant.CodeGeexChatURL,
		ChatCompleteMethod: http.MethodPost,
		Headers: map[string]string{
			constant.Accept:    constant.ContentTypeEventStream,
			constant.CodeToken: config.Config.Token,
			constant.UserAgent: constant.DefaultUserAgent,
		},
		ModelUrl:    constant.CodeGeexModelURL,
		ModelMethod: http.MethodGet,
	}

	http.HandleFunc("/", handler.NotFoundHandler)
	http.HandleFunc("/v1/models", provider.ModelHandle)
	http.HandleFunc("/v1/chat/completions", provider.ChatCompleteHandle)

	slog.Info("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
