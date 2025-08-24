#!/bin/bash

# Plugin Migration Script
# Migrates existing plugins to use the new generic build system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

print_header() {
    echo -e "${CYAN}$1${NC}"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [plugin-name] [options]"
    echo ""
    echo "Arguments:"
    echo "  plugin-name           Name of plugin to migrate (optional - migrates all if not specified)"
    echo ""
    echo "Options:"
    echo "  --dry-run            Show what would be changed without making changes"
    echo "  --backup             Create backup of original files"
    echo "  --force              Overwrite existing files without prompting"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                           # Migrate all plugins"
    echo "  $0 github-dashboard          # Migrate specific plugin"
    echo "  $0 --dry-run                 # Show what would be migrated"
    echo "  $0 github-dashboard --backup # Migrate with backup"
}

# Default values
PLUGIN_NAME=""
DRY_RUN=false
BACKUP=false
FORCE=false

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --backup)
            BACKUP=true
            shift
            ;;
        --force)
            FORCE=true
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

# Function to backup a file
backup_file() {
    local file="$1"
    if [[ -f "$file" && "$BACKUP" == true ]]; then
        local backup_name="${file}.backup.$(date +%Y%m%d-%H%M%S)"
        cp "$file" "$backup_name"
        print_info "Backed up $file to $backup_name"
    fi
}

# Function to check if user wants to proceed
confirm_action() {
    local message="$1"
    if [[ "$FORCE" == true ]]; then
        return 0
    fi

    echo -n "$message (y/N): "
    read -r response
    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Function to validate plugin.json structure
validate_plugin_json() {
    local plugin_dir="$1"
    local plugin_json="$plugin_dir/plugin.json"

    if [[ ! -f "$plugin_json" ]]; then
        return 1
    fi

    # Check if jq is available
    if ! command -v jq &> /dev/null; then
        print_warning "jq not found. Cannot validate plugin.json structure."
        return 0
    fi

    # Check required fields
    local id=$(jq -r '.info.id // empty' "$plugin_json" 2>/dev/null)
    local name=$(jq -r '.info.name // empty' "$plugin_json" 2>/dev/null)
    local version=$(jq -r '.info.version // empty' "$plugin_json" 2>/dev/null)

    if [[ -z "$id" || -z "$name" || -z "$version" ]]; then
        print_warning "plugin.json missing required fields (info.id, info.name, info.version)"
        return 1
    fi

    return 0
}

# Function to update plugin.json if needed
update_plugin_json() {
    local plugin_dir="$1"
    local plugin_json="$plugin_dir/plugin.json"
    local plugin_name=$(basename "$plugin_dir")

    if [[ ! -f "$plugin_json" ]]; then
        print_warning "No plugin.json found in $plugin_dir"
        return 1
    fi

    # Check if jq is available
    if ! command -v jq &> /dev/null; then
        print_warning "jq not found. Cannot update plugin.json automatically."
        return 0
    fi

    local needs_update=false
    local temp_json=$(mktemp)

    # Read current plugin.json
    local current_json=$(cat "$plugin_json")

    # Check if runtime section exists
    if ! echo "$current_json" | jq -e '.runtime' >/dev/null 2>&1; then
        needs_update=true
        print_info "Adding missing runtime section to plugin.json"

        # Add runtime section
        current_json=$(echo "$current_json" | jq '. + {
            "runtime": {
                "executable": "'$plugin_name'",
                "frontend_entry": "component.js"
            }
        }')
    fi

    # Check if build_info section exists and update it
    local go_version=$(go version 2>/dev/null | awk '{print $3}' || echo "unknown")
    local commit_hash=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
    local build_time=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

    current_json=$(echo "$current_json" | jq '. + {
        "build_info": {
            "built_at": "'$build_time'",
            "go_version": "'$go_version'",
            "commit": "'$commit_hash'"
        }
    }')
    needs_update=true

    if [[ "$needs_update" == true ]]; then
        if [[ "$DRY_RUN" == true ]]; then
            print_info "[DRY RUN] Would update plugin.json in $plugin_dir"
        else
            backup_file "$plugin_json"
            echo "$current_json" | jq '.' > "$temp_json"
            mv "$temp_json" "$plugin_json"
            print_success "Updated plugin.json in $plugin_dir"
        fi
    fi

    rm -f "$temp_json"
    return 0
}

# Function to create new build.sh
create_new_build_script() {
    local plugin_dir="$1"
    local plugin_name=$(basename "$plugin_dir")
    local build_script="$plugin_dir/build.sh"

    # Check if build.sh already uses the new system
    if [[ -f "$build_script" ]] && grep -q "scripts/build-plugin.sh" "$build_script"; then
        print_info "$plugin_name already uses the new build system"
        return 0
    fi

    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would create/update build.sh for $plugin_name"
        return 0
    fi

    # Backup existing build.sh if it exists
    if [[ -f "$build_script" ]]; then
        backup_file "$build_script"
    fi

    # Create new build.sh based on plugin name
    local script_content=""

    case "$plugin_name" in
        "github-dashboard")
            script_content='#!/bin/bash

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
    echo "   Make sure you'\''re running this from the github-dashboard plugin directory"
    exit 1
fi

# Colors for output
GREEN='"'"'\033[0;32m'"'"'
BLUE='"'"'\033[0;34m'"'"'
CYAN='"'"'\033[0;36m'"'"'
NC='"'"'\033[0m'"'"'

echo -e "${CYAN}🐙 GitHub Dashboard Plugin Builder${NC}"
echo -e "${BLUE}Building GitHub Dashboard plugin...${NC}"
echo ""

# Call the generic build script with any passed arguments
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
fi'
            ;;
        "json-linter-formatter")
            script_content='#!/bin/bash

