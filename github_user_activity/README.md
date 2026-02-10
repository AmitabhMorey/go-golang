# GitHub User Activity CLI 🚀

A simple command-line interface (CLI) application to fetch and display recent GitHub activity for any user. Built with Go.

**Project URL:** https://roadmap.sh/projects/github-user-activity

## Features ✨

- 📊 Fetch recent activity using the GitHub API
- 🎯 Display various event types (pushes, stars, forks, issues, pull requests, etc.)
- 🎨 Clean and readable terminal output
- ⚠️ Error handling for invalid usernames
- 🌐 Real-time data from GitHub

## Requirements 📋

- Go 1.21 or higher
- Internet connection to access GitHub API

## Installation 🔧

1. Clone or download this project
2. Navigate to the project directory:
```bash
cd github_user_activity
```

3. Build the application:
```bash
go build -o github-activity
```

## Usage 💻

### Fetch user activity

```bash
./github-activity <username>
```

### Example

```bash
./github-activity torvalds
# Output: Recent activity for torvalds:
```

```bash
./github-activity gaearon
```

### Output Example

```
Recent activity for gaearon:

- Commented on issue #4624 in bluesky-social/atproto
- Closed issue #4624 in bluesky-social/atproto
- Starred treethought/obsidian-atmosphere
- Merged pull request #499 in bluesky-social/atproto-website
- Pushed commits to gaearon/atproto
- Opened pull request #4610 in bluesky-social/atproto
- Forked bluesky-social/atproto-website
```

## Supported Event Types 📌

- **PushEvent**: Commits pushed to a repository
- **IssuesEvent**: Issues opened, closed, or reopened
- **WatchEvent**: Repository starred
- **ForkEvent**: Repository forked
- **CreateEvent**: Repository, branch, or tag created
- **DeleteEvent**: Branch or tag deleted
- **PullRequestEvent**: Pull request opened, closed, or merged
- **IssueCommentEvent**: Comment on an issue
- **PullRequestReviewEvent**: Pull request reviewed
- **PullRequestReviewCommentEvent**: Comment on pull request review

## Example Workflow 🔄

```bash
# Check your own activity
./github-activity yourusername

# Check activity of popular developers
./github-activity torvalds
./github-activity gaearon
./github-activity tj

# Check activity of organizations
./github-activity github
```

## Error Handling ⚡

The application handles common errors gracefully:
- Invalid or non-existent usernames
- Network connectivity issues
- GitHub API rate limiting
- JSON parsing errors

## Project Structure 📁

```
github_user_activity/
├── main.go          # Main application code
├── go.mod           # Go module file
├── README.md        # This file
└── github-activity  # Compiled binary (created after build)
```

## License 📄

This project is open source and available for educational purposes.
