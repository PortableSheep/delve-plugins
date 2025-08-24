#!/bin/bash

# Generic Plugin Build Script
# Usage: ./scripts/build-plugin.sh [plugin-name] [options]
# If no plugin name provided, builds the plugin in current directory

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  [INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅ [SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️  [WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}❌ [ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}🔨 $1${NC}"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [plugin-name] [options]"
    echo ""
    echo "Arguments:"
    echo "  plugin-name           Name of plugin to build (optional if run from plugin directory)"
    echo ""
    echo "Options:"
    echo "  -v, --version <ver>   Override version from plugin.json"
    echo "  -p, --platforms <list> Comma-separated list of platforms (default: all)"
    echo "  -o, --output <dir>    Output directory (default: releases/\$VERSION)"
    echo "  -c, --clean          Clean previous builds"
    echo "  --no-frontend        Skip frontend build"
    echo "  --no-package         Skip creating release package"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Platform options:"
    echo "  darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Build current directory plugin"
    echo "  $0 github-dashboard                  # Build specific plugin"
    echo "  $0 github-dashboard -v v1.1.0        # Build with custom version"
    echo "  $0 github-dashboard -p darwin/amd64  # Build for specific platform"
    echo "  $0 github-dashboard --clean          # Clean and build"
}

# Default values
PLUGIN_NAME=""
PLUGIN_DIR=""
VERSION=""
PLATFORMS="darwin/amd64,darwin/arm64,linux/amd64,linux/arm64,windows/amd64"
OUTPUT_DIR=""
CLEAN=false
BUILD_FRONTEND=true
CREATE_PACKAGE=true
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -p|--platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        --no-frontend)
            BUILD_FRONTEND=false
            shift
            ;;
        --no-package)
            CREATE_PACKAGE=false
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        -*)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
        *)
            if [[ -z "$PLUGIN_NAME" ]]; then
                PLUGIN_NAME="$1"
            else
                print_error "Too many arguments"
                show_usage
                exit 1
            fi
            shift
            ;;
    esac
done

# Determine plugin directory
if [[ -n "$PLUGIN_NAME" ]]; then
    PLUGIN_DIR="$ROOT_DIR/$PLUGIN_NAME"
    if [[ ! -d "$PLUGIN_DIR" ]]; then
        print_error "Plugin directory not found: $PLUGIN_DIR"
        exit 1
    fi
else
    # Check if we're in a plugin directory
    if [[ -f "plugin.json" && -f "main.go" ]]; then
        PLUGIN_DIR="$(pwd)"
        PLUGIN_NAME="$(basename "$PLUGIN_DIR")"
    else
        print_error "No plugin specified and current directory is not a plugin"
        print_error "Run from a plugin directory or specify plugin name"
        show_usage
        exit 1
    fi
fi

print_info "Building plugin: $PLUGIN_NAME"
print_info "Plugin directory: $PLUGIN_DIR"

cd "$PLUGIN_DIR"

# Validate plugin structure
if [[ ! -f "plugin.json" ]]; then
    print_error "plugin.json not found in $PLUGIN_DIR"
    exit 1
fi

if [[ ! -f "main.go" ]]; then
    print_error "main.go not found in $PLUGIN_DIR"
    exit 1
fi

# Check dependencies
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go to build this plugin."
    exit 1
fi

if ! command -v jq &> /dev/null; then
    print_error "jq is not installed. Please install jq to parse plugin.json."
    exit 1
fi

# Parse plugin metadata
PLUGIN_INFO=$(jq -r '.info' plugin.json 2>/dev/null)
if [[ "$PLUGIN_INFO" == "null" ]]; then
    print_error "Invalid plugin.json: missing 'info' section"
    exit 1
fi

PLUGIN_ID=$(echo "$PLUGIN_INFO" | jq -r '.id // empty')
PLUGIN_DISPLAY_NAME=$(echo "$PLUGIN_INFO" | jq -r '.name // empty')
PLUGIN_VERSION=$(echo "$PLUGIN_INFO" | jq -r '.version // empty')

