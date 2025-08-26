package qwencode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/config"
	"github.com/WenChunTech/OpenAICompatible/src/request"
)

const (
	qwenOAuthTokenEndpoint = "https://chat.qwen.ai/api/v1/oauth2/token"
	qwenOAuthClientID      = "f0304373b74a44d2b584a3fb70ca9e56"
)

func isTokenValid(creds *config.QwenCodeToken) bool {
	if creds == nil {
		return false
	}
	return time.Now().Before(creds.Expiry.Add(-30 * time.Second))
}

func refreshAccessToken(ctx context.Context, token *config.QwenCodeToken) (*config.QwenCodeToken, error) {
	if token == nil || token.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	formData := map[string][]string{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.RefreshToken},
		"client_id":     {qwenOAuthClientID},
	}

	resp, err := request.NewRequestBuilder(qwenOAuthTokenEndpoint, http.MethodPost).
		WithForm(formData).
		Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to refresh access token", "error", err)
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err == nil {
			slog.Error("Failed to read response body for token refresh failed", "error", err)
			return nil, fmt.Errorf("token refresh failed with status code: %d, error message: %s", resp.StatusCode, err)
		}
		slog.Error("Failed to refresh token", "error", string(errMsg))
		return nil, fmt.Errorf("token refresh failed with status code: %d, error message: %s", resp.StatusCode, string(errMsg))
	}

	var newToken config.QwenCodeToken
	if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
		slog.Error("Failed to decode token refresh response", "error", err)
		return nil, fmt.Errorf("failed to decode token refresh response: %w", err)
	}

	newToken.Expiry = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = token.RefreshToken
	}
	if newToken.ResourceURL == "" {
		newToken.ResourceURL = token.ResourceURL
	}
	token = &newToken
	return token, nil
}

func GetQwenCodeToken(ctx context.Context, token *config.QwenCodeToken) (*config.QwenCodeToken, error) {
	if isTokenValid(token) {
		return token, nil
	}

	return refreshAccessToken(ctx, token)
}
