# Theme Inheritance Guide for GitHub Dashboard Plugin

## Overview

The GitHub Dashboard plugin has been updated to inherit themes from the parent Delve application using CSS custom properties (variables). This ensures the plugin seamlessly integrates with any theme the main app uses.

## How Theme Inheritance Works

### CSS Custom Properties
The plugin uses CSS custom properties with fallback values:

```css
.github-dashboard {
    background-color: var(--bg-color, #ffffff);
    color: var(--text-color, #24292f);
    border: 1px solid var(--border-color, #d0d7de);
}
```

### Theme Variables Used

| Variable | Usage | Fallback |
|----------|-------|----------|
| `--bg-color` | Main background | `#ffffff` |
| `--text-color` | Primary text color | `#24292f` |
| `--border-color` | Border color | `#d0d7de` |
| `--card-bg` | Card background | `#ffffff` |
| `--primary-color` | Accent/focus color | `#0969da` |
| `--primary-color-hover` | Primary hover state | `#0860ca` |
| `--primary-bg` | Primary background | `#f6f8ff` |
| `--error-color` | Error text | `#cf222e` |
| `--error-bg` | Error background | `#fff8f6` |
| `--error-border` | Error border | `#ffcecb` |
| `--success-color` | Success text | `#166534` |
| `--success-bg` | Success background | `#dcfce7` |
| `--info-color` | Info text | `#0969da` |
| `--info-bg` | Info background | `#ddf4ff` |
| `--info-border` | Info border | `#54aeff` |
| `--info-text` | Info content text | `#0969da` |
| `--muted-color` | Secondary text | `#656d76` |
| `--muted-bg` | Muted backgrounds | `#f6f8fa` |
| `--shadow-color` | Box shadows | `rgba(0, 0, 0, 0.1)` |

### Spacing Variables

| Variable | Usage | Fallback |
|----------|-------|----------|
| `--spacing-xs` | Extra small spacing | `0.25rem` |
| `--spacing-sm` | Small spacing | `0.5rem` |
| `--spacing-md` | Medium spacing | `1rem` |
| `--spacing-lg` | Large spacing | `1.5rem` |
| `--spacing-xl` | Extra large spacing | `2rem` |

### Typography Variables

| Variable | Usage | Fallback |
|----------|-------|----------|
| `--font-family` | Main font family | `-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif` |
| `--font-size-xs` | Extra small text | `0.75rem` |
| `--font-size-sm` | Small text | `0.875rem` |
| `--font-size-md` | Medium text | `1rem` |
| `--font-size-lg` | Large text | `1.25rem` |
| `--font-size-xl` | Extra large text | `1.5rem` |

### Border Radius Variables

| Variable | Usage | Fallback |
|----------|-------|----------|
| `--border-radius` | Standard border radius | `6px` |
| `--border-radius-sm` | Small border radius | `4px` |

## Setting Up Theme Integration in Main App

### Option 1: CSS Variables in Main App
Define theme variables in your main application's CSS:

```css
:root {
    /* Light theme */
    --bg-color: #ffffff;
    --text-color: #24292f;
    --border-color: #d0d7de;
    --card-bg: #ffffff;
    --primary-color: #0969da;
    --primary-color-hover: #0860ca;
    --error-color: #cf222e;
    --success-color: #166534;
    --muted-color: #656d76;
    --muted-bg: #f6f8fa;
}

[data-theme="dark"] {
    /* Dark theme */
    --bg-color: #0d1117;
    --text-color: #f0f6fc;
    --border-color: #30363d;
    --card-bg: #161b22;
    --primary-color: #58a6ff;
    --primary-color-hover: #4493f8;
    --error-color: #f85149;
    --success-color: #56d364;
    --muted-color: #8b949e;
    --muted-bg: #21262d;
}
```

### Option 2: JavaScript Theme Injection
Dynamically set theme variables:

```javascript
function applyTheme(theme) {
    const root = document.documentElement;
    
    if (theme === 'dark') {
        root.style.setProperty('--bg-color', '#0d1117');
        root.style.setProperty('--text-color', '#f0f6fc');
        root.style.setProperty('--border-color', '#30363d');
        root.style.setProperty('--card-bg', '#161b22');
        root.style.setProperty('--primary-color', '#58a6ff');
        root.style.setProperty('--error-color', '#f85149');
        root.style.setProperty('--success-color', '#56d364');
        root.style.setProperty('--muted-color', '#8b949e');
        root.style.setProperty('--muted-bg', '#21262d');
    } else {
        root.style.setProperty('--bg-color', '#ffffff');
        root.style.setProperty('--text-color', '#24292f');
        root.style.setProperty('--border-color', '#d0d7de);
        root.style.setProperty('--card-bg', '#ffffff');
        root.style.setProperty('--primary-color', '#0969da');
        root.style.setProperty('--error-color', '#cf222e');
        root.style.setProperty('--success-color', '#166534');
        root.style.setProperty('--muted-color', '#656d76');
        root.style.setProperty('--muted-bg', '#f6f8fa');
    }
}
```

