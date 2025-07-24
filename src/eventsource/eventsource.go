package eventsource

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/WenChunTech/OpenAICompatible/src/converter"
)

type EventSourceResponse[T converter.ChatCompletionConverter] struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Data  *T     `json:"data"`
	Retry int    `json:"retry,omitempty"`
}

func (p *EventSourceResponse[T]) Json2EventSource() (string, error) {
	data := *p.Data
	openapiData, err := data.Convert()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(openapiData)
	if err != nil {
		slog.Error("marshal json data error", slog.Any("openapiData", openapiData))
	}

	return fmt.Sprintf("data: %s\n\n", string(jsonData)), nil
}
