package config

import (
	"os"
	"path/filepath"
)

// Version can be overridden at build time via ldflags.
var Version = "0.1.0"

const (
	BaseURL         = "https://injectionator.com"
	ConfigDirName   = ".n8r"
	CredentialsFile = "credentials.json"
)

// ConfigDir returns the absolute path to ~/.n8r, creating it if needed.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ConfigDirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

// CredentialsPath returns the absolute path to the credentials file.
func CredentialsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CredentialsFile), nil
}
