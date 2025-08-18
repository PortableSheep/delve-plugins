package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Plugin represents the GitHub Dashboard plugin
type Plugin struct {
	config Config
	client *http.Client
}

// Config holds the plugin configuration
type Config struct {
	GitHubToken     string   `json:"github_token"`
	Repositories    []string `json:"repositories"`
	RefreshInterval int      `json:"refresh_interval"`
}

// Repository represents a GitHub repository
type Repository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	OpenIssues  int       `json:"open_issues_count"`
	Language    string    `json:"language"`
	UpdatedAt   time.Time `json:"updated_at"`
	HTMLURL     string    `json:"html_url"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	HTMLURL   string    `json:"html_url"`
	User      struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`
}

// Initialize sets up the plugin with configuration
func (p *Plugin) Initialize(config map[string]interface{}) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := json.Unmarshal(configBytes, &p.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if p.config.GitHubToken == "" {
		return fmt.Errorf("github_token is required")
	}

	if p.config.RefreshInterval <= 0 {
		p.config.RefreshInterval = 300 // Default to 5 minutes
	}

	p.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return nil
}

// Start begins the plugin operation
func (p *Plugin) Start(ctx context.Context) error {
	// Start background refresh goroutine
	go p.backgroundRefresh(ctx)
	return nil
}

// Stop shuts down the plugin
func (p *Plugin) Stop() error {
	return nil
}

// GetInfo returns plugin information
func (p *Plugin) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":         "GitHub Dashboard",
		"version":      "1.0.0",
		"status":       "running",
		"repositories": len(p.config.Repositories),
	}
}

// HandleRequest processes HTTP requests from the frontend
func (p *Plugin) HandleRequest(method, path string, body []byte) ([]byte, error) {
	switch {
	case method == "GET" && path == "/repositories":
		return p.getRepositories()
	case method == "GET" && strings.HasPrefix(path, "/repositories/") && strings.HasSuffix(path, "/pulls"):
		repo := strings.TrimPrefix(path, "/repositories/")
		repo = strings.TrimSuffix(repo, "/pulls")
		return p.getPullRequests(repo)
	case method == "POST" && path == "/refresh":
		return p.forceRefresh()
	default:
		return nil, fmt.Errorf("endpoint not found: %s %s", method, path)
	}
}

// HealthCheck performs health checks
func (p *Plugin) HealthCheck(checkName string) error {
	switch checkName {
	case "github_api_connectivity":
		return p.checkGitHubConnectivity()
	case "token_validity":
		return p.checkTokenValidity()
	default:
		return fmt.Errorf("unknown health check: %s", checkName)
	}
}

// getRepositories fetches repository information
func (p *Plugin) getRepositories() ([]byte, error) {
	var repositories []Repository

	for _, repoName := range p.config.Repositories {
		repo, err := p.fetchRepository(repoName)
		if err != nil {
			continue // Skip failed repositories
		}
		repositories = append(repositories, repo)
	}

	return json.Marshal(repositories)
}

// getPullRequests fetches pull requests for a repository
func (p *Plugin) getPullRequests(repoName string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+p.config.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var pullRequests []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, err
	}

	return json.Marshal(pullRequests)
}

// forceRefresh triggers an immediate refresh
func (p *Plugin) forceRefresh() ([]byte, error) {
	// Trigger refresh logic here
	return json.Marshal(map[string]string{"status": "refreshed"})
}

// fetchRepository gets repository data from GitHub API
func (p *Plugin) fetchRepository(repoName string) (Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Repository{}, err
	}

	req.Header.Set("Authorization", "token "+p.config.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return Repository{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Repository{}, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return Repository{}, err
	}

	return repo, nil
}

// backgroundRefresh runs periodic refreshes
func (p *Plugin) backgroundRefresh(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(p.config.RefreshInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Perform background refresh
			_, _ = p.getRepositories()
		}
	}
}

// checkGitHubConnectivity verifies GitHub API connectivity
func (p *Plugin) checkGitHubConnectivity() error {
	req, err := http.NewRequest("GET", "https://api.github.com", nil)
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API not reachable: %d", resp.StatusCode)
	}

	return nil
}

// checkTokenValidity verifies the GitHub token is valid
func (p *Plugin) checkTokenValidity() error {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+p.config.GitHubToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid GitHub token")
	}

	return nil
}

// GetMenuItems returns menu items for the plugin UI
func (p *Plugin) GetMenuItems() ([]map[string]interface{}, error) {
	// Return menu items for the GitHub Dashboard
	return []map[string]interface{}{
		{
			"id":    "dashboard",
			"text":  "GitHub Dashboard",
			"type":  "view",
			"icon":  "github",
			"group": "tools",
			"order": 1,
		},
		{
			"id":    "refresh",
			"text":  "Refresh Data",
			"type":  "action",
			"icon":  "refresh",
			"group": "actions",
			"order": 1,
			"data":  map[string]interface{}{"action": "refresh"},
		},
		{
			"id":       "repositories",
			"text":     "Repositories",
			"type":     "menu",
			"icon":     "folder",
			"group":    "views",
			"order":    2,
			"children": []map[string]interface{}{},
		},
	}, nil
}

// NewPlugin creates a new plugin instance
func NewPlugin() PluginAPI {
	return &Plugin{}
}

// Entry point for the plugin
func main() {
	plugin := NewPlugin()
	fmt.Println("Starting GitHub Dashboard plugin...")
	if err := RunPluginServer(plugin); err != nil {
		fmt.Printf("Plugin server error: %v\n", err)
		os.Exit(1)
	}
}
