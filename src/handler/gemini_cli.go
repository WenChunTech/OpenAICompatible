package handler

import (
	gemini "github.com/WenChunTech/OpenAICompatible/src/provider/geminicli"
	"golang.org/x/oauth2"
)

func NewGeminiCliHandler(projectID string, token *oauth2.Token) *Handler[*gemini.GeminiCliProvider] {
	geminicliHandler := &Handler[*gemini.GeminiCliProvider]{}
	geminicliHandler.P = gemini.NewGeminiCliProvider(projectID, token)
	return geminicliHandler
}
