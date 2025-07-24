package converter

import "github.com/WenChunTech/OpenAICompatible/src/model"

type Converter[T any] interface {
	Convert() (T, error)
}

type ChatCompletionConverter interface {
	Convert() (*model.OpenAPIChatCompletionStreamResponse, error)
}

type ModelConverter interface {
	Convert() (*model.OpenAIModelList, error)
}
