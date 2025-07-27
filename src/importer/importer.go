package importer

import (
	"context"

	"github.com/WenChunTech/OpenAICompatible/src/model"
)

type Importer[T any] interface {
	Import(T) error
}

// OpenAIChatCompletionImporter is an interface for importing OpenAI chat completion requests.
type OpenAIChatCompletionImporter interface {
	ImportOpenAIChatCompletionRequest(ctx context.Context, req *model.OpenAIChatCompletionRequest) error
}
