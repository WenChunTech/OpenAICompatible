package qwencode

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type QwenCodeChatCompleteRequest struct {
	*openai.OpenAIChatCompletionRequest
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Option struct {
	IncludeUser bool `json:"include_user"`
}

func (q *QwenCodeChatCompleteRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error {
	q.OpenAIChatCompletionRequest = req
	return nil
}
