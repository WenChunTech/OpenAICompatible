package dangbei

import (
	"context"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

//	{
//	    "role": "assistant",
//	    "type": "answer",
//	    "content": "帮",
//	    "content_type": "text",
//	    "id": "362912660282413254",
//	    "parentMsgId": "362912660282413253",
//	    "conversation_id": "362544896454295941",
//	    "created_at": 1755186644,
//	    "requestId": "a19f3010-8b6e-413e-9a10-0b2ce88ec033",
//	    "supportDownload": false
//	}
type DangBeiChatCompleteResponse struct {
	Role            string `json:"role"`
	Type            string `json:"type"`
	Content         string `json:"content"`
	ContentType     string `json:"content_type"`
	ID              string `json:"id"`
	ParentMsgID     string `json:"parentMsgId"`
	ConversationID  string `json:"conversation_id"`
	CreatedAt       int    `json:"created_at"`
	RequestID       string `json:"requestId"`
	SupportDownload bool   `json:"supportDownload"`
}

func (c DangBeiChatCompleteResponse) Convert(ctx context.Context) (*openai.OpenAPIChatCompletionStreamResponse, error) {

	return &openai.OpenAPIChatCompletionStreamResponse{
		ID:      c.ID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   "",
		Choices: []openai.OpenAIStreamChoice{{
			Index: 0,
			Delta: &openai.Delta{
				Role:    c.Role,
				Content: c.Content,
			},
		}},
		Usage: nil,
	}, nil
}
