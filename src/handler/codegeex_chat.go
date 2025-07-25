package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model"
	"github.com/WenChunTech/OpenAICompatible/src/parser"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

// ChatProxyChatHandler handles the proxying of chat completion requests to the CodeGeex API.
func ChatProxyChatHandler(w http.ResponseWriter, r *http.Request) {
	openAPIReq, err := CheckParamterAndDecodeRequest(r, w)
	if err != nil {
		slog.Error("Failed to check and decode request", "error", err)
		http.Error(w, "Failed to check and decode request", http.StatusBadRequest)
		return
	}

	// Extract the last user message as the prompt.
	geexReqBody, err := CreateRequestBody(openAPIReq, w)
	if err != nil {
		slog.Error("Failed to create request body", "error", err)
		http.Error(w, "Failed to create request body", http.StatusInternalServerError)
		return
	}

	// Create and send the request to CodeGeex.
	builder := BuildRequest(geexReqBody)
	resp, err := builder.Do(nil)
	if err != nil {
		slog.Error("Request to CodeGeex failed", "error", err)
		http.Error(w, "Failed to proxy request", http.StatusInternalServerError)
		return
	}

	// Set response headers for SSE.
	err = HandleResponse(w, resp)
	if err != nil {
		slog.Error("Failed to handle response", "error", err)
		return
	}

}

func HandleResponse(w http.ResponseWriter, resp *request.Response) error {
	w.Header().Set(constant.ContentType, constant.ContentTypeEventStream)
	w.Header().Set(constant.CacheControl, constant.CacheControlNoCache)
	w.Header().Set(constant.Connection, constant.ConnectionKeepAlive)

	sseParser := parser.NewParser()
	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Streaming unsupported")
		w.Write([]byte("Streaming unsupported"))
		return errors.New("streaming unsupported")
	}
	dataChan, errChan := resp.EventStream()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				slog.Error("Error from event stream", "error", err)
				http.Error(w, "Error reading from upstream", http.StatusInternalServerError)
				return err // Return the error
			}
		case buf, ok := <-dataChan:
			if !ok {
				return nil // Channel closed, successful completion
			}
			events := parser.Parse[model.CodeGeexEventSourceData](sseParser, buf)
			for _, event := range events {
				openAPIData, err := event.Json2EventSource()
				if err != nil {
					slog.Error("Failed to convert SSE data", "error", err)
					continue
				}
				_, writeErr := w.Write([]byte(openAPIData))
				if writeErr != nil {
					slog.Error("Failed to write to response", "error", writeErr)
					// The connection is likely closed, so we should stop.
					return writeErr
				}
				flusher.Flush()
			}
		}
	}
}

func BuildRequest(geexReqBody []byte) *request.RequestBuilder {
	builder := request.NewRequestBuilder(constant.CodeGeexChatURL, http.MethodPost)
	builder.SetJson(bytes.NewReader(geexReqBody))

	headers := map[string]string{
		constant.Accept:    constant.ContentTypeEventStream,
		constant.CodeToken: config.Config.Token,
		constant.UserAgent: constant.DefaultUserAgent,
	}
	builder.AddHeaders(headers)
	return builder
}

func CreateRequestBody(openAPIReq *model.OpenAIChatCompletionRequest, w http.ResponseWriter) ([]byte, error) {
	var promptBuilder strings.Builder
	for _, message := range openAPIReq.Messages {
		messageContent, err := json.Marshal(message)
		if err != nil {
			slog.Error("Failed to marshal message", "error", err)
			return nil, err
		}
		promptBuilder.WriteString(string(messageContent))
		promptBuilder.WriteString("\n")
	}

	// Construct the CodeGeex request.
	// These values can be hardcoded for now but should be configurable in the future.
	geexReqPayload := model.CodeGeexChatRequest{
		UserID:        config.Config.UserID,
		UserRole:      config.Config.UserRole,
		IDE:           config.Config.IDE,
		IDEVersion:    config.Config.IDEVersion,
		PluginVersion: config.Config.PluginVersion,
		Prompt:        promptBuilder.String(),
		MachineID:     config.Config.MachineID,
		TalkID:        config.Config.TalkID,
		Locale:        config.Config.Locale,
		Model:         openAPIReq.Model,
	}

	geexReqBody, err := json.Marshal(geexReqPayload)
	if err != nil {
		slog.Error("Failed to marshal CodeGeex request", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, err
	}
	return geexReqBody, nil
}

func CheckParamterAndDecodeRequest(r *http.Request, w http.ResponseWriter) (*model.OpenAIChatCompletionRequest, error) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, errors.New("method not allowed")
	}

	var openAPIReq model.OpenAIChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&openAPIReq); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return nil, errors.New("failed to decode request body")
	}
	return &openAPIReq, nil
}
