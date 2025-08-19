import { defineComponent, ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue';

// Resource cleanup manager
class ResourceManager {
    constructor() {
        this.cleanupHandlers = new Set();
        this.isShuttingDown = false;
        this.intervals = new Set();
        this.timeouts = new Set();
        this.eventListeners = new Map();
        this.abortControllers = new Set();
    }

    onCleanup(handler) {
        if (typeof handler === 'function') {
            this.cleanupHandlers.add(handler);
        }
    }

    setManagedInterval(handler, delay) {
        const id = setInterval(handler, delay);
        this.intervals.add(id);
        return id;
    }

    setManagedTimeout(handler, delay) {
        const id = setTimeout(handler, delay);
        this.timeouts.add(id);
        return id;
    }

    addManagedEventListener(target, event, handler, options = {}) {
        target.addEventListener(event, handler, options);
        if (!this.eventListeners.has(target)) {
            this.eventListeners.set(target, []);
        }
        this.eventListeners.get(target).push({ event, handler, options });
    }

    managedFetch(url, options = {}) {
        const controller = new AbortController();
        this.abortControllers.add(controller);

        return fetch(url, {
            ...options,
            signal: controller.signal
        }).finally(() => {
            this.abortControllers.delete(controller);
        });
    }

    cleanup() {
        if (this.isShuttingDown) return;
        this.isShuttingDown = true;

        console.log('[GitHub Dashboard] Starting cleanup');

        // Clear intervals
        this.intervals.forEach(id => clearInterval(id));
        this.intervals.clear();

        // Clear timeouts
        this.timeouts.forEach(id => clearTimeout(id));
        this.timeouts.clear();

        // Remove event listeners
        this.eventListeners.forEach((listeners, target) => {
            listeners.forEach(({ event, handler, options }) => {
                target.removeEventListener(event, handler, options);
            });
        });
        this.eventListeners.clear();

        // Abort pending requests
        this.abortControllers.forEach(controller => {
            if (!controller.signal.aborted) {
                controller.abort();
            }
        });
        this.abortControllers.clear();

        // Run custom cleanup handlers
        this.cleanupHandlers.forEach(handler => {
            try {
                handler();
            } catch (error) {
                console.warn('[GitHub Dashboard] Error in cleanup handler:', error);
            }
        });
        this.cleanupHandlers.clear();

        console.log('[GitHub Dashboard] Cleanup completed');
    }
}

const resourceManager = new ResourceManager();

// GitHub API helper
class GitHubAPI {
    constructor() {
        this.baseURL = 'https://api.github.com';
        this.rateLimitRemaining = 60;
        this.rateLimitReset = Date.now();
    }

    async makeRequest(endpoint, options = {}) {
        const token = await this.getStoredToken();
        const headers = {
            'Accept': 'application/vnd.github.v3+json',
            'User-Agent': 'Delve-GitHub-Dashboard/1.0',
            ...options.headers
        };

        if (token) {
            headers['Authorization'] = `token ${token}`;
        }

        try {
            const response = await resourceManager.managedFetch(`${this.baseURL}${endpoint}`, {
                ...options,
                headers
            });

            // Update rate limit info
            this.rateLimitRemaining = parseInt(response.headers.get('X-RateLimit-Remaining') || '60');
            this.rateLimitReset = parseInt(response.headers.get('X-RateLimit-Reset') || Date.now() / 1000) * 1000;

            if (!response.ok) {
                throw new Error(`GitHub API error: ${response.status} ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            if (error.name === 'AbortError') {
                console.log('Request aborted during cleanup');
                return null;
            }
            throw error;
        }
    }

    async getStoredToken() {
        try {
            if (window.delvePlugin && window.delvePlugin.loadConfig) {
                const config = await window.delvePlugin.loadConfig('dashboard_settings', '1.0.0');
                return config?.github_token || '';
            }
            return localStorage.getItem('github_dashboard_token') || '';
        } catch (error) {
            console.warn('Failed to load GitHub token:', error);
            return '';
        }
    }

    async getRepositories(repos = []) {
        if (!repos.length) {
            // Return demo repositories if none configured
            return [
                {
                    name: 'delve',
                    full_name: 'PortableSheep/delve',
                    description: 'Advanced debugging tool for Go',
                    stargazers_count: 1500,
                    forks_count: 200,
                    open_issues_count: 45,
                    language: 'Go',
                    updated_at: new Date().toISOString(),
                    html_url: 'https://github.com/PortableSheep/delve'
                }
            ];
        }

        const repositories = [];
        for (const repoName of repos) {
            try {
                const repo = await this.makeRequest(`/repos/${repoName}`);
                if (repo) repositories.push(repo);
            } catch (error) {
                console.warn(`Failed to fetch repository ${repoName}:`, error);
            }
        }
        return repositories;
    }

    async getPullRequests(repoName) {
        try {
            return await this.makeRequest(`/repos/${repoName}/pulls?state=open`);
        } catch (error) {
            console.warn(`Failed to fetch pull requests for ${repoName}:`, error);
            return [];
        }
    }

    async getUserInfo() {
        try {
            return await this.makeRequest('/user');
        } catch (error) {
            console.warn('Failed to fetch user info:', error);
            return null;
        }
    }
}

const githubAPI = new GitHubAPI();

// Main component
export const GitHubDashboardComponent = defineComponent({
    name: 'GitHubDashboard',
    setup() {
        // Reactive state
        const repositories = ref([]);
        const selectedRepo = ref(null);
        const pullRequests = ref([]);
        const loading = ref(false);
        const prLoading = ref(false);
        const error = ref('');
        const user = ref(null);
        const config = reactive({
            githubToken: '',
            repositories: [],
            refreshInterval: 300,
            showStars: true,
            showForks: true,
            showIssues: true
        });

        // Computed properties
        const hasToken = computed(() => config.githubToken.length > 0);
        const rateLimitInfo = computed(() => ({
            remaining: githubAPI.rateLimitRemaining,
            resetTime: new Date(githubAPI.rateLimitReset).toLocaleTimeString()
        }));

        // Load stored configuration
        const loadConfig = async () => {
            try {
                if (window.delvePlugin && window.delvePlugin.loadConfig) {
                    const storedConfig = await window.delvePlugin.loadConfig('dashboard_settings', '1.0.0');
                    if (storedConfig) {
                        Object.assign(config, storedConfig);
                    }
                } else {
                    // Fallback to localStorage
                    const stored = localStorage.getItem('github_dashboard_config');
                    if (stored) {
                        const parsed = JSON.parse(stored);
                        Object.assign(config, parsed);
                    }
                }
            } catch (err) {
                console.warn('Failed to load config:', err);
            }
        };

        // Save configuration
        const saveConfig = async () => {
            try {
                if (window.delvePlugin && window.delvePlugin.storeConfig) {
                    await window.delvePlugin.storeConfig('dashboard_settings', config, '1.0.0');
                } else {
                    localStorage.setItem('github_dashboard_config', JSON.stringify(config));
                }
            } catch (err) {
                console.warn('Failed to save config:', err);
            }
        };

        // Load repositories
        const loadRepositories = async () => {
            if (loading.value) return;

            try {
                loading.value = true;
                error.value = '';

                const repos = await githubAPI.getRepositories(config.repositories);
                repositories.value = repos || [];

                if (hasToken.value && !user.value) {
                    user.value = await githubAPI.getUserInfo();
                }
            } catch (err) {
                error.value = err.message;
                console.error('Failed to load repositories:', err);
            } finally {
                loading.value = false;
            }
        };

        // Load pull requests for selected repository
        const loadPullRequests = async (repo) => {
            if (prLoading.value) return;

            try {
                prLoading.value = true;
                selectedRepo.value = repo;

                const prs = await githubAPI.getPullRequests(repo.full_name);
                pullRequests.value = prs || [];
            } catch (err) {
                error.value = err.message;
                console.error('Failed to load pull requests:', err);
                pullRequests.value = [];
            } finally {
                prLoading.value = false;
            }
        };

        // Refresh all data
        const refreshAll = async () => {
            await loadRepositories();
            if (selectedRepo.value) {
                await loadPullRequests(selectedRepo.value);
            }
        };

        // Format date helper
        const formatDate = (dateString) => {
            const date = new Date(dateString);
            const now = new Date();
            const diff = now - date;
            const days = Math.floor(diff / (1000 * 60 * 60 * 24));

            if (days === 0) return 'Today';
            if (days === 1) return 'Yesterday';
            if (days < 7) return `${days} days ago`;
            return date.toLocaleDateString();
        };

        // Format number helper
        const formatNumber = (num) => {
            if (num >= 1000) {
                return (num / 1000).toFixed(1) + 'k';
            }
            return num.toString();
        };

        // Setup auto-refresh
        let refreshTimer = null;
        const setupAutoRefresh = () => {
            if (refreshTimer) {
                clearInterval(refreshTimer);
            }
            if (config.refreshInterval > 0) {
                refreshTimer = resourceManager.setManagedInterval(refreshAll, config.refreshInterval * 1000);
            }
        };

        // Watch for config changes
        watch(() => config.refreshInterval, setupAutoRefresh);
        watch(() => config.repositories, loadRepositories, { deep: true });
        watch(config, saveConfig, { deep: true });

        // Lifecycle hooks
        onMounted(async () => {
            console.log('GitHub Dashboard mounted');
            await loadConfig();
            await loadRepositories();
            setupAutoRefresh();
        });

        onUnmounted(() => {
            console.log('GitHub Dashboard unmounted');
            if (refreshTimer) {
                clearInterval(refreshTimer);
            }
        });

        // Setup cleanup
        resourceManager.onCleanup(() => {
            if (refreshTimer) {
                clearInterval(refreshTimer);
            }
        });

        return {
            repositories,
            selectedRepo,
            pullRequests,
            loading,
            prLoading,
            error,
            user,
            config,
            hasToken,
            rateLimitInfo,
            loadRepositories,
            loadPullRequests,
            refreshAll,
            formatDate,
            formatNumber,
            saveConfig
        };
    },

    template: `
        <div class="github-dashboard">
            <div class="dashboard-header">
                <div class="header-left">
                    <h1 class="dashboard-title">
                        <span class="icon">🐙</span>
                        GitHub Dashboard
                    </h1>
                    <div v-if="user" class="user-info">
                        Welcome, {{ user.login }}
                    </div>
                </div>
                <div class="header-right">
                    <div class="rate-limit" v-if="hasToken">
                        API: {{ rateLimitInfo.remaining }}/5000
                        <small>Reset: {{ rateLimitInfo.resetTime }}</small>
                    </div>
                    <button @click="refreshAll" :disabled="loading" class="refresh-btn">
                        <span class="btn-icon" :class="{ spinning: loading }">🔄</span>
                        Refresh
                    </button>
                </div>
            </div>

            <div v-if="error" class="error-message">
                <span class="error-icon">⚠️</span>
                {{ error }}
                <button @click="error = ''" class="error-close">×</button>
            </div>

            <div v-if="!hasToken" class="setup-notice">
                <div class="notice-content">
                    <h3>🔑 Setup Required</h3>
                    <p>Configure your GitHub token in settings to access real repository data.</p>
                    <p>Currently showing demo data.</p>
                </div>
            </div>

            <div class="dashboard-content">
                <div class="repositories-section">
                    <div class="section-header">
                        <h2>Repositories</h2>
                        <span class="repo-count">{{ repositories.length }}</span>
                    </div>

                    <div v-if="loading && repositories.length === 0" class="loading-state">
                        <div class="loading-spinner"></div>
                        <p>Loading repositories...</p>
                    </div>

                    <div v-else class="repositories-grid">
                        <div
                            v-for="repo in repositories"
                            :key="repo.full_name"
                            @click="loadPullRequests(repo)"
                            class="repository-card"
                            :class="{ active: selectedRepo && selectedRepo.full_name === repo.full_name }"
                        >
                            <div class="repo-header">
                                <h3 class="repo-name">{{ repo.name }}</h3>
                                <span class="repo-language" v-if="repo.language">{{ repo.language }}</span>
                            </div>

                            <p class="repo-description" v-if="repo.description">
                                {{ repo.description }}
                            </p>

                            <div class="repo-stats">
                                <span v-if="config.showStars" class="stat">
                                    <span class="stat-icon">⭐</span>
                                    {{ formatNumber(repo.stargazers_count) }}
                                </span>
                                <span v-if="config.showForks" class="stat">
                                    <span class="stat-icon">🍴</span>
                                    {{ formatNumber(repo.forks_count) }}
                                </span>
                                <span v-if="config.showIssues" class="stat">
                                    <span class="stat-icon">🐛</span>
                                    {{ repo.open_issues_count }}
                                </span>
                            </div>

                            <div class="repo-footer">
                                <span class="update-time">Updated {{ formatDate(repo.updated_at) }}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="selectedRepo" class="pull-requests-section">
                    <div class="section-header">
                        <h2>Pull Requests</h2>
                        <span class="repo-name">{{ selectedRepo.full_name }}</span>
                    </div>

                    <div v-if="prLoading" class="loading-state">
                        <div class="loading-spinner"></div>
                        <p>Loading pull requests...</p>
                    </div>

                    <div v-else-if="pullRequests.length === 0" class="empty-state">
                        <div class="empty-icon">📝</div>
                        <p>No open pull requests</p>
                    </div>

                    <div v-else class="pull-requests-list">
                        <div
                            v-for="pr in pullRequests"
                            :key="pr.number"
                            class="pull-request-card"
                        >
                            <div class="pr-header">
                                <div class="pr-number">#{{ pr.number }}</div>
                                <div class="pr-state" :class="pr.state">{{ pr.state }}</div>
                            </div>

                            <h3 class="pr-title">
                                <a :href="pr.html_url" target="_blank" rel="noopener">
                                    {{ pr.title }}
                                </a>
                            </h3>

                            <div class="pr-meta">
                                <div class="pr-author">
                                    <img
                                        :src="pr.user.avatar_url"
                                        :alt="pr.user.login"
                                        class="author-avatar"
                                    />
                                    <span>{{ pr.user.login }}</span>
                                </div>
                                <div class="pr-dates">
                                    <span>Created {{ formatDate(pr.created_at) }}</span>
                                    <span>Updated {{ formatDate(pr.updated_at) }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `
});

// Create and register the custom element
const GitHubDashboardElement = Vue.defineCustomElement(GitHubDashboardComponent);
customElements.define('github-dashboard', GitHubDashboardElement);

// Setup cleanup on page unload
window.addEventListener('beforeunload', () => {
    resourceManager.cleanup();
});

console.log('GitHub Dashboard web component registered successfully');
