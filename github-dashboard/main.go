package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	sdk "github.com/PortableSheep/delve_sdk"
)

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

// Config holds the plugin configuration
type Config struct {
	GitHubToken     string   `json:"github_token"`
	Repositories    []string `json:"repositories"`
	RefreshInterval int      `json:"refresh_interval"`
}

// GitHubDashboard represents the GitHub Dashboard plugin
type GitHubDashboard struct {
	config Config
	client *http.Client
}

var plugin *sdk.Plugin
var dashboard *GitHubDashboard

// handleHostMessage processes messages from the host application
func handleHostMessage(messageType int, data []byte) {
	log.Printf("GitHub Dashboard received message: Type=%d, Data=%s", messageType, string(data))

	switch messageType {
	case 1: // Get repositories request
		repos, err := dashboard.getRepositories()
		if err != nil {
			log.Printf("Error fetching repositories: %v", err)
			return
		}
		log.Printf("Fetched %d repositories", len(repos))

	case 2: // Get pull requests for repository
		var repoName string
		if err := json.Unmarshal(data, &repoName); err != nil {
			log.Printf("Error unmarshaling repo name: %v", err)
			return
		}
		prs, err := dashboard.getPullRequests(repoName)
		if err != nil {
			log.Printf("Error fetching pull requests for %s: %v", repoName, err)
			return
		}
		log.Printf("Fetched %d pull requests for %s", len(prs), repoName)

	case 3: // Force refresh
		log.Println("Force refresh requested")
		// Trigger refresh logic here
	}
}

