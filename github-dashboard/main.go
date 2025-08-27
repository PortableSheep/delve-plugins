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

	sdk "github.com/PortableSheep/delve-sdk"
)

// Message types for plugin communication
const (
	MessageTypeGetRepositories = 1
	MessageTypeGetPullRequests = 2
	MessageTypeRefresh         = 3
	MessageTypeGetConfig       = 4
	MessageTypeSetConfig       = 5
	MessageTypeHealthCheck     = 6
)

// API Response structures
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Language        string    `json:"language"`
	UpdatedAt       time.Time `json:"updated_at"`
	HTMLURL         string    `json:"html_url"`
	Private         bool      `json:"private"`
	Fork            bool      `json:"fork"`
	Archived        bool      `json:"archived"`
	Disabled        bool      `json:"disabled"`
	CreatedAt       time.Time `json:"created_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Size            int       `json:"size"`
	DefaultBranch   string    `json:"default_branch"`
	Topics          []string  `json:"topics"`
	License         *struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SPDXID string `json:"spdx_id"`
	} `json:"license"`
	Owner *struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
		HTMLURL   string `json:"html_url"`
		Type      string `json:"type"`
	} `json:"owner"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID             int        `json:"id"`
	Number         int        `json:"number"`
	Title          string     `json:"title"`
	Body           string     `json:"body"`
	State          string     `json:"state"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ClosedAt       *time.Time `json:"closed_at"`
	MergedAt       *time.Time `json:"merged_at"`
	HTMLURL        string     `json:"html_url"`
	Draft          bool       `json:"draft"`
	Mergeable      *bool      `json:"mergeable"`
	MergeableState string     `json:"mergeable_state"`
	Additions      int        `json:"additions"`
	Deletions      int        `json:"deletions"`
	ChangedFiles   int        `json:"changed_files"`
	User           *struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
		HTMLURL   string `json:"html_url"`
		Type      string `json:"type"`
	} `json:"user"`
	Head *struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		SHA   string `json:"sha"`
	} `json:"head"`
	Base *struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		SHA   string `json:"sha"`
	} `json:"base"`
	RequestedReviewers []struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"requested_reviewers"`
}

// Config holds the plugin configuration
type Config struct {
	GitHubToken      string   `json:"github_token"`
	Repositories     []string `json:"repositories"`
	RefreshInterval  int      `json:"refresh_interval"`
	CompactView      bool     `json:"compact_view"`
	ShowPrivateRepos bool     `json:"show_private_repos"`
	MaxReposPerPage  int      `json:"max_repos_per_page"`
	MaxPRsPerRepo    int      `json:"max_prs_per_repo"`
	CacheTimeout     int      `json:"cache_timeout"`
}

// APIRequest represents incoming API requests
type APIRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// GitHubDashboard represents the GitHub Dashboard plugin
type GitHubDashboard struct {
	config          Config
	httpClient      *http.Client
	plugin          *sdk.Plugin
	ctx             context.Context
	cancel          context.CancelFunc
	repoCache       map[string]*Repository
	prCache         map[string][]*PullRequest
	cacheTime       map[string]time.Time
	rateLimits      map[string]int
	rateLimitResets map[string]time.Time
}

var dashboard *GitHubDashboard

// Initialize sets up the dashboard with configuration
func (d *GitHubDashboard) Initialize() error {
	log.Println("🔧 Initializing GitHub Dashboard...")

	// Set default config
	d.config = Config{
		GitHubToken:      "",
		Repositories:     []string{},
		RefreshInterval:  300,
		CompactView:      false,
		ShowPrivateRepos: true,
		MaxReposPerPage:  50,
		MaxPRsPerRepo:    20,
		CacheTimeout:     300, // 5 minutes
	}

	// Initialize caches
	d.repoCache = make(map[string]*Repository)
	d.prCache = make(map[string][]*PullRequest)
	d.cacheTime = make(map[string]time.Time)
	d.rateLimits = make(map[string]int)
	d.rateLimitResets = make(map[string]time.Time)

	// Setup HTTP client with timeout
	d.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	// Load stored config
	if err := d.loadConfig(); err != nil {
		log.Printf("⚠️  Warning: Failed to load config: %v", err)
	}

	// Validate config
	if err := d.validateConfig(); err != nil {
		log.Printf("⚠️  Warning: Config validation failed: %v", err)
	}

	log.Println("✅ GitHub Dashboard initialized successfully")
	return nil
}

