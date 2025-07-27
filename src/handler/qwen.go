package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/provider/qwen"
)

var QwenHandler *Handler[*qwen.QwenProvider]

func init() {
	qwenProvider := qwen.NewQwenProvider(config.Config.Qwen.Token)
	QwenHandler = &Handler[*qwen.QwenProvider]{}
	QwenHandler.P = qwenProvider
}

func NewQwenHandler(token string) *Handler[*qwen.QwenProvider] {
	qwenHandler := &Handler[*qwen.QwenProvider]{}
	qwenHandler.P = qwen.NewQwenProvider(token)
	return qwenHandler
}
