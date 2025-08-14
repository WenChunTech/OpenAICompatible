package geminicli

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	geminicli "github.com/WenChunTech/OpenAICompatible/src/model/gemini_cli"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
	"github.com/WenChunTech/OpenAICompatible/src/request"
	"github.com/WenChunTech/OpenAICompatible/src/responser"
	"golang.org/x/oauth2"
)

var Provider *GeminiCliProvider
var CacheMap map[string]*oauth2.Token

type GeminiCliProvider struct {
	*provider.BaseProvider
	ProjectID string
	Token     *oauth2.Token
}

func init() {
	config := config.GetGeminiCliConfig()
	Provider = NewGeminiCliProvider(config.ProjectID, config.Token)
}

func NewGeminiCliProvider(projectID string, token *oauth2.Token) *GeminiCliProvider {
	return &GeminiCliProvider{
		BaseProvider: &provider.BaseProvider{
			ChatCompleteURL:    constant.GeminiCliChatCompleteURL,
			ChatCompleteMethod: http.MethodPost,
			ModelURL:           constant.GeminiCliModelListURL,
			ModelMethod:        http.MethodGet,
		},
		ProjectID: projectID,
		Token:     token,
	}
}

func (p *GeminiCliProvider) getTW(ctx context.Context) (*TokenWrapper, error) {
	token, ok := CacheMap[p.ProjectID]
	if !ok {
		token = p.Token
	}
	tw := NewTokenWrapper(ctx, token)
	token, err := tw.GetToken()
	if err != nil {
		slog.Error("get token failed", "err", err)
		return nil, err
	}
	CacheMap[p.ProjectID] = token
	return tw, nil
}

func (p *GeminiCliProvider) HandleChatCompleteRequest(ctx context.Context, r *openai.OpenAIChatCompletionRequest) (*request.Response, error) {
	tw, err := p.getTW(ctx)
	if err != nil {
		slog.Error("get token wrapper failed", "err", err)
		return nil, err
	}
	ctx = context.WithValue(ctx, constant.ProjectIDKey, p.ProjectID)
	var geminiCliRequest = &geminicli.GeminiCliChatCompletionRequest{}
	err = geminiCliRequest.ImportOpenAIChatCompletionRequest(ctx, r)
	if err != nil {
		slog.Error("import openai chat completion request failed", "err", err)
		return nil, err
	}
	reqBody, err := json.Marshal(geminiCliRequest)
	if err != nil {
		slog.Error("marshal request body failed", "err", err)
		return nil, err
	}
	return request.NewRequestBuilder(p.ChatCompleteURL, p.ChatCompleteMethod).WithJson(bytes.NewReader(reqBody)).Do(ctx, tw.GetClient())
}

func (p *GeminiCliProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.EventStreamHandler[*geminicli.GeminiChatCompletionResponse]{}
	return handler.Handle(ctx, w, r)
}

// HandleListModelRequest 函数用于处理列表模型请求
func (p *GeminiCliProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	// 创建一个新的请求构建器，传入模型URL和模型方法
	return nil, nil
}

func (p *GeminiCliProvider) HandleListModelResponse(ctx context.Context, r *request.Response) (*openai.OpenAIModelListResponse, error) {
	var modelList geminicli.GeminiCliModelListResponse
	if err := json.NewDecoder(strings.NewReader(geminicli.GEMINI_CLI_MODEL_LIST_RESPONSE)).Decode(&modelList); err != nil {
		slog.Error("decode model list failed", "err", err)
		return nil, err
	}

	openaiModelList, err := modelList.Convert(ctx)
	if err != nil {
		slog.Error("convert model list failed", "err", err)
		return nil, err
	}

	return openaiModelList, nil
}
