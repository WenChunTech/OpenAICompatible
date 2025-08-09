package responser

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/converter"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/parser"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

type EventStreamHandler[C converter.ChatCompletionConverter] struct {
}

func (h *EventStreamHandler[C]) Handle(ctx context.Context, w http.ResponseWriter, r *request.Response) error {
	defer r.Body.Close()
	sseParser := parser.NewParser()
	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Streaming unsupported")
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set(constant.ContentType, constant.ContentTypeEventStream)
	w.Header().Set(constant.CacheControl, constant.CacheControlNoCache)
	w.Header().Set(constant.Connection, constant.ConnectionKeepAlive)
	dataChan, errChan := r.EventStream()
	var builder strings.Builder
	for {
		select {
		case err := <-errChan:
			if err != nil {
				slog.Error("Error from event stream", "error", err)
				return fmt.Errorf("error from event stream: %w", err)
			}
		case buf, ok := <-dataChan:
			if !ok {
				return nil
			}

			if json.Valid(buf) {
				slog.Error("sse data parse failed", "err_msg", string(buf))
				return fmt.Errorf("sse data parse failed: %s", string(buf))
			}

			events := parser.Parse[C](sseParser, buf)
			for _, event := range events {
				data, err := (*event.Data).Convert(ctx)
				if err != nil {
					slog.Error("Failed to convert SSE data", "error", err)
					continue
				}

				buf, err := json.Marshal(data)
				if err != nil {
					slog.Error("marshal json data error", slog.Any("data", data))
					return fmt.Errorf("marshal json data error: %w", err)
				}
				builder.WriteString("data: ")
				builder.Write(buf)
				builder.WriteString("\n\n")
				_, err = w.Write([]byte(builder.String()))
				if err != nil {
					slog.Error("Failed write response", "err", err)
				}
				flusher.Flush()
				builder.Reset()
			}
		}
	}
}

type ModelListHandler[C converter.ModelConverter] struct {
}

func (h *ModelListHandler[C]) Handle(ctx context.Context, r *request.Response) (*openai.OpenAIModelListResponse, error) {
	defer r.Body.Close()
	var providerModel C
	if err := r.Json(&providerModel); err != nil {
		slog.Error("Failed to decode provider models response", "error", err)
		return nil, fmt.Errorf("failed to decode provider models response: %w", err)
	}

	return providerModel.Convert(ctx)

}
