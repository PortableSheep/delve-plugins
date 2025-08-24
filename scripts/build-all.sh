#!/bin/bash

# Build All Plugins Script
# Builds all plugins in the delve-plugins repository

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
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -p, --platforms <list>  Comma-separated list of platforms (default: all)"
    echo "  -c, --clean            Clean previous builds for all plugins"
    echo "  --no-frontend          Skip frontend builds for all plugins"
    echo "  --no-package           Skip creating release packages"
    echo "  --parallel             Build plugins in parallel (faster but less readable output)"
    echo "  --continue-on-error    Continue building other plugins if one fails"
    echo "  -v, --verbose          Show detailed build output"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Build all plugins with default settings"
    echo "  $0 --clean                          # Clean and build all plugins"
    echo "  $0 -p darwin/amd64,linux/amd64     # Build for specific platforms only"
    echo "  $0 --parallel                       # Build plugins in parallel"
    echo "  $0 --continue-on-error             # Don't stop on individual plugin failures"
}

# Default values
PLATFORMS=""
CLEAN=false
BUILD_FRONTEND=true
CREATE_PACKAGE=true
PARALLEL=false
CONTINUE_ON_ERROR=false
VERBOSE=false

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--platforms)
            PLATFORMS="$2"
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
        --parallel)
            PARALLEL=true
            shift
            ;;
        --continue-on-error)
            CONTINUE_ON_ERROR=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
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
            print_error "Unexpected argument: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Find all plugin directories
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

# Build a single plugin
build_plugin() {
    local plugin_name="$1"
    local build_args=()

    # Prepare build arguments
    if [[ -n "$PLATFORMS" ]]; then
        build_args+=("--platforms" "$PLATFORMS")
    fi

    if [[ "$CLEAN" == true ]]; then
        build_args+=("--clean")
    fi

    if [[ "$BUILD_FRONTEND" == false ]]; then
        build_args+=("--no-frontend")
    fi

    if [[ "$CREATE_PACKAGE" == false ]]; then
        build_args+=("--no-package")
    fi

    local log_file=""
    if [[ "$PARALLEL" == true ]]; then
        log_file="/tmp/build-${plugin_name}.log"
        build_args+=(">" "$log_file" "2>&1")
    fi

    print_info "Building plugin: $plugin_name"

    local start_time=$(date +%s)

    if [[ "$VERBOSE" == true || "$PARALLEL" == false ]]; then
        "$SCRIPT_DIR/build-plugin.sh" "$plugin_name" "${build_args[@]}"
    else
        if "$SCRIPT_DIR/build-plugin.sh" "$plugin_name" "${build_args[@]}" > "/tmp/build-${plugin_name}.log" 2>&1; then
            local end_time=$(date +%s)
            local duration=$((end_time - start_time))
            print_success "Built $plugin_name in ${duration}s"
        else
            local end_time=$(date +%s)
            local duration=$((end_time - start_time))
            print_error "Failed to build $plugin_name after ${duration}s"

            if [[ -f "/tmp/build-${plugin_name}.log" ]]; then
                print_error "Build log for $plugin_name:"
                echo "----------------------------------------"
                cat "/tmp/build-${plugin_name}.log"
                echo "----------------------------------------"
            fi

            if [[ "$CONTINUE_ON_ERROR" == false ]]; then
                return 1
            fi
            return 0
        fi
    fi

    return 0
}

