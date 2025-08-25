package qwencode

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/request"
)

// OAuth Endpoints
const (
	qwenOAuthTokenEndpoint = "https://chat.qwen.ai/api/v1/oauth2/token"
	qwenOAuthClientID      = "f0304373b74a44d2b584a3fb70ca9e56"
)

// ErrorData represents the standard error response.
type ErrorData struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *ErrorData) Error() string {
	return fmt.Sprintf("OAuth error: %s - %s", e.ErrorCode, e.ErrorDescription)
}

// DeviceTokenResponse corresponds to the device token success data.
type DeviceTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ResourceURL  string `json:"resource_url"`
}

// QwenCredentials extends DeviceTokenResponse with an exact expiry time.
type QwenCredentials struct {
	DeviceTokenResponse
	ExpiryDate time.Time `json:"expiry_date"`
}

// refreshAccessToken uses the refresh token to get a new access token.
func refreshAccessToken(ctx context.Context, creds *QwenCredentials) (*QwenCredentials, error) {
	if creds == nil || creds.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	slog.Info("Refreshing access token...")

	formData := map[string][]string{
		"grant_type":    {"refresh_token"},
		"refresh_token": {creds.RefreshToken},
		"client_id":     {qwenOAuthClientID},
	}

	resp, err := request.NewRequestBuilder(qwenOAuthTokenEndpoint, http.MethodPost).
		WithForm(formData).
		Do(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData ErrorData
		if err := json.NewDecoder(resp.Body).Decode(&errData); err == nil {
			return nil, &errData
		}
		return nil, fmt.Errorf("token refresh failed with status code: %d", resp.StatusCode)
	}

	var tokenData DeviceTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return nil, fmt.Errorf("failed to decode token refresh response: %w", err)
	}

	// Create new credentials, preserving the refresh token if a new one isn't provided.
	newCreds := &QwenCredentials{
		DeviceTokenResponse: DeviceTokenResponse{
			AccessToken:  tokenData.AccessToken,
			RefreshToken: tokenData.RefreshToken,
			TokenType:    tokenData.TokenType,
			ExpiresIn:    tokenData.ExpiresIn,
			ResourceURL:  tokenData.ResourceURL,
		},
		ExpiryDate: time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second),
	}

	// If the new response doesn't include a refresh token, reuse the old one.
	if newCreds.RefreshToken == "" {
		newCreds.RefreshToken = creds.RefreshToken
	}
	if newCreds.ResourceURL == "" {
		newCreds.ResourceURL = creds.ResourceURL
	}

	slog.Info("Successfully refreshed access token.")
	return newCreds, nil
}
