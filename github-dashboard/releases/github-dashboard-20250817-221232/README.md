# GitHub Dashboard Plugin Release

This is a release package for the GitHub Dashboard plugin for Delve.

## Files

- `github-dashboard`: The Go backend binary
- `component.js`: The frontend Vue web component
- `go.mod` & `go.sum`: Go module files

## Installation

1. Copy the `github-dashboard` binary to your plugins directory
2. Copy the `component.js` file to your plugin's frontend directory
3. Configure your GitHub token in the plugin settings
4. Add repositories to monitor in the configuration

## Configuration

The plugin requires:
- `github_token`: Your GitHub Personal Access Token
- `repositories`: Array of repository names (e.g., ["owner/repo"])
- `refresh_interval`: Refresh interval in seconds (default: 300)

## Features

- 🐙 Monitor multiple GitHub repositories
- 📊 View repository statistics (stars, forks, issues)
- 🔄 Real-time pull request monitoring
- 🎨 Theme-aware UI that inherits from parent application
- 💾 Persistent configuration storage
- 🔄 Auto-refresh capabilities
- 📱 Responsive design for mobile and desktop

## Theme Variables

The plugin uses CSS custom properties for theming:
- `--bg-color`: Background color
- `--text-color`: Primary text color
- `--border-color`: Border color
- `--primary-color`: Accent color
- `--error-color`: Error color
- `--success-color`: Success color
- And many more...

Built on Sun Aug 17 22:12:32 CDT 2025
