package importer

import "github.com/WenChunTech/OpenAICompatible/src/model"

type Importer[T any] interface {
	Import(T) error
}

type OpenAIChatCompletionImporter interface {
	ImportOpenAIChatCompletionReq(req *model.OpenAIChatCompletionRequest) error
}
