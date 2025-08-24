#!/bin/bash

# Plugin Build Script Template
# Copy this file to your plugin directory as build.sh and customize as needed

set -e

# Get the directory of this script (plugin directory)
PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_NAME="$(basename "$PLUGIN_DIR")"

# Find the root scripts directory
SCRIPTS_DIR=""
if [[ -f "../scripts/build-plugin.sh" ]]; then
    SCRIPTS_DIR="../scripts"
elif [[ -f "../../scripts/build-plugin.sh" ]]; then
    SCRIPTS_DIR="../../scripts"
else
    echo "❌ Could not find scripts/build-plugin.sh"
    echo "   Make sure you're running this from a plugin directory"
    exit 1
fi

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔨 Building plugin: $PLUGIN_NAME${NC}"
echo ""

# Call the generic build script with current directory
"$SCRIPTS_DIR/build-plugin.sh" "$@"

echo ""
echo -e "${GREEN}✅ Build completed for $PLUGIN_NAME${NC}"
echo ""
echo "📋 Plugin is ready to use!"
echo ""
echo "🚀 Next steps:"
echo "1. Test the plugin binaries on your target platforms"
echo "2. Configure the plugin settings as needed"
echo "3. Deploy or distribute the plugin"

# Plugin-specific post-build tasks can be added here
# Examples:
# - Copy additional files
# - Run tests
# - Generate documentation
# - Update version files
