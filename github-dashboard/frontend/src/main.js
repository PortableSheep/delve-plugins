import { createApp } from 'vue'
import App from './App.vue'

// Plugin API helper functions
window.pluginAPI = {
    async request(method, path, body = null) {
        try {
            // This would be implemented by the plugin system to communicate with the Go backend
            const response = await fetch(`/plugin/github-dashboard${path}`, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: body ? JSON.stringify(body) : null
            })

            if (!response.ok) {
                throw new Error(`Request failed: ${response.statusText}`)
            }

            return await response.json()
        } catch (error) {
            console.error('Plugin API request failed:', error)
            throw error
        }
    },

    async getRepositories() {
        return this.request('GET', '/repositories')
    },

    async getPullRequests(repo) {
        return this.request('GET', `/repositories/${encodeURIComponent(repo)}/pulls`)
    },

    async refresh() {
        return this.request('POST', '/refresh')
    }
}

createApp(App).mount('#app')