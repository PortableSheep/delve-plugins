#!/usr/bin/env node

/**
 * Generic Frontend Component Builder
 *
 * This script builds Vue.js components into standalone web components
 * that can be used in the Delve plugin system. It reads configuration
 * from plugin.json and package.json to determine build settings.
 */

const fs = require('fs');
const path = require('path');

// ANSI color codes for console output
const colors = {
    reset: '\x1b[0m',
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    cyan: '\x1b[36m'
};

function log(level, message) {
    const timestamp = new Date().toISOString().split('T')[1].slice(0, 8);
    const levelColors = {
        info: colors.blue,
        success: colors.green,
        warning: colors.yellow,
        error: colors.red
    };

    const color = levelColors[level] || colors.reset;
    const emoji = {
        info: 'ℹ️ ',
        success: '✅',
        warning: '⚠️ ',
        error: '❌'
    }[level] || '';

    console.log(`${color}${emoji} [${timestamp}] ${message}${colors.reset}`);
}

function readJsonFile(filePath) {
    try {
        const content = fs.readFileSync(filePath, 'utf-8');
        return JSON.parse(content);
    } catch (error) {
        if (error.code === 'ENOENT') {
            return null;
        }
        throw new Error(`Failed to parse ${filePath}: ${error.message}`);
    }
}

function findProjectRoot() {
    let currentDir = process.cwd();

    while (currentDir !== path.dirname(currentDir)) {
        if (fs.existsSync(path.join(currentDir, 'plugin.json'))) {
            return currentDir;
        }
        currentDir = path.dirname(currentDir);
    }

    return process.cwd();
}

function detectFramework(packageJson, pluginJson) {
    if (!packageJson || !packageJson.dependencies) {
        return 'vanilla';
    }

    const deps = { ...packageJson.dependencies, ...packageJson.devDependencies };

    if (deps.vue || deps['@vue/compiler-sfc']) {
        return 'vue';
    } else if (deps.react || deps['react-dom']) {
        return 'react';
    } else if (deps.svelte) {
        return 'svelte';
    } else {
        return 'vanilla';
    }
}

function buildVueComponent(config) {
    log('info', 'Building Vue.js component...');

    const VUE_CDN = 'https://unpkg.com/vue@3/dist/vue.global.js';

    // Read the main component file
    let componentPath = config.entry || 'src/component.vue';
    if (!fs.existsSync(componentPath)) {
        // Try alternative paths
        const alternatives = [
            'component.vue',
            'src/main.vue',
            'src/index.vue',
            'main.vue',
            'index.vue'
        ];

        componentPath = alternatives.find(p => fs.existsSync(p));
        if (!componentPath) {
            throw new Error('No Vue component file found. Expected component.vue or src/component.vue');
        }
    }

    log('info', `Reading component from: ${componentPath}`);
    const componentContent = fs.readFileSync(componentPath, 'utf-8');

    // Parse Vue SFC (Single File Component)
    const templateMatch = componentContent.match(/<template>([\s\S]*?)<\/template>/);
    const scriptMatch = componentContent.match(/<script>([\s\S]*?)<\/script>/);
    const styleMatch = componentContent.match(/<style[^>]*>([\s\S]*?)<\/style>/);

    const template = templateMatch ? templateMatch[1].trim() : '<div>No template found</div>';
    const script = scriptMatch ? scriptMatch[1].trim() : 'export default {}';
    const styles = styleMatch ? styleMatch[1].trim() : '';

    // Generate the component JavaScript
    const componentCode = `
// ${config.name} - Generated Vue Component
// Built on ${new Date().toISOString()}

(function() {
    'use strict';

    // Load Vue.js if not already loaded
    if (typeof Vue === 'undefined') {
        const script = document.createElement('script');
        script.src = '${VUE_CDN}';
        script.onload = initializeComponent;
        document.head.appendChild(script);
    } else {
        initializeComponent();
    }

    function initializeComponent() {
        // Component styles
        const styles = \`${styles}\`;
        if (styles.trim()) {
            const styleElement = document.createElement('style');
            styleElement.textContent = styles;
            document.head.appendChild(styleElement);
        }

        // Component definition
        const component = {
            template: \`${template}\`,
            ${script.replace(/export\s+default\s*/, '').replace(/^\{/, '').replace(/\}$/, '')}
        };

        // Register as custom element
        const componentName = '${config.customElementName || (config.id + '-component')}';

        if (customElements && !customElements.get(componentName)) {
            // Create Vue app and mount as custom element
            class ${config.className || config.id.replace(/[-_]/g, '')}Element extends HTMLElement {
                connectedCallback() {
                    const { createApp } = Vue;
                    this.app = createApp(component);
                    this.app.mount(this);
                }

                disconnectedCallback() {
                    if (this.app) {
                        this.app.unmount();
                    }
                }
            }

            customElements.define(componentName, ${config.className || config.id.replace(/[-_]/g, '')}Element);
        }

        // Also expose as global for direct use
        window.${config.globalName || config.id.replace(/[-_]/g, '')} = component;

        console.log('${config.name} component loaded successfully');
    }
})();
`;

    return componentCode;
}

