package model

import (
	"context"
	"time"
)

// ErrorMessage encapsulates an error with an associated HTTP status code.
type ErrorMessage struct {
	StatusCode int
	Error      error
}

// GCPProject represents the response structure for a Google Cloud project list request.
type GCPProject struct {
	Projects []GCPProjectProjects `json:"projects"`
}

// GCPProjectLabels defines the labels associated with a GCP project.
type GCPProjectLabels struct {
	GenerativeLanguage string `json:"generative-language"`
}

// GCPProjectProjects contains details about a single Google Cloud project.
type GCPProjectProjects struct {
	ProjectNumber  string           `json:"projectNumber"`
	ProjectID      string           `json:"projectId"`
	LifecycleState string           `json:"lifecycleState"`
	Name           string           `json:"name"`
	Labels         GCPProjectLabels `json:"labels"`
	CreateTime     time.Time        `json:"createTime"`
}

func (c *GeminiChatCompletionRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *OpenAIChatCompletionRequest) error {
	contents := []Content{}
	for _, i := range req.Messages {
		contents = append(contents, Content{
			Role: i.Role,
			Parts: []Part{
				{
					Text: i.Content[0].Text,
				},
			},
		})
	}
	projectID := ctx.Value("project_id").(string)
	*c = GeminiChatCompletionRequest{
		Model:   req.Model,
		Project: projectID,
		Request: Request{
			Contents: contents,
			GenerationConfig: GenerationConfig{
				ThinkingConfig: GenerationConfigThinkingConfig{
					IncludeThoughts: true,
				},
			},
		},
	}

	return nil
}

type GeminiChatCompletionRequest struct {
	Model   string  `json:"model"`
	Project string  `json:"project"`
	Request Request `json:"request"`
}

type Request struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
}

// Content represents a single message in a conversation, with a role and parts.
type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

// Part represents a distinct piece of content within a message, which can be
// text, inline data (like an image), a function call, or a function response.
type Part struct {
	Text             string            `json:"text,omitempty"`
	InlineData       *InlineData       `json:"inlineData,omitempty"`
	FunctionCall     *FunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *FunctionResponse `json:"functionResponse,omitempty"`
}

// InlineData represents base64-encoded data with its MIME type.
type InlineData struct {
	MimeType string `json:"mime_type,omitempty"`
	Data     string `json:"data,omitempty"`
}

// FunctionCall represents a tool call requested by the model, including the
// function name and its arguments.
type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// FunctionResponse represents the result of a tool execution, sent back to the model.
type FunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

// GenerateContentRequest is the top-level request structure for the streamGenerateContent endpoint.
type GenerateContentRequest struct {
	SystemInstruction *Content          `json:"systemInstruction,omitempty"`
	Contents          []Content         `json:"contents"`
	Tools             []ToolDeclaration `json:"tools,omitempty"`
	GenerationConfig  `json:"generationConfig"`
}

// GenerationConfig defines parameters that control the model's generation behavior.
type GenerationConfig struct {
	ThinkingConfig GenerationConfigThinkingConfig `json:"thinkingConfig,omitempty"`
	Temperature    float64                        `json:"temperature,omitempty"`
	TopP           float64                        `json:"topP,omitempty"`
	TopK           float64                        `json:"topK,omitempty"`
}

// GenerationConfigThinkingConfig specifies configuration for the model's "thinking" process.
type GenerationConfigThinkingConfig struct {
	// IncludeThoughts determines whether the model should output its reasoning process.
	IncludeThoughts bool `json:"include_thoughts,omitempty"`
}

// ToolDeclaration defines the structure for declaring tools (like functions)
// that the model can call.
type ToolDeclaration struct {
	FunctionDeclarations []interface{} `json:"functionDeclarations"`
}

func (c *GeminiChatCompletionResponse) Convert(ctx context.Context) (*OpenAPIChatCompletionStreamResponse, error) {
	return &OpenAPIChatCompletionStreamResponse{}, nil
}

type GeminiChatCompletionResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Candidates    []GeminiCandidate `json:"candidates"`
	UsageMetadata UsageMetadata     `json:"usageMetadata"`
	ModelVersion  string            `json:"modelVersion"`
	CreateTime    time.Time         `json:"createTime"`
	ResponseId    string            `json:"responseId"`
}

type GeminiCandidate struct {
	Content Content `json:"content"`
}

type UsageMetadata struct {
	TrafficType string `json:"trafficType"`
}
