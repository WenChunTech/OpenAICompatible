package geminicli

import (
	"context"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

func (c *GeminiChatCompletionResponse) Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error) {
	return nil, nil
}

type GeminiChatCompletionResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Candidates    []GeminiCandidate `json:"candidates"`
	UsageMetadata UsageMetadata     `json:"usageMetadata"`
	ModelVersion  string            `json:"modelVersion"`
	CreateTime    time.Time         `json:"createTime"`
	ResponseId    string            `json:"responseId"`
}

type GeminiCandidate struct {
	Content Content `json:"content"`
}

type UsageMetadata struct {
	TrafficType string `json:"trafficType"`
}
