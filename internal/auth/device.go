package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/injectionator/n8r/internal/config"
)

// DeviceCodeResponse is returned by the device code endpoint.
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// TokenResponse is returned when the device is authorized.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ErrorResponse represents an error from the auth API.
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

var (
	ErrAuthorizationPending = errors.New("authorization_pending")
	ErrSlowDown             = errors.New("slow_down")
	ErrExpiredToken         = errors.New("expired_token")
	ErrAccessDenied         = errors.New("access_denied")
)

// RequestDeviceCode initiates the device authorization flow.
func RequestDeviceCode() (*DeviceCodeResponse, error) {
	url := config.BaseURL + "/api/auth/device/code"
	req, err := http.NewRequest("POST", url, strings.NewReader("{}"))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", config.BaseURL)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device code request failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var dcr DeviceCodeResponse
	if err := json.Unmarshal(body, &dcr); err != nil {
		return nil, fmt.Errorf("failed to parse device code response: %w", err)
	}

	return &dcr, nil
}

// PollForToken polls the token endpoint until authorization is granted or the code expires.
func PollForToken(deviceCode string, interval int, expiresIn int) (*TokenResponse, error) {
	url := config.BaseURL + "/api/auth/device/token"
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)
	pollInterval := time.Duration(interval) * time.Second

	for {
		if time.Now().After(deadline) {
			return nil, ErrExpiredToken
		}

		time.Sleep(pollInterval)

		payload := fmt.Sprintf(`{"device_code":"%s"}`, deviceCode)
		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", config.BaseURL)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to poll for token: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			var tr TokenResponse
			if err := json.Unmarshal(body, &tr); err != nil {
				return nil, fmt.Errorf("failed to parse token response: %w", err)
			}
			return &tr, nil
		}

		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			switch errResp.Error {
			case "authorization_pending":
				continue
			case "slow_down":
				pollInterval += 5 * time.Second
				continue
			case "expired_token":
				return nil, ErrExpiredToken
			case "access_denied":
				return nil, ErrAccessDenied
			}
		}

		return nil, fmt.Errorf("unexpected response (HTTP %d): %s", resp.StatusCode, string(body))
	}
}