# Use provided version or fall back to plugin.json
if [[ -z "$VERSION" ]]; then
    VERSION="$PLUGIN_VERSION"
fi

if [[ -z "$VERSION" ]]; then
    print_error "No version found in plugin.json or provided via --version"
    exit 1
fi

# Use plugin ID for binary names, fall back to directory name
if [[ -n "$PLUGIN_ID" ]]; then
    BINARY_NAME="$PLUGIN_ID"
else
    BINARY_NAME="$PLUGIN_NAME"
fi

print_info "Plugin ID: ${PLUGIN_ID:-$PLUGIN_NAME}"
print_info "Plugin Name: ${PLUGIN_DISPLAY_NAME:-$PLUGIN_NAME}"
print_info "Version: $VERSION"

# Set output directory
if [[ -z "$OUTPUT_DIR" ]]; then
    OUTPUT_DIR="releases/$VERSION"
fi

# Clean previous builds if requested
if [[ "$CLEAN" == true ]]; then
    print_step "Cleaning previous builds..."
    rm -rf releases/
    print_success "Cleaned previous builds"
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

print_step "Building $PLUGIN_DISPLAY_NAME $VERSION..."

# Download Go dependencies
print_info "Downloading Go dependencies..."
go mod tidy

# Build for specified platforms
IFS=',' read -ra PLATFORM_ARRAY <<< "$PLATFORMS"

print_step "Building binaries for ${#PLATFORM_ARRAY[@]} platforms..."

for platform in "${PLATFORM_ARRAY[@]}"; do
    if [[ ! "$platform" =~ ^[a-z]+/[a-z0-9]+$ ]]; then
        print_warning "Invalid platform format: $platform (expected: os/arch)"
        continue
    fi

    os=$(echo "$platform" | cut -d'/' -f1)
    arch=$(echo "$platform" | cut -d'/' -f2)

    binary_name="$BINARY_NAME-$os-$arch"
    if [[ "$os" == "windows" ]]; then
        binary_name="$binary_name.exe"
    fi

    print_info "Building for $os/$arch..."

    if GOOS="$os" GOARCH="$arch" go build -ldflags="-s -w" -o "$OUTPUT_DIR/$binary_name" .; then
        file_size=$(ls -lh "$OUTPUT_DIR/$binary_name" | awk '{print $5}')
        print_success "Built $binary_name ($file_size)"
    else
        print_error "Failed to build for $os/$arch"
        exit 1
    fi
done

# Build frontend component if it exists and requested
if [[ "$BUILD_FRONTEND" == true && -d "frontend" ]]; then
    print_step "Building frontend component..."

    cd frontend

    # Check for package.json
    if [[ -f "package.json" ]]; then
        # Check for build script
        if [[ -f "build-component.js" ]]; then
            print_info "Using custom build script..."
            if command -v node &> /dev/null; then
                node build-component.js
                if [[ -f "component.js" ]]; then
                    cp component.js "../$OUTPUT_DIR/"
                    file_size=$(ls -lh "../$OUTPUT_DIR/component.js" | awk '{print $5}')
                    print_success "Frontend component built: component.js ($file_size)"
                else
                    print_warning "Frontend build script ran but component.js not found"
                fi
            else
                print_error "Node.js not found. Cannot build frontend component."
                exit 1
            fi
        else
            # Standard npm build
            print_info "Installing frontend dependencies..."
            if [[ ! -d "node_modules" ]]; then
                npm install
            fi

            # Check for build scripts in package.json
            if npm run | grep -q "build:component"; then
                print_info "Running npm run build:component..."
                npm run build:component
            elif npm run | grep -q "build"; then
                print_info "Running npm run build..."
                npm run build
            else
                print_warning "No build script found in package.json"
            fi

            # Look for common output files
            for file in component.js dist/component.js build/component.js dist/index.js; do
                if [[ -f "$file" ]]; then
                    cp "$file" "../$OUTPUT_DIR/component.js"
                    file_size=$(ls -lh "../$OUTPUT_DIR/component.js" | awk '{print $5}')
                    print_success "Frontend component built: component.js ($file_size)"
                    break
                fi
            done
        fi
    else
        print_warning "No package.json found in frontend directory"
    fi

    cd ..
