package model

// OpenAIChatMessage represents a single message in a chat completion request or response.
type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatCompletionRequest represents the request body sent to the chat completion proxy.
type OpenAIChatCompletionRequest struct {
	Model    string              `json:"model"`
	Messages []OpenAIChatMessage `json:"messages"`
}

// OpenAIChatCompletion represents the non-streaming chat completion response.
type OpenAIChatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int               `json:"index"`
		Message      OpenAIChatMessage `json:"message"`
		FinishReason string            `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// OpenAIStreamChoice represents a choice in a streaming chat completion response.
type OpenAIStreamChoice struct {
	Index int `json:"index"`
	Delta struct {
		Role    *string `json:"role,omitempty"`
		Content *string `json:"content,omitempty"`
	} `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

// OpenAPIChatCompletionStreamResponse represents a streaming chat completion response.
type OpenAPIChatCompletionStreamResponse struct {
	ID                string               `json:"id"`
	Object            string               `json:"object"`
	Created           int64                `json:"created"`
	Model             string               `json:"model"`
	SystemFingerprint string               `json:"system_fingerprint,omitempty"`
	Choices           []OpenAIStreamChoice `json:"choices"`
}

// Model represents a single model listing in the OpenAI API.
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelList represents a list of models in the OpenAI API.
type ModelList struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}
