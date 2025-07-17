package parser

import (
	"encoding/json"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/converter"
	"github.com/WenChunTech/OpenAICompatible/src/sse"
)

var reLine = regexp.MustCompile(`\n`)

type Parser struct {
	reminder *strings.Builder
}

func NewParser() *Parser {
	return &Parser{
		reminder: &strings.Builder{},
	}
}

// Parse processes a chunk of SSE data and returns any complete events.
func Parse[O any, T converter.Converter[O]](p *Parser, data []byte) []sse.SSEEventResponse[O, T] {
	p.reminder.Write(data)
	fullData := p.reminder.String()
	fullData = strings.ReplaceAll(fullData, "\r\n", "\n")
	p.reminder.Reset()

	parts := strings.Split(fullData, "\n\n")
	// The last part may be incomplete. Save it for the next call.
	lastPartIndex := len(parts) - 1
	if lastPartIndex >= 0 && fullData[len(fullData)-len("\n\n"):] != "\n\n" {
		p.reminder.WriteString(parts[lastPartIndex])
		parts = parts[:lastPartIndex]
	}

	events := make([]sse.SSEEventResponse[O, T], 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}

		lines := reLine.Split(part, -1)
		var event sse.SSEEventResponse[O, T]
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			fields := strings.SplitN(line, ":", 2)
			if len(fields) != 2 {
				continue
			}

			field := strings.TrimSpace(fields[0])
			value := strings.TrimSpace(fields[1])

			switch field {
			case "id":
				event.ID = value
			case "event":
				event.Event = value
			case "data":
				if len(value) == 0 {
					continue
				}
				err := json.Unmarshal([]byte(value), &event.Data)
				if err != nil {
					slog.Error("Failed to unmarshal data field", "error", err, "data", value, "line", line)
				}

			case "retry":
				retry, err := strconv.Atoi(value)
				if err != nil {
					slog.Error("Failed to parse retry value", "error", err, "value", value, "line", line)
					continue
				}
				event.Retry = retry
			}
		}

		if event.ID != "" {
			events = append(events, event)
		}
	}

	if len(events) == 0 {
		return nil
	}
	return events
}