function buildReactComponent(config) {
    log('info', 'Building React component...');

    // For React, we'd need a more complex build process
    // This is a simplified version
    throw new Error('React component building not yet implemented. Please use Vue.js or vanilla JavaScript.');
}

function buildVanillaComponent(config) {
    log('info', 'Building vanilla JavaScript component...');

    let componentPath = config.entry || 'src/component.js';
    if (!fs.existsSync(componentPath)) {
        const alternatives = ['component.js', 'src/main.js', 'src/index.js', 'main.js', 'index.js'];
        componentPath = alternatives.find(p => fs.existsSync(p));
        if (!componentPath) {
            throw new Error('No JavaScript component file found');
        }
    }

    const componentContent = fs.readFileSync(componentPath, 'utf-8');

    const wrappedComponent = `
// ${config.name} - Vanilla JavaScript Component
// Built on ${new Date().toISOString()}

(function() {
    'use strict';

    ${componentContent}

    console.log('${config.name} component loaded successfully');
})();
`;

    return wrappedComponent;
}

function main() {
    log('info', 'Starting frontend component build...');

    try {
        const projectRoot = findProjectRoot();
        const frontendDir = process.cwd();

        log('info', `Project root: ${projectRoot}`);
        log('info', `Frontend directory: ${frontendDir}`);

        // Read configuration files
        const pluginJson = readJsonFile(path.join(projectRoot, 'plugin.json'));
        const packageJson = readJsonFile(path.join(frontendDir, 'package.json'));

        if (!pluginJson) {
            throw new Error('plugin.json not found in project root');
        }

        // Extract configuration
        const config = {
            id: pluginJson.info?.id || path.basename(projectRoot),
            name: pluginJson.info?.name || 'Unknown Plugin',
            version: pluginJson.info?.version || '1.0.0',
            entry: pluginJson.frontend?.entry,
            customElementName: pluginJson.frontend?.customElement,
            globalName: pluginJson.frontend?.globalName,
            className: pluginJson.frontend?.className,
            outputFile: 'component.js'
        };

        log('info', `Building component: ${config.name} v${config.version}`);

        // Detect framework
        const framework = detectFramework(packageJson, pluginJson);
        log('info', `Detected framework: ${framework}`);

        // Build component based on framework
        let componentCode;

        switch (framework) {
            case 'vue':
                componentCode = buildVueComponent(config);
                break;
            case 'react':
                componentCode = buildReactComponent(config);
                break;
            case 'vanilla':
            default:
                componentCode = buildVanillaComponent(config);
                break;
        }

        // Write output file
        const outputPath = path.join(frontendDir, config.outputFile);
        fs.writeFileSync(outputPath, componentCode);

        const stats = fs.statSync(outputPath);
        const sizeKB = (stats.size / 1024).toFixed(2);

        log('success', `Component built successfully!`);
        log('info', `Output: ${outputPath} (${sizeKB} KB)`);

        // Validate the output
        if (stats.size === 0) {
            throw new Error('Generated component file is empty');
        }

        log('success', 'Frontend build completed successfully!');

    } catch (error) {
        log('error', `Build failed: ${error.message}`);
        process.exit(1);
    }
}

// Run if called directly
if (require.main === module) {
    main();
}

module.exports = { main, buildVueComponent, buildVanillaComponent };
