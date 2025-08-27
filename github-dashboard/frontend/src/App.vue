<template>
    <div class="github-dashboard">
        <!-- Header Section -->
        <header class="dashboard-header">
            <div class="header-content">
                <div class="header-title">
                    <div class="title-icon">🐙</div>
                    <div class="title-text">
                        <h1>GitHub Dashboard</h1>
                        <p class="subtitle">
                            Monitor your repositories and pull requests
                        </p>
                    </div>
                </div>
                <div class="header-actions">
                    <button
                        @click="refresh"
                        :disabled="loading"
                        class="btn btn-primary"
                        :class="{ 'btn-loading': loading }"
                    >
                        <span class="btn-icon" :class="{ rotating: loading }">
                            {{ loading ? "⟳" : "↻" }}
                        </span>
                        {{ loading ? "Refreshing..." : "Refresh" }}
                    </button>
                </div>
            </div>
        </header>

        <!-- Error Alert -->
        <div v-if="error" class="alert alert-error" role="alert">
            <div class="alert-icon">⚠️</div>
            <div class="alert-content">
                <div class="alert-title">Something went wrong</div>
                <div class="alert-message">{{ error }}</div>
            </div>
            <button
                @click="clearError"
                class="alert-dismiss"
                aria-label="Dismiss error"
            >
                ×
            </button>
        </div>

        <!-- Main Content -->
        <div class="dashboard-content">
            <!-- Repositories Section -->
            <section class="dashboard-section repositories-section">
                <div class="section-header">
                    <h2 class="section-title">
                        <span class="section-icon">📁</span>
                        Repositories
                        <span
                            v-if="repositories.length > 0"
                            class="section-count"
                            >{{ repositories.length }}</span
                        >
                    </h2>
                </div>

                <div class="section-content">
                    <!-- Loading State -->
                    <div
                        v-if="loading && repositories.length === 0"
                        class="loading-state"
                    >
                        <div class="skeleton-cards">
                            <div v-for="i in 3" :key="i" class="skeleton-card">
                                <div class="skeleton-header">
                                    <div
                                        class="skeleton-line skeleton-title"
                                    ></div>
                                    <div
                                        class="skeleton-line skeleton-stats"
                                    ></div>
                                </div>
                                <div
                                    class="skeleton-line skeleton-description"
                                ></div>
                                <div class="skeleton-footer">
                                    <div
                                        class="skeleton-line skeleton-meta"
                                    ></div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Repository List -->
                    <div v-else-if="repositories.length > 0" class="repo-grid">
                        <RepositoryCard
                            v-for="repo in repositories"
                            :key="repo.full_name"
                            :repository="repo"
                            :selected="selectedRepo === repo.full_name"
                            @select="selectRepository"
                        />
                    </div>

                    <!-- Empty State -->
                    <div v-else class="empty-state">
                        <div class="empty-icon">📂</div>
                        <h3 class="empty-title">No repositories found</h3>
                        <p class="empty-description">
                            Configure your GitHub token and repositories to get
                            started.
                        </p>
                        <button class="btn btn-outline" @click="refresh">
                            Try again
                        </button>
                    </div>
                </div>
            </section>

            <!-- Pull Requests Section -->
            <section
                v-if="selectedRepo"
                class="dashboard-section pull-requests-section"
            >
                <div class="section-header">
                    <h2 class="section-title">
                        <span class="section-icon">🔄</span>
                        Pull Requests
                        <span class="section-subtitle">{{ selectedRepo }}</span>
                        <span
                            v-if="pullRequests.length > 0"
                            class="section-count"
                            >{{ pullRequests.length }}</span
                        >
                    </h2>
                </div>

                <div class="section-content">
                    <!-- Loading State -->
                    <div v-if="pullRequestsLoading" class="loading-state">
                        <div class="skeleton-cards">
                            <div
                                v-for="i in 2"
                                :key="i"
                                class="skeleton-card skeleton-pr"
                            >
                                <div class="skeleton-header">
                                    <div
                                        class="skeleton-line skeleton-pr-title"
                                    ></div>
                                    <div class="skeleton-badge"></div>
                                </div>
                                <div class="skeleton-meta">
                                    <div class="skeleton-avatar"></div>
                                    <div
                                        class="skeleton-line skeleton-author"
                                    ></div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Pull Request List -->
                    <div v-else-if="pullRequests.length > 0" class="pr-list">
                        <PullRequestCard
                            v-for="pr in pullRequests"
                            :key="pr.number"
                            :pull-request="pr"
                        />
                    </div>

                    <!-- Empty State -->
                    <div v-else class="empty-state empty-state-compact">
                        <div class="empty-icon">✨</div>
                        <h3 class="empty-title">No open pull requests</h3>
                        <p class="empty-description">
                            This repository has no open pull requests.
                        </p>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from "vue";
