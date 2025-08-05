package openai

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type TopLogprob struct {
	Token   string `json:"token"`
	Logprob int    `json:"logprob"`
	Bytes   []int  `json:"bytes"`
}

type Content struct {
	Token       string       `json:"token"`
	Logprob     int          `json:"logprob"`
	Bytes       []int        `json:"bytes"`
	TopLogprobs []TopLogprob `json:"top_logprobs"`
}

type Logprobs struct {
	Content []Content `json:"content"`
}

type Delta struct {
	Role             string    `json:"role"`
	Content          string    `json:"content"`
	ReasoningContent *string   `json:"reasoning_content"`
	ToolCall         *ToolCall `json:"tool_call"`
}

type Message struct {
	Role        string    `json:"role"`
	Content     string    `json:"content"`
	ToolCall    *ToolCall `json:"tool_call,omitempty"`
	Refusal     *string   `json:"refusal,omitempty"`     // need to confirm type
	Annotations []string  `json:"annotations,omitempty"` // need to confirm type
}

type OpenAIStreamChoice struct {
	Index        int       `json:"index"`                   // response supports index
	Delta        *Delta    `json:"delta,omitempty"`         // stream response only supports delta
	Message      *Message  `json:"message,omitempty"`       // stream response does not support message
	FinishReason *string   `json:"finish_reason,omitempty"` // stream response does not support finish_reason
	Logprobs     *Logprobs `json:"logprobs,omitempty"`      // stream response does not support logprobs
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens,omitempty"`
	AudioTokens  int `json:"audio_tokens,omitempty"`
}

type CompletionTokensDetails struct {
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
	AudioTokens              int `json:"audio_tokens,omitempty"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
}

type Usage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	CompletionTokens        int                      `json:"completion_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"` // fucntions and logprobs are not supported
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details"`
}

type OpenAPIChatCompletionStreamResponse struct {
	ID                string               `json:"id"`
	Object            string               `json:"object"`
	Created           int64                `json:"created"`
	Model             string               `json:"model"`
	SystemFingerprint string               `json:"system_fingerprint"`
	Choices           []OpenAIStreamChoice `json:"choices"`
	Usage             *Usage               `json:"usage,omitempty"`        // stream response does not support usage
	ServiceTier       string               `json:"service_tier,omitempty"` // stream response does not support service_tier
}
