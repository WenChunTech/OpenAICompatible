package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/provider/qwen"
)

func NewQwenHandler(token string) *Handler[*qwen.QwenProvider] {
	qwenHandler := &Handler[*qwen.QwenProvider]{}
	qwenHandler.P = qwen.NewQwenProvider(token)
	return qwenHandler
}