// loadConfig attempts to load configuration from storage
func (d *GitHubDashboard) loadConfig() error {
	if d.plugin == nil {
		return fmt.Errorf("plugin not available")
	}

	storedConfig, err := d.plugin.LoadConfig("dashboard_settings")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if storedConfig != nil && storedConfig.Value != nil {
		configValue, ok := storedConfig.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid config format")
		}

		// Parse configuration
		if token, ok := configValue["github_token"].(string); ok && token != "" {
			d.config.GitHubToken = token
		}

		if repos, ok := configValue["repositories"].([]interface{}); ok {
			var repoStrings []string
			for _, repo := range repos {
				if repoStr, ok := repo.(string); ok {
					repoStrings = append(repoStrings, repoStr)
				}
			}
			d.config.Repositories = repoStrings
		}

		if interval, ok := configValue["refresh_interval"].(float64); ok {
			d.config.RefreshInterval = int(interval)
		}

		if compact, ok := configValue["compact_view"].(bool); ok {
			d.config.CompactView = compact
		}

		if showPrivate, ok := configValue["show_private_repos"].(bool); ok {
			d.config.ShowPrivateRepos = showPrivate
		}

		log.Printf("📋 Loaded config: %d repositories configured", len(d.config.Repositories))
	}

	return nil
}

// saveConfig saves current configuration to storage
func (d *GitHubDashboard) saveConfig() error {
	if d.plugin == nil {
		return fmt.Errorf("plugin not available")
	}

	configMap := map[string]interface{}{
		"github_token":       d.config.GitHubToken,
		"repositories":       d.config.Repositories,
		"refresh_interval":   d.config.RefreshInterval,
		"compact_view":       d.config.CompactView,
		"show_private_repos": d.config.ShowPrivateRepos,
		"max_repos_per_page": d.config.MaxReposPerPage,
		"max_prs_per_repo":   d.config.MaxPRsPerRepo,
		"cache_timeout":      d.config.CacheTimeout,
	}

	return d.plugin.StoreConfig("dashboard_settings", configMap, "1.0.0")
}

// validateConfig validates the current configuration
func (d *GitHubDashboard) validateConfig() error {
	if d.config.RefreshInterval < 30 {
		d.config.RefreshInterval = 30
		log.Println("⚠️  Refresh interval too low, set to 30 seconds")
	}

	if d.config.MaxReposPerPage > 100 {
		d.config.MaxReposPerPage = 100
		log.Println("⚠️  Max repos per page too high, set to 100")
	}

	if d.config.MaxPRsPerRepo > 50 {
		d.config.MaxPRsPerRepo = 50
		log.Println("⚠️  Max PRs per repo too high, set to 50")
	}

	if d.config.GitHubToken == "" {
		log.Println("⚠️  No GitHub token configured - using public API with rate limits")
	}

	return nil
}

