package qwen_code

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type QwenChatCompleteResponse struct {
	Choices []Choices             `json:"choices"`
	Usage   QwenChatCompleteUsage `json:"usage"`
}
type WebSearchInfo struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Snippet  string `json:"snippet"`
	Hostname any    `json:"hostname"`
	Hostlogo any    `json:"hostlogo"`
	Date     string `json:"date"`
}
type QwenChatCompleteExtra struct {
	FunctionID    string          `json:"function_id"`
	WebSearchInfo []WebSearchInfo `json:"web_search_info"`
}
type QwenChatCompleteDelta struct {
	Role    string                `json:"role"`
	Content string                `json:"content"`
	Name    string                `json:"name"`
	Extra   QwenChatCompleteExtra `json:"extra"`
	Phase   string                `json:"phase"`
}
type Choices struct {
	Delta QwenChatCompleteDelta `json:"delta"`
}
type OutputTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}
type QwenChatCompleteUsage struct {
	InputTokens         int                 `json:"input_tokens"`
	OutputTokens        int                 `json:"output_tokens"`
	TotalTokens         int                 `json:"total_tokens"`
	OutputTokensDetails OutputTokensDetails `json:"output_tokens_details"`
}

func (q QwenChatCompleteResponse) Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error) {
	choices := make([]openai.OpenAIStreamChoice, 0, len(q.Choices))
	for _, choice := range q.Choices {
		choices = append(choices, openai.OpenAIStreamChoice{
			Delta: &openai.Delta{
				Role:    choice.Delta.Role,
				Content: choice.Delta.Content,
			},
		})
	}

	return &openai.OpenAPIChatCompletionStreamResponse{
		Choices: choices,
		Usage: &openai.Usage{
			PromptTokens:     q.Usage.InputTokens,
			CompletionTokens: q.Usage.OutputTokens,
			TotalTokens:      q.Usage.TotalTokens,
		},
	}, nil
}

type QwenChatIDRequest struct {
	Title     string   `json:"title"`
	Models    []string `json:"models"`
	ChatMode  string   `json:"chat_mode"`
	ChatType  string   `json:"chat_type"`
	Timestamp int64    `json:"timestamp"`
}

type QwenChatIDResponse struct {
	Success   bool   `json:"success"`
	RequestID string `json:"request_id"`
	Data      struct {
		ID string `json:"id"`
	} `json:"data"`
}
