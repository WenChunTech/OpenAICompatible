package qwen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/model/qwen"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
	"github.com/WenChunTech/OpenAICompatible/src/request"
	"github.com/WenChunTech/OpenAICompatible/src/responser"
)

var Provider *QwenProvider

type QwenProvider struct {
	*provider.BaseProvider

	// ChatID string
	// Model  string
	// Token  string
}

func init() {
	config := config.GetQwenConfig()
	Provider = NewQwenProvider(config.Token)
}

func NewQwenProvider(token string) *QwenProvider {
	headers := map[string]string{
		constant.Authorization: fmt.Sprintf("Bearer %s", token),
	}
	return &QwenProvider{
		BaseProvider: &provider.BaseProvider{
			ChatCompleteURL:    constant.QwenChatCompleteURL,
			ChatCompleteMethod: http.MethodPost,
			ModelURL:           constant.QwenModelListURL,
			ModelMethod:        http.MethodGet,
			Headers:            headers,
		},
	}
}

func (p *QwenProvider) HandleChatCompleteRequest(ctx context.Context, r *openai.OpenAIChatCompletionRequest) (*request.Response, error) {
	req := qwen.QwenChatIDRequest{
		ChatMode:  "normal",
		ChatType:  "search",
		Timestamp: time.Now().Unix(),
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal chatid request body", "error", err)
	}
	resp, err := request.NewRequestBuilder(constant.QwenChatID, http.MethodPost).WithHeaders(p.Headers).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
	}
	var chatIDResp qwen.QwenChatIDResponse
	err = resp.Json(&chatIDResp)
	if err != nil {
		slog.Error("Failed to decode response body", "error", err)
	}
	resp.Body.Close()

	chatCompleteReq := new(qwen.QwenChatCompleteRequest)
	ctx = context.WithValue(ctx, constant.ChatIDKey, chatIDResp.Data.ID)
	if err := chatCompleteReq.ImportOpenAIChatCompletionRequest(ctx, r); err != nil {
		slog.Error("Failed to import openai chat completion request", "error", err)
		return nil, err
	}

	reqBody, err = json.Marshal(chatCompleteReq)
	if err != nil {
		slog.Error("Failed to marshal chatcomplete request body", "error", err)
	}
	return request.NewRequestBuilder(constant.QwenChatCompleteURL, http.MethodPost).WithQuery(string(constant.ChatIDKey), chatIDResp.Data.ID).WithHeaders(p.Headers).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
}

func (p *QwenProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.EventStreamHandler[qwen.QwenChatCompleteResponse]{}
	return handler.Handle(ctx, w, r)
}

// HandleListModelRequest 函数用于处理列表模型请求
func (p *QwenProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	// 创建一个新的请求构建器，传入模型URL和模型方法
	return request.NewRequestBuilder(p.ModelURL, p.ModelMethod).WithHeaders(p.Headers).Do(ctx, nil)
}

func (p *QwenProvider) HandleListModelResponse(ctx context.Context, r *request.Response) (*openai.OpenAIModelListResponse, error) {
	handler := responser.ModelListHandler[*qwen.QwenModelListResponse]{}
	return handler.Handle(ctx, r)
}
