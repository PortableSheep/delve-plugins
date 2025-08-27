<template>
    <article
        class="pr-card"
        :class="`pr-card--${pullRequest.state}`"
        role="button"
        tabindex="0"
        @click="openPullRequest"
        @keydown.enter="openPullRequest"
        @keydown.space.prevent="openPullRequest"
    >
        <header class="pr-card__header">
            <div class="pr-card__title-section">
                <h3 class="pr-card__title">
                    <span class="pr-card__number"
                        >#{{ pullRequest.number }}</span
                    >
                    <span class="pr-card__title-text">{{
                        pullRequest.title
                    }}</span>
                </h3>
                <div class="pr-card__state">
                    <span
                        class="pr-state-badge"
                        :class="`pr-state-badge--${pullRequest.state}`"
                        :aria-label="`Pull request is ${pullRequest.state}`"
                    >
                        <span class="pr-state-badge__icon">{{
                            getStateIcon(pullRequest.state)
                        }}</span>
                        <span class="pr-state-badge__text">{{
                            getStateText(pullRequest.state)
                        }}</span>
                    </span>
                </div>
            </div>
        </header>

        <div class="pr-card__meta">
            <div class="pr-card__author">
                <img
                    :src="pullRequest.user?.avatar_url || getDefaultAvatar()"
                    :alt="`${pullRequest.user?.login || 'Unknown user'} avatar`"
                    class="pr-card__avatar"
                    loading="lazy"
                    @error="handleAvatarError"
                />
                <div class="pr-card__author-info">
                    <span class="pr-card__author-name">
                        {{ pullRequest.user?.login || "Unknown user" }}
                    </span>
                    <div class="pr-card__timestamps">
                        <span class="pr-card__created">
                            opened {{ formatDate(pullRequest.created_at) }}
                        </span>
                        <span class="pr-card__separator">•</span>
                        <span class="pr-card__updated">
                            updated {{ formatDate(pullRequest.updated_at) }}
                        </span>
                    </div>
                </div>
            </div>

            <div class="pr-card__actions">
                <button
                    class="pr-card__link-btn"
                    @click.stop="openPullRequest"
                    :aria-label="`Open pull request #${pullRequest.number} on GitHub`"
                    title="View on GitHub"
                >
                    <span class="link-icon">↗</span>
                </button>
            </div>
        </div>

        <div class="pr-card__footer" v-if="showAdditionalInfo">
            <div class="pr-card__additional-info">
                <div v-if="pullRequest.draft" class="pr-card__draft-indicator">
                    <span class="draft-icon">📝</span>
                    <span class="draft-text">Draft</span>
                </div>
                <div
                    v-if="pullRequest.mergeable_state"
                    class="pr-card__mergeable-state"
                >
                    <span class="mergeable-icon">{{
                        getMergeableIcon(pullRequest.mergeable_state)
                    }}</span>
                    <span class="mergeable-text">{{
                        getMergeableText(pullRequest.mergeable_state)
                    }}</span>
                </div>
            </div>
        </div>
    </article>
</template>

