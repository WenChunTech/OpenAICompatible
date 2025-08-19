package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/provider"
	"github.com/WenChunTech/OpenAICompatible/src/provider/codegeex"
	"github.com/WenChunTech/OpenAICompatible/src/provider/dangbei"
	"github.com/WenChunTech/OpenAICompatible/src/provider/geminicli"
	"github.com/WenChunTech/OpenAICompatible/src/provider/qwen"
)

const Object = "model"

type ProviderManager struct {
	PrefixMap map[string]func() provider.Provider
	ModelList *openai.OpenAIModelListResponse
}

func newCodeGeexProvider() provider.Provider {
	config := config.NextCodeGeexConfig()
	return codegeex.NewCodeGeexProvider(config.Token)
}

func newQwenProvider() provider.Provider {
	config := config.NextQwenConfig()
	return qwen.NewQwenProvider(config.Token)
}

func newGeminiCliProvider() provider.Provider {
	config := config.NextGeminiCliConfig()
	return geminicli.NewGeminiCliProvider(config.ProjectID, config.Token)
}

func newDangBeiProvider() provider.Provider {
	return dangbei.NewDangBeiProvider()
}

func InitProviderManager() *ProviderManager {
	var err error
	context := context.Background()
	manager := NewProviderManager()

	if config.GetCodeGeexConfig() != nil {
		err = manager.RegisterProvider(context, constant.CodeGeexPrefix, newCodeGeexProvider)
		if err != nil {
			slog.Error("Failed to register codegeex provider", "error", err)
		}
	} else {
		slog.Warn("CodeGeex config is nil, skipping CodeGeex provider registration")
	}

	if config.GetQwenConfig() != nil {
		err = manager.RegisterProvider(context, constant.QwenPrefix, newQwenProvider)
		if err != nil {
			slog.Error("Failed to register qwen provider", "error", err)
		}
	} else {
		slog.Warn("Qwen config is nil, skipping Qwen provider registration")
	}

	if config.GetGeminiCliConfig() != nil {
		err = manager.RegisterProvider(context, constant.GeminiCliPrefix, newGeminiCliProvider)
		if err != nil {
			slog.Error("Failed to register gemini_cli provider", "error", err)
		}
	} else {
		slog.Warn("GeminiCli config is nil, skipping GeminiCli provider registration")
	}

	err = manager.RegisterProvider(context, constant.DangBeiPrefix, newDangBeiProvider)
	if err != nil {
		slog.Error("Failed to register dangbei provider", "error", err)
	}
	return manager
}

func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		PrefixMap: make(map[string]func() provider.Provider),
		ModelList: &openai.OpenAIModelListResponse{
			Object: Object,
		},
	}
}

func (m *ProviderManager) RegisterProvider(ctx context.Context, prefix string, providerFunc func() provider.Provider) error {
	provider := providerFunc()
	resp, err := provider.HandleListModelRequest(ctx)
	if err != nil {
		slog.Error("Failed to handle list model request", "error", err)
		return fmt.Errorf("failed to handle list model request: %w", err)
	}

	providerModelList, err := provider.HandleListModelResponse(ctx, resp)
	if err != nil {
		slog.Error("Failed to decode list model response", "error", err)
		return fmt.Errorf("failed to decode list model response: %w", err)
	}

	for _, model := range providerModelList.Data {
		prefixModel := fmt.Sprintf("%s/%s", prefix, model.ID)
		model.ID = prefixModel
		m.PrefixMap[prefixModel] = providerFunc
	}
	m.ModelList.Data = append(m.ModelList.Data, providerModelList.Data...)
	return nil
}

func (m *ProviderManager) validateRequest(r *http.Request) (*openai.OpenAIChatCompletionRequest, error) {
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

func (m *ProviderManager) HandleChatComplete(w http.ResponseWriter, r *http.Request) {
	slog.Info("Start request chat complete by handler manager")
	ctx := r.Context()
	reqBody, err := m.validateRequest(r)
	if err != nil {
		if errors.Is(err, errors.New("method not allowed")) {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if providerFunc, ok := m.PrefixMap[reqBody.Model]; ok {
		provider := providerFunc()
		index := strings.Index(reqBody.Model, "/")
		reqBody.Model = reqBody.Model[index+1:]
		resp, err := provider.HandleChatCompleteRequest(ctx, reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = provider.HandleChatCompleteResponse(ctx, w, resp)
		if err != nil {
			slog.Error("HandleChatCompleteResponse failed", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	slog.Error("Model not found", "model", reqBody.Model)
	http.Error(w, "Model not found", http.StatusNotFound)
}

func (m *ProviderManager) HandleListModel(w http.ResponseWriter, r *http.Request) {
	slog.Info("Start request list model by handler manager")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m.ModelList); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
