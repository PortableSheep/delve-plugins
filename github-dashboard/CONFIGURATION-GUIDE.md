# GitHub Dashboard Plugin Configuration Guide

## Overview

The GitHub Dashboard plugin provides a comprehensive interface for monitoring GitHub repositories and pull requests within the Delve application. This guide covers all configuration options and setup procedures.

## Quick Start

1. **Install the plugin** in your Delve plugins directory
2. **Configure your GitHub token** (required for API access)
3. **Add repositories** to monitor
4. **Customize display options** as needed

## Configuration Options

### Required Settings

#### GitHub Personal Access Token
- **Key**: `github_token`
- **Type**: `string`
- **Required**: Yes (for real data)
- **Description**: Your GitHub Personal Access Token for API authentication

**How to create a GitHub token:**
1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Click "Generate new token (classic)"
3. Select these scopes:
   - `repo` (for private repositories)
   - `public_repo` (for public repositories)
   - `read:user` (for user information)
4. Copy the generated token

#### Repositories to Monitor
- **Key**: `repositories`
- **Type**: `array` of strings
- **Required**: No (demo data shown if empty)
- **Format**: `["owner/repo", "owner/another-repo"]`
- **Example**: `["PortableSheep/delve", "golang/go", "microsoft/vscode"]`

### Optional Settings

#### Refresh Interval
- **Key**: `refresh_interval`
- **Type**: `integer`
- **Default**: `300` (5 minutes)
- **Unit**: Seconds
- **Description**: How often to automatically refresh repository data

#### Display Options
- **Key**: `show_stars`
- **Type**: `boolean`
- **Default**: `true`
- **Description**: Show star counts on repository cards

- **Key**: `show_forks`
- **Type**: `boolean`
- **Default**: `true`
- **Description**: Show fork counts on repository cards

- **Key**: `show_issues`
- **Type**: `boolean`
- **Default**: `true`
- **Description**: Show open issue counts on repository cards

## Configuration Methods

### Method 1: Plugin Settings UI (Recommended)
If your Delve application provides a plugin settings interface:

1. Open plugin settings
2. Find "GitHub Dashboard" plugin
3. Enter your configuration values
4. Save settings

### Method 2: Direct Storage (Advanced)
Using the plugin's storage system:

```javascript
// Example configuration object
const config = {
    github_token: "ghp_your_token_here",
    repositories: [
        "PortableSheep/delve",
        "golang/go",
        "microsoft/vscode"
    ],
    refresh_interval: 300,
    show_stars: true,
    show_forks: true,
    show_issues: true
};

// Store configuration (if storage API is available)
if (window.delvePlugin && window.delvePlugin.storeConfig) {
    await window.delvePlugin.storeConfig('dashboard_settings', config, '1.0.0');
}
```

### Method 3: Environment Variables (Development)
For development or testing:

```bash
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_REPOS="PortableSheep/delve,golang/go"
export REFRESH_INTERVAL="600"
```

## Configuration Examples

### Basic Setup
Monitor a single repository with default settings:

```json
{
    "github_token": "ghp_your_token_here",
    "repositories": ["PortableSheep/delve"]
}
```

### Multiple Repositories
Monitor multiple repositories with custom refresh:

```json
{
    "github_token": "ghp_your_token_here",
    "repositories": [
        "PortableSheep/delve",
        "golang/go",
        "microsoft/vscode",
        "facebook/react",
        "vuejs/vue"
    ],
    "refresh_interval": 180,
    "show_stars": true,
    "show_forks": true,
    "show_issues": false
}
```

### Minimal Display
Focus on essential information only:

```json
{
    "github_token": "ghp_your_token_here",
    "repositories": ["your-org/critical-repo"],
    "refresh_interval": 60,
    "show_stars": false,
    "show_forks": false,
    "show_issues": true
}
```

### Demo Mode
No token required, shows sample data:

```json
{
    "repositories": [],
    "show_stars": true,
    "show_forks": true,
    "show_issues": true
}
```

## Security Considerations

