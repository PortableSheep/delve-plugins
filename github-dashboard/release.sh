#!/bin/bash

set -e

echo "🚀 Building GitHub Dashboard Plugin for Release..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go to build this plugin."
    exit 1
fi

# Get plugin directory
PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PLUGIN_DIR"

VERSION=${1:-"v1.0.0"}
echo "📋 Building version: $VERSION"

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf releases/$VERSION
mkdir -p releases/$VERSION

# Build frontend first
echo "🎨 Building frontend component..."
cd frontend

# Check if Node.js is installed and build component
if command -v node &> /dev/null && [ -f "package.json" ]; then
    # Install dependencies if needed
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing frontend dependencies..."
        npm install
    fi

    # Build the component
    echo "🏗️  Building Vue component..."
    npm run build:component
else
    echo "⚠️  Node.js not found or package.json missing. Using static build..."
    if [ -f "build-component.js" ]; then
        node build-component.js
    fi
fi

cd ..

# Ensure component.js exists
if [ ! -f "frontend/component.js" ]; then
    echo "❌ Frontend component.js not found after build"
    exit 1
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod tidy

# Define target platforms
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

echo "🔨 Building binaries for ${#PLATFORMS[@]} platforms..."

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    BINARY_NAME="github-dashboard"
    if [ "$GOOS" = "windows" ]; then
        BINARY_NAME="github-dashboard.exe"
    fi

    OUTPUT_DIR="releases/$VERSION"
    PLATFORM_DIR="$OUTPUT_DIR/${GOOS}-${GOARCH}"

    echo "  📦 Building for $GOOS/$GOARCH..."

    # Create platform directory
    mkdir -p "$PLATFORM_DIR"

    # Build binary
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X main.Version=$VERSION" \
        -o "$PLATFORM_DIR/$BINARY_NAME" \
        main.go

    # Copy frontend component
    cp frontend/component.js "$PLATFORM_DIR/"

    # Copy metadata
    cp build/plugin.json "$PLATFORM_DIR/"

    # Create platform-specific archive
    cd "$OUTPUT_DIR"
    if [ "$GOOS" = "windows" ]; then
        zip -r "github-dashboard-${GOOS}-${GOARCH}.zip" "${GOOS}-${GOARCH}/"
    else
        tar -czf "github-dashboard-${GOOS}-${GOARCH}.tar.gz" "${GOOS}-${GOARCH}/"
    fi
    cd - > /dev/null

    # Generate checksums
    if [ "$GOOS" = "windows" ]; then
        ARCHIVE="$OUTPUT_DIR/github-dashboard-${GOOS}-${GOARCH}.zip"
    else
        ARCHIVE="$OUTPUT_DIR/github-dashboard-${GOOS}-${GOARCH}.tar.gz"
    fi

    CHECKSUM=$(shasum -a 256 "$ARCHIVE" | cut -d' ' -f1)
    SIZE=$(stat -f%z "$ARCHIVE" 2>/dev/null || stat -c%s "$ARCHIVE" 2>/dev/null)

    echo "    ✓ ${GOOS}-${GOARCH}: $(basename "$ARCHIVE") (${SIZE} bytes, sha256:${CHECKSUM:0:8}...)"

    # Store checksums for registry update
    echo "${GOOS}-${GOARCH}:${CHECKSUM}:${SIZE}" >> "$OUTPUT_DIR/checksums.txt"
done

# Create combined frontend archive
echo "🎨 Creating frontend distribution..."
FRONTEND_DIR="$OUTPUT_DIR/frontend"
mkdir -p "$FRONTEND_DIR"

# Copy all frontend assets
cp frontend/component.js "$FRONTEND_DIR/"
cp frontend/dist/index.html "$FRONTEND_DIR/" 2>/dev/null || cp frontend/index.html "$FRONTEND_DIR/"

# Create frontend archive
cd "$OUTPUT_DIR"
tar -czf "frontend.tar.gz" frontend/
FRONTEND_CHECKSUM=$(shasum -a 256 "frontend.tar.gz" | cut -d' ' -f1)
echo "frontend:${FRONTEND_CHECKSUM}" >> checksums.txt
cd - > /dev/null

# Create plugin metadata with checksums
echo "📝 Creating plugin metadata..."
cat > "$OUTPUT_DIR/plugin.json" << EOF
{
  "info": {
    "id": "github-dashboard",
    "name": "GitHub Dashboard",
    "version": "$VERSION",
    "description": "Monitor GitHub repositories and pull requests with real-time updates",
    "author": "Michael Gunderson",
    "license": "MIT",
    "homepage": "https://github.com/PortableSheep/delve-plugins/tree/main/github-dashboard",
    "repository": "https://github.com/PortableSheep/delve-plugins",
    "icon": "🐙",
    "tags": ["github", "dashboard", "monitoring", "git", "repositories"],
    "category": "development-tools",
    "min_delve_version": "v0.1.0"
  },
  "runtime": {
    "executable": "github-dashboard",
    "frontend_entry": "component.js",
    "permissions": [
      "network.http",
      "storage.local"
    ]
  },
  "config_schema": {
    "github_token": {
      "type": "string",
      "required": true,
      "description": "GitHub Personal Access Token",
      "sensitive": true
    },
    "repositories": {
      "type": "array",
      "items": {
        "type": "string",
        "pattern": "^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$"
      },
      "description": "List of repositories to monitor (owner/repo format)",
      "default": []
    },
    "refresh_interval": {
      "type": "integer",
      "minimum": 30,
      "maximum": 3600,
      "default": 300,
      "description": "Data refresh interval in seconds"
    },
    "compact_view": {
      "type": "boolean",
      "default": false,
      "description": "Use compact view for repository list"
    }
  },
  "health_checks": [
    {
      "name": "github_api_connectivity",
      "interval": "60s",
      "timeout": "10s",
      "description": "Check GitHub API connectivity"
    },
    {
      "name": "token_validity",
      "interval": "300s",
      "timeout": "5s",
      "description": "Validate GitHub token"
    }
  ],
  "build_info": {
    "built_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "go_version": "$(go version | cut -d' ' -f3)",
    "commit": "$(git rev-parse HEAD 2>/dev/null || echo 'unknown')"
  }
}
EOF

