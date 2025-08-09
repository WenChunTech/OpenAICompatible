package codegeex

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
	"github.com/WenChunTech/OpenAICompatible/src/util"
)

type Candidate struct {
	CandidateMsgID    string `json:"candidate_msg_id"`
	CandidateType     string `json:"candidate_type"`
	SelectedCandidate string `json:"selected_candidate"`
}

type History struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
	ID     string `json:"id"`
}

type CodeGeexChatRequest struct {
	UserID        string     `json:"user_id"`
	UserRole      int        `json:"user_role"`
	IDE           string     `json:"ide"`
	IDEVersion    string     `json:"ide_version"`
	PluginVersion string     `json:"plugin_version"`
	Prompt        string     `json:"prompt"`
	MachineID     string     `json:"machineId"`
	TalkID        string     `json:"talkId"`
	Locale        string     `json:"locale"`
	Model         string     `json:"model"`
	Agent         *string    `json:"agent"`
	Candidates    *Candidate `json:"candidates"`
	History       []History  `json:"history"`
}

func (c *CodeGeexChatRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error {
	var promptBuilder strings.Builder
	for _, message := range req.Messages {
		messageContent, err := json.Marshal(message)
		if err != nil {
			slog.Error("Failed to marshal message", "error", err)
			return err
		}
		promptBuilder.WriteString(string(messageContent))
		promptBuilder.WriteString("\n")
	}
	talkID := config.Config.CodeGeex.TalkID
	if talkID == "" {
		id, err := util.GenerateUUID()
		if err != nil {
			slog.Error("Failed to generate UUID", "error", err)
			return err
		}
		talkID = id
	}

	*c = CodeGeexChatRequest{
		UserID:        config.Config.CodeGeex.UserID,
		UserRole:      config.Config.CodeGeex.UserRole,
		IDE:           config.Config.CodeGeex.IDE,
		IDEVersion:    config.Config.CodeGeex.IDEVersion,
		PluginVersion: config.Config.CodeGeex.PluginVersion,
		Prompt:        promptBuilder.String(),
		MachineID:     config.Config.CodeGeex.MachineID,
		TalkID:        talkID,
		Locale:        config.Config.CodeGeex.Locale,
		Model:         req.Model,
	}
	return nil
}
