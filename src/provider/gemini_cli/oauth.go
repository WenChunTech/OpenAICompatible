package gemini

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Token 表示 OAuth2 令牌
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

// Config 表示 OAuth2 配置
type Config struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
	Scopes       []string
}

// launchBrowser 在系统默认浏览器中打开 URL
func launchBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // Linux
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

// 生成随机状态字符串
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetTokenFromWeb 从 Web 浏览器获取新的令牌
func GetTokenFromWeb(ctx context.Context, config *Config) (*Token, error) {
	// 创建通道用于接收授权码
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// 设置重定向 URL
	if config.RedirectURL == "" {
		config.RedirectURL = "http://localhost:8085/oauth2callback"
	}

	// 创建新的 HTTP 服务器
	server := &http.Server{Addr: ":8085"}

	// 生成状态令牌
	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("生成状态令牌失败: %w", err)
	}

	// 设置回调处理
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		queryState := r.URL.Query().Get("state")
		if queryState != state {
			err := errors.New("状态不匹配，可能存在安全风险")
			log.Printf("认证失败: %v", err)
			errChan <- err
			http.Error(w, "认证失败: 状态不匹配", http.StatusBadRequest)
			return
		}

		if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
			err := fmt.Errorf("认证失败: %s", errorMsg)
			log.Printf("认证失败: %v", err)
			errChan <- err
			http.Error(w, fmt.Sprintf("认证失败: %s", errorMsg), http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			err := errors.New("未收到授权码")
			log.Printf("认证失败: %v", err)
			errChan <- err
			http.Error(w, "认证失败: 未收到授权码", http.StatusBadRequest)
			return
		}

		// 返回成功页面
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `&lt;!DOCTYPE html&gt;
&lt;html&gt;
&lt;head&gt;
    &lt;title&gt;认证成功&lt;/title&gt;
    &lt;style&gt;
        body { font-family: Arial, sans-serif; text-align: center; padding-top: 50px; }
        h1 { color: #4CAF50; }
    &lt;/style&gt;
&lt;/head&gt;
&lt;body&gt;
    &lt;h1&gt;认证成功！&lt;/h1&gt;
    &lt;p&gt;您可以关闭此窗口了。&lt;/p&gt;
&lt;/body&gt;
&lt;/html&gt;`)

		codeChan <- code
	})

	// 启动服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务器启动失败: %v", err)
			errChan <- fmt.Errorf("服务器启动失败: %w", err)
		}
	}()

	// 构建授权 URL
	authURL := fmt.Sprintf("%s?%s",
		config.AuthURL,
		url.Values{
			"client_id":     {config.ClientID},
			"redirect_uri":  {config.RedirectURL},
			"response_type": {"code"},
			"scope":         {strings.Join(config.Scopes, " ")},
			"state":         {state},
			"access_type":   {"offline"},
			"prompt":        {"consent"},
		}.Encode(),
	)

	// 打开浏览器
	if err := launchBrowser(authURL); err != nil {
		log.Printf("无法自动打开浏览器: %v", err)
		fmt.Printf("请手动访问此 URL 进行授权：%s\n", authURL)
	}

	// 等待授权码或错误
	var authCode string
	select {
	case code := <-codeChan:
		authCode = code
		// 关闭服务器
		go server.Shutdown(context.Background())
	case err := <-errChan:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, errors.New("认证超时")
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// 交换授权码获取令牌
	token, err := exchangeToken(ctx, config, authCode)
	if err != nil {
		return nil, fmt.Errorf("交换令牌失败: %w", err)
	}

	return token, nil
}

// exchangeToken 使用授权码交换访问令牌
func exchangeToken(ctx context.Context, config *Config, authCode string) (*Token, error) {
	data := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"code":          {authCode},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {config.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取令牌失败: %s", body)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}
