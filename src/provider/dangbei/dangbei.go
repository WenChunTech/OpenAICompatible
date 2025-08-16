package dangbei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/dangbei"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
	"github.com/WenChunTech/OpenAICompatible/src/request"
	"github.com/WenChunTech/OpenAICompatible/src/responser"
	"github.com/WenChunTech/OpenAICompatible/src/util"
)

var Provider *DangBeiProvider

const createRequest = `{"conversationList":[{"metaData":{"chatModelConfig":{},"superAgentPath":"/chat"},"shareId":"","isAnonymous":false,"source":""}]}`
const SignURL = "https://ai-dangbei.deno.dev/"

type DangBeiProvider struct {
	*provider.BaseProvider
}

type SignResponse struct {
	Success bool `json:"success"`
	Data    Sign `json:"data"`
}
type Sign struct {
	Nonce     string `json:"nonce"`
	Sign      string `json:"sign"`
	Timestamp int64  `json:"timestamp"`
}

func init() {
	Provider = NewDangBeiProvider()
}

func NewDangBeiProvider() *DangBeiProvider {
	return &DangBeiProvider{
		BaseProvider: &provider.BaseProvider{
			ChatCompleteURL:    constant.DangBeiChatCompleteURL,
			ChatCompleteMethod: http.MethodPost,
			ModelURL:           constant.DangBeiModelListURL,
			ModelMethod:        http.MethodGet,
		},
	}
}

func (p *DangBeiProvider) HandleChatCompleteRequest(ctx context.Context, r *openai.OpenAIChatCompletionRequest) (*request.Response, error) {
	deviceID := util.GenerateUUID()
	timestamp := time.Now().Unix()
	sign := util.Sign(timestamp, createRequest)
	headers := map[string]string{
		"deviceId":  deviceID,
		"nonce":     util.Nonce,
		"sign":      sign,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}

	resp, err := request.NewRequestBuilder(constant.DangBeiChatCreateURL, http.MethodPost).WithHeaders(headers).WithJson(strings.NewReader(createRequest)).Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to create chat", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	var createResp dangbei.DangBeiChatCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	if err != nil {
		slog.Error("Failed to decode create chat response", "error", err)
		return nil, err
	}

	payload := fmt.Sprintf(`{"timestamp":%d}`, time.Now().UnixMilli())
	sign = util.Sign(timestamp, payload)
	headers = map[string]string{
		"deviceId":  deviceID,
		"nonce":     util.Nonce,
		"sign":      sign,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}
	resp, err = request.NewRequestBuilder(constant.DangBeiGenerateIdURL, http.MethodPost).WithHeaders(headers).WithJson(strings.NewReader(payload)).Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to generate chat ID", "error", err)
		return nil, err
	}
	ctx = context.WithValue(ctx, constant.ConversationID, createResp.Data.ConversationID)

	var generateResp dangbei.DangBeiGenerateIdResponse
	err = json.NewDecoder(resp.Body).Decode(&generateResp)
	if err != nil {
		slog.Error("Failed to decode generate chat ID response", "error", err)
		return nil, err
	}
	ctx = context.WithValue(ctx, constant.DangbeiChatID, generateResp.Data)

	var bangbeiChatCompleteRequest dangbei.DangBeiChatCompleteRequest
	err = bangbeiChatCompleteRequest.ImportOpenAIChatCompletionRequest(ctx, r)
	if err != nil {
		slog.Error("Failed to import openai chat completion request", "error", err)
		return nil, err
	}

	reqBody, err := json.Marshal(bangbeiChatCompleteRequest)
	if err != nil {
		slog.Error("Failed to marshal dangbei chat complete request", "error", err)
		return nil, err
	}

	slog.Info(string(reqBody))
	resp, err = request.NewRequestBuilder(SignURL, http.MethodPost).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to sign request", "error", err)
		return nil, err
	}
	var signResp SignResponse
	err = json.NewDecoder(resp.Body).Decode(&signResp)
	if err != nil {
		slog.Error("Failed to decode sign response", "error", err)
		return nil, err
	}

	slog.Info("Sign response", "data", signResp.Data)
	headers = map[string]string{
		"content-type": "application/json",
		"deviceId":     deviceID,
		"nonce":        signResp.Data.Nonce,
		"sign":         signResp.Data.Sign,
		"timestamp":    strconv.FormatInt(signResp.Data.Timestamp, 10),
	}

	return request.NewRequestBuilder(p.ChatCompleteURL, p.ChatCompleteMethod).WithHeaders(headers).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
}

func (p *DangBeiProvider) HandleChatCompleteResponse(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	handler := responser.EventStreamHandler[dangbei.DangBeiChatCompleteResponse]{}
	return handler.Handle(ctx, w, r)
}

func (p *DangBeiProvider) HandleListModelRequest(ctx context.Context) (*request.Response, error) {
	deviceID := util.GenerateUUID()
	timestamp := time.Now().Unix()
	sign := util.Sign(timestamp, "")
	headers := map[string]string{
		"deviceId":  deviceID,
		"nonce":     util.Nonce,
		"sign":      sign,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}
	return request.NewRequestBuilder(p.ModelURL, p.ModelMethod).WithHeaders(headers).Do(ctx, nil)
}

func (p *DangBeiProvider) HandleListModelResponse(ctx context.Context, r *request.Response) (*openai.OpenAIModelListResponse, error) {
	handler := responser.ModelListHandler[*dangbei.DangBeiModelListResponse]{}
	return handler.Handle(ctx, r)
}
