package qwencode

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

var qwenCodeModels = []*QwenModelListResponse{
	{
		ID:      "qwen3-coder-plus",
		Object:  "model",
		Created: 1754686206,
		OwnedBy: "qwen",
	},
	{
		ID:      "qwen3-coder-turbo",
		Object:  "model",
		Created: 1754686206,
		OwnedBy: "qwen",
	},
	{
		ID:      "qwen3-plus",
		Object:  "model",
		Created: 1754686206,
		OwnedBy: "qwen",
	},
	{
		ID:      "qwen3-turbo",
		Object:  "model",
		Created: 1754686206,
		OwnedBy: "qwen",
	},
}

type QwenModelListResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

func (q *QwenModelListResponse) Convert(ctx context.Context) (*openai.OpenAIModelListResponse, error) {
	models := make([]*openai.Model, len(qwenCodeModels))
	for i, model := range qwenCodeModels {
		models[i] = &openai.Model{
			ID:      model.ID,
			Object:  model.Object,
			Created: model.Created,
			OwnedBy: model.OwnedBy,
		}
	}

	return &openai.OpenAIModelListResponse{
		Object: "models",
		Data:   models,
	}, nil
}
