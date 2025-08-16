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
type DangBeiChatCompleteRequest struct {
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

func (c *DangBeiChatCompleteRequest) ImportOpenAIChatCompletionRequest(ctx context.Context, req *openai.OpenAIChatCompletionRequest) error {
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

	*c = DangBeiChatCompleteRequest{
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

type DangBeiChatCreateResponse struct {
	Success    bool           `json:"success"`
	ErrCode    any            `json:"errCode"`
	ErrMessage any            `json:"errMessage"`
	RequestID  string         `json:"requestId"`
	Data       ChatCreateData `json:"data"`
}

type ChatModelConfig struct {
	Model   string `json:"model"`
	Options []any  `json:"options"`
}

type ActiveConversation struct {
	ConversationID string `json:"conversationId"`
	IsVisible      bool   `json:"isVisible"`
	Order          int    `json:"order"`
}

type MultiModelLayout struct {
	ActiveConversation []ActiveConversation `json:"activeConversation"`
	LastUpdateTime     string               `json:"lastUpdateTime"`
}

type MetaData struct {
	WriteCode        string           `json:"writeCode"`
	ChatModelConfig  ChatModelConfig  `json:"chatModelConfig"`
	MultiModelLayout MultiModelLayout `json:"multiModelLayout"`
	SuperAgentPath   string           `json:"superAgentPath"`
	PageType         any              `json:"pageType"`
}

type ConversationList struct {
	ConversationID   string   `json:"conversationId"`
	ConversationType int      `json:"conversationType"`
	Title            string   `json:"title"`
	UserID           any      `json:"userId"`
	DeviceID         string   `json:"deviceId"`
	TitleSummaryFlag int      `json:"titleSummaryFlag"`
	MetaData         MetaData `json:"metaData"`
	IsAnonymous      bool     `json:"isAnonymous"`
	AnonymousKey     string   `json:"anonymousKey"`
	LastChatModel    any      `json:"lastChatModel"`
	ConversationList any      `json:"conversationList"`
}

type ChatCreateData struct {
	ConversationID   string             `json:"conversationId"`
	ConversationType int                `json:"conversationType"`
	Title            string             `json:"title"`
	UserID           any                `json:"userId"`
	DeviceID         string             `json:"deviceId"`
	TitleSummaryFlag int                `json:"titleSummaryFlag"`
	MetaData         MetaData           `json:"metaData"`
	IsAnonymous      bool               `json:"isAnonymous"`
	AnonymousKey     string             `json:"anonymousKey"`
	LastChatModel    any                `json:"lastChatModel"`
	ConversationList []ConversationList `json:"conversationList"`
}

type DangBeiGenerateIdResponse struct {
	Success    bool   `json:"success"`
	ErrCode    any    `json:"errCode"`
	ErrMessage any    `json:"errMessage"`
	RequestID  string `json:"requestId"`
	Data       string `json:"data"`
}
