package model

import (
	"context"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
)

type QwenListModelResponse struct {
	Data []Data `json:"data"`
}
type Capabilities struct {
	Document  bool `json:"document"`
	Vision    bool `json:"vision"`
	Video     bool `json:"video"`
	Audio     bool `json:"audio"`
	Citations bool `json:"citations"`
}
type Abilities struct {
	Document  int `json:"document"`
	Vision    int `json:"vision"`
	Video     int `json:"video"`
	Audio     int `json:"audio"`
	Mcp       int `json:"mcp"`
	Citations int `json:"citations"`
}
type Meta struct {
	ProfileImageURL     string       `json:"profile_image_url"`
	Description         string       `json:"description"`
	Capabilities        Capabilities `json:"capabilities"`
	ShortDescription    string       `json:"short_description"`
	MaxContextLength    int          `json:"max_context_length"`
	MaxGenerationLength int          `json:"max_generation_length"`
	Abilities           Abilities    `json:"abilities"`
	ChatType            []string     `json:"chat_type"`
	Mcp                 []string     `json:"mcp"`
	Modality            []string     `json:"modality"`
}
type Info struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	BaseModelID     any    `json:"base_model_id"`
	Name            string `json:"name"`
	Meta            Meta   `json:"meta"`
	AccessControl   any    `json:"access_control"`
	IsActive        bool   `json:"is_active"`
	IsVisitorActive bool   `json:"is_visitor_active"`
	UpdatedAt       int    `json:"updated_at"`
	CreatedAt       int    `json:"created_at"`
}
type Data struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Object    string `json:"object"`
	OwnedBy   string `json:"owned_by"`
	Info      Info   `json:"info"`
	Preset    bool   `json:"preset"`
	ActionIds []any  `json:"action_ids"`
}

func (q *QwenListModelResponse) Convert(ctx context.Context) (*OpenAIModelList, error) {
	models := make([]Model, 0, len(q.Data))
	for _, data := range q.Data {
		models = append(models, Model{
			ID:      data.ID,
			Object:  data.Object,
			Created: int64(data.Info.CreatedAt),
			OwnedBy: data.Info.UserID,
		})
	}

	return &OpenAIModelList{
		Data: models,
	}, nil
}

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

func (q *QwenChatCompleteRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *OpenAIChatCompletionRequest) error {
	messages := make([]Messages, 0, len(req.Messages))
	for _, message := range req.Messages {
		messages = append(messages, Messages{
			Role:     message.Role,
			Content:  message.Content,
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
		Stream:            req.Stream,
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

func (q QwenChatCompleteResponse) Convert(ctx context.Context) (*OpenAPIChatCompletionStreamResponse, error) {
	choices := make([]OpenAIStreamChoice, 0, len(q.Choices))
	for _, choice := range q.Choices {
		choices = append(choices, OpenAIStreamChoice{
			Delta: Delta{
				Role:    choice.Delta.Role,
				Content: choice.Delta.Content,
			},
		})
	}

	return &OpenAPIChatCompletionStreamResponse{
		Choices: choices,
		Usage: &Usage{
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
