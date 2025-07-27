package parser

import (
	"encoding/json"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/converter"
	"github.com/WenChunTech/OpenAICompatible/src/eventsource"
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
func Parse[T converter.ChatCompletionConverter](p *Parser, data []byte) []eventsource.EventSourceResponse[T] {
	p.reminder.Write(data)
	fullData := p.reminder.String()
	fullData = strings.ReplaceAll(fullData, "\r\n", "\n")
	p.reminder.Reset()

	parts := strings.Split(fullData, "\n\n")
	// The last part may be incomplete. Save it for the next call.
	lastPartIndex := len(parts) - 1
	if lastPartIndex >= 0 && !strings.HasSuffix(fullData, "\n\n") {
		p.reminder.WriteString(parts[lastPartIndex])
		parts = parts[:lastPartIndex]
	}

	events := make([]eventsource.EventSourceResponse[T], 0, len(parts))
	for _, part := range parts {
		if event, ok := parsePart[T](part); ok {
			events = append(events, event)
		}
	}

	if len(events) == 0 {
		return nil
	}
	return events
}

// parsePart parses a single event part and returns the event and a boolean indicating if it was successful.
func parsePart[T converter.ChatCompletionConverter](part string) (eventsource.EventSourceResponse[T], bool) {
	if strings.TrimSpace(part) == "" {
		return eventsource.EventSourceResponse[T]{}, false
	}

	lines := reLine.Split(part, -1)
	var event eventsource.EventSourceResponse[T]
	var data strings.Builder

	for _, line := range lines {
		parseField(line, &event, &data)
	}

	dataStr := data.String()
	if dataStr != "" {
		parseDataField(dataStr, &event)
	}

	// An event is valid if it has any of the fields set.
	if event.ID != "" || event.Event != "" || dataStr != "" || event.Retry != 0 {
		return event, true
	}

	return eventsource.EventSourceResponse[T]{}, false
}

// parseField parses a single line of an event and updates the event and data builder.
func parseField[T converter.ChatCompletionConverter](line string, event *eventsource.EventSourceResponse[T], data *strings.Builder) {
	if len(line) == 0 {
		return
	}
	fields := strings.SplitN(line, ":", 2)
	if len(fields) != 2 {
		// According to the SSE spec, lines without a colon are ignored.
		return
	}

	field := strings.TrimSpace(fields[0])
	value := ""
	if len(fields) > 1 {
		value = strings.TrimPrefix(fields[1], " ")
	}

	switch field {
	case "id":
		event.ID = value
	case "event":
		event.Event = value
	case "data":
		if data.Len() > 0 {
			data.WriteRune('\n')
		}
		data.WriteString(value)
	case "retry":
		retry, err := strconv.Atoi(value)
		if err != nil {
			slog.Error("Failed to parse retry value", "error", err, "value", value, "line", line)
			return
		}
		event.Retry = retry
	}
}

// parseDataField parses the data field of an event.
func parseDataField[T converter.ChatCompletionConverter](dataStr string, event *eventsource.EventSourceResponse[T]) {
	if dataStr == "[DONE]" {
		event.Data = nil
	} else {
		// The data field is a JSON object, unmarshal it.
		var data T
		err := json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			slog.Error("Failed to unmarshal data field", "error", err, "data", dataStr)
			return
		}
		event.Data = &data
	}
}
