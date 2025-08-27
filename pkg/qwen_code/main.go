package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/request"
)

const QwenOAuthFile = "qwen_code.json"

// OAuth Endpoints
const (
	qwenOAuthBaseURL            = "https://chat.qwen.ai"
	qwenOAuthDeviceCodeEndpoint = qwenOAuthBaseURL + "/api/v1/oauth2/device/code"
	qwenOAuthTokenEndpoint      = qwenOAuthBaseURL + "/api/v1/oauth2/token"
)

// OAuth Client Configuration
const (
	qwenOAuthClientID   = "f0304373b74a44d2b584a3fb70ca9e56"
	qwenOAuthScope      = "openid profile email model.completion"
	qwenOAuthGrantType  = "urn:ietf:params:oauth:grant-type:device_code"
	codeChallengeMethod = "S256"
)

// ErrorData represents the standard error response.
type ErrorData struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *ErrorData) Error() string {
	return fmt.Sprintf("OAuth error: %s - %s", e.ErrorCode, e.ErrorDescription)
}

// DeviceAuthorizationResponse corresponds to the device authorization success data.
type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
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
	Expiry time.Time `json:"expiry"`
}

// generateCodeVerifier creates a random string for PKCE.
func generateCodeVerifier() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		slog.Error("Failed to generate code verifier", "error", err)
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// isTokenValid checks if the credentials are still valid, considering a small buffer.
func isTokenValid(creds *QwenCredentials) bool {
	if creds == nil {
		return false
	}
	// Check if the token is expiring within the next 30 seconds.
	return time.Now().Before(creds.Expiry.Add(-30 * time.Second))
}

// loadCredentials reads and parses credentials from the QwenOAuthFile.
func loadCredentials() (*QwenCredentials, error) {
	data, err := os.ReadFile(QwenOAuthFile)
	if err != nil {
		return nil, err // The caller can check for os.IsNotExist(err)
	}

	var creds QwenCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to decode credentials from %s: %w", QwenOAuthFile, err)
	}
	slog.Info("Successfully loaded credentials from", "file", QwenOAuthFile)
	return &creds, nil
}

// saveCredentials serializes and writes credentials to the QwenOAuthFile.
func saveCredentials(creds *QwenCredentials) error {
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode credentials: %w", err)
	}

	if err := os.WriteFile(QwenOAuthFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file %s: %w", QwenOAuthFile, err)
	}
	slog.Info("Successfully saved credentials to", "file", QwenOAuthFile)
	return nil
}

// generateCodeChallenge creates a SHA-256 hash of the code verifier.
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// requestDeviceAuthorization initiates the device authorization flow.
func requestDeviceAuthorization(ctx context.Context) (*DeviceAuthorizationResponse, string, error) {
	verifier := generateCodeVerifier()
	challenge := generateCodeChallenge(verifier)

	formData := map[string][]string{
		"client_id":             {qwenOAuthClientID},
		"scope":                 {qwenOAuthScope},
		"code_challenge":        {challenge},
		"code_challenge_method": {codeChallengeMethod},
	}

	// Generate a random request ID
	reqIDBytes := make([]byte, 16)
	if _, err := rand.Read(reqIDBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate request ID: %w", err)
	}
	reqID := fmt.Sprintf("%x", reqIDBytes)

	headers := map[string]string{
		"x-request-id": reqID,
		"Accept":       "application/json",
	}

	resp, err := request.NewRequestBuilder(qwenOAuthDeviceCodeEndpoint, http.MethodPost).
		WithHeaders(headers).
		WithForm(formData).
		Do(ctx, nil)

	if err != nil {
		return nil, "", fmt.Errorf("device authorization request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData ErrorData
		if err := json.NewDecoder(resp.Body).Decode(&errData); err == nil {
			return nil, "", fmt.Errorf("device authorization failed: %w", &errData)
		}
		return nil, "", fmt.Errorf("device authorization failed with status code: %d", resp.StatusCode)
	}

	var authResp DeviceAuthorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode device authorization response: %w", err)
	}

	slog.Info("Device authorization successful", "response", authResp)
	return &authResp, verifier, nil
}

