<template>
    <div class="json-linter-container">
        <div class="header">
            <h3>JSON Linter & Formatter</h3>
            <div class="toolbar">
                <button
                    @click="formatJSON"
                    :disabled="!hasInput"
                    class="format-btn"
                >
                    Format JSON
                </button>
                <button @click="clearAll" class="clear-btn">Clear</button>
                <button
                    @click="copyFormatted"
                    :disabled="!isValid"
                    class="copy-btn"
                >
                    Copy Formatted
                </button>
            </div>
        </div>

        <div class="content-area">
            <div class="input-section">
                <div class="input-header">
                    <label for="json-input">JSON Editor:</label>
                    <div class="status" :class="statusClass">
                        <span class="status-icon">{{ statusIcon }}</span>
                        <span class="status-text">{{ statusText }}</span>
                    </div>
                </div>
                <div
                    id="json-input"
                    ref="jsonEditor"
                    contenteditable="true"
                    @input="onContentEditableInput"
                    @paste="onPaste"
                    @keydown="onKeyDown"
                    class="json-editor"
                    :class="{
                        error: !isValid && hasInput,
                        valid: isValid && hasInput,
                    }"
                    data-placeholder="Paste or type JSON here..."
                ></div>
            </div>

            <div
                class="output-section"
                v-if="!isValid && hasInput && errorMessage"
            >
                <div class="output-header">
                    <label>Error Details:</label>
                </div>
                <pre class="json-output error" v-html="errorMessage"></pre>
            </div>
        </div>

        <div class="error-section" v-if="!isValid && hasInput && errorMessage">
            <div class="error-message">
                <strong>Error:</strong> {{ errorMessage }}
                <span v-if="lineNumber" class="line-info"
                    >(Line {{ lineNumber }})</span
                >
            </div>
        </div>

        <div class="stats">
            <div class="stat">
                <strong>Characters:</strong> {{ characterCount }}
            </div>
            <div class="stat"><strong>Lines:</strong> {{ lineCount }}</div>
            <div class="stat" v-if="isValid">
                <strong>Size:</strong> {{ formatSize(characterCount) }}
            </div>
            <div class="stat">
                <strong>Validations:</strong> {{ validationCount }}
            </div>
            <div class="stat">
                <strong>Mount Count:</strong> {{ mountCount }}
            </div>
        </div>

        <div class="persistence-info">
            <small>State: {{ persistenceStatus }}</small>
        </div>
    </div>
</template>

<script setup>
import {
    ref,
    computed,
    onMounted,
    onBeforeUnmount,
    getCurrentInstance,
    nextTick,
    watch,
} from "vue";
import delveSDK from "./delve-sdk.js";

// State storage key - unique per plugin
const STORAGE_KEY = "json-linter-formatter-state";

// Global state store - persists across component instances
if (!window.jsonLinterGlobalState) {
    window.jsonLinterGlobalState = {
        inputJSON: "",
        formattedJSON: "",
        isValid: false,
        errorMessage: "",
        lineNumber: 0,
        validationCount: 0,
        mountCount: 0,
        lastUpdated: null,
    };

    // Try to restore from localStorage on first load
    try {
        const saved = localStorage.getItem(STORAGE_KEY);
        if (saved) {
            const parsed = JSON.parse(saved);
            Object.assign(window.jsonLinterGlobalState, parsed);
            console.log(
                "🔄 Restored state from localStorage:",
                window.jsonLinterGlobalState,
            );
        }
    } catch (error) {
        console.warn("Failed to restore state from localStorage:", error);
    }
}

// Reactive state - bound to global state
const inputJSON = ref(window.jsonLinterGlobalState.inputJSON);
const formattedJSON = ref(window.jsonLinterGlobalState.formattedJSON);
const isValid = ref(window.jsonLinterGlobalState.isValid);
const errorMessage = ref(window.jsonLinterGlobalState.errorMessage);
const lineNumber = ref(window.jsonLinterGlobalState.lineNumber);
const validationCount = ref(window.jsonLinterGlobalState.validationCount);
const mountCount = ref(window.jsonLinterGlobalState.mountCount);
const persistenceStatus = ref("Loading...");

