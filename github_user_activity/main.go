package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Type      string    `json:"type"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Action       string   `json:"action"`
	RefType      string   `json:"ref_type"`
	Ref          string   `json:"ref"`
	PushID       int64    `json:"push_id"`
	Size         int      `json:"size"`
	Commits      []Commit `json:"commits"`
	Issue        *Issue   `json:"issue"`
	PullRequest  *Issue   `json:"pull_request"`
}

type Commit struct {
	Message string `json:"message"`
}

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	fetchAndDisplayActivity(username)
}

func fetchAndDisplayActivity(username string) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("User '%s' not found\n", username)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Error: GitHub API returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Printf("No recent activity found for user '%s'\n", username)
		return
	}

	fmt.Printf("Recent activity for %s:\n\n", username)
	displayEvents(events)
}

func displayEvents(events []Event) {
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commitCount := event.Payload.Size
			if commitCount > 0 {
				fmt.Printf("- Pushed %d commit(s) to %s\n", commitCount, event.Repo.Name)
			} else {
				fmt.Printf("- Pushed commits to %s\n", event.Repo.Name)
			}
		case "IssuesEvent":
			if event.Payload.Issue != nil {
				fmt.Printf("- %s issue #%d in %s\n", 
					capitalize(event.Payload.Action), 
					event.Payload.Issue.Number, 
					event.Repo.Name)
			}
		case "WatchEvent":
			fmt.Printf("- Starred %s\n", event.Repo.Name)
		case "ForkEvent":
			fmt.Printf("- Forked %s\n", event.Repo.Name)
		case "CreateEvent":
			if event.Payload.RefType == "repository" {
				fmt.Printf("- Created repository %s\n", event.Repo.Name)
			} else {
				fmt.Printf("- Created %s in %s\n", event.Payload.RefType, event.Repo.Name)
			}
		case "DeleteEvent":
			fmt.Printf("- Deleted %s in %s\n", event.Payload.RefType, event.Repo.Name)
		case "PullRequestEvent":
			if event.Payload.PullRequest != nil {
				fmt.Printf("- %s pull request #%d in %s\n", 
					capitalize(event.Payload.Action), 
					event.Payload.PullRequest.Number, 
					event.Repo.Name)
			}
		case "IssueCommentEvent":
			fmt.Printf("- Commented on issue #%d in %s\n", 
				event.Payload.Issue.Number, 
				event.Repo.Name)
		case "PullRequestReviewEvent":
			fmt.Printf("- Reviewed pull request in %s\n", event.Repo.Name)
		case "PullRequestReviewCommentEvent":
			fmt.Printf("- Commented on pull request in %s\n", event.Repo.Name)
		default:
			fmt.Printf("- %s in %s\n", formatEventType(event.Type), event.Repo.Name)
		}
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

func formatEventType(eventType string) string {
	// Remove "Event" suffix and add spaces before capitals
	result := ""
	for i, char := range eventType {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result += " "
		}
		result += string(char)
	}
	// Remove " Event" at the end
	if len(result) > 6 && result[len(result)-6:] == " Event" {
		result = result[:len(result)-6]
	}
	return result
}
