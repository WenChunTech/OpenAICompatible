package provider

import (
	"context"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

type BaseProvider struct {
	ChatCompleteURL    string
	ChatCompleteMethod string
	ModelURL           string
	ModelMethod        string
	Headers            map[string]string
}

type Provider interface {
	HandleChatCompleteRequest(ctx context.Context, r *openai.OpenAIChatCompletionRequest) (*request.Response, error)
	HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error

	HandleListModelRequest(ctx context.Context) (*request.Response, error)
	HandleListModelResponse(ctx context.Context, r *request.Response) (*openai.OpenAIModelListResponse, error)
}