// Computed properties
const hasInput = computed(() => inputJSON.value.trim().length > 0);
const characterCount = computed(() => inputJSON.value.length);
const lineCount = computed(() =>
    inputJSON.value ? inputJSON.value.split("\n").length : 0,
);

const statusClass = computed(() => {
    if (!hasInput.value) return "status-neutral";
    return isValid.value ? "status-valid" : "status-error";
});

const statusIcon = computed(() => {
    if (!hasInput.value) return "○";
    return isValid.value ? "✓" : "✗";
});

const statusText = computed(() => {
    if (!hasInput.value) return "Ready";
    return isValid.value ? "Valid JSON" : "Invalid JSON";
});

const formattedOutput = computed(() => {
    if (!hasInput.value) {
        return '<span class="placeholder">Formatted JSON will appear here...</span>';
    }

    if (!isValid.value) {
        return `<span class="error-text">${errorMessage.value}</span>`;
    }

    return syntaxHighlight(formattedJSON.value);
});

// State persistence functions
const saveState = () => {
    try {
        const state = {
            inputJSON: inputJSON.value,
            formattedJSON: formattedJSON.value,
            isValid: isValid.value,
            errorMessage: errorMessage.value,
            lineNumber: lineNumber.value,
            validationCount: validationCount.value,
            mountCount: mountCount.value,
            lastUpdated: new Date().toISOString(),
        };

        // Update global state
        Object.assign(window.jsonLinterGlobalState, state);

        // Save to localStorage
        localStorage.setItem(STORAGE_KEY, JSON.stringify(state));

        persistenceStatus.value = `Saved at ${new Date().toLocaleTimeString()}`;
        console.log("💾 State saved:", state);
    } catch (error) {
        console.error("Failed to save state:", error);
        persistenceStatus.value = "Save failed";
    }
};

const loadState = () => {
    try {
        // Load from global state (which may have been updated by another instance)
        const state = window.jsonLinterGlobalState;

        inputJSON.value = state.inputJSON || "";
        formattedJSON.value = state.formattedJSON || "";
        isValid.value = state.isValid || false;
        errorMessage.value = state.errorMessage || "";
        lineNumber.value = state.lineNumber || 0;
        validationCount.value = state.validationCount || 0;
        // Don't restore mountCount - it should increment

        persistenceStatus.value = state.lastUpdated
            ? `Restored from ${new Date(state.lastUpdated).toLocaleTimeString()}`
            : "No saved state";

        console.log("📥 State loaded:", state);

        // Update editor content
        nextTick(() => {
            if (jsonEditor.value) {
                if (isValid.value && formattedJSON.value) {
                    jsonEditor.value.innerHTML = syntaxHighlight(
                        formattedJSON.value,
                    );
                } else {
                    jsonEditor.value.textContent = inputJSON.value;
                }
            }

            // Re-validate if we have content
            if (hasInput.value) {
                validateAndFormat();
            }
        });
    } catch (error) {
        console.error("Failed to load state:", error);
        persistenceStatus.value = "Load failed";
    }
};

// Auto-save with debouncing
let saveTimeout = null;
const debouncedSave = () => {
    clearTimeout(saveTimeout);
    saveTimeout = setTimeout(saveState, 500);
};

// Watch for changes and auto-save
watch([inputJSON, isValid, errorMessage], () => {
    debouncedSave();
});

