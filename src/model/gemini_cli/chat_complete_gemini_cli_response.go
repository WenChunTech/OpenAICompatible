package geminicli

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

func (c *GeminiChatCompletionResponse) Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error) {
	choices := make([]openai.OpenAIStreamChoice, 0)
	for i, candidate := range c.Response.Candidates {
		// 创建 OpenAI 格式的 Delta 对象
		delta := &openai.Delta{
			Role: candidate.Content.Role,
		}

		// 处理 Content 中的各个部分
		var contentText string
		for _, part := range candidate.Content.Parts {
			// 目前仅处理文本内容，其他类型（如函数调用）可以根据需要扩展
			if part.Text != "" {
				contentText += part.Text
			}

			// 处理函数调用（如果有）
			if part.FunctionCall != nil {
				argsJSON, err := json.Marshal(part.FunctionCall.Args)
				if err == nil {
					delta.ToolCall = &openai.ToolCall{
						ID:   fmt.Sprintf("call_%d", i),
						Type: "function",
						Function: openai.FunctionCall{
							Name:      part.FunctionCall.Name,
							Arguments: string(argsJSON),
						},
					}
				}
			}
		}

		// 设置内容文本
		delta.Content = contentText

		// 创建 OpenAI 格式的 Choice 对象
		choice := openai.OpenAIStreamChoice{
			Index: i,
			Delta: delta,
		}

		choices = append(choices, choice)
	}

	// 创建 Usage 对象（如果有相关信息）
	var usage *openai.Usage
	// 这里可以根据 Gemini 响应中的信息填充 Usage，如果有的话

	resp := &openai.OpenAPIChatCompletionStreamResponse{
		ID:      c.Response.ResponseId,
		Object:  "chat.completion.chunk",
		Created: c.Response.CreateTime.Unix(),
		Model:   c.Response.ModelVersion,
		Choices: choices,
		Usage:   usage,
	}
	return resp, nil
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
