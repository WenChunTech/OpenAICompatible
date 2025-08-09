package codegeex

import (
	"context"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type Option struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description *Description `json:"description"`
	Host        string       `json:"host"`
}

type Description struct {
	Zh string `json:"zh"`
	En string `json:"en"`
}

type Promote struct {
	Key  string `json:"key"`
	Text struct {
		Zh string `json:"zh"`
		En string `json:"en"`
	} `json:"text"`
}

type CodeGeexModelListResponse struct {
	Options []Option `json:"options"`
	Default string   `json:"default"`
	Promote *Promote `json:"promote"`
	IP      string   `json:"ip"`
}

func (c *CodeGeexModelListResponse) Convert(ctx context.Context) (*openai.OpenAIModelListResponse, error) {
	models := make([]*openai.Model, len(c.Options))
	for i, option := range c.Options {
		models[i] = &openai.Model{
			ID:      option.ID,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "codegeex",
		}
	}
	return &openai.OpenAIModelListResponse{
		Object: "list",
		Data:   models,
	}, nil
}
