package qwen

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type QwenModelListResponse struct {
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

func (q *QwenModelListResponse) Convert(ctx context.Context) (*openai.OpenAIModelListResponse, error) {
	models := make([]*openai.Model, 0, len(q.Data))
	for _, data := range q.Data {
		models = append(models, &openai.Model{
			ID:      data.ID,
			Object:  data.Object,
			Created: int64(data.Info.CreatedAt),
			OwnedBy: data.Info.UserID,
		})
	}

	return &openai.OpenAIModelListResponse{
		Data: models,
	}, nil
}
