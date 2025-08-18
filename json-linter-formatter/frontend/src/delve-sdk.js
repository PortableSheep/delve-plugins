/**
 * Delve Plugin SDK - JavaScript/Frontend utilities
 * Provides standardized cleanup and lifecycle management for plugin frontends
 */

class DelvePluginSDK {
    constructor() {
        this.cleanupHandlers = new Set();
        this.isShuttingDown = false;
        this.intervals = new Set();
        this.timeouts = new Set();
        this.eventListeners = new Map();
        this.websockets = new Set();
    }

    /**
     * Register a cleanup handler to be called when the plugin shuts down
     * @param {Function} handler - Function to call during cleanup
     */
    onCleanup(handler) {
        if (typeof handler === 'function') {
            this.cleanupHandlers.add(handler);
        }
    }

    /**
     * Create a managed interval that will be automatically cleared on cleanup
     * @param {Function} callback - Function to call on interval
     * @param {number} delay - Delay in milliseconds
     * @returns {number} - Interval ID
     */
    setManagedInterval(callback, delay) {
        const intervalId = setInterval(callback, delay);
        this.intervals.add(intervalId);
        return intervalId;
    }

    /**
     * Create a managed timeout that will be automatically cleared on cleanup
     * @param {Function} callback - Function to call after timeout
     * @param {number} delay - Delay in milliseconds
     * @returns {number} - Timeout ID
     */
    setManagedTimeout(callback, delay) {
        const timeoutId = setTimeout(callback, delay);
        this.timeouts.add(timeoutId);
        return timeoutId;
    }

    /**
     * Add a managed event listener that will be automatically removed on cleanup
     * @param {EventTarget} target - Element to attach listener to
     * @param {string} event - Event name
     * @param {Function} handler - Event handler function
     * @param {object} options - Event listener options
     */
    addManagedEventListener(target, event, handler, options = {}) {
        target.addEventListener(event, handler, options);

        if (!this.eventListeners.has(target)) {
            this.eventListeners.set(target, []);
        }
        this.eventListeners.get(target).push({ event, handler, options });
    }

    /**
     * Register a WebSocket for automatic cleanup
     * @param {WebSocket} websocket - WebSocket connection to manage
     */
    addManagedWebSocket(websocket) {
        this.websockets.add(websocket);

        // Auto-remove when connection closes naturally
        websocket.addEventListener('close', () => {
            this.websockets.delete(websocket);
        });
    }

    /**
     * Clear a managed interval
     * @param {number} intervalId - Interval ID to clear
     */
    clearManagedInterval(intervalId) {
        clearInterval(intervalId);
        this.intervals.delete(intervalId);
    }

    /**
     * Clear a managed timeout
     * @param {number} timeoutId - Timeout ID to clear
     */
    clearManagedTimeout(timeoutId) {
        clearTimeout(timeoutId);
        this.timeouts.delete(timeoutId);
    }

    /**
     * Remove a managed event listener
     * @param {EventTarget} target - Element to remove listener from
     * @param {string} event - Event name
     * @param {Function} handler - Event handler function
     */
    removeManagedEventListener(target, event, handler) {
        target.removeEventListener(event, handler);

        if (this.eventListeners.has(target)) {
            const listeners = this.eventListeners.get(target);
            const index = listeners.findIndex(l => l.event === event && l.handler === handler);
            if (index !== -1) {
                listeners.splice(index, 1);
            }
            if (listeners.length === 0) {
                this.eventListeners.delete(target);
            }
        }
    }

    /**
     * Perform cleanup of all managed resources
     * This should be called when the plugin is being shut down
     */
    cleanup() {
        if (this.isShuttingDown) {
            return; // Prevent multiple cleanup calls
        }

        this.isShuttingDown = true;
        console.log('[Delve SDK] Starting plugin cleanup');

        // Clear all intervals
        this.intervals.forEach(intervalId => {
            clearInterval(intervalId);
        });
        this.intervals.clear();

        // Clear all timeouts
        this.timeouts.forEach(timeoutId => {
            clearTimeout(timeoutId);
        });
        this.timeouts.clear();

        // Remove all event listeners
        this.eventListeners.forEach((listeners, target) => {
            listeners.forEach(({ event, handler, options }) => {
                target.removeEventListener(event, handler, options);
            });
        });
        this.eventListeners.clear();

        // Close all WebSockets
        this.websockets.forEach(websocket => {
            if (websocket.readyState === WebSocket.OPEN || websocket.readyState === WebSocket.CONNECTING) {
                websocket.close(1000, 'Plugin shutting down');
            }
        });
        this.websockets.clear();

        // Call all custom cleanup handlers
        this.cleanupHandlers.forEach(handler => {
            try {
                handler();
            } catch (error) {
                console.warn('[Delve SDK] Error in cleanup handler:', error);
            }
        });
        this.cleanupHandlers.clear();

        console.log('[Delve SDK] Plugin cleanup completed');
    }

    /**
     * Setup automatic cleanup integration with Vue component
     * Call this in your Vue component's onMounted hook
     * @param {object} vueInstance - Vue component instance
     */
    setupVueIntegration(vueInstance) {
        // Expose cleanup method for the host app
        if (vueInstance && vueInstance.exposed) {
            vueInstance.exposed.cleanup = () => this.cleanup();
        }

        // Also attach to the root element for direct access
        if (vueInstance && vueInstance.vnode && vueInstance.vnode.el) {
            vueInstance.vnode.el.cleanup = () => this.cleanup();
        }
    }

    /**
     * Create a fetch wrapper that can be cancelled during cleanup
     * @param {string} url - URL to fetch
     * @param {object} options - Fetch options
     * @returns {Promise} - Fetch promise with abort support
     */
    managedFetch(url, options = {}) {
        const controller = new AbortController();

        // Add abort controller to cleanup
        this.onCleanup(() => {
            controller.abort();
        });

        return fetch(url, {
            ...options,
            signal: controller.signal
        });
    }
}

// Export singleton instance
const delveSDK = new DelvePluginSDK();

export default delveSDK;

// Also provide named exports for convenience
export const {
    onCleanup,
    setManagedInterval,
    setManagedTimeout,
    addManagedEventListener,
    addManagedWebSocket,
    clearManagedInterval,
    clearManagedTimeout,
    removeManagedEventListener,
    cleanup,
    setupVueIntegration,
    managedFetch
} = delveSDK;