// Methods
const validateAndFormat = () => {
    validationCount.value++;

    if (!hasInput.value) {
        resetValidationState();
        return;
    }

    try {
        const parsed = JSON.parse(inputJSON.value);
        const formatted = JSON.stringify(parsed, null, 2);

        formattedJSON.value = formatted;
        isValid.value = true;
        errorMessage.value = "";
        lineNumber.value = 0;
    } catch (error) {
        isValid.value = false;
        errorMessage.value = error.message;
        formattedJSON.value = "";

        // Try to extract line number from error message
        const lineMatch = error.message.match(/line (\d+)/i);
        if (lineMatch) {
            lineNumber.value = parseInt(lineMatch[1]);
        } else {
            const posMatch = error.message.match(/position (\d+)/i);
            if (posMatch) {
                const position = parseInt(posMatch[1]);
                const textUpToPosition = inputJSON.value.substring(0, position);
                lineNumber.value = textUpToPosition.split("\n").length;
            }
        }
    }
};

const resetValidationState = () => {
    formattedJSON.value = "";
    isValid.value = false;
    errorMessage.value = "";
    lineNumber.value = 0;
};

const jsonEditor = ref(null);

const onContentEditableInput = (event) => {
    const content = event.target.textContent || "";
    inputJSON.value = content;

    // Debounce validation and auto-format
    clearTimeout(window.jsonLinterValidationTimeout);
    window.jsonLinterValidationTimeout = setTimeout(() => {
        validateAndFormat();
        if (isValid.value && formattedJSON.value !== content) {
            // Auto-format valid JSON
            autoFormatEditor();
        }
    }, 500);
};

const onPaste = (event) => {
    event.preventDefault();
    const paste = (event.clipboardData || window.clipboardData).getData("text");

    // Insert plain text
    const selection = window.getSelection();
    if (selection.rangeCount) {
        const range = selection.getRangeAt(0);
        range.deleteContents();
        range.insertNode(document.createTextNode(paste));
        range.collapse(false);
        selection.removeAllRanges();
        selection.addRange(range);
    }

    // Trigger input event
    inputJSON.value = jsonEditor.value?.textContent || "";
    onContentEditableInput({ target: jsonEditor.value });
};

const onKeyDown = (event) => {
    // Handle Tab key for proper indentation
    if (event.key === "Tab") {
        event.preventDefault();
        document.execCommand("insertText", false, "  ");
    }
};

const autoFormatEditor = () => {
    if (jsonEditor.value && formattedJSON.value) {
        const currentContent = jsonEditor.value.textContent;
        if (currentContent !== formattedJSON.value) {
            // Store cursor position
            const selection = window.getSelection();
            const range =
                selection.rangeCount > 0 ? selection.getRangeAt(0) : null;

            // Update content with syntax highlighting
            jsonEditor.value.innerHTML = syntaxHighlight(formattedJSON.value);

            // Try to restore cursor position (simplified)
            if (range) {
                try {
                    selection.removeAllRanges();
                    selection.addRange(range);
                } catch (e) {
                    // Fallback: place cursor at end
                    const newRange = document.createRange();
                    newRange.selectNodeContents(jsonEditor.value);
                    newRange.collapse(false);
                    selection.removeAllRanges();
                    selection.addRange(newRange);
                }
            }
        }
    }
};

const onInput = () => {
    // Legacy method for compatibility
    onContentEditableInput({ target: { textContent: inputJSON.value } });
};

const formatJSON = () => {
    if (isValid.value) {
        inputJSON.value = formattedJSON.value;
        if (jsonEditor.value) {
            jsonEditor.value.innerHTML = syntaxHighlight(formattedJSON.value);
        }
    }
};

const clearAll = () => {
    inputJSON.value = "";
    if (jsonEditor.value) {
        jsonEditor.value.textContent = "";
    }
    resetValidationState();
    saveState(); // Immediately save cleared state
};

