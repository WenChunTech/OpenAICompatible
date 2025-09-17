package geminicli

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/translator"
)

func (c *GeminiCliChatCompletionRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error {
	var systemInstruction *Content
	contents := []Content{}
	toolItems := make(map[string]*FunctionResponse)
	for _, msg := range req.Messages {
		content := Content{}
		switch msg.Role {
		case "system":
			if text, ok := msg.Content.(string); ok {
				systemInstruction = &Content{
					Role:  "user",
					Parts: []Part{{Text: text}},
				}
			}
			continue
		case "user":
			content.Role = "user"
			switch v := msg.Content.(type) {
			case string:
				content.Parts = append(content.Parts, Part{Text: v})
			case []interface{}:
				for _, part := range v {
					if p, ok := part.(map[string]interface{}); ok {
						switch p["type"] {
						case "text":
							if text, ok := p["text"].(string); ok {
								content.Parts = append(content.Parts, Part{Text: text})
							}
						case "file":
							if fileInfo, ok := p["file"].(map[string]interface{}); ok {
								if filename, ok := fileInfo["filename"].(string); ok {
									if fileData, ok := fileInfo["file_data"].(string); ok {
										ext := ""
										if split := strings.Split(filename, "."); len(split) > 1 {
											ext = split[len(split)-1]
										}
										if mimeType, ok := translator.MimeTypes[ext]; ok {
											content.Parts = append(content.Parts, Part{
												InlineData: &InlineData{
													MimeType: mimeType,
													Data:     fileData,
												},
											})
										} else {
											slog.Warn("Unknown file name extension, skipping file", "ext", ext)
										}
									}
								}
							}
						case "image_url":
							if imageUrl, ok := p["image_url"].(map[string]interface{}); ok {
								if url, ok := imageUrl["url"].(string); ok && len(url) > 5 {
									parts := strings.SplitN(url[5:], ";", 2)
									if len(parts) == 2 && len(parts[1]) > 7 {
										content.Parts = append(content.Parts, Part{
											InlineData: &InlineData{
												MimeType: parts[0],
												Data:     parts[1][7:],
											},
										})
									}
								}
							}
						}
					}
				}
			}
		case "assistant":
			content.Role = "model"
			if text, ok := msg.Content.(string); ok && text != "" {
				content.Parts = append(content.Parts, Part{Text: text})
			}

			if len(msg.ToolCalls) > 0 {
				functionCallParts := []Part{}
				toolResponseParts := []Part{}

				for _, toolCall := range msg.ToolCalls {
					args := map[string]interface{}{}
					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err == nil {
						functionCallParts = append(functionCallParts, Part{
							FunctionCall: &FunctionCall{
								Name: toolCall.Function.Name,
								Args: args,
							},
						})
					}

					if response, ok := toolItems[toolCall.ID]; ok {
						response.Name = toolCall.Function.Name
						toolResponseParts = append(toolResponseParts, Part{FunctionResponse: response})
					}
				}

				if len(functionCallParts) > 0 {
					content.Parts = append(content.Parts, functionCallParts...)
				}

				if len(content.Parts) > 0 {
					contents = append(contents, content)
				}

				if len(toolResponseParts) > 0 {
					contents = append(contents, Content{
						Role:  "tool",
						Parts: toolResponseParts,
					})
				}
			}
		case "tool":
			if msg.ToolCallID != nil && msg.Content != nil {
				if s, ok := msg.Content.(string); ok {
					var responseData interface{}
					var jsonObj map[string]interface{}
					if err := json.Unmarshal([]byte(s), &jsonObj); err == nil {
						responseData = jsonObj
					} else {
						responseData = s
					}
					toolItems[*msg.ToolCallID] = &FunctionResponse{
						Response: map[string]interface{}{"result": responseData},
					}
				}
			}
		}

		if len(content.Parts) > 0 {
			contents = append(contents, content)
		}
	}

	projectID := ""
	if pID := ctx.Value(constant.ProjectIDKey); pID != nil {
		projectID = pID.(string)
	}

	c.Model = req.Model
	c.Project = projectID
	c.Request.SystemInstruction = systemInstruction
	c.Request.Contents = contents

	if req.ReasoningEffort != nil {
		switch *req.ReasoningEffort {
		case "none":
			c.Request.GenerationConfig.ThinkingConfig.IncludeThoughts = false
			c.Request.GenerationConfig.ThinkingConfig.ThinkingBudget = 0
		case "auto":
			c.Request.GenerationConfig.ThinkingConfig.ThinkingBudget = -1
			c.Request.GenerationConfig.ThinkingConfig.IncludeThoughts = true
		case "low":
			c.Request.GenerationConfig.ThinkingConfig.ThinkingBudget = 1024
			c.Request.GenerationConfig.ThinkingConfig.IncludeThoughts = true
		case "medium":
			c.Request.GenerationConfig.ThinkingConfig.ThinkingBudget = 8192
			c.Request.GenerationConfig.ThinkingConfig.IncludeThoughts = true
		case "high":
			c.Request.GenerationConfig.ThinkingConfig.ThinkingBudget = 24576
			c.Request.GenerationConfig.ThinkingConfig.IncludeThoughts = true
		default:
			c.Request.GenerationConfig.ThinkingConfig.ThinkingBudget = -1
			c.Request.GenerationConfig.ThinkingConfig.IncludeThoughts = true
		}
	}

	if req.Temperature != nil {
		c.Request.GenerationConfig.Temperature = *req.Temperature
	}
	if req.TopP != nil {
		c.Request.GenerationConfig.TopP = *req.TopP
	}

	if len(req.Tools) > 0 {
		declarations := []interface{}{}
		for _, tool := range req.Tools {
			if tool.Type == "function" {
				declarations = append(declarations, convertOpenAIFunctionToGemini(tool.Function))
			}
		}
		if len(declarations) > 0 {
			c.Request.Tools = []ToolDeclaration{
				{FunctionDeclarations: declarations},
			}
		}
	}

	return nil
}

type GeminiCliChatCompletionRequest struct {
	Model   string  `json:"model"`
	Project string  `json:"project"`
	Request Request `json:"request"`
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
type Request struct {
	SystemInstruction *Content          `json:"systemInstruction,omitempty"`
	Contents          []Content         `json:"contents"`
	Tools             []ToolDeclaration `json:"tools,omitempty"`
	GenerationConfig  *GenerationConfig `json:"generationConfig,omitempty"`
}

// GenerationConfig defines parameters that control the model's generation behavior.
type GenerationConfig struct {
	ThinkingConfig *GenerationConfigThinkingConfig `json:"thinkingConfig,omitempty"`
	Temperature    float64                         `json:"temperature,omitempty"`
	TopP           float64                         `json:"topP,omitempty"`
	TopK           float64                         `json:"topK,omitempty"`
}

// GenerationConfigThinkingConfig specifies configuration for the model's "thinking" process.
type GenerationConfigThinkingConfig struct {
	// IncludeThoughts determines whether the model should output its reasoning process.
	IncludeThoughts bool `json:"include_thoughts,omitempty"`
	ThinkingBudget  int  `json:"thinking_budget,omitempty"`
}

// ToolDeclaration defines the structure for declaring tools (like functions)
// that the model can call.
type ToolDeclaration struct {
	FunctionDeclarations []interface{} `json:"functionDeclarations"`
}

func convertOpenAIFunctionToGemini(openAIFunc openai.FunctionDefinition) map[string]interface{} {
	geminiFunc := map[string]interface{}{
		"name": openAIFunc.Name,
	}
	if openAIFunc.Description != nil {
		geminiFunc["description"] = *openAIFunc.Description
	}
	if openAIFunc.Parameters != nil {
		geminiFunc["parameters"] = convertSchema(openAIFunc.Parameters)
	}
	return geminiFunc
}

func convertSchema(schema interface{}) interface{} {
	if schema == nil {
		return nil
	}

	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return schema // Not a map, return as is
	}

	geminiSchema := make(map[string]interface{})
	for key, value := range schemaMap {
		switch key {
		case "type":
			if typeStr, ok := value.(string); ok {
				geminiSchema[key] = strings.ToUpper(typeStr)
			}
		case "description", "required":
			geminiSchema[key] = value
		case "properties":
			if propertiesMap, ok := value.(map[string]interface{}); ok {
				newProperties := make(map[string]interface{})
				for propKey, propValue := range propertiesMap {
					if strings.HasPrefix(propKey, "$") {
						continue // Skip $ref for now
					}
					newProperties[propKey] = convertSchema(propValue)
				}
				geminiSchema[key] = newProperties
			}
		case "items":
			geminiSchema[key] = convertSchema(value)
		default:
			// Ignore other fields like $schema, additionalProperties
		}
	}
	return geminiSchema
}
