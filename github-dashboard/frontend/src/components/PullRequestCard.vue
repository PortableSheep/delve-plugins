<template>
  <div class="pr-card">
    <div class="pr-header">
      <div class="pr-title">
        <a :href="pullRequest.html_url" target="_blank" class="pr-link">
          #{{ pullRequest.number }} {{ pullRequest.title }}
        </a>
      </div>
      <div class="pr-state" :class="pullRequest.state">
        {{ pullRequest.state }}
      </div>
    </div>

    <div class="pr-meta">
      <div class="pr-author">
        <img
            :src="pullRequest.user.avatar_url"
            :alt="pullRequest.user.login"
            class="avatar"
        />
        <span>{{ pullRequest.user.login }}</span>
      </div>

      <div class="pr-dates">
        <span>Created {{ formatDate(pullRequest.created_at) }}</span>
        <span>Updated {{ formatDate(pullRequest.updated_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PullRequestCard',
  props: {
    pullRequest: {
      type: Object,
      required: true
    }
  },
  methods: {
    formatDate(dateString) {
      const date = new Date(dateString)
      const now = new Date()
      const diff = now - date
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))

      if (days === 0) return 'today'
      if (days === 1) return 'yesterday'
      if (days < 7) return `${days} days ago`
      return date.toLocaleDateString()
    }
  }
}
</script>

<style scoped>
.pr-card {
  border: 1px solid #d1d5da;
  border-radius: 6px;
  padding: 16px;
  background: white;
  transition: box-shadow 0.2s;
}

.pr-card:hover {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.pr-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  gap: 12px;
}

.pr-title {
  flex: 1;
}

.pr-link {
  color: #0366d6;
  text-decoration: none;
  font-weight: 500;
  line-height: 1.4;
}

.pr-link:hover {
  text-decoration: underline;
}

.pr-state {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  text-transform: capitalize;
  white-space: nowrap;
}

.pr-state.open {
  background: #dcffe4;
  color: #28a745;
}

.pr-state.closed {
  background: #ffeaea;
  color: #d73a49;
}

.pr-state.merged {
  background: #f0f4ff;
  color: #6f42c1;
}

.pr-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #586069;
}

.pr-author {
  display: flex;
  align-items: center;
  gap: 8px;
}

.avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
}

.pr-dates {
  display: flex;
  flex-direction: column;
  gap: 2px;
  text-align: right;
}

@media (max-width: 600px) {
  .pr-header {
    flex-direction: column;
    align-items: stretch;
  }

  .pr-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .pr-dates {
    text-align: left;
  }
}
</style>