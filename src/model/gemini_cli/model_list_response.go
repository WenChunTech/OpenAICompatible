package geminicli

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

func (g *GeminiCliModelListResponse) Convert(ctx context.Context) (*openai.OpenAIModelListResponse, error) {
	models := make([]*openai.Model, 0, len(g.Models))
	for _, model := range g.Models {
		models = append(models, &openai.Model{
			ID:      model.Name[len("models/"):],
			Object:  model.Name,
			OwnedBy: "google",
		})
	}
	return &openai.OpenAIModelListResponse{
		Object: "model",
		Data:   models,
	}, nil
}

type GeminiCliModelListResponse struct {
	Models        []GeminiCliModel `json:"models"`
	NextPageToken string           `json:"nextPageToken"`
}

type GeminiCliModel struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	DisplayName      string   `json:"displayName"`
	Description      string   `json:"description"`
	InputTokenLimit  int      `json:"inputTokenLimit"`
	OutputTokenLimit int      `json:"outputTokenLimit"`
	SupportedMethods []string `json:"supportedGenerationMethods"`
	Temperature      float64  `json:"temperature"`
	TopP             float64  `json:"topP"`
	TopK             int      `json:"topK"`
	MaxTemperature   float64  `json:"maxTemperature"`
	Thinking         bool     `json:"thinking"`
}

const GEMINI_CLI_MODEL_LIST_RESPONSE = `{
  "models": [
    {
      "name": "models/gemini-2.5-flash",
      "version": "001",
      "displayName": "Gemini 2.5 Flash",
      "description": "Stable version of Gemini 2.5 Flash, our mid-size multimodal model that supports up to 1 million tokens, released in June of 2025.",
      "inputTokenLimit": 1048576,
      "outputTokenLimit": 65536,
      "supportedGenerationMethods": [
        "generateContent",
        "countTokens",
        "createCachedContent",
        "batchGenerateContent"
      ],
      "temperature": 1,
      "topP": 0.95,
      "topK": 64,
      "maxTemperature": 2,
      "thinking": true
    },
    {
      "name": "models/gemini-2.5-pro",
      "version": "2.5",
      "displayName": "Gemini 2.5 Pro",
      "description": "Stable release (June 17th, 2025) of Gemini 2.5 Pro",
      "inputTokenLimit": 1048576,
      "outputTokenLimit": 65536,
      "supportedGenerationMethods": [
        "generateContent",
        "countTokens",
        "createCachedContent",
        "batchGenerateContent"
      ],
      "temperature": 1,
      "topP": 0.95,
      "topK": 64,
      "maxTemperature": 2,
      "thinking": true
    }
  ],
  "nextPageToken": ""
}`
