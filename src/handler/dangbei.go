package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/provider/dangbei"
)

func NewDangBeiHandler() *Handler[*dangbei.DangBeiProvider] {
	dangbeiHandler := &Handler[*dangbei.DangBeiProvider]{}
	dangbeiHandler.P = dangbei.NewDangBeiProvider()
	return dangbeiHandler
}
