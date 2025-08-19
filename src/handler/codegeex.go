package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/provider/codegeex"
)

func NewCodeGeexHandler(token string) *Handler[*codegeex.CodeGeexProvider] {
	codegeexHandler := &Handler[*codegeex.CodeGeexProvider]{}
	codegeexHandler.P = codegeex.NewCodeGeexProvider(token)
	return codegeexHandler
}