// makeGitHubRequest performs a request to GitHub API with proper error handling
func (d *GitHubDashboard) makeGitHubRequest(endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("https://api.github.com%s", endpoint)

	req, err := http.NewRequestWithContext(d.ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Delve-GitHub-Dashboard/1.0")

	// Add authentication if available
	if d.config.GitHubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", d.config.GitHubToken))
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Update rate limit info
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		// Store rate limit info for monitoring
	}

	// Handle rate limiting
	if resp.StatusCode == 403 && strings.Contains(resp.Header.Get("X-RateLimit-Remaining"), "0") {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GitHub API error: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// getRepositories fetches repository information
func (d *GitHubDashboard) getRepositories() ([]*Repository, error) {
	log.Println("📡 Fetching repositories...")

	var repositories []*Repository

	// If no repositories configured, return demo data
	if len(d.config.Repositories) == 0 {
		log.Println("📊 No repositories configured, returning demo data")
		return d.getDemoRepositories(), nil
	}

	// Fetch configured repositories
	for _, repoName := range d.config.Repositories {
		// Check cache first
		cacheKey := fmt.Sprintf("repo:%s", repoName)
		if cached, exists := d.repoCache[cacheKey]; exists {
			if cacheTime, exists := d.cacheTime[cacheKey]; exists {
				if time.Since(cacheTime) < time.Duration(d.config.CacheTimeout)*time.Second {
					repositories = append(repositories, cached)
					continue
				}
			}
		}

		// Fetch from API
		repo, err := d.fetchRepository(repoName)
		if err != nil {
			log.Printf("❌ Failed to fetch repository %s: %v", repoName, err)
			continue
		}

		if repo != nil {
			repositories = append(repositories, repo)
			// Cache the result
			d.repoCache[cacheKey] = repo
			d.cacheTime[cacheKey] = time.Now()
		}
	}

	log.Printf("✅ Successfully fetched %d repositories", len(repositories))
	return repositories, nil
}

// fetchRepository gets repository data from GitHub API
func (d *GitHubDashboard) fetchRepository(repoName string) (*Repository, error) {
	endpoint := fmt.Sprintf("/repos/%s", repoName)

	resp, err := d.makeGitHubRequest(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("failed to decode repository data: %w", err)
	}

	return &repo, nil
}

// getPullRequests fetches pull requests for a repository
func (d *GitHubDashboard) getPullRequests(repoName string) ([]*PullRequest, error) {
	log.Printf("📡 Fetching pull requests for %s...", repoName)

	// Check cache first
	cacheKey := fmt.Sprintf("prs:%s", repoName)
	if cached, exists := d.prCache[cacheKey]; exists {
		if cacheTime, exists := d.cacheTime[cacheKey]; exists {
			if time.Since(cacheTime) < time.Duration(d.config.CacheTimeout)*time.Second {
				log.Printf("📋 Using cached data for %s", repoName)
				return cached, nil
			}
		}
	}

	// If no token, return demo data
	if d.config.GitHubToken == "" {
		log.Println("📊 No GitHub token configured, returning demo data")
		return d.getDemoPullRequests(repoName), nil
	}

	endpoint := fmt.Sprintf("/repos/%s/pulls?state=open&per_page=%d", repoName, d.config.MaxPRsPerRepo)

	resp, err := d.makeGitHubRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
	}
	defer resp.Body.Close()

	var pullRequests []*PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, fmt.Errorf("failed to decode pull requests data: %w", err)
	}

	// Cache the results
	d.prCache[cacheKey] = pullRequests
	d.cacheTime[cacheKey] = time.Now()

	log.Printf("✅ Successfully fetched %d pull requests for %s", len(pullRequests), repoName)
	return pullRequests, nil
}

