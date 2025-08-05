package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	u := "https://chat.qwen.ai/api/v2/chat/completions?chat_id=1cf2fb60-9e99-45db-9719-34ccc607e58f"
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		t.Fatal(err)
	}

	if parsedURL.Host != "chat.qwen.ai" {
		t.Errorf("Expected host to be 'chat.qwen.ai', got '%s'", parsedURL.Host)
	}
	if parsedURL.Path != "/api/v2/chat/completions" {
		t.Errorf("Expected path to be '/api/v2/chat/completions', got '%s'", parsedURL.Path)
	}
	if parsedURL.Query().Get("chat_id") != "1cf2fb60-9e99-45db-9719-34ccc607e58f" {
		t.Errorf("Expected chat_id to be '1cf2fb60-9e99-45db-9719-34ccc607e58f', got '%s'", parsedURL.Query().Get("chat_id"))
	}
}

type Content struct {
	Content []ContentItem `json:"content"`
}

type ContentItem struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL struct {
		URL string `json:"url,omitempty"`
	} `json:"image_url,omitempty"`
}

// test the json unmarshalling
func TestUnmarshal(t *testing.T) {
	s := `{"content": [
          {
            "type": "text",
            "text": "What is in this image?"
          },
          {
            "type": "image_url",
            "image_url": {
              "url": "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg"
            }
          }
        ]}`
	var content Content
	err := json.Unmarshal([]byte(s), &content)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(content.Content)
}

// ChatCompletionRequest represents the request body for the OpenAI Chat Completions API.
// Based on documentation from https://platform.openai.com/docs/api-reference/chat/create
type ChatCompletionRequest struct {
	Model               string                  `json:"model"`
	Messages            []ChatCompletionMessage `json:"messages"`
	Audio               *AudioOutputConfig      `json:"audio,omitempty"` // For audio output requests
	FrequencyPenalty    *float64                `json:"frequency_penalty,omitempty"`
	LogitBias           map[string]int          `json:"logit_bias,omitempty"`
	LogProbs            *bool                   `json:"logprobs,omitempty"`
	TopLogProbs         *int                    `json:"top_logprobs,omitempty"`
	MaxCompletionTokens *int                    `json:"max_completion_tokens,omitempty"` // Preferred over MaxTokens
	Metadata            map[string]string       `json:"metadata,omitempty"`              // Up to 16 key-value pairs
	Modalities          []string                `json:"modalities,omitempty"`            // e.g., ["text"], ["text", "audio"]
	N                   *int                    `json:"n,omitempty"`
	ParallelToolCalls   *bool                   `json:"parallel_tool_calls,omitempty"`
	Prediction          *json.RawMessage        `json:"prediction,omitempty"` // string or PredictionObject
	PresencePenalty     *float64                `json:"presence_penalty,omitempty"`
	PromptCacheKey      *string                 `json:"prompt_cache_key,omitempty"`
	ReasoningEffort     *string                 `json:"reasoning_effort,omitempty"` // o-series models
	ResponseFormat      *json.RawMessage        `json:"response_format,omitempty"`  // Can be TextResponseFormat, JSONObjectResponseFormat, JSONSchemaResponseFormat
	SafetyIdentifier    *string                 `json:"safety_identifier,omitempty"`
	Seed                *int                    `json:"seed,omitempty"`
	ServiceTier         *string                 `json:"service_tier,omitempty"`
	Stop                *json.RawMessage        `json:"stop,omitempty"` // string or []string
	Store               *bool                   `json:"store,omitempty"`
	Stream              *bool                   `json:"stream,omitempty"`
	StreamOptions       *StreamOptions          `json:"stream_options,omitempty"`
	Temperature         *float64                `json:"temperature,omitempty"`
	TopP                *float64                `json:"top_p,omitempty"`
	Tools               []Tool                  `json:"tools,omitempty"`
	ToolChoice          *json.RawMessage        `json:"tool_choice,omitempty"` // string or ToolChoiceObject
	WebSearchOptions    *WebSearchOptions       `json:"web_search_options,omitempty"`
}

// ChatCompletionMessage represents a single message in the conversation.
type ChatCompletionMessage struct {
	Role         string        `json:"role"`    // "system", "user", "assistant", "developer", "tool"
	Content      interface{}   `json:"content"` // string or []ContentPart
	Name         *string       `json:"name,omitempty"`
	Audio        *AudioData    `json:"audio,omitempty"`         // For previous assistant audio response
	Refusal      *string       `json:"refusal,omitempty"`       // For assistant refusal
	ToolCallID   *string       `json:"tool_call_id,omitempty"`  // Required for tool messages
	ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`    // For assistant tool calls
	FunctionCall *FunctionCall `json:"function_call,omitempty"` // Deprecated
}

// ContentPart represents a part of the message content (text, image, audio, file).
// You might need to implement custom marshaling/unmarshaling for this.
// type ContentPart interface{} // Could be TextContentPart, ImageContentPart, etc.

// TextContentPart represents textual content within a message.
type TextContentPart struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

// AudioOutputConfig represents parameters for audio output.
type AudioOutputConfig struct {
	Format string `json:"format"` // "wav", "mp3", "flac", "opus", "pcm16"
	Voice  string `json:"voice"`  // "alloy", "ash", etc.
}

// AudioData represents data about a previous audio response (placeholder).
type AudioData struct {
	ID string `json:"id"`
	// Add other fields as needed from the documentation
}

// Tool represents a function the model can call.
type Tool struct {
	Type     string             `json:"type"` // Currently, only "function" is supported
	Function FunctionDefinition `json:"function"`
}

// FunctionDefinition defines a function available for the model to call.
type FunctionDefinition struct {
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"` // JSON Schema object
	Strict      *bool       `json:"strict,omitempty"`
}

// ToolCall represents a call to a tool/function generated by the model.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // Currently, only "function" is supported
	Function FunctionCall `json:"function"`
}

// FunctionCall represents the name and arguments of a function call (generated or for tool messages).
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON formatted string
}

type PredictionObject struct {
	Content struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Type string `json:"type"`
}

// StreamOptions represents options for streaming responses.
type StreamOptions struct {
	// Add fields as defined in the API docs if needed
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

// WebSearchOptions represents options for the web search tool.
type WebSearchOptions struct {
	SearchContextSize *string       `json:"search_context_size,omitempty"` // "low", "medium", "high"
	UserLocation      *UserLocation `json:"user_location,omitempty"`
}

// UserLocation represents approximate location for the search.
type UserLocation struct {
	Approximate ApproximateLocation `json:"approximate"`
}

// ApproximateLocation represents the approximate location parameters.
type ApproximateLocation struct {
	// Define fields for approximate location (e.g., city, country, coordinates)
	// The API docs don't specify the exact structure, so you'll need to adapt this.
	Type string `json:"type"` // Always "approximate"
}

// Example Response Formats (you would use one of these for ResponseFormat field)
type TextResponseFormat struct {
	Type string `json:"type"` // Always "text"
}

type JSONObjectResponseFormat struct {
	Type string `json:"type"` // Always "json_object"
}

type JSONSchemaResponseFormat struct {
	Type       string `json:"type"` // Always "json_schema"
	JSONSchema struct {
		Name        string      `json:"name"`
		Description *string     `json:"description,omitempty"`
		Schema      interface{} `json:"schema,omitempty"` // JSON Schema object
		Strict      *bool       `json:"strict,omitempty"`
	} `json:"json_schema"`
}
