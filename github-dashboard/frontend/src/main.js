import { createApp, defineCustomElement } from "vue";
import App from "./App.vue";

// GitHub Dashboard Web Component
class GitHubDashboardElement extends HTMLElement {
  constructor() {
    super();
    this.app = null;
    this.isConnected = false;
    this.shadow = null;
    this.cleanup = new Set();
  }

  connectedCallback() {
    if (this.isConnected) return;

    try {
      // Create shadow DOM for style encapsulation
      this.shadow = this.attachShadow({ mode: "open" });

      // Create container for Vue app
      const container = document.createElement("div");
      container.className = "github-dashboard-root";
      this.shadow.appendChild(container);

      // Add CSS reset and base styles to shadow DOM
      const styles = document.createElement("style");
      styles.textContent = `
        :host {
          display: block;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Noto Sans', Helvetica, Arial, sans-serif;
          font-size: 14px;
          line-height: 1.5;
          color: #1f2328;
        }

        * {
          box-sizing: border-box;
        }

        .github-dashboard-root {
          min-height: 400px;
          background: #ffffff;
        }
      `;
      this.shadow.appendChild(styles);

      // Set up plugin API if available
      this.setupPluginAPI();

      // Create and mount Vue app
      this.app = createApp(App, {
        // Pass any props from attributes
        ...this.getPropsFromAttributes(),
      });

      // Global error handling
      this.app.config.errorHandler = (err, instance, info) => {
        console.error("GitHub Dashboard Error:", err, info);
      };

      // Mount the app
      this.app.mount(container);
      this.isConnected = true;

      // Set up cleanup handlers
      this.setupCleanup();
    } catch (error) {
      console.error("Failed to initialize GitHub Dashboard:", error);
      this.renderError(error.message);
    }
  }

  disconnectedCallback() {
    this.cleanup.forEach((handler) => {
      try {
        handler();
      } catch (error) {
        console.warn("Cleanup error:", error);
      }
    });
    this.cleanup.clear();

    if (this.app) {
      this.app.unmount();
      this.app = null;
    }

    this.isConnected = false;
  }

  setupPluginAPI() {
    // Set up communication with the plugin backend
    if (!window.pluginAPI) {
      window.pluginAPI = {};
    }

    // Mock API for development/fallback
    if (!window.pluginAPI.getRepositories) {
      window.pluginAPI.getRepositories = async () => {
        // Return demo data if no real API is available
        return [
          {
            name: "delve",
            full_name: "PortableSheep/delve",
            description: "Advanced plugin system for Go applications",
            stargazers_count: 1250,
            forks_count: 180,
            open_issues_count: 23,
            language: "Go",
            updated_at: new Date().toISOString(),
            html_url: "https://github.com/PortableSheep/delve",
          },
        ];
      };
    }

    if (!window.pluginAPI.getPullRequests) {
      window.pluginAPI.getPullRequests = async (repoName) => {
        // Return demo data if no real API is available
        return [
          {
            number: 42,
            title: "Add new dashboard features",
            state: "open",
            created_at: new Date(Date.now() - 86400000).toISOString(),
            updated_at: new Date(Date.now() - 3600000).toISOString(),
            html_url: `https://github.com/${repoName}/pull/42`,
            user: {
              login: "developer",
              avatar_url: "https://github.com/identicons/developer.png",
            },
          },
        ];
      };
    }
  }

  setupCleanup() {
    // Set up cleanup for intervals, event listeners, etc.
    const originalSetInterval = window.setInterval;
    const originalSetTimeout = window.setTimeout;
    const originalAddEventListener = EventTarget.prototype.addEventListener;

    const intervals = new Set();
    const timeouts = new Set();
    const listeners = new WeakMap();

    // Track intervals
    window.setInterval = (handler, delay) => {
      const id = originalSetInterval(handler, delay);
      intervals.add(id);
      return id;
    };

    // Track timeouts
    window.setTimeout = (handler, delay) => {
      const id = originalSetTimeout(handler, delay);
      timeouts.add(id);
      return id;
    };

    // Add cleanup handlers
    this.cleanup.add(() => {
      intervals.forEach((id) => clearInterval(id));
      timeouts.forEach((id) => clearTimeout(id));
      intervals.clear();
      timeouts.clear();
    });

    // Restore original functions on cleanup
    this.cleanup.add(() => {
      window.setInterval = originalSetInterval;
      window.setTimeout = originalSetTimeout;
    });
  }

  getPropsFromAttributes() {
    const props = {};

    // Get configuration from attributes
    if (this.hasAttribute("github-token")) {
      props.githubToken = this.getAttribute("github-token");
    }

    if (this.hasAttribute("repositories")) {
      try {
        props.repositories = JSON.parse(this.getAttribute("repositories"));
      } catch (e) {
        console.warn("Invalid repositories attribute:", e);
      }
    }

    if (this.hasAttribute("refresh-interval")) {
      props.refreshInterval = parseInt(
        this.getAttribute("refresh-interval"),
        10,
      );
    }

    if (this.hasAttribute("compact-view")) {
      props.compactView = this.getAttribute("compact-view") === "true";
    }

    return props;
  }

  renderError(message) {
    if (this.shadow) {
      this.shadow.innerHTML = `
        <style>
          .error-container {
            padding: 20px;
            background: #fff8f8;
            border: 1px solid #ffebee;
            border-radius: 6px;
            color: #d73a49;
            text-align: center;
          }
          .error-icon {
            font-size: 48px;
            margin-bottom: 16px;
          }
          .error-message {
            font-size: 16px;
            margin-bottom: 8px;
          }
          .error-details {
            font-size: 12px;
            color: #586069;
          }
        </style>
        <div class="error-container">
          <div class="error-icon">⚠️</div>
          <div class="error-message">Failed to load GitHub Dashboard</div>
          <div class="error-details">${message}</div>
        </div>
      `;
    }
  }

  // Attribute change handling
  static get observedAttributes() {
    return ["github-token", "repositories", "refresh-interval", "compact-view"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue !== newValue && this.isConnected && this.app) {
      // Re-initialize with new attributes if needed
      const props = this.getPropsFromAttributes();
      // Update app props if possible, or trigger a re-render
      console.log("Attribute changed:", name, "from", oldValue, "to", newValue);
    }
  }
}

// Load Vue from CDN if not available
function initializeComponent() {
  // Register the custom element
  if (!customElements.get("github-dashboard")) {
    customElements.define("github-dashboard", GitHubDashboardElement);
    console.log("GitHub Dashboard web component registered");
  }
}

// Check if Vue is available
if (typeof window !== "undefined") {
  if (window.Vue) {
    // Vue is already loaded
    initializeComponent();
  } else {
    // Load Vue from CDN
    const script = document.createElement("script");
    script.src = "https://unpkg.com/vue@3.4.0/dist/vue.global.js";
    script.onload = () => {
      console.log("Vue loaded from CDN");
      initializeComponent();
    };
    script.onerror = () => {
      console.error("Failed to load Vue from CDN");
    };
    document.head.appendChild(script);
  }
} else {
  // Server-side or non-browser environment
  console.warn("GitHub Dashboard component requires a browser environment");
}

// Export for module systems
export default GitHubDashboardElement;
export { GitHubDashboardElement };