// getDemoRepositories returns demo repository data
func (d *GitHubDashboard) getDemoRepositories() []*Repository {
	return []*Repository{
		{
			ID:              1,
			Name:            "delve",
			FullName:        "PortableSheep/delve",
			Description:     "Advanced plugin system for Go applications - Configure your GitHub token for real data",
			StargazersCount: 1250,
			ForksCount:      180,
			OpenIssuesCount: 23,
			Language:        "Go",
			UpdatedAt:       time.Now().Add(-2 * time.Hour),
			HTMLURL:         "https://github.com/PortableSheep/delve",
			Private:         false,
			Fork:            false,
			Archived:        false,
			Disabled:        false,
			CreatedAt:       time.Now().Add(-365 * 24 * time.Hour),
			PushedAt:        time.Now().Add(-2 * time.Hour),
			Size:            15420,
			DefaultBranch:   "main",
			Topics:          []string{"go", "plugins", "development"},
		},
		{
			ID:              2,
			Name:            "delve-plugins",
			FullName:        "PortableSheep/delve-plugins",
			Description:     "Official plugin collection for Delve - Monitor GitHub repositories and pull requests",
			StargazersCount: 89,
			ForksCount:      34,
			OpenIssuesCount: 8,
			Language:        "Go",
			UpdatedAt:       time.Now().Add(-1 * time.Hour),
			HTMLURL:         "https://github.com/PortableSheep/delve-plugins",
			Private:         false,
			Fork:            false,
			Archived:        false,
			Disabled:        false,
			CreatedAt:       time.Now().Add(-180 * 24 * time.Hour),
			PushedAt:        time.Now().Add(-1 * time.Hour),
			Size:            8935,
			DefaultBranch:   "main",
			Topics:          []string{"plugins", "extensions", "github"},
		},
	}
}

// getDemoPullRequests returns demo pull request data
func (d *GitHubDashboard) getDemoPullRequests(repoName string) []*PullRequest {
	return []*PullRequest{
		{
			ID:             42,
			Number:         42,
			Title:          "Add enhanced GitHub dashboard features",
			Body:           "This PR adds new features to the GitHub dashboard including better UI, caching, and error handling.",
			State:          "open",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			UpdatedAt:      time.Now().Add(-1 * time.Hour),
			HTMLURL:        fmt.Sprintf("https://github.com/%s/pull/42", repoName),
			Draft:          false,
			MergeableState: "clean",
			Additions:      245,
			Deletions:      89,
			ChangedFiles:   12,
			User: &struct {
				Login     string `json:"login"`
				ID        int    `json:"id"`
				AvatarURL string `json:"avatar_url"`
				HTMLURL   string `json:"html_url"`
				Type      string `json:"type"`
			}{
				Login:     "developer",
				ID:        12345,
				AvatarURL: "https://github.com/identicons/developer.png",
				HTMLURL:   "https://github.com/developer",
				Type:      "User",
			},
			Head: &struct {
				Label string `json:"label"`
				Ref   string `json:"ref"`
				SHA   string `json:"sha"`
			}{
				Label: "developer:feature/enhanced-dashboard",
				Ref:   "feature/enhanced-dashboard",
				SHA:   "abc123def456",
			},
			Base: &struct {
				Label string `json:"label"`
				Ref   string `json:"ref"`
				SHA   string `json:"sha"`
			}{
				Label: "PortableSheep:main",
				Ref:   "main",
				SHA:   "def456ghi789",
			},
		},
		{
			ID:             38,
			Number:         38,
			Title:          "Fix responsive design issues",
			Body:           "Addresses mobile layout problems and improves accessibility.",
			State:          "open",
			CreatedAt:      time.Now().Add(-48 * time.Hour),
			UpdatedAt:      time.Now().Add(-2 * time.Hour),
			HTMLURL:        fmt.Sprintf("https://github.com/%s/pull/38", repoName),
			Draft:          true,
			MergeableState: "draft",
			Additions:      156,
			Deletions:      23,
			ChangedFiles:   8,
			User: &struct {
				Login     string `json:"login"`
				ID        int    `json:"id"`
				AvatarURL string `json:"avatar_url"`
				HTMLURL   string `json:"html_url"`
				Type      string `json:"type"`
			}{
				Login:     "designer",
				ID:        67890,
				AvatarURL: "https://github.com/identicons/designer.png",
				HTMLURL:   "https://github.com/designer",
				Type:      "User",
			},
		},
	}
}

