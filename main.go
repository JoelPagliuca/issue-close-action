package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type IssueEvent struct {
	Action string `json:"action"`
	Issue  struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	} `json:"issue"`
}

func main() {
	// Parse the event payload
	payloadPath := os.Getenv("GITHUB_EVENT_PATH")
	payload, err := parsePayload(payloadPath)
	if err != nil {
		fmt.Printf("Error parsing payload: %v", err)
		return
	}

	// Check if the event was a closed issue
	if payload.Action != "closed" {
		fmt.Printf("Not a closed issue event, action=%s\n", payload.Action)
		return
	}

	// Send the issue title to a Discord webhook
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL != "" {
		err = sendDiscordWebhook(webhookURL, payload.Issue.Title)
		if err != nil {
			fmt.Printf("Error sending Discord webhook: %v", err)
			return
		}
		fmt.Println("Sent issue title to Discord!")
	}

	// Print a message to the console
	fmt.Printf("Issue #%d closed: %s\n", payload.Issue.Number, payload.Issue.Title)
}

// Helper function to parse the Github event payload
func parsePayload(path string) (*IssueEvent, error) {
	payloadFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open payload file: %v", err)
	}
	defer payloadFile.Close()

	var payload IssueEvent
	err = json.NewDecoder(payloadFile).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload file: %v", err)
	}

	return &payload, nil
}

// Helper function to send a message to a Discord webhook
func sendDiscordWebhook(webhookURL string, message string) error {
	// Create the payload data
	payload := map[string]string{
		"content": message,
	}

	// Convert the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %v", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", webhookURL, strings.NewReader(string(payloadJSON)))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	return nil
}