else
    if [[ "$BUILD_FRONTEND" == true ]]; then
        print_info "No frontend directory found, skipping frontend build"
    else
        print_info "Frontend build disabled"
    fi
fi

# Copy plugin metadata
if [[ -f "plugin.json" ]]; then
    cp plugin.json "$OUTPUT_DIR/"
    print_success "Copied plugin.json"
fi

# Copy additional files if they exist
for file in README.md LICENSE go.mod go.sum; do
    if [[ -f "$file" ]]; then
        cp "$file" "$OUTPUT_DIR/"
        print_info "Copied $file"
    fi
done

# Create release package if requested
if [[ "$CREATE_PACKAGE" == true ]]; then
    print_step "Creating release package..."

    cd releases

    # Create different archive formats based on content
    if [[ -f "$VERSION/component.js" ]]; then
        # Create tar.gz for plugins with frontend components
        ARCHIVE_NAME="$BINARY_NAME-$VERSION.tar.gz"
        tar -czf "$ARCHIVE_NAME" "$VERSION"
        archive_size=$(ls -lh "$ARCHIVE_NAME" | awk '{print $5}')
        print_success "Created release archive: $ARCHIVE_NAME ($archive_size)"
    fi

    # Also create a zip file for Windows compatibility
    if command -v zip &> /dev/null; then
        ARCHIVE_NAME="$BINARY_NAME-$VERSION.zip"
        zip -r "$ARCHIVE_NAME" "$VERSION" >/dev/null
        archive_size=$(ls -lh "$ARCHIVE_NAME" | awk '{print $5}')
        print_success "Created release archive: $ARCHIVE_NAME ($archive_size)"
    fi

    cd ..
fi

# Generate checksums
print_step "Generating checksums..."
cd "$OUTPUT_DIR"

CHECKSUM_FILE="checksums.txt"
echo "# Checksums for $BINARY_NAME $VERSION" > "$CHECKSUM_FILE"
echo "# Generated on $(date -u '+%Y-%m-%d %H:%M:%S UTC')" >> "$CHECKSUM_FILE"
echo "" >> "$CHECKSUM_FILE"

for file in *; do
    if [[ -f "$file" && "$file" != "$CHECKSUM_FILE" ]]; then
        if command -v sha256sum &> /dev/null; then
            checksum=$(sha256sum "$file" | cut -d' ' -f1)
        elif command -v shasum &> /dev/null; then
            checksum=$(shasum -a 256 "$file" | cut -d' ' -f1)
        else
            print_warning "No SHA256 tool found, skipping checksums"
            break
        fi
        echo "$checksum  $file" >> "$CHECKSUM_FILE"
    fi
done

if [[ -f "$CHECKSUM_FILE" && -s "$CHECKSUM_FILE" ]]; then
    print_success "Generated checksums: $CHECKSUM_FILE"
fi

cd - >/dev/null

# Display build summary
echo ""
print_success "🎉 Build completed successfully!"
echo ""
echo "📋 Build Summary:"
echo "├── Plugin: $PLUGIN_DISPLAY_NAME"
echo "├── Version: $VERSION"
echo "├── Platforms: ${#PLATFORM_ARRAY[@]}"

# List built files
echo "├── Files:"
for file in "$OUTPUT_DIR"/*; do
    if [[ -f "$file" ]]; then
        filename=$(basename "$file")
        file_size=$(ls -lh "$file" | awk '{print $5}')
        echo "│   ├── $filename ($file_size)"
    fi
done

echo "└── Output: $OUTPUT_DIR"
echo ""

# Show next steps
echo "🚀 Next steps:"
echo "1. Test the plugin binaries on target platforms"
echo "2. Update plugin documentation if needed"
echo "3. Commit and tag the release"
if [[ -d "$ROOT_DIR/../delve-registry" ]]; then
    echo "4. Run registry update to publish the plugin"
fi

# Return to original directory
cd "$SCRIPT_DIR"