import RepositoryCard from "./components/RepositoryCard.vue";
import PullRequestCard from "./components/PullRequestCard.vue";

export default {
    name: "GitHubDashboard",
    components: {
        RepositoryCard,
        PullRequestCard,
    },
    props: {
        githubToken: {
            type: String,
            default: "",
        },
        repositories: {
            type: Array,
            default: () => [],
        },
        refreshInterval: {
            type: Number,
            default: 300000, // 5 minutes
        },
        compactView: {
            type: Boolean,
            default: false,
        },
    },
    setup(props) {
        const repositories = ref([]);
        const pullRequests = ref([]);
        const selectedRepo = ref(null);
        const loading = ref(false);
        const pullRequestsLoading = ref(false);
        const error = ref(null);
        let refreshTimer = null;

        const loadRepositories = async () => {
            try {
                loading.value = true;
                error.value = null;

                if (window.pluginAPI && window.pluginAPI.getRepositories) {
                    const data = await window.pluginAPI.getRepositories();
                    repositories.value = Array.isArray(data) ? data : [];
                } else {
                    throw new Error("Plugin API not available");
                }
            } catch (err) {
                error.value = err.message || "Failed to load repositories";
                console.error("Failed to load repositories:", err);
                repositories.value = [];
            } finally {
                loading.value = false;
            }
        };

        const selectRepository = async (repo) => {
            if (selectedRepo.value === repo.full_name) return;

            selectedRepo.value = repo.full_name;
            pullRequests.value = [];

            try {
                pullRequestsLoading.value = true;
                error.value = null;

                if (window.pluginAPI && window.pluginAPI.getPullRequests) {
                    const data = await window.pluginAPI.getPullRequests(
                        repo.full_name,
                    );
                    pullRequests.value = Array.isArray(data) ? data : [];
                } else {
                    throw new Error("Plugin API not available");
                }
            } catch (err) {
                error.value = `Failed to load pull requests for ${repo.full_name}: ${err.message}`;
                console.error("Failed to load pull requests:", err);
                pullRequests.value = [];
            } finally {
                pullRequestsLoading.value = false;
            }
        };

        const refresh = async () => {
            await loadRepositories();
            if (selectedRepo.value) {
                const repo = repositories.value.find(
                    (r) => r.full_name === selectedRepo.value,
                );
                if (repo) {
                    await selectRepository(repo);
                } else {
                    // Selected repo no longer exists
                    selectedRepo.value = null;
                    pullRequests.value = [];
                }
            }
        };

        const clearError = () => {
            error.value = null;
        };

        const startAutoRefresh = () => {
            if (props.refreshInterval > 0) {
                refreshTimer = setInterval(refresh, props.refreshInterval);
            }
        };

        const stopAutoRefresh = () => {
            if (refreshTimer) {
                clearInterval(refreshTimer);
                refreshTimer = null;
            }
        };

        onMounted(() => {
            loadRepositories();
            startAutoRefresh();
        });

        onUnmounted(() => {
            stopAutoRefresh();
        });

        return {
            repositories,
            pullRequests,
            selectedRepo,
            loading,
            pullRequestsLoading,
            error,
            selectRepository,
            refresh,
            clearError,
        };
    },
};
</script>