### Token Security
- **Never commit tokens** to version control
- **Use environment variables** in production
- **Rotate tokens regularly** (recommended every 90 days)
- **Use minimal scopes** (only what's needed)
- **Store securely** using your platform's secret management

### Token Scopes
Minimum required scopes:
- `public_repo`: Access to public repositories
- `read:user`: Basic user information

Additional scopes for enhanced features:
- `repo`: Access to private repositories
- `read:org`: Organization information
- `notifications`: Notification access (future feature)

### Rate Limiting
- GitHub API has rate limits (5000 requests/hour for authenticated users)
- Plugin automatically handles rate limiting
- Longer refresh intervals reduce API usage
- Rate limit information is displayed in the UI

## Troubleshooting

### Common Issues

#### "No GitHub token configured" Warning
**Problem**: Plugin shows demo data and warning message
**Solution**: Configure your GitHub Personal Access Token

#### "API rate limit exceeded" Error
**Problem**: Too many API requests
**Solutions**:
- Increase `refresh_interval`
- Reduce number of monitored repositories
- Wait for rate limit reset (shown in UI)

#### "Repository not found" Error
**Problem**: Cannot access specified repository
**Possible causes**:
- Repository name is incorrect (use `owner/repo` format)
- Repository is private and token lacks `repo` scope
- Repository has been deleted or moved
- Token has expired or been revoked

#### Empty Repository List
**Problem**: No repositories shown despite having token
**Solutions**:
- Check repository names format (`owner/repo`)
- Verify token has correct scopes
- Check network connectivity
- Look for errors in browser console

#### Pull Requests Not Loading
**Problem**: Repository loads but pull requests don't
**Possible causes**:
- Network timeout (repository has many PRs)
- API rate limit reached
- Repository has disabled pull requests
- Token lacks sufficient permissions

### Debug Mode
Enable debug logging in browser console:

```javascript
// Enable debug logging
localStorage.setItem('github_dashboard_debug', 'true');

// Disable debug logging
localStorage.removeItem('github_dashboard_debug');
```

### Network Issues
Check these if experiencing connectivity problems:
- Internet connection stability
- Corporate firewall blocking GitHub API
- Proxy configuration
- DNS resolution for `api.github.com`

## Performance Optimization

### Repository Selection
- Monitor only actively developed repositories
- Avoid archived or inactive repositories
- Group related repositories logically

### Refresh Intervals
- Use longer intervals for stable repositories
- Use shorter intervals for active development
- Consider your team's workflow patterns

### API Usage
- Monitor rate limit usage in the UI
- Adjust monitoring frequency based on usage
- Use webhooks for real-time updates (future feature)

## Advanced Configuration

### Custom Themes
The plugin supports theme inheritance. Configure theme variables in your main application:

```css
:root {
    --primary-color: #your-brand-color;
    --bg-color: #your-background;
    --text-color: #your-text-color;
}
```

### Storage Backend
The plugin supports multiple storage backends:
- Delve SDK storage (preferred)
- Browser localStorage (fallback)
- Custom storage adapters (advanced)

### Integration with Other Tools
The plugin can be integrated with:
- CI/CD pipelines (via webhooks)
- Issue tracking systems
- Team communication tools
- Project management platforms

## Migration Guide

### From Version 1.0 to 2.0
1. Update plugin binary
2. Update frontend component
3. Configuration format remains compatible
4. New theme system available

### From Other GitHub Tools
- Export repository lists from existing tools
- Convert webhook configurations
- Migrate team settings and preferences

## Support and Community

### Getting Help
- Check this configuration guide
- Review troubleshooting section
- Search existing issues
- Create new issue with details

### Reporting Issues
Include this information:
- Plugin version
- Configuration (without sensitive data)
- Error messages
- Browser console logs
- Steps to reproduce

### Contributing
- Submit bug reports
- Suggest new features
- Contribute code improvements
- Update documentation

## API Reference

### Configuration Schema
```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "github_token": {
            "type": "string",
            "description": "GitHub Personal Access Token"
        },
        "repositories": {
            "type": "array",
            "items": {
                "type": "string",
                "pattern": "^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$"
            },
            "description": "List of repositories in owner/repo format"
        },
        "refresh_interval": {
            "type": "integer",
            "minimum": 30,
            "maximum": 3600,
            "default": 300,
            "description": "Refresh interval in seconds"
        },
        "show_stars": {
            "type": "boolean",
            "default": true,
            "description": "Show star counts"
        },
        "show_forks": {
            "type": "boolean",
            "default": true,
            "description": "Show fork counts"
        },
        "show_issues": {
            "type": "boolean",
            "default": true,
            "description": "Show issue counts"
        }
    }
}
```

### Storage Keys
- `dashboard_settings`: Main configuration
- `repository_data`: Cached repository information
- `user_preferences`: UI preferences
- `api_cache`: API response cache

## Changelog

### Version 2.0.0
- New web component architecture
- Theme inheritance support
- Improved performance
- Enhanced error handling
- Mobile responsive design

### Version 1.0.0
- Initial release
- Basic repository monitoring
- Pull request viewing
- Configuration management