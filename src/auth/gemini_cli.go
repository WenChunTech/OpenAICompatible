package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	Port              = 8090
	OauthJSONPath     = "oauth_creds.json"
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

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func newTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	codeChan := make(chan string)
	errChan := make(chan error)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "path", r.URL.Path)
	})

	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		if err := r.URL.Query().Get("error"); err != "" {
			slog.Error("Authentication failed", "error", err)
			errChan <- fmt.Errorf("authentication failed via callback: %s", err)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			slog.Error("Authentication failed", "error", "code not found in callback")
			errChan <- fmt.Errorf("code not found in callback")
			return
		}
		slog.Info("Authentication successful")
		codeChan <- code
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			errChan <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))

	var err error
	err = openBrowser(authURL)
	if err != nil {
		fmt.Printf("无法自动打开浏览器，请手动访问此 URL：%s\n", authURL)
	}

	var authCode string
	select {
	case code := <-codeChan:
		authCode = code
	case err = <-errChan:
		return nil, err
	case <-time.After(1 * time.Minute):
		slog.Error("Authentication timeout")
		return nil, fmt.Errorf("authentication timeout")
	}

	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	return token, nil
}

func saveGeminiCLiToken(token *oauth2.Token) error {
	file, err := os.Create(OauthJSONPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(token); err != nil {
		return fmt.Errorf("failed to encode token: %w", err)
	}

	return nil
}

func newConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     OauthClientID,
		ClientSecret: OauthClientSecret,
		RedirectURL:  OauthRedirectURL,
		Scopes:       OauthScopes,
		Endpoint:     google.Endpoint,
	}
}

func StartGeminiCli() {
	slog.Info("Starting Auth Gemini CLI Token...")
	ctx := context.Background()
	config := newConfig()
	token, err := newTokenFromWeb(ctx, config)
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return
	}
	if err := saveGeminiCLiToken(token); err != nil {
		slog.Error("Failed to save token", "error", err)
		return
	}
	slog.Info("Token saved successfully")
}