### Option 3: Vue Theme Provider Component
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
        '--text-color': '#24292f',
        '--border-color': '#d0d7de',
        '--primary-color': '#0969da',
        '--error-color': '#cf222e',
        '--success-color': '#166534',
        '--muted-color': '#656d76',
        '--muted-bg': '#f6f8fa',
    },
    dark: {
        '--bg-color': '#0d1117',
        '--text-color': '#f0f6fc',
        '--border-color': '#30363d',
        '--primary-color': '#58a6ff',
        '--error-color': '#f85149',
        '--success-color': '#56d364',
        '--muted-color': '#8b949e',
        '--muted-bg': '#21262d',
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
    --card-bg: #ffffff;
    --primary-color: #0969da;
    --error-color: #cf222e;
    --success-color: #1a7f37;
    --muted-color: #656d76;
    --muted-bg: #f6f8fa;
}
```

### GitHub Dark Theme
```css
:root {
    --bg-color: #0d1117;
    --text-color: #f0f6fc;
    --border-color: #30363d;
    --card-bg: #161b22;
    --primary-color: #58a6ff;
    --error-color: #f85149;
    --success-color: #56d364;
    --muted-color: #8b949e;
    --muted-bg: #21262d;
}
```

### VS Code Dark Theme
```css
:root {
    --bg-color: #1e1e1e;
    --text-color: #d4d4d4;
    --border-color: #3e3e3e;
    --card-bg: #252526;
    --primary-color: #007acc;
    --error-color: #f44747;
    --success-color: #4ec9b0;
    --muted-color: #969696;
    --muted-bg: #2d2d30;
}
```

### Solarized Light Theme
```css
:root {
    --bg-color: #fdf6e3;
    --text-color: #657b83;
    --border-color: #eee8d5;
    --card-bg: #fdf6e3;
    --primary-color: #268bd2;
    --error-color: #dc322f;
    --success-color: #859900;
    --muted-color: #93a1a1;
    --muted-bg: #eee8d5;
}
```

## Plugin-Specific Components

### Repository Cards
Repository cards use these specific variables:
- `--card-bg`: Card background color
- `--border-color`: Card border color
- `--primary-bg`: Active card background
- `--shadow-color`: Hover shadow effect

### Pull Request Cards
Pull request cards use:
- `--card-bg`: Card background
- `--success-bg` and `--success-color`: Open PR status
- `--muted-color`: PR metadata text

### Error and Info Messages
- `--error-bg`, `--error-color`, `--error-border`: Error messages
- `--info-bg`, `--info-color`, `--info-border`: Setup notices

## Benefits of Theme Inheritance

1. **Consistent UI**: Plugin matches main app appearance automatically
2. **Accessibility**: Respects user's theme preferences (dark mode, high contrast)
3. **Customization**: Easy to create custom themes that apply everywhere
4. **Maintainability**: Theme changes in main app automatically apply to plugins
5. **Professional Appearance**: Seamless integration with any design system
6. **GitHub-like Experience**: Default styling matches GitHub's interface

## Testing Theme Integration

1. **Set theme variables** in main app
2. **Load the GitHub Dashboard plugin**
3. **Verify colors match** main app theme
4. **Switch themes** and confirm plugin adapts
5. **Test both light and dark modes**
6. **Check responsive behavior** on different screen sizes

## Fallback Behavior

If no theme variables are provided, the plugin uses GitHub-inspired defaults that work well in most contexts. This ensures the plugin is functional and attractive even without explicit theme setup.

## Advanced Customization

For advanced theme customization, you can override specific plugin styles:

```css
github-dashboard {
    --primary-color: #your-custom-color;
    --error-color: #your-error-color;
    --card-bg: #your-card-background;
}
```

This allows fine-tuning of plugin appearance while maintaining overall theme coherence.

## Mobile Responsiveness

The plugin includes responsive design that adapts to mobile screens:
- Stacked layout on narrow screens
- Adjusted spacing and typography
- Touch-friendly interactive elements
- Optimized for both portrait and landscape orientations

All responsive features respect the theme variables, ensuring consistent appearance across all device sizes.