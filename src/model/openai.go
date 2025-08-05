package model

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Index    string   `json:"index"`
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function,omitempty"`
}

type Logprobs struct {
	TokenIds      []int    `json:"token_ids"`
	Tokens        []string `json:"tokens"`
	TokenLogprobs []int    `json:"token_logprobs"`
}

// type StreamOption struct {
// 	IncludeUsage bool `json:"include_usage"`
// }

type StructureContent struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageUrl *ImageUrl `json:"image_url,omitempty"`
}

type ImageUrl struct {
	Url string `json:"url"`
}

// OpenAIChatMessage represents a single message in a chat completion request or response.
type OpenAIChatMessage struct {
	Role    string             `json:"role"`
	Content []StructureContent `json:"content"`
}

// "tools": [
//
//	  {
//	    "type": "function",
//	    "function": {
//	      "name": "get_current_weather",
//	      "description": "Get the current weather in a given location",
//	      "parameters": {
//	        "type": "object",
//	        "properties": {
//	          "location": {
//	            "type": "string",
//	            "description": "The city and state, e.g. San Francisco, CA"
//	          },
//	          "unit": {
//	            "type": "string",
//	            "enum": ["celsius", "fahrenheit"]
//	          }
//	        },
//	        "required": ["location"]
//	      }
//	    }
//	  }
//	],
//	"tool_choice": "auto"
type Tool struct {
	Type     string          `json:"type"`
	Function RequestFunction `json:"function"`
}

type RequestFunction struct {
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Parameters  map[string]map[string]interface{} `json:"parameters"`
	Required    []string                          `json:"required"`
}

// OpenAIChatCompletionRequest represents the request body sent to the chat completion proxy.
type OpenAIChatCompletionRequest struct {
	Model       string              `json:"model"`
	Messages    []OpenAIChatMessage `json:"messages,omitempty"`
	MaxTokens   int                 `json:"max_tokens"`
	Stream      bool                `json:"stream"`
	Tools       []Tool              `json:"tools,omitempty"`
	ToolChoice  string              `json:"tool_choice,omitempty"`
	Logprobs    bool                `json:"logprobs"`
	TopLogprobs int                 `json:"top_logprobs"`
}

type Choice struct {
	Text         string            `json:"text"`
	Index        int               `json:"index"`
	Seed         int               `json:"seed"`
	Message      OpenAIChatMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

type Warning struct {
	Message string `json:"message"`
}

// OpenAIChatCompletion represents the non-streaming chat completion response.
type OpenAIChatCompletion struct {
	ID       string    `json:"id"`
	Object   string    `json:"object"`
	Created  int64     `json:"created"`
	Model    string    `json:"model"`
	Warnings []Warning `json:"warnings"`
	Choices  []Choice  `json:"choices"`
	Usage    *Usage    `json:"usage"`
}

type Delta struct {
	Role             string    `json:"role"`
	Content          string    `json:"content"`
	ReasoningContent *string   `json:"reasoning_content"`
	ToolCall         *ToolCall `json:"tool_call"`
}

// OpenAIStreamChoice represents a choice in a streaming chat completion response.
type OpenAIStreamChoice struct {
	Index        int       `json:"index"`
	Delta        Delta     `json:"delta"`
	FinishReason *string   `json:"finish_reason"`
	Logprobs     *Logprobs `json:"logprobs"`
	MatchedSotp  *string   `json:"matched_stop"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAPIChatCompletionStreamResponse represents a streaming chat completion response.
type OpenAPIChatCompletionStreamResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
	Choices []OpenAIStreamChoice `json:"choices"`
	Usage   *Usage               `json:"usage"`
}

// Model represents a single model listing in the OpenAI API.
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// OpenAIModelListResponse represents a list of models in the OpenAI API.
type OpenAIModelListResponse struct {
	Object string   `json:"object"`
	Data   []*Model `json:"data"`
}
