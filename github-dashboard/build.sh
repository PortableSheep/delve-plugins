#!/bin/bash

# GitHub Dashboard Plugin Build Script
# Uses the generic plugin build system

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
    echo "   Make sure you're running this from the github-dashboard plugin directory"
    exit 1
fi

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}🐙 GitHub Dashboard Plugin Builder${NC}"
echo -e "${BLUE}Building GitHub Dashboard plugin...${NC}"
echo ""

# Call the generic build script with any passed arguments
# This allows for custom build options like:
# ./build.sh --clean
# ./build.sh --platforms darwin/amd64,linux/amd64
# ./build.sh --no-frontend
"$SCRIPTS_DIR/build-plugin.sh" "$@"

build_exit_code=$?

echo ""
if [[ $build_exit_code -eq 0 ]]; then
    echo -e "${GREEN}✅ GitHub Dashboard plugin built successfully!${NC}"
    echo ""
    echo "📋 Plugin is ready to use!"
    echo ""
    echo "🚀 Next steps:"
    echo "1. Configure your GitHub token in the plugin settings"
    echo "2. Add repositories to monitor (format: owner/repo)"
    echo "3. Set your preferred refresh interval (30-3600 seconds)"
    echo "4. Test the plugin in your Delve environment"
    echo ""
    echo "📖 For configuration help, see:"
    echo "   • CONFIGURATION-GUIDE.md"
    echo "   • README.md"
else
    echo "❌ Build failed. Check the output above for details."
    exit $build_exit_code
fi
