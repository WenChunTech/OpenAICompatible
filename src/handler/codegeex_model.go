package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

// ChatProxyModelHandler handles the request for listing models and returns them in the OpenAI format.
func ChatProxyModelHandler(w http.ResponseWriter, r *http.Request) {
	// Build the request to CodeGeex
	builder := request.NewRequestBuilder(constant.CodeGeexModelURL, http.MethodGet)
	resp, err := builder.Do(nil)
	if err != nil {
		slog.Error("Request to CodeGeex failed", "error", err)
		http.Error(w, "Failed to fetch models from CodeGeex", http.StatusInternalServerError)
		return
	}

	// Decode the response
	var codegeexModels model.CodeGeexModelOptions
	if err := resp.Json(&codegeexModels); err != nil {
		slog.Error("Failed to decode CodeGeex models response", "error", err)
		http.Error(w, "Failed to decode response from CodeGeex", http.StatusInternalServerError)
		return
	}

	// Convert to OpenAI format
	openaiModelList, err := codegeexModels.Convert()
	if err != nil {
		slog.Error("Failed to convert to OpenAI model list", "error", err)
		http.Error(w, "Failed to convert to OpenAI model list", http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openaiModelList); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
