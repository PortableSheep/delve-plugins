package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	sdk "github.com/PortableSheep/delve_sdk"
)

// JSONValidationResult represents the result of JSON validation and formatting
type JSONValidationResult struct {
	IsValid       bool   `json:"isValid"`
	FormattedJSON string `json:"formattedJson"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
	LineNumber    int    `json:"lineNumber,omitempty"`
}

var plugin *sdk.Plugin

// handleHostMessage processes messages from the host application
func handleHostMessage(messageType int, data []byte) {
	log.Printf("JSON Linter received message: Type=%d, Data=%s", messageType, string(data))

	switch messageType {
	case 1: // JSON validation request
		result := validateAndFormatJSON(string(data))
		log.Printf("Validation result: %+v", result)
	}
}

// validateAndFormatJSON takes a JSON string, validates it, and returns formatted JSON
func validateAndFormatJSON(jsonStr string) JSONValidationResult {
	result := JSONValidationResult{
		IsValid: false,
	}

	// Trim whitespace
	jsonStr = strings.TrimSpace(jsonStr)

	if jsonStr == "" {
		result.ErrorMessage = "Empty JSON input"
		return result
	}

	// Parse the JSON to check if it's valid
	var parsed interface{}
	err := json.Unmarshal([]byte(jsonStr), &parsed)

	if err != nil {
		result.ErrorMessage = err.Error()

		// Try to extract line number from error message if available
		if strings.Contains(err.Error(), "line") {
			// This is a simple approach - real implementation might need more sophisticated parsing
			parts := strings.Split(err.Error(), "line")
			if len(parts) > 1 {
				// Extract line number - this is simplified
				result.LineNumber = 1 // Default to line 1 if we can't parse it properly
			}
		}
		return result
	}

	// If valid, format it nicely
	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		result.ErrorMessage = "Failed to format JSON: " + err.Error()
		return result
	}

	result.IsValid = true
	result.FormattedJSON = string(formatted)
	return result
}

// demonstrateStorage shows how to use the plugin storage system for JSON snippets
// This is optional and can be disabled by setting DISABLE_STORAGE_DEMO=true
func demonstrateStorage(p *sdk.Plugin) {
	// Check if storage demo is disabled
	if os.Getenv("DISABLE_STORAGE_DEMO") == "true" {
		log.Println("Storage demonstration disabled (DISABLE_STORAGE_DEMO=true)")
		return
	}

	go func() {
		log.Println("=== JSON Linter Plugin Storage Demonstration (Background) ===")

		// Store some configuration for the JSON linter
		config := map[string]interface{}{
			"theme":           "monokai",
			"autoFormat":      true,
			"validateOnType":  true,
			"tabSize":         2,
			"autoSave":        true,
			"saveInterval":    500,
			"persistentState": true,
			"syntaxHighlight": true,
		}

		// Try to store config with timeout
		err := p.StoreConfig("linter_settings", config, "1.0.0")
		if err != nil {
			log.Printf("⚠️  Storage not available - linter config not stored: %v", err)
		} else {
			log.Println("✓ Stored JSON linter configuration")
		}

		// Store some sample JSON snippets
		sampleJSON := map[string]interface{}{
			"snippets": []map[string]interface{}{
				{"name": "API Response", "json": `{"status": "success", "data": {"id": 1, "name": "John Doe"}}`},
				{"name": "Config File", "json": `{"server": {"host": "localhost", "port": 3000}, "database": {"url": "mongodb://localhost:27017"}}`},
				{"name": "User Profile", "json": `{"user": {"id": 123, "username": "john_doe", "email": "john@example.com", "preferences": {"theme": "dark", "notifications": true}}}`},
			},
			"createdAt": time.Now(),
			"version":   "1.0.0",
		}

		err = p.StoreData("json_snippets", sampleJSON, "1.0.0")
		if err != nil {
			log.Printf("⚠️  Storage not available - JSON snippets not stored: %v", err)
		} else {
			log.Println("✓ Stored JSON snippets")
		}

		// Show storage statistics
		stats, err := p.GetStats()
		if err != nil {
			log.Printf("⚠️  Storage stats not available: %v", err)
		} else {
			log.Printf("\n=== Storage Statistics ===")
			for storageType, count := range stats {
				log.Printf("%s: %d items", storageType, count)
			}
		}
	}()
}

func main() {
	log.Printf("JSON Linter Plugin launched with arguments: %v", os.Args)

	// Define the plugin's metadata
	pluginInfo := &sdk.RegisterRequest{
		Name:             "json-linter-formatter",
		Description:      "A plugin for validating, linting, and formatting JSON with syntax highlighting and persistent state",
		UiComponentPath:  "frontend/component.js",
		CustomElementTag: "json-linter-formatter",
	}

	// Connect to the host and register
	var err error
	plugin, err = sdk.Start(pluginInfo)
	if err != nil {
		log.Fatalf("Failed to start JSON linter plugin: %v", err)
	}

	// Demonstrate storage functionality (non-blocking, optional)
	demonstrateStorage(plugin)

	// Test the JSON validation functionality
	log.Println("=== Testing JSON Validation ===")

	testCases := []string{
		`{"valid": "json", "number": 42, "boolean": true, "null_value": null}`,
		`{"invalid": "json", "missing": }`,
		`{"nested": {"object": {"with": ["arrays", "and", "values"]}}}`,
		`[{"id": 1, "name": "Item 1"}, {"id": 2, "name": "Item 2"}]`,
		`invalid json string`,
		``,
		`"just a string"`,
		`42`,
		`true`,
		`null`,
	}

	for i, testCase := range testCases {
		log.Printf("Test case %d: %s", i+1, testCase)
		result := validateAndFormatJSON(testCase)
		if result.IsValid {
			log.Printf("✓ Valid JSON, formatted:\n%s", result.FormattedJSON)
		} else {
			log.Printf("✗ Invalid JSON: %s", result.ErrorMessage)
		}
		log.Println("---")
	}

	// Start listening for events from the host
	log.Println("JSON Linter Plugin is running and listening for host events.")
	log.Println("Plugin supports persistent state using global window state.")
	log.Println("Note: Storage operations run in background and may timeout if storage service is unavailable.")
	log.Println("To disable storage demo entirely, set environment variable DISABLE_STORAGE_DEMO=true")
	plugin.Listen(handleHostMessage)

	log.Println("JSON Linter Plugin shutting down.")
}
