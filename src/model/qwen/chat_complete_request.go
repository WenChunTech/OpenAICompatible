package qwen

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type QwenChatCompleteRequest struct {
	Stream            bool       `json:"stream"`
	IncrementalOutput bool       `json:"incremental_output"`
	ChatID            string     `json:"chat_id"`
	ChatMode          string     `json:"chat_mode"`
	Model             string     `json:"model"`
	ParentID          any        `json:"parent_id"`
	Messages          []Messages `json:"messages"`
	Timestamp         int        `json:"timestamp"`
	Size              string     `json:"size"`
}
type FeatureConfig struct {
	ThinkingEnabled bool   `json:"thinking_enabled"`
	OutputSchema    string `json:"output_schema"`
	SearchVersion   string `json:"search_version"`
	ThinkingBudget  int    `json:"thinking_budget"`
}
type QwenChatCompleteMeta struct {
	SubChatType string `json:"subChatType"`
}
type Extra struct {
	Meta QwenChatCompleteMeta `json:"meta"`
}
type Messages struct {
	Fid           string        `json:"fid"`
	ParentID      any           `json:"parentId"`
	ChildrenIds   []string      `json:"childrenIds"`
	Role          string        `json:"role"`
	Content       string        `json:"content"`
	UserAction    string        `json:"user_action"`
	Files         []any         `json:"files"`
	Timestamp     int           `json:"timestamp"`
	Models        []string      `json:"models"`
	ChatType      string        `json:"chat_type"`
	FeatureConfig FeatureConfig `json:"feature_config"`
	Extra         Extra         `json:"extra"`
	SubChatType   string        `json:"sub_chat_type"`
	ParentID0     any           `json:"parent_id"`
}

func (q *QwenChatCompleteRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error {
	messages := make([]Messages, 0, len(req.Messages))

	for _, message := range req.Messages {
		var content string
		switch message.Content.(type) {
		case string:
			content = message.Content.(string)
		case map[string]interface{}:
			contentMap := message.Content.(map[string]interface{})
			jsonContent, err := json.Marshal(contentMap)
			if err != nil {
				slog.Error("Failed to marshal content map", "error", err)
				return fmt.Errorf("failed to marshal content map: %w", err)
			}
			content = string(jsonContent)
		}

		messages = append(messages, Messages{
			Role:     message.Role,
			Content:  content,
			Models:   []string{req.Model},
			ChatType: "search",
			FeatureConfig: FeatureConfig{
				ThinkingEnabled: true,
				SearchVersion:   "v2",
				ThinkingBudget:  38912,
			},
			Extra: Extra{
				Meta: QwenChatCompleteMeta{
					SubChatType: "search",
				},
			},
			SubChatType: "search",
		})
	}

	var chatID string
	if id := ctx.Value(constant.ChatIDKey); id != nil {
		if idStr, ok := id.(string); ok {
			chatID = idStr
		}
	}

	*q = QwenChatCompleteRequest{
		Stream:            req.Stream != nil && *req.Stream,
		IncrementalOutput: true,
		ChatID:            chatID,
		ChatMode:          "normal",
		Model:             req.Model,
		ParentID:          nil,
		Messages:          messages,
		Timestamp:         int(time.Now().Unix()),
	}

	return nil
}
