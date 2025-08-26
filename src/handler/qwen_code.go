package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/config"
	qwencode "github.com/WenChunTech/OpenAICompatible/src/provider/qwen_code"
)

func NewQwenCodeHandler(token *config.QwenCodeToken) *Handler[*qwencode.QwenCodeProvider] {
	qwenCodeHandler := &Handler[*qwencode.QwenCodeProvider]{}
	qwenCodeHandler.P = qwencode.NewQwenCodeProvider(token)
	return qwenCodeHandler
}
