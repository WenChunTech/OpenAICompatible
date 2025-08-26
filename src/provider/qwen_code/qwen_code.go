package qwencode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/model/qwencode"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

type QwenCodeProvider struct {
	Token *config.QwenCodeToken
}

func NewQwenCodeProvider(token *config.QwenCodeToken) *QwenCodeProvider {
	return &QwenCodeProvider{Token: token}
}

func (p *QwenCodeProvider) HandleChatCompleteRequest(ctx context.Context, r *openai.OpenAIChatCompletionRequest) (*request.Response, error) {
	var qwenCodeChatCompleteRequest qwencode.QwenCodeChatCompleteRequest
	qwenCodeChatCompleteRequest.ImportOpenAIChatCompletionRequest(ctx, r)
	reqBody, err := json.Marshal(qwenCodeChatCompleteRequest)
	if err != nil {
		slog.Error("")
		return nil, err
	}

	creds, err := GetQwenCodeToken(ctx, p.Token)
	if err != nil {
		slog.Error("Failed to get Qwen code token", "error", err)
		return nil, err
	}

	return request.NewRequestBuilder(fmt.Sprintf("https://%s/v1/chat/completions", creds.ResourceURL), http.MethodPost).
		WithHeaders(map[string]string{
			constant.Authorization: fmt.Sprintf("Bearer %s", creds.AccessToken),
		}).
		WithJson(bytes.NewReader(reqBody)).
		Do(ctx, nil)
}

func (p *QwenCodeProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	defer r.Body.Close()
	w.WriteHeader(r.StatusCode)
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	_, err := io.Copy(w, r.Body)
	if err != nil {
		slog.Error("Failed to copy response body", "error", err)
	}
	return err
}

func (p *QwenCodeProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	return nil, nil
}

func (p *QwenCodeProvider) HandleListModelResponse(ctx context.Context, r *request.Response) (*openai.OpenAIModelListResponse, error) {
	var modelList *qwencode.QwenModelListResponse
	return modelList.Convert(ctx)
}
