#!/usr/bin/env node

import { readFileSync, writeFileSync, existsSync, mkdirSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { build } from "vite";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Configuration
const config = {
  input: join(__dirname, "src", "main.js"),
  output: join(__dirname, "component.js"),
  distDir: join(__dirname, "dist"),
  tempDir: join(__dirname, ".temp"),
  vueVersion: "3.4.21",
};

const VUE_CDN = `https://unpkg.com/vue@${config.vueVersion}/dist/vue.global.js`;

console.log("🔨 Building GitHub Dashboard Web Component for Distribution...");

// Ensure directories exist
if (!existsSync(config.distDir)) {
  mkdirSync(config.distDir, { recursive: true });
}

// Build with Vite first
console.log("📦 Building with Vite...");

try {
  await build({
    configFile: join(__dirname, "vite.config.js"),
    build: {
      outDir: config.distDir,
      emptyOutDir: true,
      lib: {
        entry: config.input,
        name: "GitHubDashboard",
        fileName: "component",
        formats: ["iife"],
      },
      rollupOptions: {
        output: {
          format: "iife",
          inlineDynamicImports: true,
          manualChunks: undefined,
          assetFileNames: "assets/[name].[ext]",
          chunkFileNames: "[name].js",
          entryFileNames: "component.js",
        },
      },
      minify: "terser",
      terserOptions: {
        compress: {
          drop_console: false,
          drop_debugger: true,
        },
        mangle: {
          keep_classnames: true,
          keep_fnames: true,
        },
      },
      target: "es2015",
      cssCodeSplit: false,
    },
    define: {
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: false,
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
    },
  });

  console.log("✅ Vite build completed");
} catch (error) {
  console.error("❌ Vite build failed:", error);
  process.exit(1);
}

// Read the built file
const builtFilePath = join(config.distDir, "component.js");
let builtContent = "";

try {
  builtContent = readFileSync(builtFilePath, "utf-8");
  console.log("📖 Read built component file");
} catch (error) {
  console.error("❌ Failed to read built file:", error);
  process.exit(1);
}

// Create the final standalone component
const standaloneComponent = `/*!
 * GitHub Dashboard Web Component v1.0.0
 * Delve Plugin for monitoring GitHub repositories and pull requests
 *
 * Built: ${new Date().toISOString()}
 * Vue Version: ${config.vueVersion}
 *
 * Copyright (c) 2024 Michael Gunderson
 * Licensed under MIT License
 */

(function (global) {
  'use strict';

  // Component state management
  let isVueLoaded = false;
  let isComponentRegistered = false;
  let pendingInitializations = [];

  // Error handling
  function handleError(error, context) {
    console.error(\`[GitHub Dashboard] Error in \${context}:\`, error);

    // Create error display element
    const errorElement = document.createElement('div');
    errorElement.innerHTML = \`
      <div style="
        padding: 20px;
        background: #fff8f8;
        border: 1px solid #ffcdd2;
        border-radius: 6px;
        color: #d1242f;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
        margin: 10px 0;
      ">
        <div style="font-weight: 600; margin-bottom: 8px;">
          ⚠️ GitHub Dashboard Error
        </div>
        <div style="font-size: 14px; margin-bottom: 8px;">
          Failed to initialize component: \${context}
        </div>
        <div style="font-size: 12px; color: #636c76;">
          \${error.message || 'Unknown error occurred'}
        </div>
      </div>
    \`;
    return errorElement;
  }

  // Vue loader with fallback
  function loadVue() {
    return new Promise((resolve, reject) => {
      // Check if Vue is already available
      if (typeof global.Vue !== 'undefined') {
        isVueLoaded = true;
        console.log('[GitHub Dashboard] Vue already available');
        resolve(global.Vue);
        return;
      }

      // Load Vue from CDN
      const script = document.createElement('script');
      script.src = '${VUE_CDN}';
      script.async = true;
      script.crossOrigin = 'anonymous';

      script.onload = () => {
        if (typeof global.Vue !== 'undefined') {
          isVueLoaded = true;
          console.log('[GitHub Dashboard] Vue loaded from CDN');
          resolve(global.Vue);
        } else {
          reject(new Error('Vue failed to load properly'));
        }
      };

      script.onerror = (error) => {
        console.error('[GitHub Dashboard] Failed to load Vue from CDN:', error);
        reject(new Error('Failed to load Vue from CDN'));
      };

      // Timeout fallback
      setTimeout(() => {
        if (!isVueLoaded) {
          reject(new Error('Vue loading timeout'));
        }
      }, 10000); // 10 second timeout

      document.head.appendChild(script);
    });
  }

  // Component registration
  function registerComponent() {
    try {
      if (isComponentRegistered) {
        console.log('[GitHub Dashboard] Component already registered');
        return Promise.resolve();
      }

      // Built component code
      ${builtContent}

      // Mark as registered
      isComponentRegistered = true;
      console.log('[GitHub Dashboard] Component registered successfully');

      // Process any pending initializations
      pendingInitializations.forEach(callback => {
        try {
          callback();
        } catch (error) {
          console.error('[GitHub Dashboard] Pending initialization failed:', error);
        }
      });
      pendingInitializations = [];

      return Promise.resolve();
    } catch (error) {
      console.error('[GitHub Dashboard] Component registration failed:', error);
      return Promise.reject(error);
    }
  }

  // Initialize the component
  async function initialize() {
    try {
      console.log('[GitHub Dashboard] Initializing web component...');

      // Load Vue if needed
      if (!isVueLoaded) {
        await loadVue();
      }

      // Register the component
      await registerComponent();

      console.log('[GitHub Dashboard] Component initialization complete');
    } catch (error) {
      handleError(error, 'initialization');
      throw error;
    }
  }

  // Plugin API setup
  function setupPluginAPI() {
    if (!global.pluginAPI) {
      global.pluginAPI = {};
    }

    // Default implementations for development/testing
    if (!global.pluginAPI.getRepositories) {
      global.pluginAPI.getRepositories = async () => {
        console.log('[GitHub Dashboard] Using mock repository data');
        return [
          {
            name: 'delve',
            full_name: 'PortableSheep/delve',
            description: 'Advanced plugin system for Go applications - Configure your GitHub token for real data',
            stargazers_count: 1250,
            forks_count: 180,
            open_issues_count: 23,
            language: 'Go',
            updated_at: new Date(Date.now() - 86400000).toISOString(),
            html_url: 'https://github.com/PortableSheep/delve',
          },
          {
            name: 'delve-plugins',
            full_name: 'PortableSheep/delve-plugins',
            description: 'Official plugin collection for Delve',
            stargazers_count: 89,
            forks_count: 34,
            open_issues_count: 8,
            language: 'Go',
            updated_at: new Date(Date.now() - 3600000).toISOString(),
            html_url: 'https://github.com/PortableSheep/delve-plugins',
          }
        ];
      };
    }

    if (!global.pluginAPI.getPullRequests) {
      global.pluginAPI.getPullRequests = async (repoName) => {
        console.log(\`[GitHub Dashboard] Using mock PR data for \${repoName}\`);
        return [
          {
            number: 42,
            title: 'Add enhanced GitHub dashboard features',
            state: 'open',
            created_at: new Date(Date.now() - 86400000).toISOString(),
            updated_at: new Date(Date.now() - 3600000).toISOString(),
            html_url: \`https://github.com/\${repoName}/pull/42\`,
            user: {
              login: 'developer',
              avatar_url: 'https://github.com/identicons/developer.png',
            },
            draft: false,
          },
          {
            number: 38,
            title: 'Fix responsive design issues',
            state: 'open',
            created_at: new Date(Date.now() - 172800000).toISOString(),
            updated_at: new Date(Date.now() - 7200000).toISOString(),
            html_url: \`https://github.com/\${repoName}/pull/38\`,
            user: {
              login: 'designer',
              avatar_url: 'https://github.com/identicons/designer.png',
            },
            draft: true,
          }
        ];
      };
    }
  }

  // DOM ready handler
  function onDOMReady(callback) {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', callback);
    } else {
      callback();
    }
  }

  // Main entry point
  function main() {
    setupPluginAPI();

    onDOMReady(async () => {
      try {
        await initialize();

        // Auto-initialize any existing elements
        const existingElements = document.querySelectorAll('github-dashboard');
        if (existingElements.length > 0) {
          console.log(\`[GitHub Dashboard] Found \${existingElements.length} existing elements\`);
        }
      } catch (error) {
        console.error('[GitHub Dashboard] Failed to initialize:', error);
      }
    });
  }

  // Export for module systems
  if (typeof module !== 'undefined' && module.exports) {
    module.exports = { initialize, setupPluginAPI };
  } else if (typeof define === 'function' && define.amd) {
    define([], () => ({ initialize, setupPluginAPI }));
  } else {
    // Browser global
    global.GitHubDashboard = { initialize, setupPluginAPI };
  }

  // Auto-initialize in browser environment
  if (typeof window !== 'undefined') {
    main();
  }

})(typeof window !== 'undefined' ? window : typeof global !== 'undefined' ? global : this);
`;

// Write the final component file
try {
  writeFileSync(config.output, standaloneComponent, "utf-8");
  console.log(`✅ Standalone component created: ${config.output}`);
} catch (error) {
  console.error("❌ Failed to write component file:", error);
  process.exit(1);
}

// Create a minified version for production
console.log("🗜️  Creating minified version...");

try {
  const { minify } = await import("terser");
  const minified = await minify(standaloneComponent, {
    compress: {
      drop_console: false,
      drop_debugger: true,
      pure_funcs: ["console.log"],
    },
    mangle: {
      keep_classnames: true,
      keep_fnames: true,
    },
    format: {
      comments: /^!/,
    },
  });

  const minifiedPath = config.output.replace(".js", ".min.js");
  writeFileSync(minifiedPath, minified.code, "utf-8");
  console.log(`✅ Minified component created: ${minifiedPath}`);
} catch (error) {
  console.warn("⚠️  Failed to create minified version:", error.message);
}

// Generate component stats
const stats = {
  buildTime: new Date().toISOString(),
  vueVersion: config.vueVersion,
  outputSize: Math.round(standaloneComponent.length / 1024),
  features: [
    "Vue 3 Web Component",
    "Shadow DOM encapsulation",
    "GitHub API integration",
    "Responsive design",
    "Error handling",
    "Mock data fallback",
    "CDN fallback loading",
  ],
};

writeFileSync(
  join(config.distDir, "build-stats.json"),
  JSON.stringify(stats, null, 2),
);

console.log("\n📊 Build Statistics:");
console.log(`   Component size: ~${stats.outputSize}KB`);
console.log(`   Vue version: ${stats.vueVersion}`);
console.log(`   Features: ${stats.features.length}`);
console.log(
  "\n🎉 GitHub Dashboard Web Component build completed successfully!",
);
console.log(
  "💡 The component is now ready for distribution via the plugin registry.",
);
