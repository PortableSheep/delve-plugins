<template>
  <div class="github-dashboard">
    <header class="dashboard-header">
      <h1>🐙 GitHub Dashboard</h1>
      <button @click="refresh" :disabled="loading" class="refresh-btn">
        <span v-if="loading">🔄</span>
        <span v-else>↻</span>
        Refresh
      </button>
    </header>

    <div v-if="error" class="error-message">
      ❌ {{ error }}
    </div>

    <div v-if="loading && repositories.length === 0" class="loading">
      Loading repositories...
    </div>

    <div v-else class="dashboard-content">
      <div class="repositories-grid">
        <RepositoryList
            :repositories="repositories"
            @select-repository="selectRepository"
            :selected="selectedRepo"
        />
      </div>

      <div v-if="selectedRepo" class="pull-requests-panel">
        <h2>Pull Requests - {{ selectedRepo }}</h2>
        <div v-if="pullRequestsLoading" class="loading">
          Loading pull requests...
        </div>
        <div v-else class="pull-requests">
          <PullRequestCard
              v-for="pr in pullRequests"
              :key="pr.number"
              :pull-request="pr"
          />
          <div v-if="pullRequests.length === 0" class="no-prs">
            No open pull requests
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import RepositoryList from './components/RepositoryList.vue'
import PullRequestCard from './components/PullRequestCard.vue'

export default {
  name: 'GitHubDashboard',
  components: {
    RepositoryList,
    PullRequestCard
  },
  setup() {
    const repositories = ref([])
    const pullRequests = ref([])
    const selectedRepo = ref(null)
    const loading = ref(false)
    const pullRequestsLoading = ref(false)
    const error = ref(null)

    const loadRepositories = async () => {
      try {
        loading.value = true
        error.value = null
        const data = await window.pluginAPI.getRepositories()
        repositories.value = data || []
      } catch (err) {
        error.value = err.message
        console.error('Failed to load repositories:', err)
      } finally {
        loading.value = false
      }
    }

    const selectRepository = async (repo) => {
      selectedRepo.value = repo.full_name
      try {
        pullRequestsLoading.value = true
        const data = await window.pluginAPI.getPullRequests(repo.full_name)
        pullRequests.value = data || []
      } catch (err) {
        error.value = err.message
        console.error('Failed to load pull requests:', err)
        pullRequests.value = []
      } finally {
        pullRequestsLoading.value = false
      }
    }

    const refresh = async () => {
      await loadRepositories()
      if (selectedRepo.value) {
        const repo = repositories.value.find(r => r.full_name === selectedRepo.value)
        if (repo) {
          await selectRepository(repo)
        }
      }
    }

    onMounted(() => {
      loadRepositories()
    })

    return {
      repositories,
      pullRequests,
      selectedRepo,
      loading,
      pullRequestsLoading,
      error,
      selectRepository,
      refresh
    }
  }
}
</script>

<style scoped>
.github-dashboard {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e1e5e9;
}

.dashboard-header h1 {
  margin: 0;
  color: #24292e;
  font-size: 24px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #0366d6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: #0256cc;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  background: #ffeaea;
  color: #d73a49;
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 20px;
  border: 1px solid #fdaeb7;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #586069;
}

.dashboard-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.pull-requests-panel h2 {
  margin: 0 0 16px 0;
  color: #24292e;
  font-size: 18px;
}

.pull-requests {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.no-prs {
  text-align: center;
  color: #586069;
  padding: 40px;
  font-style: italic;
}

@media (max-width: 768px) {
  .dashboard-content {
    grid-template-columns: 1fr;
  }

  .dashboard-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
}
</style>