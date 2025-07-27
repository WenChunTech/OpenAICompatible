package codegeex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
	"github.com/WenChunTech/OpenAICompatible/src/request"
	"github.com/WenChunTech/OpenAICompatible/src/responser"
)

type CodeGeexProvider struct {
	*provider.BaseProvider
}

func NewCodeGeexProvider(token string) *CodeGeexProvider {
	headers := map[string]string{
		constant.Accept:    constant.ContentTypeEventStream,
		constant.CodeToken: token,
		constant.UserAgent: constant.DefaultUserAgent,
	}

	return &CodeGeexProvider{
		BaseProvider: &provider.BaseProvider{
			ChatCompleteURL:    constant.CodeGeexChatURL,
			ChatCompleteMethod: http.MethodPost,
			ModelURL:           constant.CodeGeexModelURL,
			ModelMethod:        http.MethodGet,
			Headers:            headers,
		},
	}
}

func (p *CodeGeexProvider) HandleChatCompleteRequest(ctx context.Context, r *model.OpenAIChatCompletionRequest) (*request.Response, error) {
	codegeex := model.CodeGeexChatRequest{}
	err := codegeex.ImportOpenAIChatCompletionRequest(ctx, r)
	if err != nil {
		slog.Error("Failed to import OpenAI chat completion request", "error", err)
		return nil, fmt.Errorf("failed to import OpenAI chat completion request: %w", err)
	}
	requestBody, err := json.Marshal(codegeex)
	if err != nil {
		slog.Error("Failed to marshal request body", "error", err)
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return request.NewRequestBuilder(p.ChatCompleteURL, p.ChatCompleteMethod).AddHeaders(p.Headers).SetJson(bytes.NewReader(requestBody)).Do(ctx, nil)
}
func (p *CodeGeexProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.EventStreamHandler[model.CodeGeexEventSourceData]{}
	return handler.Handle(ctx, w, r)
}

func (p *CodeGeexProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	return request.NewRequestBuilder(p.ModelURL, p.ModelMethod).Do(ctx, nil)
}

func (p *CodeGeexProvider) HandleListModelResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.ModelListHandler[*model.CodeGeexModelOptions]{}
	return handler.Handle(ctx, w, r)
}
