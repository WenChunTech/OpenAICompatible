package geminicli

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	OauthClientID     = "681255809395-oo8ft2oprdrnp9e3aqf6av3hmdib135j.apps.googleusercontent.com"
	OauthClientSecret = "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"
	OauthRedirectURL  = "http://localhost:8085/oauth2callback"
)

var (
	OauthScopes = []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
)

func newConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     OauthClientID,
		ClientSecret: OauthClientSecret,
		RedirectURL:  OauthRedirectURL,
		Scopes:       OauthScopes,
		Endpoint:     google.Endpoint,
	}
}

// TokenWrapper 包装TokenSource以便获取最新token
type TokenWrapper struct {
	source oauth2.TokenSource
	client *http.Client
}

// NewTokenWrapper 创建TokenWrapper
func NewTokenWrapper(ctx context.Context, token *oauth2.Token) *TokenWrapper {
	config := newConfig()
	// 使用ReuseTokenSource来缓存token，避免频繁刷新
	source := oauth2.ReuseTokenSource(token, config.TokenSource(ctx, token))
	client := oauth2.NewClient(ctx, source)
	return &TokenWrapper{
		source: source,
		client: client,
	}
}

// GetClient 返回HTTP客户端
func (tw *TokenWrapper) GetClient() *http.Client {
	return tw.client
}

// GetAccessToken 获取当前访问令牌
func (tw *TokenWrapper) GetAccessToken() (string, error) {
	token, err := tw.source.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// GetToken 获取完整的token信息
func (tw *TokenWrapper) GetToken() (*oauth2.Token, error) {
	return tw.source.Token()
}