# Generate registry entry
echo "📊 Generating registry entry..."
cat > "$OUTPUT_DIR/registry-entry.yml" << 'EOL'
github-dashboard:
  name: "GitHub Dashboard"
  description: "Monitor GitHub repositories and pull requests with real-time updates"
  author: "Michael Gunderson"
  license: "MIT"
  repository: "https://github.com/PortableSheep/delve-plugins"
  homepage: "https://github.com/PortableSheep/delve-plugins/tree/main/github-dashboard"
  tags: ["github", "dashboard", "monitoring", "git", "repositories"]
  category: "development-tools"
  min_delve_version: "v0.1.0"
  versions:
EOL

# Add version entry with checksums
echo "    - version: \"$VERSION\"" >> "$OUTPUT_DIR/registry-entry.yml"
echo "      released: \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" >> "$OUTPUT_DIR/registry-entry.yml"
echo "      compatibility: [\"v0.1.0\", \"v0.2.0\"]" >> "$OUTPUT_DIR/registry-entry.yml"
echo "      assets:" >> "$OUTPUT_DIR/registry-entry.yml"

# Add asset entries from checksums
while IFS=':' read -r platform checksum size; do
    if [ "$platform" != "frontend" ]; then
        GOOS=${platform%-*}
        GOARCH=${platform#*-}

        if [ "$GOOS" = "windows" ]; then
            ASSET_NAME="github-dashboard-${platform}.zip"
        else
            ASSET_NAME="github-dashboard-${platform}.tar.gz"
        fi

        cat >> "$OUTPUT_DIR/registry-entry.yml" << EOL
        ${platform}:
          url: "github-dashboard/releases/$VERSION/$ASSET_NAME"
          checksum: "sha256:$checksum"
          size: $size
EOL
    fi
done < "$OUTPUT_DIR/checksums.txt"

# Add frontend entry
FRONTEND_CHECKSUM=$(grep "frontend:" "$OUTPUT_DIR/checksums.txt" | cut -d':' -f2)
cat >> "$OUTPUT_DIR/registry-entry.yml" << EOL
      frontend:
        entry: "component.js"
        checksum: "sha256:$FRONTEND_CHECKSUM"
      plugin_metadata:
        url: "github-dashboard/plugin.json"
        checksum: "sha256:$(shasum -a 256 "$OUTPUT_DIR/plugin.json" | cut -d' ' -f1)"
EOL

# Create release summary
echo "📋 Creating release summary..."
cat > "$OUTPUT_DIR/RELEASE.md" << EOF
# GitHub Dashboard Plugin - Release $VERSION

Generated on: $(date)

## 📦 Assets

### Platform Binaries
$(ls -la "$OUTPUT_DIR"/*.tar.gz "$OUTPUT_DIR"/*.zip 2>/dev/null | awk '{printf "- %s (%s bytes)\n", $9, $5}' | sed 's|.*/||')

### Frontend
- frontend.tar.gz (Vue component and assets)

### Metadata
- plugin.json (Plugin configuration and schema)
- registry-entry.yml (For registry integration)

## 🔍 Checksums

$(cat "$OUTPUT_DIR/checksums.txt" | while IFS=':' read -r platform checksum size; do
    echo "- $platform: sha256:${checksum:0:16}... ($size bytes)"
done)

## 🚀 Installation

1. Download the appropriate binary for your platform
2. Extract the archive to your Delve plugins directory
3. Configure your GitHub token in the plugin settings
4. Add repositories to monitor

## 📖 Configuration

Required:
- \`github_token\`: Your GitHub Personal Access Token with repo access

Optional:
- \`repositories\`: Array of "owner/repo" strings to monitor
- \`refresh_interval\`: Update frequency in seconds (30-3600, default: 300)
- \`compact_view\`: Use compact display mode (default: false)

## 🔗 Links

- Repository: https://github.com/PortableSheep/delve-plugins
- Documentation: https://github.com/PortableSheep/delve-plugins/tree/main/github-dashboard
- Issues: https://github.com/PortableSheep/delve-plugins/issues

EOF

# Display build summary
echo ""
echo "🎉 Release build completed successfully!"
echo ""
echo "📊 Build Summary:"
echo "├── Version: $VERSION"
echo "├── Platforms: ${#PLATFORMS[@]} ($(echo "${PLATFORMS[@]}" | tr ' ' ', '))"
echo "├── Total archives: $(ls "$OUTPUT_DIR"/*.tar.gz "$OUTPUT_DIR"/*.zip 2>/dev/null | wc -l | tr -d ' ')"
echo "└── Output directory: releases/$VERSION"
echo ""
echo "📁 Release artifacts:"
ls -la "$OUTPUT_DIR" | grep -E '\.(tar\.gz|zip|json|yml|md|txt)$' | awk '{printf "├── %s (%s bytes)\n", $9, $5}'
echo ""
echo "🚀 Next steps:"
echo "1. Test the release artifacts"
echo "2. Create git tag: git tag $VERSION"
echo "3. Push to repository: git push origin $VERSION"
echo "4. Update delve-registry with registry-entry.yml content"
echo "5. Publish release artifacts to distribution server"
