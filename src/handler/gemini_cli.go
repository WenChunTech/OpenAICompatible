package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/config"
	gemini "github.com/WenChunTech/OpenAICompatible/src/provider/geminicli"
	"golang.org/x/oauth2"
)

var GeminiCliHandler *Handler[*gemini.GeminiCliProvider]

func init() {
	config := config.GetGeminiCliConfig()
	geminiProvider := gemini.NewGeminiCliProvider(config.ProjectID, config.Token)
	GeminiCliHandler = &Handler[*gemini.GeminiCliProvider]{}
	GeminiCliHandler.P = geminiProvider
}

func NewGeminiCliHandler(projectID string, token *oauth2.Token) *Handler[*gemini.GeminiCliProvider] {
	geminicliHandler := &Handler[*gemini.GeminiCliProvider]{}
	geminicliHandler.P = gemini.NewGeminiCliProvider(projectID, token)
	return geminicliHandler
}