// handleHostMessage processes messages from the host application
func handleHostMessage(messageType int, data []byte) {
	log.Printf("📨 GitHub Dashboard received message: Type=%d, Data=%s", messageType, string(data))

	var response APIResponse

	switch messageType {
	case MessageTypeGetRepositories:
		repos, err := dashboard.getRepositories()
		if err != nil {
			log.Printf("❌ Error fetching repositories: %v", err)
			response = APIResponse{
				Success: false,
				Error:   err.Error(),
			}
		} else {
			response = APIResponse{
				Success: true,
				Data:    repos,
			}
		}

	case MessageTypeGetPullRequests:
		var request APIRequest
		if err := json.Unmarshal(data, &request); err != nil {
			log.Printf("❌ Error unmarshaling PR request: %v", err)
			response = APIResponse{
				Success: false,
				Error:   "Invalid request format",
			}
		} else {
			repoName, ok := request.Params["repository"].(string)
			if !ok {
				log.Printf("❌ Missing repository parameter in PR request")
				response = APIResponse{
					Success: false,
					Error:   "Missing repository parameter",
				}
			} else {
				prs, err := dashboard.getPullRequests(repoName)
				if err != nil {
					log.Printf("❌ Error fetching pull requests: %v", err)
					response = APIResponse{
						Success: false,
						Error:   err.Error(),
					}
				} else {
					response = APIResponse{
						Success: true,
						Data:    prs,
					}
				}
			}
		}

	case MessageTypeRefresh:
		// Clear caches to force refresh
		dashboard.repoCache = make(map[string]*Repository)
		dashboard.prCache = make(map[string][]*PullRequest)
		dashboard.cacheTime = make(map[string]time.Time)

		log.Println("🔄 Cache cleared, data will be refreshed on next request")
		response = APIResponse{
			Success: true,
			Data:    "Cache cleared successfully",
		}

	case MessageTypeGetConfig:
		response = APIResponse{
			Success: true,
			Data:    dashboard.config,
		}

	case MessageTypeSetConfig:
		var newConfig Config
		if err := json.Unmarshal(data, &newConfig); err != nil {
			log.Printf("❌ Error unmarshaling config: %v", err)
			response = APIResponse{
				Success: false,
				Error:   "Invalid config format",
			}
		} else {
			dashboard.config = newConfig
			if err := dashboard.validateConfig(); err != nil {
				log.Printf("⚠️  Config validation warning: %v", err)
			}
			if err := dashboard.saveConfig(); err != nil {
				log.Printf("❌ Error saving config: %v", err)
				response = APIResponse{
					Success: false,
					Error:   "Failed to save config",
				}
			} else {
				log.Println("✅ Configuration updated successfully")
				response = APIResponse{
					Success: true,
					Data:    "Configuration updated",
				}
			}
		}

	case MessageTypeHealthCheck:
		// Perform health check
		healthy := true
		message := "Plugin is healthy"

		// Check GitHub API connectivity if token is available
		if dashboard.config.GitHubToken != "" {
			_, err := dashboard.makeGitHubRequest("/user")
			if err != nil {
				healthy = false
				message = fmt.Sprintf("GitHub API connectivity failed: %v", err)
			}
		}

		response = APIResponse{
			Success: healthy,
			Data: map[string]interface{}{
				"healthy":          healthy,
				"message":          message,
				"repos_configured": len(dashboard.config.Repositories),
				"cache_entries":    len(dashboard.repoCache) + len(dashboard.prCache),
				"has_github_token": dashboard.config.GitHubToken != "",
			},
		}

	default:
		log.Printf("❓ Unknown message type: %d", messageType)
		response = APIResponse{
			Success: false,
			Error:   "Unknown message type",
		}
	}

	// Send response back to host (if needed - depends on SDK implementation)
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ Error marshaling response: %v", err)
		return
	}

	log.Printf("📤 Sending response: %s", string(responseData))
}

