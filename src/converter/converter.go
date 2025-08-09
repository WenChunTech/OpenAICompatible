package converter

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type Converter[T any] interface {
	Convert() (T, error)
}

// ChatCompletionConverter is an interface for converting chat completion responses.
type ChatCompletionConverter interface {
	Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error)
}

// ModelConverter is an interface for converting model list responses.
type ModelConverter interface {
	Convert(ctx context.Context) (*openai.OpenAIModelListResponse, error)
}
