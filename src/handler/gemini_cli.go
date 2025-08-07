package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/provider/codegeex"
	gemini "github.com/WenChunTech/OpenAICompatible/src/provider/gemini_cli"
)

var GeminiCliHandler *Handler[*gemini.GeminiProvider]

func init() {
	codegeexProvider := codegeex.NewCodeGeexProvider(config.Config.CodeGeex.Token)
	CodeGeexHandler = &Handler[*codegeex.CodeGeexProvider]{}
	CodeGeexHandler.P = codegeexProvider
}

func NewGeminiCliHandler(token string, projectID string) *Handler[*gemini.GeminiProvider] {
	geminiHandler := &Handler[*gemini.GeminiProvider]{}
	geminiHandler.P = gemini.NewGeminiProvider(token, projectID)
	return geminiHandler
}
