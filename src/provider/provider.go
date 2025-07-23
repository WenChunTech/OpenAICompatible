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

type Provider[P importer.OpenAIChatCompletionImporter, C converter.Converter[O], O any] struct {
	Provider P
	Url      string
	Method   string
	Headers  map[string]string
}

func (p *Provider[P, C, O]) validateRequest(r *http.Request) (*model.OpenAIChatCompletionRequest, error) {
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

func (p *Provider[P, C, O]) serializeAndConvert(reqBody *model.OpenAIChatCompletionRequest) ([]byte, error) {
	provider := p.Provider
	provider.ImportOpenAIChatCompletionReq(reqBody)
	providerReqBody, err := json.Marshal(provider)
	if err != nil {
		slog.Error("Failed to marshal provider request", "error", err)
		return nil, fmt.Errorf("failed to marshal provider request: %w", err)
	}
	return providerReqBody, nil
}

func (p *Provider[P, C, O]) doRequest(requestBody []byte) (*request.Response, error) {
	builder := request.NewRequestBuilder(p.Url, p.Method)
	builder.SetJson(bytes.NewReader(requestBody))
	builder.AddHeaders(p.Headers)
	resp, err := builder.Do(nil)
	if err != nil {
		slog.Error("Request to provider failed", "error", err)
		return nil, fmt.Errorf("failed to proxy request: %w", err)
	}
	return resp, nil
}

func (p *Provider[P, C, O]) handleSSEResponse(w http.ResponseWriter, resp *request.Response) {
	w.Header().Set(constant.ContentType, constant.ContentTypeEventStream)
	w.Header().Set(constant.CacheControl, constant.CacheControlNoCache)
	w.Header().Set(constant.Connection, constant.ConnectionKeepAlive)

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

func (p *Provider[P, C, O]) Handle(w http.ResponseWriter, r *http.Request) {
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
