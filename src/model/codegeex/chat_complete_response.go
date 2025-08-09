package codegeex

import (
	"context"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type CodeGeexChatCompleteResponse struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Model        string `json:"model"`
}

func (c CodeGeexChatCompleteResponse) Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error) {
	choice := openai.OpenAIStreamChoice{
		Index: 0,
		Delta: &openai.Delta{
			Content: c.Text,
			Role:    constant.AssistantRole,
		},
	}
	if c.FinishReason != "" {
		finishReason := c.FinishReason
		choice.FinishReason = &finishReason
	}
	return &openai.OpenAPIChatCompletionStreamResponse{
		ID:      c.ID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   c.Model,
		Choices: []openai.OpenAIStreamChoice{choice},
		Usage:   nil,
	}, nil
}
