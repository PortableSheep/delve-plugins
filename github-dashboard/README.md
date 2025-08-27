# 🐙 GitHub Dashboard Plugin

A comprehensive GitHub dashboard plugin for Delve that provides real-time monitoring of your repositories, pull requests, and development activity with a beautiful, responsive interface.

![Version](https://img.shields.io/badge/version-1.0.1-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)

## ✨ Features

- **📊 Repository Monitoring**: View key metrics including stars, forks, issues, and recent activity
- **🔄 Pull Request Tracking**: Monitor open PRs with detailed information and status indicators
- **🎨 Modern UI**: Built with GitHub's Primer design system for familiar, polished experience
- **📱 Responsive Design**: Optimized for all screen sizes with mobile-first approach
- **⚡ Intelligent Caching**: Reduces API calls and improves performance with smart data caching
- **🛡️ Error Recovery**: Robust error handling with graceful fallbacks and user-friendly messages
- **♿ Accessibility**: Full keyboard navigation and screen reader support
- **🎭 Demo Mode**: Works without GitHub token using realistic demo data
- **⚙️ Flexible Configuration**: Extensive customization options for different workflows

## 🚀 Installation

### Via Plugin Manager (Recommended)

1. Open Delve application
2. Navigate to **Plugins** → **Browse Registry**
3. Search for "GitHub Dashboard"
4. Click **Install**

### Manual Installation

1. Download the latest release from the [plugin registry](https://github.com/PortableSheep/delve-registry)
2. Extract to your Delve plugins directory
3. Restart Delve

## 🔧 Configuration

### Basic Setup

1. Open Delve Settings → Plugins → GitHub Dashboard
2. **Optional**: Add your GitHub Personal Access Token for:
   - Higher rate limits (5000/hour vs 60/hour)
   - Access to private repositories
   - Real-time data instead of demo data

### GitHub Token Setup

1. Go to [GitHub Settings → Personal Access Tokens](https://github.com/settings/tokens)
2. Click **Generate new token** → **Generate new token (classic)**
3. Select scopes:
   - `public_repo` (for public repositories)
   - `repo` (for private repositories)
   - `read:user` (for user information)
4. Copy the token and paste it in the plugin configuration

### Repository Configuration

Add repositories to monitor in the format `owner/repository-name`:

```
PortableSheep/delve
golang/go
microsoft/vscode
facebook/react
```

### Advanced Configuration

| Setting | Description | Default | Range |
|---------|-------------|---------|-------|
| `refresh_interval` | Background refresh interval (seconds) | 300 | 30-3600 |
| `compact_view` | Use compact layout for repository cards | false | boolean |
| `show_private_repos` | Include private repositories | true | boolean |
| `max_repos_per_page` | Maximum repositories to display | 50 | 10-100 |
| `max_prs_per_repo` | Maximum PRs per repository | 20 | 5-50 |
| `cache_timeout` | Data cache timeout (seconds) | 300 | 60-1800 |

## 📖 Usage

### Dashboard Overview

The GitHub Dashboard provides two main sections:

1. **Repositories**: Cards showing repository information, statistics, and recent activity
2. **Pull Requests**: Detailed view of open PRs for the selected repository

### Navigation

- **Click** any repository card to view its pull requests
- **Click** the external link icon (↗) to open repositories/PRs on GitHub
- **Use keyboard navigation** with Tab/Enter/Space for accessibility

### Visual Indicators

- **Green badge**: Open pull requests
- **Red badge**: Closed pull requests  
- **Purple badge**: Merged pull requests
- **Draft indicator**: Shows draft PRs with 📝 icon
- **Mergeable status**: Clean ✅, Conflicts ⚠️, Blocked 🚫

## 🛠️ Development

### Prerequisites

- Node.js 18+ 
- npm 9+
- Go 1.21+

### Frontend Development

```bash
cd frontend/
npm install
npm run dev
```

### Backend Development

```bash
go mod tidy
go run main.go --ws-port=8080
```

### Building

```bash
# Frontend component
cd frontend/
npm run build:component

# Backend plugin
go build -o github-dashboard main.go
```

### Testing

The plugin includes comprehensive error handling and works in multiple modes:

- **With GitHub Token**: Full functionality with real data
- **Without Token**: Demo mode with realistic sample data
- **Offline Mode**: Graceful degradation with cached data

## 🎨 UI Components

### Repository Card

```vue
<RepositoryCard
  :repository="repo"
  :selected="selectedRepo === repo.full_name"
  @select="selectRepository"
/>
```

### Pull Request Card

```vue
<PullRequestCard
  :pull-request="pr"
  :show-additional-info="true"
/>
```

## 🔌 API Reference

### Plugin Messages

| Message Type | Description | Payload |
|--------------|-------------|---------|
| `1` | Get repositories | - |
| `2` | Get pull requests | `{"method": "getPullRequests", "params": {"repository": "owner/repo"}}` |
| `3` | Force refresh | - |
| `4` | Get configuration | - |
| `5` | Update configuration | `Config` object |
| `6` | Health check | - |

### Response Format

```json
{
  "success": true,
  "data": {...},
  "error": "Error message if success is false"
}
```

## 📱 Screenshots

### Dashboard Overview
![Dashboard Overview](screenshots/dashboard-overview.png)

### Repository Cards
![Repository Cards](screenshots/repository-cards.png)

### Pull Requests
![Pull Requests](screenshots/pull-requests.png)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Vue 3 Composition API patterns
- Use GitHub's Primer design tokens
- Ensure accessibility with ARIA labels
- Add comprehensive error handling
- Include responsive design considerations
- Write meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/PortableSheep/delve-plugins/issues)
- **Documentation**: [Plugin Documentation](https://github.com/PortableSheep/delve-plugins/tree/main/github-dashboard)
- **Community**: [Delve Community Forum](https://community.delve.dev)

## 📝 Changelog

### v1.0.1 (2025-08-26)
- ✨ Complete UI redesign with GitHub Primer design system
- 📱 Improved responsive design and accessibility
- ⚡ Added intelligent caching and error recovery
- 🎨 Better loading states and visual hierarchy
- 📦 Enhanced plugin distribution compatibility
- 🎭 Added demo mode for users without GitHub tokens
- ⚙️ Improved configuration options and validation

### v1.0.0 (2025-08-23)
- 🎉 Initial release
- 📊 Basic GitHub repository and PR monitoring
- 🔌 WebSocket communication with host
- ⚙️ Configuration management

## 🙏 Acknowledgments

- **GitHub Primer**: Design system and components
- **Vue.js**: Reactive UI framework
- **Vite**: Build tool and development server
- **Delve Team**: Plugin system and SDK

---

<div align="center">
Made with ❤️ for the Delve community
</div>