<script>
export default {
    name: "PullRequestCard",
    props: {
        pullRequest: {
            type: Object,
            required: true,
        },
        showAdditionalInfo: {
            type: Boolean,
            default: false,
        },
    },
    emits: ["select"],
    methods: {
        formatDate(dateString) {
            if (!dateString) return "unknown time";

            try {
                const date = new Date(dateString);
                const now = new Date();
                const diffInSeconds = Math.floor((now - date) / 1000);

                // Less than a minute
                if (diffInSeconds < 60) {
                    return "just now";
                }

                // Less than an hour
                if (diffInSeconds < 3600) {
                    const minutes = Math.floor(diffInSeconds / 60);
                    return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
                }

                // Less than a day
                if (diffInSeconds < 86400) {
                    const hours = Math.floor(diffInSeconds / 3600);
                    return `${hours} hour${hours > 1 ? "s" : ""} ago`;
                }

                // Less than a week
                if (diffInSeconds < 604800) {
                    const days = Math.floor(diffInSeconds / 86400);
                    return `${days} day${days > 1 ? "s" : ""} ago`;
                }

                // Less than a month
                if (diffInSeconds < 2592000) {
                    const weeks = Math.floor(diffInSeconds / 604800);
                    return `${weeks} week${weeks > 1 ? "s" : ""} ago`;
                }

                // More than a month, show date
                return date.toLocaleDateString("en-US", {
                    month: "short",
                    day: "numeric",
                    year:
                        date.getFullYear() !== now.getFullYear()
                            ? "numeric"
                            : undefined,
                });
            } catch (error) {
                return "unknown time";
            }
        },

        getStateIcon(state) {
            const stateIcons = {
                open: "🟢",
                closed: "🔴",
                merged: "🟣",
            };
            return stateIcons[state] || "❓";
        },

        getStateText(state) {
            const stateTexts = {
                open: "Open",
                closed: "Closed",
                merged: "Merged",
            };
            return stateTexts[state] || "Unknown";
        },

        getMergeableIcon(state) {
            const mergeableIcons = {
                clean: "✅",
                dirty: "⚠️",
                unstable: "🟡",
                blocked: "🚫",
            };
            return mergeableIcons[state] || "❓";
        },

        getMergeableText(state) {
            const mergeableTexts = {
                clean: "Ready to merge",
                dirty: "Conflicts",
                unstable: "Checks pending",
                blocked: "Blocked",
            };
            return mergeableTexts[state] || "Unknown status";
        },

        getDefaultAvatar() {
            return "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHZpZXdCb3g9IjAgMCAyMCAyMCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMTAiIGN5PSIxMCIgcj0iMTAiIGZpbGw9IiNkMWQ5ZTAiLz4KPHN2ZyB3aWR0aD0iMTIiIGhlaWdodD0iMTIiIHZpZXdCb3g9IjAgMCAxMiAxMiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4PSI0IiB5PSI0Ij4KPHBhdGggZD0iTTYgNkM3LjY1Njg1IDYgOSA0LjY1Njg1IDkgM0M5IDEuMzQzMTUgNy42NTY4NSAwIDYgMEM0LjM0MzE1IDAgMyAxLjM0MzE1IDMgM0MzIDQuNjU2ODUgNC4zNDMxNSA2IDYgNloiIGZpbGw9IiM2MzZjNzYiLz4KPHBhdGggZD0iTTEyIDEwLjVDMTIgOC4wMTQ3MiAxMC4yMDkxIDYgOCA2SDRDMS43OTA4NiA2IDAgOC4wMTQ3MiAwIDEwLjVWMTJIMTJWMTAuNVoiIGZpbGw9IiM2MzZjNzYiLz4KPC9zdmc+Cjwvc3ZnPgo=";
        },

        handleAvatarError(event) {
            event.target.src = this.getDefaultAvatar();
        },

        openPullRequest(event) {
            event.preventDefault();
            const url = this.pullRequest.html_url;
            if (url) {
                window.open(url, "_blank", "noopener,noreferrer");
            }
            this.$emit("select", this.pullRequest);
        },
    },
};
</script>

<style scoped>
/* Pull Request Card Component */
.pr-card {
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
}

