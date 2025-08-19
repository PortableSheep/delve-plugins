# GitHub Dashboard Plugin for Delve

A comprehensive GitHub repository and pull request monitoring plugin for the Delve application, featuring real-time updates, theme inheritance, and a responsive design.

## 🚀 Features

- **📊 Repository Monitoring**: Track multiple GitHub repositories with real-time statistics
- **🔄 Pull Request Tracking**: View and monitor open pull requests across repositories
- **🎨 Theme Inheritance**: Seamlessly integrates with parent application themes
- **💾 Persistent Storage**: Configuration and state persist across sessions
- **📱 Responsive Design**: Works perfectly on desktop and mobile devices
- **⚡ Auto-refresh**: Configurable automatic data refresh intervals
- **🔒 Secure Authentication**: Uses GitHub Personal Access Tokens
- **📈 Rate Limit Monitoring**: Built-in GitHub API rate limit tracking

## 🏗️ Architecture

This plugin uses the new Delve SDK architecture with:
- **Go Backend**: Handles GitHub API communication and data processing
- **Vue Web Component**: Modern frontend built as a custom element
- **Theme System**: CSS custom properties for seamless theme integration
- **Storage System**: Persistent configuration and state management

## 📦 Installation

### Prerequisites

- Go 1.24.2 or later
- Node.js 16+ (optional, for development)
- GitHub Personal Access Token

### Build from Source

```bash
# Clone or navigate to the plugin directory
cd delve-plugins/github-dashboard

# Build the plugin (builds both Go backend and frontend component)
./build.sh

# The build creates:
# - github-dashboard (Go binary)
# - frontend/component.js (Vue web component)
# - releases/github-dashboard-*.tar.gz (release package)
```

### Using Pre-built Release

1. Download the latest release package
2. Extract the files to your Delve plugins directory
3. Ensure the `github-dashboard` binary is executable

## ⚙️ Configuration

### Required Settings

#### GitHub Personal Access Token
Create a token at [GitHub Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens):

Required scopes:
- `repo` (for private repositories)
- `public_repo` (for public repositories)  
- `read:user` (for user information)

#### Repository List
Specify repositories to monitor in `owner/repo` format:
```json
{
  "repositories": [
    "PortableSheep/delve",
    "golang/go",
    "microsoft/vscode"
  ]
}
```

### Optional Settings

```json
{
  "github_token": "ghp_your_token_here",
  "repositories": ["owner/repo1", "owner/repo2"],
  "refresh_interval": 300,
  "show_stars": true,
  "show_forks": true,
  "show_issues": true
}
```

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `github_token` | string | "" | GitHub Personal Access Token |
| `repositories` | array | [] | List of repositories to monitor |
| `refresh_interval` | integer | 300 | Auto-refresh interval in seconds |
| `show_stars` | boolean | true | Display star counts |
| `show_forks` | boolean | true | Display fork counts |
| `show_issues` | boolean | true | Display issue counts |

## 🎨 Theme Integration

The plugin inherits themes from the parent Delve application using CSS custom properties:

### Theme Variables

```css
:root {
  /* Core colors */
  --bg-color: #ffffff;
  --text-color: #24292f;
  --border-color: #d0d7de;
  --primary-color: #0969da;
  
  /* Status colors */
  --error-color: #cf222e;
  --success-color: #166534;
  --info-color: #0969da;
  
  /* Interactive elements */
  --card-bg: #ffffff;
  --muted-color: #656d76;
  --shadow-color: rgba(0, 0, 0, 0.1);
}
```

### Dark Theme Example

```css
[data-theme="dark"] {
  --bg-color: #0d1117;
  --text-color: #f0f6fc;
  --border-color: #30363d;
  --card-bg: #161b22;
  --primary-color: #58a6ff;
  --error-color: #f85149;
  --success-color: #56d364;
  --muted-color: #8b949e;
}
```

See [THEME-INHERITANCE.md](THEME-INHERITANCE.md) for complete theme customization guide.

## 🔧 Usage

1. **Configure Token**: Set your GitHub Personal Access Token in plugin settings
2. **Add Repositories**: Specify repositories to monitor
3. **View Dashboard**: 
   - Repository cards show stats, language, and last update
   - Click any repository to view its pull requests
   - Use the refresh button for manual updates
4. **Monitor Rate Limits**: API usage is displayed in the header

### Demo Mode

Without a GitHub token, the plugin shows demo data to demonstrate functionality.

## 🛠️ Development

### Project Structure

```
github-dashboard/
├── main.go                     # Go backend
├── go.mod                      # Go dependencies
├── build.sh                    # Build script
├── frontend/
│   ├── component.js           # Built web component
│   ├── build-component.js     # Build script
│   ├── package.json          # Node dependencies
│   └── src/
│       └── component.js      # Vue component source
├── releases/                  # Build artifacts
├── THEME-INHERITANCE.md      # Theme guide
├── CONFIGURATION-GUIDE.md    # Config guide
└── README.md                 # This file
```

### Building

```bash
# Full build (Go + frontend)
./build.sh

# Go backend only
go build -o github-dashboard main.go

# Frontend component only
cd frontend && npm run build:component
```

### Development Server

```bash
cd frontend
npm install
npm run dev
```

## 🔒 Security

- **Token Storage**: Tokens are stored securely using the Delve storage system
- **API Limits**: Built-in rate limiting prevents API abuse
- **Minimal Scopes**: Use only required GitHub API scopes
- **No Logging**: Sensitive data is never logged

## 🐛 Troubleshooting

### Common Issues

**"No GitHub token configured"**
- Configure your Personal Access Token in plugin settings

**"API rate limit exceeded"**  
- Increase refresh interval or wait for reset (time shown in UI)

**"Repository not found"**
- Check repository name format (`owner/repo`)
- Verify token has access to the repository

**Empty repository list**
- Ensure repository names are correct
- Check token scopes and permissions

### Debug Mode

Enable debug logging:
```javascript
localStorage.setItem('github_dashboard_debug', 'true')
```

See [CONFIGURATION-GUIDE.md](CONFIGURATION-GUIDE.md) for detailed troubleshooting.

## 📊 Performance

- **Efficient API Usage**: Batch requests and intelligent caching
- **Rate Limit Aware**: Automatic throttling and monitoring
- **Memory Management**: Proper cleanup and resource management
- **Responsive Loading**: Progressive data loading with loading states

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

### Development Guidelines

- Follow Go conventions for backend code
- Use Vue 3 Composition API for frontend
- Maintain theme variable compatibility
- Include comprehensive error handling
- Add appropriate documentation

## 📄 License

This plugin is part of the Delve project and follows the same licensing terms.

## 🆕 Changelog

### Version 2.0.0 (Current)
- ✨ New web component architecture
- 🎨 Theme inheritance support  
- 📱 Mobile responsive design
- ⚡ Improved performance and error handling
- 💾 Enhanced storage system integration

### Version 1.0.0
- 🚀 Initial release
- 📊 Basic repository monitoring
- 🔄 Pull request viewing
- ⚙️ Configuration management

## 🔗 Related Documentation

- [Theme Inheritance Guide](THEME-INHERITANCE.md)
- [Configuration Guide](CONFIGURATION-GUIDE.md)
- [Delve SDK Documentation](../delve_sdk/)
- [GitHub API Documentation](https://docs.github.com/en/rest)

## 💡 Tips

- Use longer refresh intervals for stable repositories
- Monitor API rate limits in the UI header
- Group related repositories for better organization
- Configure only actively developed repositories
- Use dark theme for better GitHub-like experience

---

**Made with ❤️ for the Delve community**