<style scoped>
/* Design Tokens - GitHub Primer inspired */
:root {
    --color-canvas-default: #ffffff;
    --color-canvas-subtle: #f6f8fa;
    --color-canvas-inset: #f6f8fa;
    --color-border-default: #d1d9e0;
    --color-border-muted: #d8dee4;
    --color-fg-default: #1f2328;
    --color-fg-muted: #636c76;
    --color-fg-subtle: #8c959f;
    --color-accent-fg: #0969da;
    --color-accent-emphasis: #0969da;
    --color-danger-fg: #d1242f;
    --color-success-fg: #1a7f37;
    --color-attention-fg: #9a6700;

    --space-1: 4px;
    --space-2: 8px;
    --space-3: 16px;
    --space-4: 24px;
    --space-5: 32px;
    --space-6: 40px;

    --font-size-small: 12px;
    --font-size-normal: 14px;
    --font-size-medium: 16px;
    --font-size-large: 20px;
    --font-size-xl: 24px;

    --border-radius: 6px;
    --border-radius-medium: 8px;

    --shadow-small: 0 1px 0 rgba(31, 35, 40, 0.04);
    --shadow-medium: 0 3px 6px rgba(31, 35, 40, 0.15);
    --shadow-large: 0 8px 24px rgba(31, 35, 40, 0.2);

    --transition-duration: 0.2s;
}

/* Base Styles */
.github-dashboard {
    font-family:
        -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Helvetica,
        Arial, sans-serif;
    font-size: var(--font-size-normal);
    line-height: 1.5;
    color: var(--color-fg-default);
    background: var(--color-canvas-default);
    min-height: 400px;
}

/* Header */
.dashboard-header {
    background: var(--color-canvas-default);
    border-bottom: 1px solid var(--color-border-default);
    padding: var(--space-4) var(--space-4) var(--space-3);
    margin-bottom: var(--space-4);
}

.header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    max-width: 1200px;
    margin: 0 auto;
}

.header-title {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

.title-icon {
    font-size: var(--font-size-xl);
}

.title-text h1 {
    margin: 0;
    font-size: var(--font-size-large);
    font-weight: 600;
    color: var(--color-fg-default);
}

.subtitle {
    margin: var(--space-1) 0 0;
    font-size: var(--font-size-small);
    color: var(--color-fg-muted);
}

/* Buttons */
.btn {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-3);
    font-size: var(--font-size-normal);
    font-weight: 500;
    line-height: 1.25;
    border: 1px solid;
    border-radius: var(--border-radius);
    cursor: pointer;
    transition: all var(--transition-duration) ease;
    text-decoration: none;
    white-space: nowrap;
}

.btn:focus {
    outline: 2px solid var(--color-accent-fg);
    outline-offset: -2px;
}

.btn-primary {
    color: #ffffff;
    background: var(--color-accent-emphasis);
    border-color: var(--color-accent-emphasis);
}

.btn-primary:hover:not(:disabled) {
    background: #0860ca;
    border-color: #0860ca;
}

.btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-outline {
    color: var(--color-accent-fg);
    background: var(--color-canvas-default);
    border-color: var(--color-border-default);
}

.btn-outline:hover:not(:disabled) {
    background: var(--color-canvas-subtle);
    border-color: var(--color-border-muted);
}

.btn-icon {
    font-size: var(--font-size-medium);
    transition: transform 0.6s ease;
}

.btn-icon.rotating {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

/* Alert */
.alert {
    display: flex;
    align-items: flex-start;
    gap: var(--space-3);
    padding: var(--space-3);
    margin: 0 var(--space-4) var(--space-4);
    border: 1px solid;
    border-radius: var(--border-radius);
    position: relative;
}

.alert-error {
    background: #fff8f8;
    border-color: #ffcdd2;
    color: var(--color-danger-fg);
}

.alert-icon {
    font-size: var(--font-size-medium);
    flex-shrink: 0;
}

.alert-content {
    flex: 1;
}

.alert-title {
    font-weight: 600;
    margin-bottom: var(--space-1);
}

.alert-message {
    font-size: var(--font-size-small);
    color: var(--color-fg-muted);
}

.alert-dismiss {
    background: none;
    border: none;
    font-size: var(--font-size-large);
    color: var(--color-fg-muted);
    cursor: pointer;
    padding: 0;
    line-height: 1;
}

.alert-dismiss:hover {
    color: var(--color-fg-default);
}

/* Dashboard Content */
.dashboard-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 var(--space-4);
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--space-5);
    align-items: start;
}

