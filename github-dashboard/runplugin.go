package main

import (
	"context"
)

// This file contains a stub for backward compatibility during transition to gRPC

// PluginAPI defines the interface that plugins must implement
type PluginAPI interface {
	Initialize(config map[string]interface{}) error
	Start(ctx context.Context) error
	Stop() error
	GetInfo() map[string]interface{}
	HandleRequest(method, path string, body []byte) ([]byte, error)
	HealthCheck(checkName string) error
	GetMenuItems() ([]map[string]interface{}, error)
}

// RunPluginServer is a stub function for backward compatibility
// This will be deprecated and removed after all plugins migrate to gRPC
func RunPluginServer(plugin PluginAPI) error {
	// This function is only kept for backward compatibility
	// It actually doesn't get called because main() now uses proto.RunGrpcPluginServer
	return nil
}
