package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GeminiResponse represents the expected JSON response from the Gemini API.
type GeminiResponse struct {
	Message string `json:"message"`
}

// GenerateWelcomeMessage calls the Gemini API to generate a custom welcome message.
// In a real implementation, ensure to add authentication, error handling, and proper API endpoint configuration.
func GenerateWelcomeMessage(username string) string {
	// Construct the prompt for Gemini. You can adjust the text as needed.
	prompt := fmt.Sprintf("Generate a creative welcome message for %s joining the group.", username)
	payload := map[string]string{
		"prompt": prompt,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fallbackMessage(username)
	}

	// Replace with your actual Gemini API endpoint.
	apiURL := "https://api.gemini.example.com/generate"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fallbackMessage(username)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fallbackMessage(username)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return fallbackMessage(username)
	}

	return geminiResp.Message
}

// fallbackMessage returns a simple welcome message if Gemini generation fails.
func fallbackMessage(username string) string {
	return fmt.Sprintf("Welcome to the group, %s!", username)
}