# JSON Linter & Formatter Plugin Build Script
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
    echo "   Make sure you'\''re running this from the json-linter-formatter plugin directory"
    exit 1
fi

# Colors for output
GREEN='"'"'\033[0;32m'"'"'
BLUE='"'"'\033[0;34m'"'"'
CYAN='"'"'\033[0;36m'"'"'
NC='"'"'\033[0m'"'"'

echo -e "${CYAN}🔧 JSON Linter & Formatter Plugin Builder${NC}"
echo -e "${BLUE}Building JSON processing plugin...${NC}"
echo ""

# Call the generic build script with any passed arguments
"$SCRIPTS_DIR/build-plugin.sh" "$@"

build_exit_code=$?

echo ""
if [[ $build_exit_code -eq 0 ]]; then
    echo -e "${GREEN}✅ JSON Linter & Formatter plugin built successfully!${NC}"
    echo ""
    echo "📋 Plugin is ready to use!"
    echo ""
    echo "🚀 Next steps:"
    echo "1. Test JSON linting and formatting functionality"
    echo "2. Configure validation rules if needed"
    echo "3. Test the plugin in your Delve environment"
else
    echo "❌ Build failed. Check the output above for details."
    exit $build_exit_code
fi'
            ;;
        *)
            # Generic template for other plugins
            script_content='#!/bin/bash

# '"$plugin_name"' Plugin Build Script
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
    echo "   Make sure you'\''re running this from a plugin directory"
    exit 1
fi

# Colors for output
GREEN='"'"'\033[0;32m'"'"'
BLUE='"'"'\033[0;34m'"'"'
NC='"'"'\033[0m'"'"'

echo -e "${BLUE}🔨 Building plugin: $PLUGIN_NAME${NC}"
echo ""

# Call the generic build script with any passed arguments
"$SCRIPTS_DIR/build-plugin.sh" "$@"

build_exit_code=$?

echo ""
if [[ $build_exit_code -eq 0 ]]; then
    echo -e "${GREEN}✅ Build completed for $PLUGIN_NAME${NC}"
    echo ""
    echo "📋 Plugin is ready to use!"
    echo ""
    echo "🚀 Next steps:"
    echo "1. Test the plugin binaries on your target platforms"
    echo "2. Configure the plugin settings as needed"
    echo "3. Deploy or distribute the plugin"
else
    echo "❌ Build failed. Check the output above for details."
    exit $build_exit_code
fi'
            ;;
    esac

    echo "$script_content" > "$build_script"
    chmod +x "$build_script"
    print_success "Created new build.sh for $plugin_name"
}