// Initialize sets up the dashboard with configuration
func (d *GitHubDashboard) Initialize() error {
	// Try to load config from storage
	config := map[string]interface{}{
		"github_token":     "",
		"repositories":     []string{},
		"refresh_interval": 300,
	}

	// Load stored config
	storedConfig, err := plugin.LoadConfig("dashboard_settings")
	if err == nil && storedConfig != nil {
		configValue := storedConfig.Value.(map[string]interface{})
		if token, ok := configValue["github_token"].(string); ok && token != "" {
			config["github_token"] = token
		}
		if repos, ok := configValue["repositories"].([]interface{}); ok {
			var repoStrings []string
			for _, repo := range repos {
				if repoStr, ok := repo.(string); ok {
					repoStrings = append(repoStrings, repoStr)
				}
			}
			config["repositories"] = repoStrings
		}
		if interval, ok := configValue["refresh_interval"].(float64); ok {
			config["refresh_interval"] = int(interval)
		}
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := json.Unmarshal(configBytes, &d.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// GitHub token is required for API access
	if d.config.GitHubToken == "" {
		log.Println("⚠️  No GitHub token configured - API access will be limited")
	}

	if d.config.RefreshInterval <= 0 {
		d.config.RefreshInterval = 300 // Default to 5 minutes
	}

	d.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return nil
}

// getRepositories fetches repository information
func (d *GitHubDashboard) getRepositories() ([]Repository, error) {
	var repositories []Repository

	if len(d.config.Repositories) == 0 {
		// Return some default repositories for demo
		return []Repository{
			{
				Name:        "delve",
				FullName:    "PortableSheep/delve",
				Description: "Advanced debugging tool for the Go programming language",
				Stars:       1500,
				Forks:       200,
				OpenIssues:  45,
				Language:    "Go",
				UpdatedAt:   time.Now().Add(-2 * time.Hour),
				HTMLURL:     "https://github.com/PortableSheep/delve",
			},
		}, nil
	}

	for _, repoName := range d.config.Repositories {
		repo, err := d.fetchRepository(repoName)
		if err != nil {
			log.Printf("Failed to fetch repository %s: %v", repoName, err)
			continue // Skip failed repositories
		}
		repositories = append(repositories, repo)
	}

	return repositories, nil
}

// getPullRequests fetches pull requests for a repository
func (d *GitHubDashboard) getPullRequests(repoName string) ([]PullRequest, error) {
	if d.config.GitHubToken == "" {
		// Return mock data for demo
		return []PullRequest{
			{
				Number:    123,
				Title:     "Add new feature for dashboard improvements",
				State:     "open",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-2 * time.Hour),
				HTMLURL:   fmt.Sprintf("https://github.com/%s/pull/123", repoName),
				User: struct {
					Login     string `json:"login"`
					AvatarURL string `json:"avatar_url"`
				}{
					Login:     "developer",
					AvatarURL: "https://github.com/identicons/developer.png",
				},
			},
		}, nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+d.config.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := d.client.Do(req)
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

	return pullRequests, nil
}

// fetchRepository gets repository data from GitHub API
func (d *GitHubDashboard) fetchRepository(repoName string) (Repository, error) {
	if d.config.GitHubToken == "" {
		// Return mock data for demo
		return Repository{
			Name:        strings.Split(repoName, "/")[1],
			FullName:    repoName,
			Description: "Demo repository (configure GitHub token for real data)",
			Stars:       100,
			Forks:       20,
			OpenIssues:  5,
			Language:    "Go",
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
			HTMLURL:     fmt.Sprintf("https://github.com/%s", repoName),
		}, nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Repository{}, err
	}

	req.Header.Set("Authorization", "token "+d.config.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := d.client.Do(req)
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

// checkGitHubConnectivity verifies GitHub API connectivity
func (d *GitHubDashboard) checkGitHubConnectivity() error {
	req, err := http.NewRequest("GET", "https://api.github.com", nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
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
func (d *GitHubDashboard) checkTokenValidity() error {
	if d.config.GitHubToken == "" {
		return fmt.Errorf("no GitHub token configured")
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+d.config.GitHubToken)

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid GitHub token")
	}

	return nil
}

// backgroundRefresh runs periodic refreshes
func (d *GitHubDashboard) backgroundRefresh(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(d.config.RefreshInterval) * time.Second)
	defer ticker.Stop()

	log.Printf("Started background refresh with interval: %d seconds", d.config.RefreshInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Background refresh stopped")
			return
		case <-ticker.C:
			// Perform background refresh
			log.Println("Performing background refresh...")
			repos, err := d.getRepositories()
			if err != nil {
				log.Printf("Background refresh failed: %v", err)
			} else {
				log.Printf("Background refresh completed: %d repositories", len(repos))
			}
		}
	}
}

// demonstrateStorage shows how to use the plugin storage system
func demonstrateStorage(p *sdk.Plugin) {
	// Check if storage demo is disabled
	if os.Getenv("DISABLE_STORAGE_DEMO") == "true" {
		log.Println("Storage demonstration disabled (DISABLE_STORAGE_DEMO=true)")
		return
	}

	go func() {
		log.Println("=== GitHub Dashboard Plugin Storage Demonstration (Background) ===")

		// Store some configuration for the GitHub dashboard
		config := map[string]interface{}{
			"theme":               "github",
			"refreshInterval":     300,
			"autoRefresh":         true,
			"showStars":           true,
			"showForks":           true,
			"showIssues":          true,
			"defaultView":         "repositories",
			"compactView":         false,
			"persistentState":     true,
			"enableNotifications": true,
		}

		// Try to store config with timeout
		err := p.StoreConfig("dashboard_settings", config, "1.0.0")
		if err != nil {
			log.Printf("⚠️  Storage not available - dashboard config not stored: %v", err)
		} else {
			log.Println("✓ Stored GitHub dashboard configuration")
		}

		// Store some sample repository data
		sampleRepos := map[string]interface{}{
			"watched_repositories": []map[string]interface{}{
				{"name": "PortableSheep/delve", "lastChecked": time.Now()},
				{"name": "golang/go", "lastChecked": time.Now()},
				{"name": "microsoft/vscode", "lastChecked": time.Now()},
			},
			"user_preferences": map[string]interface{}{
				"defaultRepo":     "PortableSheep/delve",
				"sortBy":          "updated",
				"showArchived":    false,
				"maxReposPerPage": 20,
			},
			"createdAt": time.Now(),
			"version":   "1.0.0",
		}

		err = p.StoreData("repository_data", sampleRepos, "1.0.0")
		if err != nil {
			log.Printf("⚠️  Storage not available - repository data not stored: %v", err)
		} else {
			log.Println("✓ Stored repository data")
		}

		// Show storage statistics
		stats, err := p.GetStats()
		if err != nil {
			log.Printf("⚠️  Storage stats not available: %v", err)
		} else {
			log.Printf("\n=== Storage Statistics ===")
			for storageType, count := range stats {
				log.Printf("%s: %d items", storageType, count)
			}
		}
	}()
}

func main() {
	log.Printf("GitHub Dashboard Plugin launched with arguments: %v", os.Args)

	// Define the plugin's metadata
	pluginInfo := &sdk.RegisterRequest{
		Name:             "github-dashboard",
		Description:      "A comprehensive GitHub dashboard for monitoring repositories, pull requests, and issues with real-time updates",
		UiComponentPath:  "frontend/component.js",
		CustomElementTag: "github-dashboard",
	}

	// Connect to the host and register
	var err error
	plugin, err = sdk.Start(pluginInfo)
	if err != nil {
		log.Fatalf("Failed to start GitHub Dashboard plugin: %v", err)
	}

	// Initialize the dashboard
	dashboard = &GitHubDashboard{}
	if err := dashboard.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize dashboard: %v", err)
	}

	// Demonstrate storage functionality (non-blocking, optional)
	demonstrateStorage(plugin)

	// Start background refresh
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go dashboard.backgroundRefresh(ctx)

	// Test the GitHub API functionality
	log.Println("=== Testing GitHub API Functionality ===")

	// Test connectivity
	if err := dashboard.checkGitHubConnectivity(); err != nil {
		log.Printf("⚠️  GitHub connectivity check failed: %v", err)
	} else {
		log.Println("✓ GitHub API is reachable")
	}

	// Test token validity if configured
	if dashboard.config.GitHubToken != "" {
		if err := dashboard.checkTokenValidity(); err != nil {
			log.Printf("⚠️  Token validation failed: %v", err)
		} else {
			log.Println("✓ GitHub token is valid")
		}
	}

	// Test repository fetching
	repos, err := dashboard.getRepositories()
	if err != nil {
		log.Printf("⚠️  Failed to fetch repositories: %v", err)
	} else {
		log.Printf("✓ Successfully fetched %d repositories", len(repos))
		for _, repo := range repos {
			log.Printf("  - %s (%d stars, %d forks)", repo.FullName, repo.Stars, repo.Forks)
		}
	}

	// Start listening for events from the host
	log.Println("GitHub Dashboard Plugin is running and listening for host events.")
	log.Println("Plugin supports persistent state and real-time GitHub data.")
	log.Println("Note: Storage operations run in background and may timeout if storage service is unavailable.")
	log.Println("To disable storage demo entirely, set environment variable DISABLE_STORAGE_DEMO=true")
	plugin.Listen(handleHostMessage)

	log.Println("GitHub Dashboard Plugin shutting down.")
}