.pr-card:hover {
    border-color: var(--color-border-muted, #d8dee4);
    box-shadow: var(--shadow-medium, 0 3px 6px rgba(31, 35, 40, 0.15));
    transform: translateY(-1px);
}

.pr-card:focus {
    outline: 2px solid var(--color-accent-emphasis, #0969da);
    outline-offset: -2px;
}

/* State-specific styling */
.pr-card--open {
    border-left: 3px solid #1a7f37;
}

.pr-card--closed {
    border-left: 3px solid #d1242f;
}

.pr-card--merged {
    border-left: 3px solid #8250df;
}

/* Header */
.pr-card__header {
    display: flex;
    flex-direction: column;
    gap: var(--space-2, 8px);
}

.pr-card__title-section {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--space-3, 16px);
}

.pr-card__title {
    margin: 0;
    font-size: var(--font-size-normal, 14px);
    font-weight: 600;
    line-height: 1.4;
    color: var(--color-fg-default, #1f2328);
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: baseline;
    gap: var(--space-2, 8px);
}

.pr-card__number {
    color: var(--color-fg-muted, #636c76);
    font-weight: 400;
    flex-shrink: 0;
}

.pr-card__title-text {
    color: var(--color-accent-fg, #0969da);
    text-decoration: none;
    word-break: break-word;
}

.pr-card__title:hover .pr-card__title-text {
    text-decoration: underline;
}

.pr-card__state {
    flex-shrink: 0;
}

/* State Badge */
.pr-state-badge {
    display: inline-flex;
    align-items: center;
    gap: var(--space-1, 4px);
    padding: var(--space-1, 4px) var(--space-2, 8px);
    border-radius: 12px;
    font-size: var(--font-size-small, 12px);
    font-weight: 500;
    text-transform: capitalize;
    white-space: nowrap;
}

.pr-state-badge--open {
    background: #dafbe1;
    color: #1a7f37;
}

.pr-state-badge--closed {
    background: #ffebe9;
    color: #d1242f;
}

.pr-state-badge--merged {
    background: #fbefff;
    color: #8250df;
}

.pr-state-badge__icon {
    font-size: 10px;
}

/* Meta Section */
.pr-card__meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3, 16px);
}

.pr-card__author {
    display: flex;
    align-items: center;
    gap: var(--space-2, 8px);
    flex: 1;
    min-width: 0;
}

.pr-card__avatar {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    flex-shrink: 0;
    background: var(--color-canvas-subtle, #f6f8fa);
}

.pr-card__author-info {
    display: flex;
    flex-direction: column;
    gap: var(--space-1, 4px);
    min-width: 0;
    flex: 1;
}

.pr-card__author-name {
    font-size: var(--font-size-small, 12px);
    font-weight: 500;
    color: var(--color-fg-default, #1f2328);
}

.pr-card__timestamps {
    display: flex;
    align-items: center;
    gap: var(--space-1, 4px);
    font-size: var(--font-size-small, 12px);
    color: var(--color-fg-muted, #636c76);
}

.pr-card__separator {
    color: var(--color-fg-subtle, #8c959f);
}

.pr-card__created,
.pr-card__updated {
    white-space: nowrap;
}

.pr-card__actions {
    display: flex;
    align-items: center;
}

.pr-card__link-btn {
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

.pr-card__link-btn:hover {
    background: var(--color-canvas-subtle, #f6f8fa);
    border-color: var(--color-border-default, #d1d9e0);
    color: var(--color-fg-default, #1f2328);
}

.pr-card__link-btn:focus {
    outline: 2px solid var(--color-accent-emphasis, #0969da);
    outline-offset: -2px;
}

.link-icon {
    font-size: var(--font-size-normal, 14px);
    font-weight: bold;
}

/* Footer */
.pr-card__footer {
    padding-top: var(--space-2, 8px);
    border-top: 1px solid var(--color-border-muted, #d8dee4);
}

.pr-card__additional-info {
    display: flex;
    gap: var(--space-3, 16px);
    align-items: center;
}

.pr-card__draft-indicator,
.pr-card__mergeable-state {
    display: flex;
    align-items: center;
    gap: var(--space-1, 4px);
    font-size: var(--font-size-small, 12px);
    color: var(--color-fg-muted, #636c76);
}

.draft-icon,
.mergeable-icon {
    font-size: 10px;
}

/* Responsive Design */
@media (max-width: 768px) {
    .pr-card__title-section {
        flex-direction: column;
        align-items: stretch;
        gap: var(--space-2, 8px);
    }

    .pr-card__state {
        align-self: flex-start;
    }

    .pr-card__meta {
        flex-direction: column;
        align-items: stretch;
        gap: var(--space-2, 8px);
    }

    .pr-card__actions {
        align-self: flex-end;
    }

    .pr-card__timestamps {
        flex-direction: column;
        align-items: flex-start;
        gap: var(--space-1, 4px);
    }

    .pr-card__separator {
        display: none;
    }
}

@media (max-width: 480px) {
    .pr-card {
        padding: var(--space-2, 8px);
        gap: var(--space-2, 8px);
    }

    .pr-card__title {
        font-size: var(--font-size-small, 12px);
        flex-direction: column;
        align-items: flex-start;
        gap: var(--space-1, 4px);
    }

    .pr-card__additional-info {
        flex-direction: column;
        align-items: flex-start;
        gap: var(--space-2, 8px);
    }
}

/* Accessibility */
@media (prefers-reduced-motion: reduce) {
    .pr-card {
        transition: none;
    }

    .pr-card:hover {
        transform: none;
    }
}

/* High Contrast Mode */
@media (prefers-contrast: high) {
    .pr-card {
        border-width: 2px;
    }

    .pr-card:hover {
        border-width: 3px;
    }

    .pr-state-badge {
        border: 1px solid currentColor;
    }
}

/* Dark mode support */
@media (prefers-color-scheme: dark) {
    .pr-state-badge--open {
        background: #1a7f3730;
        color: #4ade80;
    }

    .pr-state-badge--closed {
        background: #d1242f30;
        color: #f87171;
    }

    .pr-state-badge--merged {
        background: #8250df30;
        color: #c084fc;
    }
}
</style>
