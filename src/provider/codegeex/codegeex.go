package codegeex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
	"github.com/WenChunTech/OpenAICompatible/src/request"
	"github.com/WenChunTech/OpenAICompatible/src/responser"
)

var Provider *CodeGeexProvider

type CodeGeexProvider struct {
	*provider.BaseProvider
}

func init() {
	Provider = NewCodeGeexProvider(config.Config.CodeGeex.Token)
}

func NewCodeGeexProvider(token string) *CodeGeexProvider {
	headers := map[string]string{
		constant.Accept:    constant.ContentTypeEventStream,
		constant.CodeToken: token,
		constant.UserAgent: constant.DefaultUserAgent,
	}

	return &CodeGeexProvider{
		BaseProvider: &provider.BaseProvider{
			ChatCompleteURL:    constant.CodeGeexChatCompleteURL,
			ChatCompleteMethod: http.MethodPost,
			ModelURL:           constant.CodeGeexModelListURL,
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
	return request.NewRequestBuilder(p.ChatCompleteURL, p.ChatCompleteMethod).WithHeaders(p.Headers).WithJson(bytes.NewReader(requestBody)).Do(ctx, nil)
}

func (p *CodeGeexProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.EventStreamHandler[model.CodeGeexEventSourceData]{}
	return handler.Handle(ctx, w, r)
}

func (p *CodeGeexProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	return request.NewRequestBuilder(p.ModelURL, p.ModelMethod).Do(ctx, nil)
}

func (p *CodeGeexProvider) HandleListModelResponse(ctx context.Context, r *request.Response) (*model.OpenAIModelListResponse, error) {
	handler := responser.ModelListHandler[*model.CodeGeexModelOptions]{}
	return handler.Handle(ctx, r)
}
