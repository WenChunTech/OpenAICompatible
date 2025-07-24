package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/converter"
	"github.com/WenChunTech/OpenAICompatible/src/importer"
	"github.com/WenChunTech/OpenAICompatible/src/model"
	"github.com/WenChunTech/OpenAICompatible/src/parser"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

type BaseProvider[P importer.OpenAIChatCompletionImporter, M converter.ModelConverter, C converter.ChatCompletionConverter] struct {
	ChatCompleteResp   P
	ModelResp          M
	ChatCompleteUrl    string
	ChatCompleteMethod string
	Headers            map[string]string
	ModelUrl           string
	ModelMethod        string
}

type Provider interface {
}

func (p *BaseProvider[P, M, C]) validateRequest(r *http.Request) (*model.OpenAIChatCompletionRequest, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var reqBody model.OpenAIChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	return &reqBody, nil
}

func (p *BaseProvider[P, M, C]) serializeAndConvert(reqBody *model.OpenAIChatCompletionRequest) ([]byte, error) {
	provider := p.ChatCompleteResp
	provider.ImportOpenAIChatCompletionReq(reqBody)
	providerReqBody, err := json.Marshal(provider)
	if err != nil {
		slog.Error("Failed to marshal provider request", "error", err)
		return nil, fmt.Errorf("failed to marshal provider request: %w", err)
	}
	return providerReqBody, nil
}

func (p *BaseProvider[P, M, C]) doRequest(requestBody []byte) (*request.Response, error) {
	builder := request.NewRequestBuilder(p.ChatCompleteUrl, p.ChatCompleteMethod)
	builder.SetJson(bytes.NewReader(requestBody))
	if p.Headers != nil {
		builder.AddHeaders(p.Headers)
	}
	resp, err := builder.Do(nil)
	if err != nil {
		slog.Error("Request to provider failed", "error", err)
		return nil, fmt.Errorf("failed to proxy request: %w", err)
	}
	return resp, nil
}

func (p *BaseProvider[P, M, C]) handleSSEResponse(w http.ResponseWriter, resp *request.Response) {
	w.Header().Set(constant.ContentType, constant.ContentTypeEventStream)
	w.Header().Set(constant.CacheControl, constant.CacheControlNoCache)
	w.Header().Set(constant.Connection, constant.ConnectionKeepAlive)
	defer resp.Body.Close()

	sseParser := parser.NewParser()
	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Streaming unsupported")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := resp.EventStream()
	for buf := range ch {
		events := parser.Parse[C](sseParser, buf)
		for _, event := range events {
			openAPIData, err := event.Json2EventSource()
			if err != nil {
				slog.Error("Failed to convert SSE data", "error", err)
				continue
			}
			w.Write([]byte(openAPIData))
			flusher.Flush()
		}
	}
}

func (p *BaseProvider[P, M, C]) ChatCompleteHandle(w http.ResponseWriter, r *http.Request) {
	reqBody, err := p.validateRequest(r)
	if err != nil {
		if errors.Is(err, errors.New("method not allowed")) {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	providerReqBody, err := p.serializeAndConvert(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := p.doRequest(providerReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.handleSSEResponse(w, resp)
}

func (p *BaseProvider[P, M, C]) ModelHandle(w http.ResponseWriter, r *http.Request) {
	builder := request.NewRequestBuilder(p.ModelUrl, p.ModelMethod)
	resp, err := builder.Do(nil)
	if err != nil {
		slog.Error("Request to CodeGeex failed", "error", err)
		http.Error(w, "Failed to fetch models from CodeGeex", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	m := p.ModelResp
	if err := resp.Json(m); err != nil {
		slog.Error("Failed to decode CodeGeex models response", "error", err)
		http.Error(w, "Failed to decode response from CodeGeex", http.StatusInternalServerError)
		return
	}

	openaiModelList, err := m.Convert()
	if err != nil {
		slog.Error("Failed to convert to OpenAI model list", "error", err)
		http.Error(w, "Failed to convert to OpenAI model list", http.StatusInternalServerError)
		return
	}

	w.Header().Set(constant.ContentType, constant.ContentTypeJson)
	if err := json.NewEncoder(w).Encode(openaiModelList); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
