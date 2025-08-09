package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
)

type Handler[P provider.Provider] struct {
	P P
}

func (p *Handler[P]) validateRequest(r *http.Request) (*openai.OpenAIChatCompletionRequest, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var reqBody openai.OpenAIChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	return &reqBody, nil
}

func (p *Handler[P]) HandleChatComplete(w http.ResponseWriter, r *http.Request) {
	slog.Info("Start request chat complete")
	ctx := r.Context()
	reqBody, err := p.validateRequest(r)
	if err != nil {
		if errors.Is(err, errors.New("method not allowed")) {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	resp, err := p.P.HandleChatCompleteRequest(ctx, reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.P.HandleChatCompleteResponse(ctx, w, resp)
	if err != nil {
		slog.Error("HandleChatCompleteResponse failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Handler[P]) HandleListModel(w http.ResponseWriter, r *http.Request) {
	slog.Info("Start request list model")
	ctx := r.Context()
	resp, err := p.P.HandleListModelRequest(ctx)
	if err != nil {
		slog.Error("HandleListModelRequest failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	openaiModelList, err := p.P.HandleListModelResponse(ctx, resp)
	if err != nil {
		slog.Error("HandleListModelResponse failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set(constant.ContentType, constant.ContentTypeJson)
	if err := json.NewEncoder(w).Encode(openaiModelList); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
