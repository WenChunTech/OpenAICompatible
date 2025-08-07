package gemini

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

var Provider *GeminiProvider

type GeminiProvider struct {
	*provider.BaseProvider
	ProjectID string
	Token     string
}

func NewGeminiProvider(token, projectID string) *GeminiProvider {
	headers := map[string]string{
		constant.Authorization: fmt.Sprintf("Bearer %s", token),
	}
	return &GeminiProvider{
		BaseProvider: &provider.BaseProvider{
			ChatCompleteURL:    constant.GeminiCliChatCompleteURL,
			ChatCompleteMethod: http.MethodPost,
			ModelURL:           constant.GeminiCliModelListURL,
			ModelMethod:        http.MethodGet,
			Headers:            headers,
		},
		ProjectID: projectID,
	}
}

func (p *GeminiProvider) HandleChatCompleteRequest(ctx context.Context, r *model.OpenAIChatCompletionRequest) (*request.Response, error) {
	reqBody, err := json.Marshal(r)
	if err != nil {
		slog.Error("marshal request body failed", "err", err)
		return nil, err
	}
	return request.NewRequestBuilder(constant.GeminiCliChatCompleteURL, p.ChatCompleteMethod).WithHeaders(p.Headers).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
}

func (p *GeminiProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.EventStreamHandler[*model.GeminiChatCompletionResponse]{}
	return handler.Handle(ctx, w, r)
}

// HandleListModelRequest 函数用于处理列表模型请求
func (p *GeminiProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	// 创建一个新的请求构建器，传入模型URL和模型方法
	return request.NewRequestBuilder(p.ModelURL, p.ModelMethod).WithHeaders(p.Headers).Do(ctx, nil)
}

func (p *GeminiProvider) HandleListModelResponse(ctx context.Context, r *request.Response) (*model.OpenAIModelListResponse, error) {
	handler := responser.ModelListHandler[*model.QwenListModelResponse]{}
	return handler.Handle(ctx, r)
}
