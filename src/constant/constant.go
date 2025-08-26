package constant

const (
	BufferSize = 1024
)

const (
	Accept              = "Accept"
	CacheControl        = "Cache-Control"
	CacheControlNoCache = "no-cache"
	UserAgent           = "User-Agent"
	Authorization       = "Authorization"
)

const (
	Connection          = "Connection"
	ConnectionKeepAlive = "keep-alive"
	CodeToken           = "Code-Token"
)

const (
	ContentType            = "Content-Type"
	ContentTypeText        = "text/plain"
	ContentTypeJson        = "application/json"
	ContentTypeForm        = "application/x-www-form-urlencoded"
	ContentTypeFormData    = "multipart/form-data"
	ContentTypeEventStream = "text/event-stream"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Code/1.102.0 Chrome/134.0.6998.205 Electron/35.6.0 Safari/537.36"

const (
	UserRole      = "user"
	AssistantRole = "assistant"
)

type ContextKey string

const (
	ChatIDKey      ContextKey = "chat_id"
	ProjectIDKey   ContextKey = "project_id"
	DangbeiChatID  ContextKey = "chatId"
	ConversationID ContextKey = "conversationId"
)

const (
	HTTPSScheme = "https"
	HTTPScheme  = "http"
)

const (
	CodeGeexPrefix           = "codegeex"
	CodeGeexHost             = "codegeex.cn"
	CodeGeexChatCompletePath = "/prod/code/chatCodeSseV3/chat"
	CodeGeexModelListPath    = "/prod/v3/chat/model"

	CodeGeexChatCompleteURL = "https://codegeex.cn/prod/code/chatCodeSseV3/chat"
	CodeGeexModelListURL    = "https://codegeex.cn/prod/v3/chat/model"
)

const (
	QwenPrefix           = "qwen"
	QwenHost             = "chat.qwen.ai"
	QwenModelListPath    = "/api/models"
	QwenChatIDPath       = "/api/v2/chats/new"
	QwenChatCompletePath = "/api/v2/chat/completions"

	QwenChatCompleteURL = "https://chat.qwen.ai/api/v2/chat/completions"
	QwenChatID          = "https://chat.qwen.ai/api/v2/chats/new"
	QwenModelListURL    = "https://chat.qwen.ai/api/models"
)

const (
	QwenCodePrefix = "qwen_code"
)

const (
	GeminiCliPrefix           = "gemini_cli"
	GeminiCliHost             = "cloudcode-pa.googleapis.com"
	GeminiCliChatCompletePath = "/v1internal:streamGenerateContent"
	GeminiCliModelListPath    = ""

	GeminiCliChatCompleteURL = "https://cloudcode-pa.googleapis.com/v1internal:streamGenerateContent?alt=sse"
	GeminiCliModelListURL    = ""
)

const (
	DangBeiPrefix           = "dangbei"
	DangBeiHost             = "ai-api.dangbei.net"
	DangBeiChatCompletePath = "/ai-search/chatApi/v2/chat"
	DangBeiModelListPath    = "/ai-search/configApi/v1/getChatModelConfig"

	DangBeiChatCreateURL   = "https://ai-api.dangbei.net/ai-search/conversationApi/v1/batch/create"
	DangBeiGenerateIdURL   = "https://ai-api.dangbei.net/ai-search/commonApi/v1/generateId"
	DangBeiChatCompleteURL = "https://ai-api.dangbei.net/ai-search/chatApi/v2/chat"
	DangBeiModelListURL    = "https://ai-api.dangbei.net/ai-search/configApi/v1/getChatModelConfig"
)