# Function to migrate a single plugin
migrate_plugin() {
    local plugin_dir="$1"
    local plugin_name=$(basename "$plugin_dir")

    print_header "Migrating plugin: $plugin_name"
    print_header "================================="

    # Check if it's a valid plugin directory
    if [[ ! -f "$plugin_dir/main.go" ]]; then
        print_warning "$plugin_name is not a valid plugin (no main.go found)"
        return 1
    fi

    if [[ ! -f "$plugin_dir/plugin.json" ]]; then
        print_warning "$plugin_name is missing plugin.json"
        return 1
    fi

    # Validate plugin.json structure
    if ! validate_plugin_json "$plugin_dir"; then
        print_warning "$plugin_name has invalid plugin.json structure"
        if [[ "$DRY_RUN" == false ]] && confirm_action "Continue with migration anyway?"; then
            # Continue with migration
            :
        else
            return 1
        fi
    fi

    # Update plugin.json if needed
    update_plugin_json "$plugin_dir"

    # Create new build script
    create_new_build_script "$plugin_dir"

    # Check for frontend directory and update if needed
    if [[ -d "$plugin_dir/frontend" ]]; then
        print_info "Frontend directory found in $plugin_name"

        # Check if frontend has custom build script
        if [[ -f "$plugin_dir/frontend/build-component.js" ]]; then
            print_info "Custom frontend build script found"
        elif [[ -f "$plugin_dir/frontend/package.json" ]]; then
            print_info "Standard npm-based frontend found"
        else
            print_warning "Frontend directory exists but no build configuration found"
        fi
    fi

    print_success "Migration completed for $plugin_name"
    echo ""
}

# Function to find all plugins
find_plugins() {
    local plugins=()

    for dir in "$ROOT_DIR"/*/; do
        if [[ -d "$dir" && -f "$dir/plugin.json" && -f "$dir/main.go" ]]; then
            plugin_name=$(basename "$dir")
            # Skip hidden directories and scripts directory
            if [[ ! "$plugin_name" =~ ^\. && "$plugin_name" != "scripts" ]]; then
                plugins+=("$plugin_name")
            fi
        fi
    done

    echo "${plugins[@]}"
}

# Main execution
main() {
    print_header "🔄 Delve Plugin Migration Tool"
    print_header "==============================="
    echo ""

    if [[ "$DRY_RUN" == true ]]; then
        print_info "Running in DRY RUN mode - no changes will be made"
        echo ""
    fi

    # Determine which plugins to migrate
    local plugins_to_migrate=()

    if [[ -n "$PLUGIN_NAME" ]]; then
        # Migrate specific plugin
        local plugin_dir="$ROOT_DIR/$PLUGIN_NAME"
        if [[ ! -d "$plugin_dir" ]]; then
            print_error "Plugin directory not found: $plugin_dir"
            exit 1
        fi
        plugins_to_migrate=("$PLUGIN_NAME")
    else
        # Migrate all plugins
        plugins_to_migrate=($(find_plugins))
    fi

    if [[ ${#plugins_to_migrate[@]} -eq 0 ]]; then
        print_error "No plugins found to migrate"
        exit 1
    fi

    print_info "Found ${#plugins_to_migrate[@]} plugin(s) to migrate:"
    for plugin in "${plugins_to_migrate[@]}"; do
        echo "  • $plugin"
    done
    echo ""

    if [[ "$DRY_RUN" == false ]] && [[ ${#plugins_to_migrate[@]} -gt 1 ]]; then
        if ! confirm_action "Proceed with migration of all plugins?"; then
            print_info "Migration cancelled"
            exit 0
        fi
        echo ""
    fi

    # Migrate each plugin
    local success_count=0
    local failed_count=0

    for plugin in "${plugins_to_migrate[@]}"; do
        if migrate_plugin "$ROOT_DIR/$plugin"; then
            ((success_count++))
        else
            ((failed_count++))
        fi
    done

    # Summary
    print_header "📊 Migration Summary"
    print_header "==================="
    echo ""
    echo "  • Total plugins: ${#plugins_to_migrate[@]}"
    echo "  • Successfully migrated: $success_count"
    echo "  • Failed: $failed_count"
    echo ""

    if [[ $failed_count -eq 0 ]]; then
        print_success "🎉 All plugins migrated successfully!"
        echo ""
        if [[ "$DRY_RUN" == false ]]; then
            print_info "🚀 Next steps:"
            echo "1. Test the new build scripts with existing plugins"
            echo "2. Run './scripts/build-all.sh' to test all plugins"
            echo "3. Update any custom build logic if needed"
            echo "4. Commit the migrated files"
        fi
    else
        print_warning "Some plugins failed to migrate. Check the output above for details."
        exit 1
    fi
}

# Check dependencies
if ! command -v jq &> /dev/null; then
    print_warning "jq not installed. Some plugin.json validation features will be limited."
    print_info "Install jq for full functionality: brew install jq (macOS) or apt-get install jq (Ubuntu)"
    echo ""
fi

# Run main function
main "$@"
