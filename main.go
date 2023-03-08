package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type IssueEvent struct {
	Action string `json:"action"`
	Issue  struct {
		Number int    `json:"number"`
		State  string `json:"state"`
		Title  string `json:"title"`
	} `json:"issue"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
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

	// Print the issue title
	fmt.Printf("Issue #%d closed: %s\n", payload.Issue.Number, payload.Issue.Title)

	// Perform your action here...
	fmt.Println("Do something with the closed issue...")
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
