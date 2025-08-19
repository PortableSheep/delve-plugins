# GitHub Dashboard Plugin - Release v1.0.0

Generated on: Mon Aug 18 21:58:22 CDT 2025

## 📦 Assets

### Platform Binaries
frontend.tar.gz (6990 bytes)
github-dashboard-darwin-amd64.tar.gz (2601874 bytes)
github-dashboard-darwin-arm64.tar.gz (2424205 bytes)
github-dashboard-linux-amd64.tar.gz (2563600 bytes)
github-dashboard-linux-arm64.tar.gz (2337589 bytes)
github-dashboard-windows-amd64.zip (2615317 bytes)

### Frontend
- frontend.tar.gz (Vue component and assets)

### Metadata
- plugin.json (Plugin configuration and schema)
- registry-entry.yml (For registry integration)

## 🔍 Checksums

- darwin-amd64: sha256:390ca2a3ece66850... (2601874 bytes)
- darwin-arm64: sha256:35285a05e9fa84c3... (2424205 bytes)
- linux-amd64: sha256:aaced8bcffe33170... (2563600 bytes)
- linux-arm64: sha256:336e845d5d97f2f0... (2337589 bytes)
- windows-amd64: sha256:0d7283aa5a0f767b... (2615317 bytes)
- frontend: sha256:b3c4ff24e5ab6067... ( bytes)

## 🚀 Installation

1. Download the appropriate binary for your platform
2. Extract the archive to your Delve plugins directory
3. Configure your GitHub token in the plugin settings
4. Add repositories to monitor

## 📖 Configuration

Required:
- `github_token`: Your GitHub Personal Access Token with repo access

Optional:
- `repositories`: Array of "owner/repo" strings to monitor
- `refresh_interval`: Update frequency in seconds (30-3600, default: 300)
- `compact_view`: Use compact display mode (default: false)

## 🔗 Links

- Repository: https://github.com/PortableSheep/delve-plugins
- Documentation: https://github.com/PortableSheep/delve-plugins/tree/main/github-dashboard
- Issues: https://github.com/PortableSheep/delve-plugins/issues

