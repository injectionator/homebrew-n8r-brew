package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/injectionator/n8r/internal/config"
)

// StoredToken is the on-disk representation of saved credentials.
type StoredToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	SavedAt     time.Time `json:"saved_at"`
}

// SaveToken persists a TokenResponse to disk.
func SaveToken(token TokenResponse) error {
	path, err := config.CredentialsPath()
	if err != nil {
		return err
	}

	stored := StoredToken{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		SavedAt:     time.Now(),
	}

	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// LoadToken reads the stored token from disk.
func LoadToken() (*StoredToken, error) {
	path, err := config.CredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var stored StoredToken
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &stored, nil
}

// DeleteToken removes the credentials file.
func DeleteToken() error {
	path, err := config.CredentialsPath()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete credentials: %w", err)
	}

	return nil
}

// IsExpired returns true if the token has expired.
func (t *StoredToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
