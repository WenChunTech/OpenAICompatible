package handler

import (
	"github.com/WenChunTech/OpenAICompatible/src/provider/dangbei"
)

var DangBeiHandler *Handler[*dangbei.DangBeiProvider]

func init() {
	dangbeiProvider := dangbei.NewDangBeiProvider()
	DangBeiHandler = &Handler[*dangbei.DangBeiProvider]{}
	DangBeiHandler.P = dangbeiProvider
}

func NewDangBeiHandler() *Handler[*dangbei.DangBeiProvider] {
	dangbeiHandler := &Handler[*dangbei.DangBeiProvider]{}
	dangbeiHandler.P = dangbei.NewDangBeiProvider()
	return dangbeiHandler
}