const copyFormatted = async () => {
    if (isValid.value && formattedJSON.value) {
        try {
            await navigator.clipboard.writeText(formattedJSON.value);
            console.log("Formatted JSON copied to clipboard");

            // Show temporary feedback
            const oldStatus = persistenceStatus.value;
            persistenceStatus.value = "Copied to clipboard!";
            setTimeout(() => {
                persistenceStatus.value = oldStatus;
            }, 2000);
        } catch (error) {
            console.error("Failed to copy to clipboard:", error);
            // Fallback for older browsers
            const textArea = document.createElement("textarea");
            textArea.value = formattedJSON.value;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand("copy");
            document.body.removeChild(textArea);
        }
    }
};

const formatSize = (bytes) => {
    if (bytes < 1024) return bytes + " bytes";
    if (bytes < 1024 * 1024) return Math.round(bytes / 1024) + " KB";
    return Math.round(bytes / (1024 * 1024)) + " MB";
};

const syntaxHighlight = (json) => {
    if (!json) return "";

    return json
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(
            /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
            function (match) {
                let cls = "number";
                if (/^"/.test(match)) {
                    if (/:$/.test(match)) {
                        cls = "key";
                    } else {
                        cls = "string";
                    }
                } else if (/true|false/.test(match)) {
                    cls = "boolean";
                } else if (/null/.test(match)) {
                    cls = "null";
                }
                return '<span class="' + cls + '">' + match + "</span>";
            },
        );
};

// Cleanup
delveSDK.onCleanup(() => {
    console.log("JSON Linter plugin shutting down - final save");
    clearTimeout(saveTimeout);
    clearTimeout(window.jsonLinterValidationTimeout);
    saveState();
});

onMounted(() => {
    mountCount.value++;
    console.log(`🚀 JSON Linter mounted (instance #${mountCount.value})`);

    // Setup SDK integration
    const instance = getCurrentInstance();
    delveSDK.setupVueIntegration(instance);

    // Load state (this handles both first load and subsequent mounts)
    loadState();

    persistenceStatus.value =
        mountCount.value === 1
            ? "First load - state initialized"
            : `Remounted #${mountCount.value} - state preserved!`;
});

onBeforeUnmount(() => {
    console.log(
        `📤 JSON Linter unmounting (instance #${mountCount.value}) - state will persist`,
    );
    clearTimeout(saveTimeout);
    clearTimeout(window.jsonLinterValidationTimeout);
    saveState();
});
</script>