# Build plugins in parallel
build_plugins_parallel() {
    local plugins=("$@")
    local pids=()
    local failed_plugins=()

    print_info "Building ${#plugins[@]} plugins in parallel..."

    # Start builds
    for plugin in "${plugins[@]}"; do
        (
            if ! build_plugin "$plugin"; then
                echo "$plugin" > "/tmp/failed-${plugin}"
            fi
        ) &
        pids+=($!)
    done

    # Wait for all builds to complete
    for i in "${!pids[@]}"; do
        local pid=${pids[$i]}
        local plugin=${plugins[$i]}

        if wait "$pid"; then
            print_success "Completed: $plugin"
        else
            print_error "Failed: $plugin"
            failed_plugins+=("$plugin")
        fi
    done

    return ${#failed_plugins[@]}
}

# Build plugins sequentially
build_plugins_sequential() {
    local plugins=("$@")
    local failed_plugins=()

    print_info "Building ${#plugins[@]} plugins sequentially..."

    for plugin in "${plugins[@]}"; do
        if ! build_plugin "$plugin"; then
            failed_plugins+=("$plugin")
            if [[ "$CONTINUE_ON_ERROR" == false ]]; then
                break
            fi
        fi
    done

    return ${#failed_plugins[@]}
}

# Main execution
main() {
    print_header "🏗️  Delve Plugin Build System"
    print_header "================================"
    echo ""

    # Find all plugins
    plugins=($(find_plugins))

    if [[ ${#plugins[@]} -eq 0 ]]; then
        print_error "No plugins found in $ROOT_DIR"
        print_info "Plugins must have both plugin.json and main.go files"
        exit 1
    fi

    print_info "Found ${#plugins[@]} plugins:"
    for plugin in "${plugins[@]}"; do
        # Get plugin info if available
        if [[ -f "$ROOT_DIR/$plugin/plugin.json" ]]; then
            plugin_name=$(jq -r '.info.name // empty' "$ROOT_DIR/$plugin/plugin.json" 2>/dev/null)
            plugin_version=$(jq -r '.info.version // empty' "$ROOT_DIR/$plugin/plugin.json" 2>/dev/null)
            if [[ -n "$plugin_name" && -n "$plugin_version" ]]; then
                echo "  • $plugin ($plugin_name v$plugin_version)"
            else
                echo "  • $plugin"
            fi
        else
            echo "  • $plugin"
        fi
    done
    echo ""

    # Build configuration summary
    print_info "Build configuration:"
    echo "  • Platforms: ${PLATFORMS:-all supported}"
    echo "  • Clean builds: $([ "$CLEAN" == true ] && echo "yes" || echo "no")"
    echo "  • Frontend builds: $([ "$BUILD_FRONTEND" == true ] && echo "yes" || echo "no")"
    echo "  • Create packages: $([ "$CREATE_PACKAGE" == true ] && echo "yes" || echo "no")"
    echo "  • Parallel builds: $([ "$PARALLEL" == true ] && echo "yes" || echo "no")"
    echo "  • Continue on error: $([ "$CONTINUE_ON_ERROR" == true ] && echo "yes" || echo "no")"
    echo ""

    local start_time=$(date +%s)
    local failed_count=0

    # Build plugins
    if [[ "$PARALLEL" == true ]]; then
        build_plugins_parallel "${plugins[@]}"
        failed_count=$?
    else
        build_plugins_sequential "${plugins[@]}"
        failed_count=$?
    fi

    local end_time=$(date +%s)
    local total_duration=$((end_time - start_time))
    local success_count=$((${#plugins[@]} - failed_count))

    echo ""
    print_header "📊 Build Summary"
    print_header "================"
    echo ""
    echo "  • Total plugins: ${#plugins[@]}"
    echo "  • Successful builds: $success_count"
    echo "  • Failed builds: $failed_count"
    echo "  • Total time: ${total_duration}s"
    echo ""

    if [[ $failed_count -eq 0 ]]; then
        print_success "🎉 All plugins built successfully!"
        echo ""
        print_info "🚀 Next steps:"
        echo "  1. Test the plugin binaries on target platforms"
        echo "  2. Update plugin documentation if needed"
        echo "  3. Commit and tag releases"
        echo "  4. Update the plugin registry"
    else
        print_warning "Some plugins failed to build"
        exit 1
    fi

    # Cleanup temporary files
    rm -f /tmp/build-*.log /tmp/failed-*
}

# Trap to cleanup on exit
trap 'rm -f /tmp/build-*.log /tmp/failed-*' EXIT

# Check dependencies
if ! command -v jq &> /dev/null; then
    print_warning "jq not installed. Plugin information display will be limited."
fi

# Run main function
main "$@"
