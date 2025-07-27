package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type NotFoundMessage struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Not Found Handler", "url", r.URL.String(), "method", r.Method, "headers", r.Header)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
	}
	defer r.Body.Close()
	slog.Info("Request body", "body", string(body))

	msg := NotFoundMessage{
		Code:        "404",
		Message:     "Not Found",
		Description: "The requested resource was not found on the server.",
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal not found message", "error", err)
	}
	w.WriteHeader(http.StatusNotFound)
	_, err = w.Write(jsonMsg)
	if err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}
