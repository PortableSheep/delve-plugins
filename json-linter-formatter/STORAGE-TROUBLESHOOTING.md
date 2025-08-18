# Storage Troubleshooting Guide for JSON Linter Plugin

## Overview

The JSON Linter & Formatter plugin may encounter storage timeout errors during startup. These are non-critical and don't affect the plugin's core functionality, but this guide explains what they are and how to resolve them.

## Common Error Messages

```
2025/08/13 12:29:13 Failed to store linter config: storage request timeout
2025/08/13 12:29:23 Failed to store JSON snippets: storage request timeout
```

## What These Errors Mean

### Storage Demo Operations
The plugin attempts to demonstrate the Delve SDK storage system by:
1. Storing sample configuration settings
2. Storing example JSON snippets
3. Retrieving storage statistics

These operations are **optional demonstrations** and don't affect the plugin's ability to:
- Validate and format JSON
- Maintain state across plugin switches
- Provide syntax highlighting
- Copy formatted JSON to clipboard

### Why Timeouts Occur

1. **Storage Service Unavailable**: The Delve storage service may not be fully initialized
2. **Network Issues**: Communication between plugin and host may be slow
3. **Resource Constraints**: System under load causing delayed responses
4. **Service Dependencies**: Storage service waiting for other components

## Solutions

### Option 1: Ignore the Warnings (Recommended)

The storage timeout errors are **non-critical**. The plugin will:
- Continue to function normally
- Use frontend localStorage for state persistence
- Provide all JSON linting and formatting features

**Action**: No action needed - the plugin works fine without backend storage.

### Option 2: Disable Storage Demo

To eliminate the warning messages entirely:

1. **Environment Variable Method**:
   ```bash
   DISABLE_STORAGE_DEMO=true ./json-linter-formatter
   ```

2. **Launch with Delve**:
   Set the environment variable before starting Delve:
   ```bash
   export DISABLE_STORAGE_DEMO=true
   ./delve
   ```

### Option 3: Wait for Storage Service

If you need the storage demo to work:

1. **Wait for Full Startup**: Allow 30-60 seconds after starting Delve before using plugins
2. **Check Storage Service**: Ensure Delve's storage service is running properly
3. **Restart if Needed**: Restart the Delve application if issues persist

## Plugin State Persistence

### Frontend State Management
The plugin uses a **dual-layer persistence system**:

1. **Primary**: Global window state (immediate persistence)
2. **Backup**: Browser localStorage (cross-session persistence)

This means your JSON content will persist even if backend storage fails.

### State Persistence Behavior

| Scenario | State Preservation | Note |
|----------|-------------------|------|
| Switch between plugins | ✅ Preserved | Via global window state |
| Refresh browser page | ✅ Preserved | Via localStorage |
| Restart Delve app | ✅ Preserved | Via localStorage |
| Backend storage timeout | ✅ Still works | Uses frontend persistence |

## Debugging Storage Issues

### Enable Verbose Logging

1. Check plugin startup logs for storage-related messages
2. Look for these log patterns:
   ```
   ✓ Stored JSON linter configuration
   ⚠️  Storage not available - linter config not stored
   ```

### Test Storage Functionality

1. **Check State Persistence**:
   - Add JSON content to the linter
   - Switch to another plugin
   - Switch back - content should be preserved

2. **Verify localStorage**:
   - Open browser dev tools
   - Check Application > Local Storage
   - Look for `json-linter-formatter-state` key

### Network Diagnostics

If storage issues persist:

1. **Check WebSocket Connection**: Ensure plugin can communicate with host
2. **Monitor Network Tab**: Look for failed requests in browser dev tools
3. **Check Delve Logs**: Review main application logs for storage service issues

## Performance Impact

### Storage Timeouts
- **Impact**: None on plugin functionality
- **Duration**: Timeouts typically occur within 10-30 seconds
- **Recovery**: Plugin continues normally after timeout

### Resource Usage
- **Memory**: Frontend state uses minimal browser memory
- **Storage**: localStorage uses small amount of disk space
- **CPU**: Background storage operations don't affect UI performance

## Best Practices

### For Users
1. **Ignore Storage Warnings**: Focus on plugin functionality
2. **Use Frontend Features**: Rely on browser-based state persistence
3. **Regular Browser Cleanup**: Clear localStorage occasionally if needed

### For Developers
1. **Test Without Storage**: Ensure plugins work without backend storage
2. **Implement Fallbacks**: Always provide alternative persistence methods
3. **Handle Timeouts Gracefully**: Don't block plugin startup on storage operations

## FAQ

### Q: Will my JSON content be lost?
**A**: No, the plugin uses frontend persistence which is independent of backend storage.

### Q: Can I disable the error messages?
**A**: Yes, set `DISABLE_STORAGE_DEMO=true` environment variable.

### Q: Does this affect other plugins?
**A**: No, each plugin manages its own storage independently.

### Q: Is this a bug?
**A**: No, this is expected behavior when storage services are slow to initialize.

### Q: Should I report these errors?
**A**: Only if they persist for extended periods (>5 minutes) or affect plugin functionality.

## Getting Help

If storage issues significantly impact your workflow:

1. **Check Delve Documentation**: Review main application storage documentation
2. **File an Issue**: Report persistent storage problems to the Delve project
3. **Community Support**: Ask in Delve community channels
4. **Workaround**: Use `DISABLE_STORAGE_DEMO=true` to eliminate warnings

## Related Files

- `main.go`: Contains storage demonstration code
- `JSONLinterComponent.vue`: Frontend state management
- `delve-sdk.js`: Plugin SDK integration
- `KEEPALIVE-FIX.md`: State persistence implementation details