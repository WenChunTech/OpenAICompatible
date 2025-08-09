package importer

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

type Importer[T any] interface {
	Import(T) error
}

// OpenAIChatCompletionImporter is an interface for importing OpenAI chat completion requests.
type OpenAIChatCompletionImporter interface {
	ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error
}
