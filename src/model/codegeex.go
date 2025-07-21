package model

import (
	"time"
)

type CodeGeexChatRequest struct {
	UserID        string  `json:"user_id"`
	UserRole      int     `json:"user_role"`
	IDE           string  `json:"ide"`
	IDEVersion    string  `json:"ide_version"`
	PluginVersion string  `json:"plugin_version"`
	Prompt        string  `json:"prompt"`
	MachineID     string  `json:"machineId"`
	TalkID        string  `json:"talkId"`
	Locale        string  `json:"locale"`
	Model         string  `json:"model"`
	Agent         *string `json:"agent"`
	Candidates    struct {
		CandidateMsgID    string `json:"candidate_msg_id"`
		CandidateType     string `json:"candidate_type"`
		SelectedCandidate string `json:"selected_candidate"`
	} `json:"candidates"`
	History []struct {
		Query  string `json:"query"`
		Answer string `json:"answer"`
		ID     string `json:"id"`
	}
}

type CodeGeexSSEData struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Model        string `json:"model"`
}

func (c CodeGeexSSEData) Convert() (*OpenAPIChatCompletionStreamResponse, error) {
	choice := OpenAIStreamChoice{
		Index: 0,
		Delta: Delta{
			Content:          c.Text,
			Role:             "assistant",
			ReasoningContent: nil,
			ToolCall:         nil,
		},
		FinishReason: nil,
		Logprobs:     nil,
		MatchedSotp:  nil,
	}
	if c.FinishReason != "" {
		finishReason := c.FinishReason
		choice.FinishReason = &finishReason
	}
	return &OpenAPIChatCompletionStreamResponse{
		ID:      c.ID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   c.Model,
		Choices: []OpenAIStreamChoice{choice},
		Usage:   nil,
	}, nil
}

type CodeGeexModelOptions struct {
	Options []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description struct {
			Zh string `json:"zh"`
			En string `json:"en"`
		} `json:"description"`
		Host string `json:"host"`
	} `json:"options"`
	Default string `json:"default"`
	Promote struct {
		Key  string `json:"key"`
		Text struct {
			Zh string `json:"zh"`
			En string `json:"en"`
		} `json:"text"`
	} `json:"promote"`
	IP string `json:"ip"`
}

func (c CodeGeexModelOptions) Convert() (*ModelList, error) {
	models := make([]Model, len(c.Options))
	for i, option := range c.Options {
		models[i] = Model{
			ID:      option.ID,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "codegeex",
		}
	}
	return &ModelList{
		Object: "list",
		Data:   models,
	}, nil
}
