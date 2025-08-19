#!/bin/bash

set -e

echo "🔨 Building GitHub Dashboard Plugin..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go to build this plugin."
    exit 1
fi

# Get plugin directory
PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PLUGIN_DIR"

echo "📂 Plugin directory: $PLUGIN_DIR"

# Initialize Go module if go.mod doesn't exist
if [ ! -f "go.mod" ]; then
    echo "📝 Initializing Go module..."
    go mod init github.com/PortableSheep/github-dashboard
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod tidy

# Build the Go backend
echo "🔧 Building Go backend..."
CGO_ENABLED=0 go build -ldflags="-s -w" -o github-dashboard main.go

# Make the binary executable
chmod +x github-dashboard

echo "✅ Go backend built successfully: github-dashboard"

# Build the frontend component
echo "🎨 Building frontend component..."

cd frontend

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "⚠️  Node.js not found. Building component with static build script..."

    # Use the custom build script
    if [ -f "build-component.js" ]; then
        node build-component.js
    else
        echo "❌ build-component.js not found. Please ensure it exists."
        exit 1
    fi
else
    # Check if package.json exists
    if [ ! -f "package.json" ]; then
        echo "❌ package.json not found in frontend directory"
        exit 1
    fi

    # Install dependencies if node_modules doesn't exist
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing frontend dependencies..."
        npm install
    fi

    # Build the component
    echo "🏗️  Building Vue component..."
    npm run build:component
fi

echo "✅ Frontend component built successfully: component.js"

cd ..

# Create releases directory if it doesn't exist
mkdir -p releases

# Create a release package
RELEASE_NAME="github-dashboard-$(date +%Y%m%d-%H%M%S)"
RELEASE_DIR="releases/$RELEASE_NAME"

echo "📦 Creating release package: $RELEASE_NAME"

mkdir -p "$RELEASE_DIR"

# Copy necessary files
cp github-dashboard "$RELEASE_DIR/"
cp frontend/component.js "$RELEASE_DIR/"
cp go.mod "$RELEASE_DIR/"
cp go.sum "$RELEASE_DIR/" 2>/dev/null || true

# Create a simple README for the release
cat > "$RELEASE_DIR/README.md" << EOF
# GitHub Dashboard Plugin Release

This is a release package for the GitHub Dashboard plugin for Delve.

## Files

- \`github-dashboard\`: The Go backend binary
- \`component.js\`: The frontend Vue web component
- \`go.mod\` & \`go.sum\`: Go module files

## Installation

1. Copy the \`github-dashboard\` binary to your plugins directory
2. Copy the \`component.js\` file to your plugin's frontend directory
3. Configure your GitHub token in the plugin settings
4. Add repositories to monitor in the configuration

## Configuration

The plugin requires:
- \`github_token\`: Your GitHub Personal Access Token
- \`repositories\`: Array of repository names (e.g., ["owner/repo"])
- \`refresh_interval\`: Refresh interval in seconds (default: 300)

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
- \`--bg-color\`: Background color
- \`--text-color\`: Primary text color
- \`--border-color\`: Border color
- \`--primary-color\`: Accent color
- \`--error-color\`: Error color
- \`--success-color\`: Success color
- And many more...

Built on $(date)
EOF

# Create a tarball
cd releases
tar -czf "$RELEASE_NAME.tar.gz" "$RELEASE_NAME"
cd ..

echo "✅ Release package created: releases/$RELEASE_NAME.tar.gz"

# Display build summary
echo ""
echo "🎉 Build completed successfully!"
echo ""
echo "📋 Build Summary:"
echo "├── Go backend: github-dashboard ($(ls -lh github-dashboard | awk '{print $5}'))"
echo "├── Frontend component: frontend/component.js ($(ls -lh frontend/component.js | awk '{print $5}'))"
echo "└── Release package: releases/$RELEASE_NAME.tar.gz ($(ls -lh releases/$RELEASE_NAME.tar.gz | awk '{print $5}'))"
echo ""
echo "🚀 Plugin is ready to use!"
echo ""
echo "📖 Next steps:"
echo "1. Configure your GitHub token in the plugin settings"
echo "2. Add repositories to monitor"
echo "3. Enjoy your GitHub dashboard!"
