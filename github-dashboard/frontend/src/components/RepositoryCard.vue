<template>
    <article
        class="repo-card"
        :class="{ 'repo-card--selected': selected }"
        @click="$emit('select', repository)"
        role="button"
        tabindex="0"
        @keydown.enter="$emit('select', repository)"
        @keydown.space.prevent="$emit('select', repository)"
    >
        <header class="repo-card__header">
            <div class="repo-card__title-section">
                <h3 class="repo-card__name">
                    <span class="repo-card__icon">📁</span>
                    {{ repository.name }}
                </h3>
                <div class="repo-card__stats">
                    <div class="repo-stat" v-if="repository.stargazers_count !== undefined">
                        <span class="repo-stat__icon">⭐</span>
                        <span class="repo-stat__value">{{ formatNumber(repository.stargazers_count) }}</span>
                    </div>
                    <div class="repo-stat" v-if="repository.forks_count !== undefined">
                        <span class="repo-stat__icon">🍴</span>
                        <span class="repo-stat__value">{{ formatNumber(repository.forks_count) }}</span>
                    </div>
                    <div class="repo-stat" v-if="repository.open_issues_count !== undefined">
                        <span class="repo-stat__icon">🐛</span>
                        <span class="repo-stat__value">{{ formatNumber(repository.open_issues_count) }}</span>
                    </div>
                </div>
            </div>
        </header>

        <div class="repo-card__content">
            <p class="repo-card__description">
                {{ repository.description || 'No description available' }}
            </p>
        </div>

        <footer class="repo-card__footer">
            <div class="repo-card__meta">
                <div v-if="repository.language" class="repo-language">
                    <span class="language-dot" :style="{ backgroundColor: getLanguageColor(repository.language) }"></span>
                    <span class="language-name">{{ repository.language }}</span>
                </div>
                <div v-if="repository.updated_at" class="repo-updated">
                    Updated {{ formatDate(repository.updated_at) }}
                </div>
            </div>
            <div class="repo-card__actions">
                <button
                    class="repo-card__link-btn"
                    @click.stop="openRepository"
                    :aria-label="`Open ${repository.name} on GitHub`"
                    title="View on GitHub"
                >
                    <span class="link-icon">↗</span>
                </button>
            </div>
        </footer>

        <div v-if="selected" class="repo-card__selected-indicator" aria-hidden="true"></div>
    </article>
</template>

<script>
export default {
    name: 'RepositoryCard',
    props: {
        repository: {
            type: Object,
            required: true,
        },
        selected: {
            type: Boolean,
            default: false,
        },
    },
    emits: ['select'],
    methods: {
        formatNumber(num) {
            if (num === undefined || num === null) return '0';
            if (num >= 1000000) {
                return (num / 1000000).toFixed(1) + 'M';
            }
            if (num >= 1000) {
                return (num / 1000).toFixed(1) + 'k';
            }
            return num.toString();
        },

        formatDate(dateString) {
            if (!dateString) return 'Unknown';

            try {
                const date = new Date(dateString);
                const now = new Date();
                const diffInDays = Math.floor((now - date) / (1000 * 60 * 60 * 24));

                if (diffInDays === 0) return 'today';
                if (diffInDays === 1) return 'yesterday';
                if (diffInDays < 7) return `${diffInDays} days ago`;
                if (diffInDays < 30) return `${Math.floor(diffInDays / 7)} weeks ago`;
                if (diffInDays < 365) return `${Math.floor(diffInDays / 30)} months ago`;
                return `${Math.floor(diffInDays / 365)} years ago`;
            } catch (error) {
                return 'Unknown';
            }
        },

        getLanguageColor(language) {
            // GitHub language colors
            const languageColors = {
                'JavaScript': '#f1e05a',
                'TypeScript': '#2b7489',
                'Python': '#3572A5',
                'Java': '#b07219',
                'Go': '#00ADD8',
                'Rust': '#dea584',
                'C++': '#f34b7d',
                'C': '#555555',
                'C#': '#239120',
                'PHP': '#4F5D95',
                'Ruby': '#701516',
                'Swift': '#ffac45',
                'Kotlin': '#F18E33',
                'Scala': '#c22d40',
                'Dart': '#00B4AB',
                'R': '#198CE7',
                'Shell': '#89e051',
                'PowerShell': '#012456',
                'HTML': '#e34c26',
                'CSS': '#563d7c',
                'Vue': '#2c3e50',
                'React': '#61dafb',
                'Angular': '#dd0031',
            };
            return languageColors[language] || '#586069';
        },

        openRepository(event) {
            event.preventDefault();
            event.stopPropagation();

            const url = this.repository.html_url || `https://github.com/${this.repository.full_name}`;
            window.open(url, '_blank', 'noopener,noreferrer');
        },
    },
};
</script>

