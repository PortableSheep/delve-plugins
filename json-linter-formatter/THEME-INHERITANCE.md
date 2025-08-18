# Theme Inheritance Guide for JSON Linter Plugin

## Overview

The JSON Linter & Formatter plugin has been updated to inherit themes from the parent Delve application using CSS custom properties (variables). This ensures the plugin seamlessly integrates with any theme the main app uses.

## How Theme Inheritance Works

### CSS Custom Properties
The plugin uses CSS custom properties with fallback values:

```css
.json-linter-container {
    background-color: var(--bg-color, #f8f9fa);
    color: var(--text-color, #343a40);
    border: 1px solid var(--border-color, #e9ecef);
}
```

### Theme Variables Used

| Variable | Usage | Fallback |
|----------|-------|----------|
| `--bg-color` | Main background | `#f8f9fa` |
| `--text-color` | Primary text color | `#343a40` |
| `--border-color` | Border color | `#e9ecef` |
| `--input-bg` | Input background | `#fff` |
| `--primary-color` | Focus/accent color | `#007bff` |
| `--primary-color-alpha` | Focus shadow | `rgba(0, 123, 255, 0.25)` |
| `--error-color` | Error borders | `#dc3545` |
| `--error-color-alpha` | Error shadow | `rgba(220, 53, 69, 0.25)` |
| `--error-bg` | Error background | `#f8d7da` |
| `--error-text` | Error text | `#721c24` |
| `--success-color` | Valid state border | `#28a745` |
| `--success-bg` | Valid state background | `#f8fff9` |
| `--stats-bg` | Statistics bar background | `#e9ecef` |
| `--muted-color` | Secondary text | `#6c757d` |
| `--muted-bg` | Muted backgrounds | `#f1f3f4` |
| `--placeholder-color` | Placeholder text | `#6c757d` |

### JSON Syntax Highlighting Colors

| Variable | Usage | Fallback |
|----------|-------|----------|
| `--json-string-color` | JSON strings | `#032f62` |
| `--json-number-color` | JSON numbers | `#005cc5` |
| `--json-boolean-color` | JSON booleans | `#d73a49` |
| `--json-null-color` | JSON null values | `#d73a49` |
| `--json-key-color` | JSON keys | `#22863a` |

## Setting Up Theme Integration in Main App

### Option 1: CSS Variables in Main App
Define theme variables in your main application's CSS:

```css
:root {
    /* Light theme */
    --bg-color: #ffffff;
    --text-color: #2d3748;
    --border-color: #e2e8f0;
    --input-bg: #ffffff;
    --primary-color: #4299e1;
    --error-color: #f56565;
    --success-color: #48bb78;
    --stats-bg: #f7fafc;
    --muted-color: #718096;
}

[data-theme="dark"] {
    /* Dark theme */
    --bg-color: #1a202c;
    --text-color: #f7fafc;
    --border-color: #4a5568;
    --input-bg: #2d3748;
    --primary-color: #63b3ed;
    --error-color: #fc8181;
    --success-color: #68d391;
    --stats-bg: #2d3748;
    --muted-color: #a0aec0;
}
```

### Option 2: JavaScript Theme Injection
Dynamically set theme variables:

```javascript
function applyTheme(theme) {
    const root = document.documentElement;
    
    if (theme === 'dark') {
        root.style.setProperty('--bg-color', '#1a202c');
        root.style.setProperty('--text-color', '#f7fafc');
        root.style.setProperty('--border-color', '#4a5568');
        // ... other properties
    } else {
        root.style.setProperty('--bg-color', '#ffffff');
        root.style.setProperty('--text-color', '#2d3748');
        root.style.setProperty('--border-color', '#e2e8f0');
        // ... other properties
    }
}
```

### Option 3: Theme Provider Component
Use a Vue theme provider:

```vue
<template>
    <div class="app" :data-theme="currentTheme">
        <!-- App content -->
        <PluginContainer />
    </div>
</template>

<script setup>
import { ref, watch } from 'vue';

const currentTheme = ref('light');

const themes = {
    light: {
        '--bg-color': '#ffffff',
        '--text-color': '#2d3748',
        '--border-color': '#e2e8f0',
        // ... other properties
    },
    dark: {
        '--bg-color': '#1a202c',
        '--text-color': '#f7fafc',
        '--border-color': '#4a5568',
        // ... other properties
    }
};

watch(currentTheme, (newTheme) => {
    const root = document.documentElement;
    const themeVars = themes[newTheme];
    
    Object.entries(themeVars).forEach(([property, value]) => {
        root.style.setProperty(property, value);
    });
});
</script>
```

## Custom Theme Examples

### GitHub Theme
```css
:root {
    --bg-color: #ffffff;
    --text-color: #24292f;
    --border-color: #d0d7de;
    --input-bg: #ffffff;
    --primary-color: #0969da;
    --error-color: #da3633;
    --success-color: #1a7f37;
    --json-string-color: #0a3069;
    --json-number-color: #0969da;
    --json-boolean-color: #8250df;
    --json-key-color: #1a7f37;
}
```

### VS Code Dark Theme
```css
:root {
    --bg-color: #1e1e1e;
    --text-color: #d4d4d4;
    --border-color: #3e3e3e;
    --input-bg: #252526;
    --primary-color: #007acc;
    --error-color: #f44747;
    --success-color: #4ec9b0;
    --json-string-color: #ce9178;
    --json-number-color: #b5cea8;
    --json-boolean-color: #569cd6;
    --json-key-color: #9cdcfe;
}
```

### Solarized Theme
```css
:root {
    --bg-color: #fdf6e3;
    --text-color: #657b83;
    --border-color: #eee8d5;
    --input-bg: #fdf6e3;
    --primary-color: #268bd2;
    --error-color: #dc322f;
    --success-color: #859900;
    --json-string-color: #2aa198;
    --json-number-color: #d33682;
    --json-boolean-color: #b58900;
    --json-key-color: #859900;
}
```

## Benefits of Theme Inheritance

1. **Consistent UI**: Plugin matches main app appearance
2. **Accessibility**: Respects user's theme preferences (dark mode, high contrast)
3. **Customization**: Easy to create custom themes
4. **Maintainability**: Theme changes in main app automatically apply to plugins
5. **Professional Appearance**: Seamless integration with any design system

## Testing Theme Integration

1. **Set theme variables** in main app
2. **Load the JSON Linter plugin**
3. **Verify colors match** main app theme
4. **Switch themes** and confirm plugin adapts
5. **Test both light and dark modes**

## Fallback Behavior

If no theme variables are provided, the plugin uses sensible defaults that work well in most contexts. This ensures the plugin is functional even without explicit theme setup.

## Advanced Customization

For advanced theme customization, you can override specific plugin styles:

```css
json-linter-formatter {
    --json-string-color: #your-custom-color;
    --error-color: #your-error-color;
}
```

This allows fine-tuning of plugin appearance while maintaining overall theme coherence.