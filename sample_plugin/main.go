package main

import (
	"log"
	"os"
	"time"

	sdk "github.com/PortableSheep/delve_sdk"
)

// handleHostMessage is the callback function that processes messages from the host.
func handleHostMessage(messageType int, data []byte) {
	log.Printf("Received message from host: Type=%d, Data=%s", messageType, string(data))
	// Here you could check the message content and trigger different actions.
	// For example, if string(data) == "shutdown" { ... }
}

// demonstrateStorage shows how to use the plugin storage system
func demonstrateStorage(plugin *sdk.Plugin) {
	log.Println("=== Plugin Storage Demonstration ===")

	// Store some configuration
	config := map[string]interface{}{
		"theme":           "dark",
		"refreshInterval": 5000,
		"apiEndpoint":     "https://api.example.com",
	}

	err := plugin.StoreConfig("settings", config, "1.0.0")
	if err != nil {
		log.Printf("Failed to store config: %v", err)
	} else {
		log.Println("✓ Stored plugin configuration")
	}

	// Store some application data
	chartData := map[string]interface{}{
		"datasets": []map[string]interface{}{
			{"label": "Dataset 1", "data": []int{10, 20, 30, 40}},
			{"label": "Dataset 2", "data": []int{15, 25, 35, 45}},
		},
		"createdAt": time.Now(),
		"version":   "1.0.0",
	}

	err = plugin.StoreData("chart_data", chartData, "1.0.0")
	if err != nil {
		log.Printf("Failed to store chart data: %v", err)
	} else {
		log.Println("✓ Stored chart data")
	}

	// Store current state
	state := map[string]interface{}{
		"currentView":   "dashboard",
		"selectedItems": []string{"item1", "item3", "item7"},
		"lastActivity":  time.Now(),
	}

	err = plugin.StoreState("current_session", state, "1.0.0")
	if err != nil {
		log.Printf("Failed to store state: %v", err)
	} else {
		log.Println("✓ Stored current state")
	}

	// Demonstrate loading data
	log.Println("\n=== Loading Stored Data ===")

	// Load configuration
	configItem, err := plugin.LoadConfig("settings")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	} else {
		log.Printf("✓ Loaded config: %+v", configItem.Value)
	}

	// Load application data
	dataItem, err := plugin.LoadData("chart_data")
	if err != nil {
		log.Printf("Failed to load chart data: %v", err)
	} else {
		log.Printf("✓ Loaded chart data (created: %s)", dataItem.CreatedAt.Format(time.RFC3339))
	}

	// Load state
	stateItem, err := plugin.LoadState("current_session")
	if err != nil {
		log.Printf("Failed to load state: %v", err)
	} else {
		log.Printf("✓ Loaded state: %+v", stateItem.Value)
	}

	// Show storage statistics
	stats, err := plugin.GetStats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		log.Printf("\n=== Storage Statistics ===")
		for storageType, count := range stats {
			log.Printf("%s: %d items", storageType, count)
		}
	}

	// List all keys in data storage
	keys, err := plugin.List("data")
	if err != nil {
		log.Printf("Failed to list data keys: %v", err)
	} else {
		log.Printf("\nData storage keys: %v", keys)
	}
}

func main() {
	log.Printf("Plugin launched with arguments: %v", os.Args)

	// Define the plugin's metadata.
	pluginInfo := &sdk.RegisterRequest{
		Name:             "super-chart-plugin",
		Description:      "A sample plugin that uses WebSockets.",
		UiComponentPath:  "frontend/component.js",
		CustomElementTag: "super-chart-plugin",
	}

	// Connect to the host and register.
	plugin, err := sdk.Start(pluginInfo)
	if err != nil {
		log.Fatalf("Failed to start plugin: %v", err)
	}

	// Demonstrate storage functionality
	demonstrateStorage(plugin)

	// Start listening for events from the host.
	// This is a blocking call and will run until the connection is closed.
	log.Println("Plugin is running and listening for host events.")
	plugin.Listen(handleHostMessage)

	log.Println("Plugin shutting down.")
}
