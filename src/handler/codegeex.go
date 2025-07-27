package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/provider/codegeex"
)

var CodeGeexHandler *Handler[*codegeex.CodeGeexProvider]

func init() {
	codegeexProvider := codegeex.NewCodeGeexProvider(config.Config.CodeGeex.Token)
	CodeGeexHandler = &Handler[*codegeex.CodeGeexProvider]{}
	CodeGeexHandler.P = codegeexProvider
}

func NewCodeGeexHandler(token string) *Handler[*codegeex.CodeGeexProvider] {
	codegeexHandler := &Handler[*codegeex.CodeGeexProvider]{}
	codegeexHandler.P = codegeex.NewCodeGeexProvider(token)
	return codegeexHandler
}
