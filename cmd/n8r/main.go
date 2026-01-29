package main

import (
	"fmt"
	"os"
	"time"

	"github.com/injectionator/n8r/internal/auth"
	"github.com/injectionator/n8r/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cmd := os.Args[1]

	switch cmd {
	case "login":
		cmdLogin()
	case "logout":
		cmdLogout()
	case "status":
		cmdStatus()
	case "version", "--version", "-v":
		fmt.Printf("n8r v%s\n", config.Version)
	case "help", "--help", "-h":
		printUsage()
	default:
		// Any other command requires authentication
		token, err := auth.LoadToken()
		if err != nil || token == nil || token.IsExpired() {
			fmt.Fprintln(os.Stderr, "Please run `n8r login` first")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`n8r v%s — Injectionator CLI

Usage:
  n8r <command>

Commands:
  login       Authenticate with Injectionator
  logout      Remove stored credentials
  status      Show authentication status
  version     Print version

Flags:
  --version   Print version
  --help      Show this help
`, config.Version)
}

func cmdLogin() {
	// Check if already logged in
	existing, _ := auth.LoadToken()
	if existing != nil && !existing.IsExpired() {
		fmt.Println("You are already authenticated.")
		fmt.Println("Run `n8r logout` first to re-authenticate.")
		return
	}

	// Request device code
	dcr, err := auth.RequestDeviceCode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Visit %s and enter code: %s\n", dcr.VerificationURI, dcr.UserCode)
	fmt.Println("Waiting for authorization...")

	// Poll for token
	token, err := auth.PollForToken(dcr.DeviceCode, dcr.Interval, dcr.ExpiresIn)
	if err != nil {
		if err == auth.ErrExpiredToken {
			fmt.Fprintln(os.Stderr, "Error: Device code expired. Please run `n8r login` again.")
		} else if err == auth.ErrAccessDenied {
			fmt.Fprintln(os.Stderr, "Error: Authorization was denied.")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	// Save token
	if err := auth.SaveToken(*token); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving credentials: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully authenticated!")
}

func cmdLogout() {
	if err := auth.DeleteToken(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Logged out. Credentials removed.")
}

func cmdStatus() {
	token, err := auth.LoadToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading credentials: %v\n", err)
		os.Exit(1)
	}

	if token == nil {
		fmt.Println("Not authenticated. Run `n8r login` to get started.")
		return
	}

	if token.IsExpired() {
		fmt.Println("Status: Token expired")
		fmt.Printf("Expired at: %s\n", token.ExpiresAt.Format(time.RFC3339))
		fmt.Println("Run `n8r login` to re-authenticate.")
		return
	}

	fmt.Println("Status: Authenticated")
	fmt.Printf("Token type: %s\n", token.TokenType)
	fmt.Printf("Expires at: %s\n", token.ExpiresAt.Format(time.RFC3339))
	fmt.Printf("Saved at:   %s\n", token.SavedAt.Format(time.RFC3339))
}