@media (max-width: 1024px) {
    .dashboard-content {
        grid-template-columns: 1fr;
        gap: var(--space-4);
    }
}

/* Sections */
.dashboard-section {
    background: var(--color-canvas-default);
}

.section-header {
    margin-bottom: var(--space-4);
}

.section-title {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    margin: 0;
    font-size: var(--font-size-medium);
    font-weight: 600;
    color: var(--color-fg-default);
}

.section-icon {
    font-size: var(--font-size-medium);
}

.section-subtitle {
    font-size: var(--font-size-small);
    font-weight: 400;
    color: var(--color-fg-muted);
    margin-left: var(--space-2);
}

.section-count {
    background: var(--color-canvas-subtle);
    color: var(--color-fg-muted);
    padding: var(--space-1) var(--space-2);
    border-radius: 12px;
    font-size: var(--font-size-small);
    font-weight: 500;
    margin-left: var(--space-2);
}

/* Repository Grid */
.repo-grid {
    display: grid;
    gap: var(--space-3);
}

/* Pull Request List */
.pr-list {
    display: grid;
    gap: var(--space-3);
}

/* Loading States */
.loading-state {
    padding: var(--space-2) 0;
}

.skeleton-cards {
    display: grid;
    gap: var(--space-3);
}

.skeleton-card {
    background: var(--color-canvas-default);
    border: 1px solid var(--color-border-default);
    border-radius: var(--border-radius);
    padding: var(--space-3);
}

.skeleton-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--space-3);
}

.skeleton-footer {
    margin-top: var(--space-3);
}

.skeleton-meta {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    margin-top: var(--space-3);
}

.skeleton-line {
    background: var(--color-canvas-subtle);
    border-radius: 3px;
    animation: skeleton-loading 1.2s ease-in-out infinite;
}

.skeleton-title {
    height: 16px;
    width: 60%;
}

.skeleton-stats {
    height: 12px;
    width: 80px;
}

.skeleton-description {
    height: 14px;
    width: 85%;
}

.skeleton-meta {
    height: 12px;
    width: 50%;
}

.skeleton-pr-title {
    height: 16px;
    width: 70%;
}

.skeleton-author {
    height: 12px;
    width: 60px;
}

.skeleton-badge {
    height: 20px;
    width: 50px;
    border-radius: 10px;
}

.skeleton-avatar {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    flex-shrink: 0;
}

@keyframes skeleton-loading {
    0% {
        opacity: 1;
    }
    50% {
        opacity: 0.4;
    }
    100% {
        opacity: 1;
    }
}

/* Empty States */
.empty-state {
    text-align: center;
    padding: var(--space-6) var(--space-4);
}

.empty-state-compact {
    padding: var(--space-5) var(--space-4);
}

.empty-icon {
    font-size: 48px;
    margin-bottom: var(--space-4);
    opacity: 0.6;
}

.empty-title {
    margin: 0 0 var(--space-2) 0;
    font-size: var(--font-size-medium);
    font-weight: 600;
    color: var(--color-fg-default);
}

.empty-description {
    margin: 0 0 var(--space-4) 0;
    font-size: var(--font-size-normal);
    color: var(--color-fg-muted);
    line-height: 1.4;
}

/* Responsive Design */
@media (max-width: 768px) {
    .header-content {
        flex-direction: column;
        align-items: stretch;
        gap: var(--space-3);
    }

    .dashboard-content {
        padding: 0 var(--space-3);
    }

    .section-title {
        flex-wrap: wrap;
    }

    .section-subtitle {
        margin-left: 0;
        margin-top: var(--space-1);
    }
}

@media (max-width: 480px) {
    .dashboard-header {
        padding: var(--space-3);
    }

    .dashboard-content {
        padding: 0 var(--space-2);
        gap: var(--space-3);
    }
}
</style>
