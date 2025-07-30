package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

var ForwardMap = map[string]string{
	constant.CodeGeexModelListPath:    constant.CodeGeexHost,
	constant.CodeGeexChatCompletePath: constant.CodeGeexHost,

	constant.QwenModelListPath:    constant.QwenHost,
	constant.QwenChatIDPath:       constant.QwenHost,
	constant.QwenChatCompletePath: constant.QwenHost,
}

func HandleForward(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.Background()
	host, ok := ForwardMap[r.URL.Path]
	if !ok {
		slog.Error("Not found", "path", r.URL.Path, "header", r.Header)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read request body", slog.Any("error", err))
		}
		slog.Error("Request body", "body", string(body))
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	r.URL.Scheme = constant.HTTPSScheme
	r.URL.Host = host
	resp, err := request.NewRequestBuilder(r.URL.String(), r.Method).WithHeader(r.Header).WithBody(r.Body).Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to build request", slog.Any("error", err))
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		slog.Error("Failed to copy response body", slog.Any("error", err))
		http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
	}
}
