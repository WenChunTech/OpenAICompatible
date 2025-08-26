package qwencode

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type QwenChatCompleteResponse struct {
	*openai.OpenAPIChatCompletionStreamResponse
}

func (q QwenChatCompleteResponse) Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error) {
	return q.OpenAPIChatCompletionStreamResponse, nil
}
