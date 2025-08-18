<template>
  <div class="repository-list">
    <h2>Repositories</h2>
    <div class="repo-cards">
      <div
          v-for="repo in repositories"
          :key="repo.full_name"
          class="repo-card"
          :class="{ selected: selected === repo.full_name }"
          @click="$emit('select-repository', repo)"
      >
        <div class="repo-header">
          <h3>{{ repo.name }}</h3>
          <div class="repo-stats">
            <span class="stat">
              ⭐ {{ repo.stargazers_count }}
            </span>
            <span class="stat">
              🍴 {{ repo.forks_count }}
            </span>
            <span class="stat">
              🐛 {{ repo.open_issues_count }}
            </span>
          </div>
        </div>

        <p class="repo-description">{{ repo.description || 'No description' }}</p>

        <div class="repo-footer">
          <span v-if="repo.language" class="language">{{ repo.language }}</span>
          <span class="updated">{{ formatDate(repo.updated_at) }}</span>
        </div>
      </div>

      <div v-if="repositories.length === 0" class="no-repos">
        No repositories configured
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'RepositoryList',
  props: {
    repositories: {
      type: Array,
      default: () => []
    },
    selected: {
      type: String,
      default: null
    }
  },
  emits: ['select-repository'],
  methods: {
    formatDate(dateString) {
      const date = new Date(dateString)
      return date.toLocaleDateString()
    }
  }
}
</script>

<style scoped>
.repository-list h2 {
  margin: 0 0 16px 0;
  color: #24292e;
  font-size: 18px;
}

.repo-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.repo-card {
  border: 1px solid #d1d5da;
  border-radius: 6px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
}

.repo-card:hover {
  border-color: #0366d6;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.repo-card.selected {
  border-color: #0366d6;
  background: #f6f8fa;
}

.repo-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.repo-header h3 {
  margin: 0;
  color: #0366d6;
  font-size: 16px;
  font-weight: 600;
}

.repo-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #586069;
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
}

.repo-description {
  margin: 0 0 12px 0;
  color: #586069;
  font-size: 14px;
  line-height: 1.4;
}

.repo-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #586069;
}

.language {
  background: #f1f8ff;
  color: #0366d6;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.no-repos {
  text-align: center;
  color: #586069;
  padding: 40px;
  font-style: italic;
}
</style>