// backgroundTasks runs periodic maintenance tasks
func (d *GitHubDashboard) backgroundTasks(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(d.config.RefreshInterval) * time.Second)
	defer ticker.Stop()

	log.Printf("🔄 Started background tasks with %d second interval", d.config.RefreshInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Background tasks stopped")
			return
		case <-ticker.C:
			// Perform background refresh of configured repositories
			if len(d.config.Repositories) > 0 {
				log.Println("🔄 Performing background repository refresh...")
				repos, err := d.getRepositories()
				if err != nil {
					log.Printf("❌ Background refresh failed: %v", err)
				} else {
					log.Printf("✅ Background refresh completed: %d repositories", len(repos))
				}
			}
		}
	}
}

// gracefulShutdown handles cleanup on shutdown
func (d *GitHubDashboard) gracefulShutdown() {
	log.Println("🛑 Starting graceful shutdown...")

	// Save current config
	if err := d.saveConfig(); err != nil {
		log.Printf("❌ Failed to save config during shutdown: %v", err)
	}

	// Cancel context to stop background tasks
	if d.cancel != nil {
		d.cancel()
	}

	log.Println("✅ Graceful shutdown completed")
}

func main() {
	log.Printf("🚀 GitHub Dashboard Plugin v1.0.0 starting with arguments: %v", os.Args)

	// Define the plugin's metadata
	pluginInfo := &sdk.RegisterRequest{
		Name:             "github-dashboard",
		Description:      "Monitor GitHub repositories and pull requests with real-time updates",
		UiComponentPath:  "frontend/component.js",
		CustomElementTag: "github-dashboard",
	}

	// Connect to the host and register
	var err error
	plugin, err := sdk.Start(pluginInfo)
	if err != nil {
		log.Fatalf("❌ Failed to start GitHub Dashboard plugin: %v", err)
	}

	// Initialize the dashboard
	ctx, cancel := context.WithCancel(context.Background())
	dashboard = &GitHubDashboard{
		plugin: plugin,
		ctx:    ctx,
		cancel: cancel,
	}

	if err := dashboard.Initialize(); err != nil {
		log.Fatalf("❌ Failed to initialize dashboard: %v", err)
	}

	// Start background tasks
	go dashboard.backgroundTasks(ctx)

	// Set up graceful shutdown
	defer dashboard.gracefulShutdown()

	// Start health monitoring
	plugin.StartHeartbeat(30*time.Second, 60*time.Second)

	// Test GitHub API connectivity on startup
	log.Println("🔍 Testing GitHub API connectivity...")
	if dashboard.config.GitHubToken != "" {
		if _, err := dashboard.makeGitHubRequest("/user"); err != nil {
			log.Printf("⚠️  GitHub API connectivity test failed: %v", err)
			log.Println("💡 Check your GitHub token configuration")
		} else {
			log.Println("✅ GitHub API connectivity confirmed")
		}
	} else {
		log.Println("⚠️  No GitHub token configured - using public API with rate limits")
	}

	// Test repository fetching on startup
	log.Println("🧪 Testing repository fetching...")
	repos, err := dashboard.getRepositories()
	if err != nil {
		log.Printf("⚠️  Repository fetch test failed: %v", err)
	} else {
		log.Printf("✅ Successfully fetched %d repositories", len(repos))
		for _, repo := range repos {
			log.Printf("   📁 %s (%d stars, %d forks)", repo.FullName, repo.StargazersCount, repo.ForksCount)
		}
	}

	// Start listening for events from the host
	log.Println("👂 GitHub Dashboard Plugin is running and listening for host events...")
	log.Println("📡 Ready to serve repository and pull request data")
	log.Printf("🔧 Configuration: %d repositories, %ds refresh interval", len(dashboard.config.Repositories), dashboard.config.RefreshInterval)

	plugin.Listen(handleHostMessage)

	log.Println("🔚 GitHub Dashboard Plugin shutting down...")
}
