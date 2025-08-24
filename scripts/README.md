# Delve Plugin Build System

This directory contains generic build scripts that can be used by any plugin in the delve-plugins repository. The build system automates the process of building Go binaries, frontend components, and creating release packages.

## Scripts Overview

### `build-plugin.sh`
The main generic build script that can build any plugin with proper structure.

**Usage:**
```bash
# Build current directory plugin
./scripts/build-plugin.sh

# Build specific plugin
./scripts/build-plugin.sh github-dashboard

# Build with custom options
./scripts/build-plugin.sh github-dashboard --clean --platforms darwin/amd64,linux/amd64
```

**Options:**
- `-v, --version <ver>`: Override version from plugin.json
- `-p, --platforms <list>`: Comma-separated list of platforms
- `-o, --output <dir>`: Custom output directory
- `-c, --clean`: Clean previous builds
- `--no-frontend`: Skip frontend build
- `--no-package`: Skip creating release packages

### `build-frontend.js`
Generic frontend component builder that supports multiple frameworks.

**Features:**
- Auto-detects framework (Vue.js, React, vanilla JS)
- Reads configuration from plugin.json
- Generates standalone web components
- Supports custom build configurations

### `build-all.sh`
Builds all plugins in the repository at once.

**Usage:**
```bash
# Build all plugins
./scripts/build-all.sh

# Build all plugins in parallel
./scripts/build-all.sh --parallel

# Build for specific platforms only
./scripts/build-all.sh --platforms darwin/amd64,linux/amd64
```

### `plugin-build-template.sh`
Template for individual plugin build scripts. Copy this to your plugin directory as `build.sh`.

## Plugin Requirements

For a plugin to work with the generic build system, it must have:

1. **plugin.json** - Plugin metadata file
2. **main.go** - Go source code entry point
3. **go.mod** - Go module file

### Required plugin.json Structure

```json
{
  "info": {
    "id": "my-plugin",
    "name": "My Plugin",
    "version": "v1.0.0",
    "description": "Description of my plugin",
    "author": "Your Name",
    "license": "MIT"
  },
  "runtime": {
    "executable": "my-plugin",
    "frontend_entry": "component.js"
  },
  "frontend": {
    "entry": "src/component.vue",
    "customElement": "my-plugin-component",
    "globalName": "MyPlugin"
  }
}
```

### Optional Frontend Structure

If your plugin has a frontend component:

```
my-plugin/
├── frontend/
│   ├── package.json
│   ├── src/
│   │   └── component.vue
│   └── build-component.js (optional custom builder)
```

## Platform Support

The build system supports building for multiple platforms:

- `darwin/amd64` - macOS Intel
- `darwin/arm64` - macOS Apple Silicon
- `linux/amd64` - Linux 64-bit
- `linux/arm64` - Linux ARM64
- `windows/amd64` - Windows 64-bit

## Output Structure

After building, each plugin will have:

```
my-plugin/
├── releases/
│   └── v1.0.0/
│       ├── my-plugin-darwin-amd64
│       ├── my-plugin-darwin-arm64
│       ├── my-plugin-linux-amd64
│       ├── my-plugin-linux-arm64
│       ├── my-plugin-windows-amd64.exe
│       ├── component.js (if frontend exists)
│       ├── plugin.json
│       ├── checksums.txt
│       └── other files...
└── my-plugin-v1.0.0.tar.gz (release package)
```

## Setting Up a New Plugin

1. **Create plugin directory structure:**
   ```bash
   mkdir my-plugin
   cd my-plugin
   ```

2. **Create plugin.json:**
   ```bash
   cp ../scripts/plugin-build-template.sh build.sh
   # Edit plugin.json with your plugin details
   ```

3. **Initialize Go module:**
   ```bash
   go mod init github.com/PortableSheep/delve-plugins/my-plugin
   ```

4. **Create main.go:**
   ```go
   package main

   import "fmt"

   func main() {
       fmt.Println("Hello from my plugin!")
   }
   ```

5. **Build the plugin:**
   ```bash
   ./build.sh
   ```

## Frontend Development

### Vue.js Components

Create a Vue single-file component:

```vue
<template>
  <div class="my-plugin">
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
  </div>
</template>

<script>
export default {
  name: 'MyPlugin',
  data() {
    return {
      title: 'My Plugin',
      description: 'This is my plugin'
    }
  }
}
</script>

<style scoped>
.my-plugin {
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
</style>
```

### Custom Build Scripts

For complex frontend builds, create a custom `frontend/build-component.js`:

```javascript
const fs = require('fs');

// Your custom build logic here
const componentCode = `
// Your generated component code
`;

fs.writeFileSync('component.js', componentCode);
console.log('Custom component built successfully');
```

## Troubleshooting

### Common Issues

1. **"jq not found"**: Install jq for JSON parsing
   ```bash
   # macOS
   brew install jq
   
   # Ubuntu/Debian
   sudo apt-get install jq
   ```

2. **"Go not found"**: Install Go from https://golang.org/

3. **Frontend build fails**: Ensure Node.js and npm are installed

4. **Permission denied**: Make build scripts executable
   ```bash
   chmod +x build.sh
   chmod +x scripts/*.sh
   ```

### Debug Mode

Run builds with verbose output:
```bash
./scripts/build-plugin.sh my-plugin --verbose
```

### Clean Builds

Remove all build artifacts:
```bash
./scripts/build-plugin.sh my-plugin --clean
```

## Contributing

When adding new features to the build system:

1. Test with existing plugins
2. Update this documentation
3. Add examples for new features
4. Ensure backward compatibility

## Dependencies

Required tools:
- **Go 1.21+**: For building plugin binaries
- **jq**: For parsing JSON files
- **Node.js** (optional): For frontend builds
- **tar/zip**: For creating release packages

Optional tools:
- **shasum/sha256sum**: For generating checksums
- **npm**: For Node.js dependency management

## Examples

### Basic Plugin Build
```bash
# Build github-dashboard plugin
./scripts/build-plugin.sh github-dashboard
```

### Custom Platform Build
```bash
# Build only for macOS
./scripts/build-plugin.sh github-dashboard --platforms darwin/amd64,darwin/arm64
```

### Clean Release Build
```bash
# Clean build with custom version
./scripts/build-plugin.sh github-dashboard --clean --version v2.0.0
```

### Build All Plugins
```bash
# Build all plugins in parallel
./scripts/build-all.sh --parallel --clean
```
