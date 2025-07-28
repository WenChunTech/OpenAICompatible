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
	CodeGeexChatURL  = "https://codegeex.cn/prod/code/chatCodeSseV3/chat"
	CodeGeexModelURL = "https://codegeex.cn/prod/v3/chat/model"
	CodeGeexPrefix   = "codegeex"
)

const (
	QwenChatURL  = "https://chat.qwen.ai/api/v2/chat/completions"
	QwenChatID   = "https://chat.qwen.ai/api/v2/chats/new"
	QwenModelURL = "https://chat.qwen.ai/api/models"
	QwenPrefix   = "qwen"
)

const (
	UserRole      = "user"
	AssistantRole = "assistant"
)

type ContextKey string

const (
	ChatIDKey ContextKey = "chat_id"
)