<style>
.json-linter-container {
    max-width: 100%;
    padding: 1rem;
    font-family: inherit;
    background-color: var(--bg-color, #f8f9fa);
    color: var(--text-color, #343a40);
    border-radius: 8px;
    border: 1px solid var(--border-color, #e9ecef);
    box-sizing: border-box;
    overflow: hidden;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #dee2e6;
}

.header h3 {
    margin: 0;
    color: #343a40;
    font-size: 1.25rem;
}

.toolbar {
    display: flex;
    gap: 0.5rem;
}

.toolbar button {
    padding: 0.375rem 0.75rem;
    border: 1px solid #ced4da;
    border-radius: 4px;
    background-color: #fff;
    color: #495057;
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.2s ease;
}

.toolbar button:hover:not(:disabled) {
    background-color: #e9ecef;
    border-color: #adb5bd;
}

.toolbar button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.format-btn:hover:not(:disabled) {
    background-color: #007bff;
    color: white;
    border-color: #007bff;
}

.clear-btn:hover:not(:disabled) {
    background-color: #dc3545;
    color: white;
    border-color: #dc3545;
}

.copy-btn:hover:not(:disabled) {
    background-color: #28a745;
    color: white;
    border-color: #28a745;
}

.content-area {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1rem;
    overflow: hidden;
}

.input-section,
.output-section {
    display: flex;
    flex-direction: column;
    min-width: 0;
    overflow: hidden;
}

.input-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
}

.input-section label,
.output-section label {
    font-weight: 600;
    color: var(--text-color, #495057);
    font-size: 0.875rem;
    margin: 0;
}

.output-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
}

.status {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.75rem;
    font-weight: 500;
    padding: 0.25rem 0.5rem;
    border-radius: 12px;
}

.status-neutral {
    background-color: #f8f9fa;
    color: #6c757d;
}

.status-valid {
    background-color: #d1edff;
    color: #0c5460;
}

.status-error {
    background-color: #f8d7da;
    color: #721c24;
}

.json-editor {
    width: 100%;
    min-height: 300px;
    max-height: 500px;
    padding: 0.75rem;
    border: 1px solid var(--border-color, #ced4da);
    border-radius: 4px;
    font-family: "Courier New", monospace;
    font-size: 0.875rem;
    line-height: 1.4;
    background-color: var(--input-bg, #eee);
    color: var(--text-color, #343a40);
    transition: border-color 0.2s ease;
    overflow: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
    box-sizing: border-box;
}

.json-editor:focus {
    outline: none;
    border-color: var(--primary-color, #007bff);
    box-shadow: 0 0 0 2px var(--primary-color-alpha, rgba(0, 123, 255, 0.25));
}

.json-editor.error {
    border-color: var(--error-color, #dc3545);
}

.json-editor.error:focus {
    border-color: var(--error-color, #dc3545);
    box-shadow: 0 0 0 2px var(--error-color-alpha, rgba(220, 53, 69, 0.25));
}

.json-editor.valid {
    border-color: var(--success-color, #28a745);
    background-color: var(--success-bg, #f8fff9);
}

.json-editor:empty:before {
    content: attr(data-placeholder);
    color: var(--placeholder-color, #6c757d);
    font-style: italic;
}

.json-output {
    width: 100%;
    min-height: 200px;
    padding: 0.75rem;
    border: 1px solid #ced4da;
    border-radius: 4px;
    font-family: "Courier New", monospace;
    font-size: 0.875rem;
    line-height: 1.4;
    background-color: #eee;
    margin: 0;
    overflow: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
    box-sizing: border-box;
}

.json-output.valid {
    border-color: #28a745;
    background-color: #f8fff9;
}

.json-output.error {
    border-color: #dc3545;
    background-color: #fff5f5;
}

.placeholder {
    color: #6c757d;
    font-style: italic;
}

.error-text {
    color: #dc3545;
}

.string {
    color: var(--json-string-color, #032f62);
}
.number {
    color: var(--json-number-color, #005cc5);
}
.boolean {
    color: var(--json-boolean-color, #d73a49);
}
.null {
    color: var(--json-null-color, #d73a49);
}
.key {
    color: var(--json-key-color, #22863a);
    font-weight: bold;
}

.error-section {
    background-color: var(--error-bg, #f8d7da);
    color: var(--error-text, #721c24);
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    border-left: 4px solid var(--error-color, #dc3545);
}

.error-message {
    font-size: 0.875rem;
    line-height: 1.4;
}

.line-info {
    font-weight: bold;
    color: #495057;
}

.stats {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.75rem;
    background-color: var(--stats-bg, #e9ecef);
    border-radius: 4px;
    font-size: 0.75rem;
    margin-bottom: 0.5rem;
}

.stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
}

.stat strong {
    color: var(--text-color, #495057);
    margin-bottom: 0.25rem;
}

.persistence-info {
    text-align: center;
    color: var(--muted-color, #6c757d);
    font-size: 0.7rem;
    padding: 0.25rem;
    background-color: var(--muted-bg, #f1f3f4);
    border-radius: 4px;
}

@media (max-width: 768px) {
    .header {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
    }

    .toolbar {
        align-self: stretch;
        justify-content: space-between;
    }

    .input-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
    }

    .stats {
        flex-direction: column;
        gap: 0.5rem;
    }

    .stat {
        flex-direction: row;
        justify-content: space-between;
    }

    .json-editor {
        min-height: 250px;
        max-height: 400px;
    }
}
</style>
