package dangbei

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

//	{
//	    "stream": true,
//	    "botCode": "AI_SEARCH",
//	    "conversationId": "361947757325324485",
//	    "question": "你是什么大模型？",
//	    "model": "deepseek",
//	    "chatOption": {
//	        "searchKnowledge": false,
//	        "searchAllKnowledge": false,
//	        "searchSharedKnowledge": false
//	    },
//	    "knowledgeList": [],
//	    "anonymousKey": "",
//	    "uuid": "361947757894304133",
//	    "chatId": "361947757894304133",
//	    "files": [],
//	    "reference": [],
//	    "role": "user",
//	    "status": "local",
//	    "content": "你是什么大模型？",
//	    "userAction": "deep,online",
//	    "agentId": ""
//	}
type DangBeiChatRequest struct {
	Stream         bool       `json:"stream"`
	BotCode        string     `json:"botCode"`
	ConversationID string     `json:"conversationId"`
	Question       string     `json:"question"`
	Model          string     `json:"model"`
	ChatOption     ChatOption `json:"chatOption"`
	KnowledgeList  []any      `json:"knowledgeList"`
	AnonymousKey   string     `json:"anonymousKey"`
	UUID           string     `json:"uuid"`
	ChatID         string     `json:"chatId"`
	Files          []any      `json:"files"`
	Reference      []any      `json:"reference"`
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	Content        string     `json:"content"`
	UserAction     string     `json:"userAction"`
	AgentID        string     `json:"agentId"`
}
type ChatOption struct {
	SearchKnowledge       bool `json:"searchKnowledge"`
	SearchAllKnowledge    bool `json:"searchAllKnowledge"`
	SearchSharedKnowledge bool `json:"searchSharedKnowledge"`
}

func (c *DangBeiChatRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error {
	var question strings.Builder
	for _, message := range req.Messages {
		switch message.Content.(type) {
		case string:
			question.WriteString(message.Content.(string))
		default:
			jsonStr, _ := json.Marshal(message.Content)
			question.WriteString(string(jsonStr))
		}
	}

	role := "user"
	if len(req.Messages) > 0 {
		role = req.Messages[0].Role
	}

	var chatID string
	if ctx.Value(constant.DangbeiChatID) != nil {
		chatID = ctx.Value(constant.DangbeiChatID).(string)
	}

	var conversationID string
	if ctx.Value(constant.ConversationID) != nil {
		conversationID = ctx.Value(constant.ConversationID).(string)
	}

	*c = DangBeiChatRequest{
		Stream:         *req.Stream,
		BotCode:        "AI_SEARCH",
		ConversationID: conversationID,
		Question:       question.String(),
		Role:           role,
		ChatID:         chatID,
		UUID:           chatID,
		Model:          req.Model,
		UserAction:     "deep,online",
		ChatOption: ChatOption{
			SearchKnowledge:       false,
			SearchAllKnowledge:    false,
			SearchSharedKnowledge: false,
		},
	}

	return nil
}
