package sse

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/WenChunTech/OpenAICompatible/src/converter"
)

type SSEEventResponse[O any, T converter.Converter[O]] struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Data  T      `json:"data"`
	Retry int    `json:"retry,omitempty"`
}

func (p *SSEEventResponse[O, T]) SSEJson2Text() (string, error) {
	openapiData, err := p.Data.Convert()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(openapiData)
	if err != nil {
		slog.Error("marshal json data error", slog.Any("openapiData", openapiData))
	}

	return fmt.Sprintf("id: %s\n"+
		"event: %s\n"+
		"data: %s\n"+
		"retry: %d\n"+
		"\n", p.ID, p.Event, string(jsonData), p.Retry), nil
}