<style scoped>
/* Repository Card Component */
.repo-card {
    background: var(--color-canvas-default, #ffffff);
    border: 1px solid var(--color-border-default, #d1d9e0);
    border-radius: var(--border-radius, 6px);
    padding: var(--space-3, 16px);
    cursor: pointer;
    transition: all var(--transition-duration, 0.2s) ease;
    position: relative;
    display: flex;
    flex-direction: column;
    gap: var(--space-3, 16px);
    min-height: 140px;
}

.repo-card:hover {
    border-color: var(--color-accent-emphasis, #0969da);
    box-shadow: var(--shadow-medium, 0 3px 6px rgba(31, 35, 40, 0.15));
    transform: translateY(-1px);
}

.repo-card:focus {
    outline: 2px solid var(--color-accent-emphasis, #0969da);
    outline-offset: -2px;
}

.repo-card--selected {
    border-color: var(--color-accent-emphasis, #0969da);
    background: var(--color-canvas-subtle, #f6f8fa);
    box-shadow: var(--shadow-small, 0 1px 0 rgba(31, 35, 40, 0.04));
}

.repo-card__selected-indicator {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: var(--color-accent-emphasis, #0969da);
    border-radius: var(--border-radius, 6px) 0 0 var(--border-radius, 6px);
}

/* Header */
.repo-card__header {
    display: flex;
    flex-direction: column;
    gap: var(--space-2, 8px);
}

.repo-card__title-section {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--space-2, 8px);
}

.repo-card__name {
    margin: 0;
    font-size: var(--font-size-medium, 16px);
    font-weight: 600;
    color: var(--color-accent-fg, #0969da);
    line-height: 1.25;
    display: flex;
    align-items: center;
    gap: var(--space-2, 8px);
    flex: 1;
    min-width: 0; /* Allow text to truncate */
}

.repo-card__icon {
    font-size: var(--font-size-normal, 14px);
    flex-shrink: 0;
}

.repo-card__stats {
    display: flex;
    gap: var(--space-3, 16px);
    align-items: center;
    flex-shrink: 0;
}

.repo-stat {
    display: flex;
    align-items: center;
    gap: var(--space-1, 4px);
    font-size: var(--font-size-small, 12px);
    color: var(--color-fg-muted, #636c76);
}

.repo-stat__icon {
    font-size: var(--font-size-small, 12px);
}

.repo-stat__value {
    font-weight: 500;
    color: var(--color-fg-default, #1f2328);
}

/* Content */
.repo-card__content {
    flex: 1;
}

.repo-card__description {
    margin: 0;
    font-size: var(--font-size-normal, 14px);
    line-height: 1.4;
    color: var(--color-fg-muted, #636c76);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

/* Footer */
.repo-card__footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-2, 8px);
    margin-top: auto;
}

.repo-card__meta {
    display: flex;
    flex-direction: column;
    gap: var(--space-1, 4px);
    flex: 1;
    min-width: 0;
}

.repo-language {
    display: flex;
    align-items: center;
    gap: var(--space-2, 8px);
    font-size: var(--font-size-small, 12px);
    color: var(--color-fg-default, #1f2328);
}

.language-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    flex-shrink: 0;
}

.language-name {
    font-weight: 500;
}

.repo-updated {
    font-size: var(--font-size-small, 12px);
    color: var(--color-fg-subtle, #8c959f);
}

.repo-card__actions {
    display: flex;
    align-items: center;
}

.repo-card__link-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: none;
    border: 1px solid transparent;
    border-radius: var(--border-radius, 6px);
    cursor: pointer;
    color: var(--color-fg-muted, #636c76);
    transition: all var(--transition-duration, 0.2s) ease;
}

.repo-card__link-btn:hover {
    background: var(--color-canvas-subtle, #f6f8fa);
    border-color: var(--color-border-default, #d1d9e0);
    color: var(--color-fg-default, #1f2328);
}

.repo-card__link-btn:focus {
    outline: 2px solid var(--color-accent-emphasis, #0969da);
    outline-offset: -2px;
}

.link-icon {
    font-size: var(--font-size-normal, 14px);
    font-weight: bold;
}

/* Responsive Design */
@media (max-width: 768px) {
    .repo-card__title-section {
        flex-direction: column;
        align-items: stretch;
        gap: var(--space-2, 8px);
    }

    .repo-card__stats {
        align-self: flex-start;
        gap: var(--space-2, 8px);
    }

    .repo-card__footer {
        flex-direction: column;
        align-items: stretch;
        gap: var(--space-2, 8px);
    }

    .repo-card__actions {
        align-self: flex-end;
    }
}

@media (max-width: 480px) {
    .repo-card {
        padding: var(--space-2, 8px);
        gap: var(--space-2, 8px);
    }

    .repo-card__name {
        font-size: var(--font-size-normal, 14px);
    }

    .repo-card__stats {
        gap: var(--space-2, 8px);
    }

    .repo-stat {
        font-size: 11px;
    }
}

/* Accessibility */
@media (prefers-reduced-motion: reduce) {
    .repo-card {
        transition: none;
    }

    .repo-card:hover {
        transform: none;
    }
}

/* High Contrast Mode */
@media (prefers-contrast: high) {
    .repo-card {
        border-width: 2px;
    }

    .repo-card:hover,
    .repo-card--selected {
        border-width: 3px;
    }
}
</style>