// pollDeviceToken polls the token endpoint to get the access token.
func pollDeviceToken(ctx context.Context, deviceCode, verifier string, interval, timeout time.Duration) (*DeviceTokenResponse, error) {
	pollingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-pollingCtx.Done():
			return nil, fmt.Errorf("polling timed out after %v", timeout)
		case <-ticker.C:
			slog.Info("Polling for device token...")
			tokenResp, err := attemptToGetToken(pollingCtx, deviceCode, verifier)
			if err != nil {
				// Check for specific polling errors
				if errData, ok := err.(*ErrorData); ok {
					switch errData.ErrorCode {
					case "authorization_pending":
						// Continue polling
						continue
					case "slow_down":
						// Increase interval and continue polling
						ticker.Reset(interval * 2) // Simple backoff
						continue
					}
				}
				// For other errors, stop polling
				return nil, err
			}
			// Success
			return tokenResp, nil
		}
	}
}

func attemptToGetToken(ctx context.Context, deviceCode, verifier string) (*DeviceTokenResponse, error) {
	formData := map[string][]string{
		"grant_type":    {qwenOAuthGrantType},
		"client_id":     {qwenOAuthClientID},
		"device_code":   {deviceCode},
		"code_verifier": {verifier},
	}

	resp, err := request.NewRequestBuilder(qwenOAuthTokenEndpoint, http.MethodPost).
		WithForm(formData).
		Do(ctx, nil)
	if err != nil {
		slog.Error("token poll request failed", "error", err)
		return nil, fmt.Errorf("token poll request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // Read body to reuse it
	if err != nil {
		slog.Error("failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errData ErrorData
		if err := json.Unmarshal(body, &errData); err == nil && errData.ErrorCode != "" {
			// Return structured error for the polling loop to handle
			return nil, &errData
		}
		return nil, fmt.Errorf("token poll failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp DeviceTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
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
		Expiry: time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second),
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

// GetToken orchestrates the entire device authorization flow.
func GetToken(ctx context.Context) (*DeviceTokenResponse, error) {
	authResp, verifier, err := requestDeviceAuthorization(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to request device authorization: %w", err)
	}

	fmt.Printf("Please go to %s and enter the code: %s\n", authResp.VerificationURI, authResp.UserCode)

	// Poll for 5 seconds interval, with a 5 minute timeout.
	tokenResp, err := pollDeviceToken(ctx, authResp.DeviceCode, verifier, 5*time.Second, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for device token: %w", err)
	}

	return tokenResp, nil
}

func main() {
	ctx := context.Background()
	var creds *QwenCredentials
	var err error

	// 1. Try to load existing credentials
	creds, err = loadCredentials()
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("No existing credentials found. Starting new device authorization flow.")
			// 4. File doesn't exist, get a new token
			tokenResp, getTokenErr := GetToken(ctx)
			if getTokenErr != nil {
				slog.Error("Failed to get new token", "error", getTokenErr)
				return
			}
			// 5. Convert to QwenCredentials
			creds = &QwenCredentials{
				DeviceTokenResponse: *tokenResp,
				Expiry:              time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
			}
		} else {
			// Other error loading credentials
			slog.Error("Failed to load credentials", "error", err)
			return
		}
	} else {
		// 2. Credentials loaded, check if token is valid
		if !isTokenValid(creds) {
			slog.Info("Access token is expired or invalid. Refreshing...")
			// 3. Token is invalid, refresh it
			refreshedCreds, refreshErr := refreshAccessToken(ctx, creds)
			if refreshErr != nil {
				slog.Error("Failed to refresh access token", "error", refreshErr)
				// If refresh fails, might need to re-authenticate
				slog.Info("Attempting to re-authenticate with device flow.")
				tokenResp, getTokenErr := GetToken(ctx)
				if getTokenErr != nil {
					slog.Error("Failed to get new token after refresh failure", "error", getTokenErr)
					return
				}
				creds = &QwenCredentials{
					DeviceTokenResponse: *tokenResp,
					Expiry:              time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
				}
			} else {
				creds = refreshedCreds
			}
		} else {
			slog.Info("Existing access token is valid.")
		}
	}

	// 6. Save the latest valid credentials
	if creds != nil {
		if err := saveCredentials(creds); err != nil {
			slog.Error("Failed to save credentials", "error", err)
			return
		}
		slog.Info("Successfully obtained and saved token.", "access_token", creds.AccessToken)
	}
}
