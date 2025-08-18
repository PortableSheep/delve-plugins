package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// PluginAPI defines the interface that plugins must implement
type PluginAPI interface {
	Initialize(config map[string]interface{}) error
	Start(ctx context.Context) error
	Stop() error
	GetInfo() map[string]interface{}
	HandleRequest(method, path string, body []byte) ([]byte, error)
	HealthCheck(checkName string) error
}

// PluginServer wraps a PluginAPI implementation to serve JSON-over-TCP requests
type PluginServer struct {
	plugin   PluginAPI
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
}

// Request represents a plugin request
type Request struct {
	Method string                 `json:"method"`
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// Response represents a plugin response
type Response struct {
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// NewPluginServer creates a new plugin server
func NewPluginServer(plugin PluginAPI) *PluginServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &PluginServer{
		plugin: plugin,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Serve starts the plugin server on the specified port
func (s *PluginServer) Serve(port int) error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	fmt.Printf("Plugin server listening on port %d\n", port)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.ctx.Err() != nil {
					return nil // Server is shutting down
				}
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}

			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles a single client connection
func (s *PluginServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			response := Response{
				Success: false,
				Error:   fmt.Sprintf("invalid request format: %v", err),
			}
			encoder.Encode(response)
			continue
		}

		response := s.handleRequest(req)
		encoder.Encode(response)
	}
}

// handleRequest processes a plugin request
func (s *PluginServer) handleRequest(req Request) Response {
	switch req.Method {
	case "health_check":
		return s.handleHealthCheck(req)
	case "execute_action":
		return s.handleExecuteAction(req)
	case "get_info":
		return s.handleGetInfo(req)
	case "initialize":
		return s.handleInitialize(req)
	case "start":
		return s.handleStart(req)
	case "stop":
		return s.handleStop(req)
	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown method: %s", req.Method),
		}
	}
}

// handleHealthCheck processes health check requests
func (s *PluginServer) handleHealthCheck(req Request) Response {
	checkName, _ := req.Data["check_name"].(string)
	
	err := s.plugin.HealthCheck(checkName)
	if err != nil {
		return Response{
			Success: true,
			Result: map[string]interface{}{
				"healthy": false,
				"message": err.Error(),
			},
		}
	}

	return Response{
		Success: true,
		Result: map[string]interface{}{
			"healthy": true,
			"message": "OK",
		},
	}
}

// handleExecuteAction processes action execution requests
func (s *PluginServer) handleExecuteAction(req Request) Response {
	data, _ := req.Data["data"]
	dataBytes, _ := json.Marshal(data)

	result, err := s.plugin.HandleRequest("POST", "/execute", dataBytes)
	if err != nil {
		return Response{
			Success: false,
			Error:   err.Error(),
		}
	}

	var resultData interface{}
	json.Unmarshal(result, &resultData)

	return Response{
		Success: true,
		Result: map[string]interface{}{
			"result": resultData,
		},
	}
}

// handleGetInfo processes info retrieval requests
func (s *PluginServer) handleGetInfo(req Request) Response {
	info := s.plugin.GetInfo()
	return Response{
		Success: true,
		Result:  info,
	}
}

// handleInitialize processes initialization requests
func (s *PluginServer) handleInitialize(req Request) Response {
	config, ok := req.Data["config"].(map[string]interface{})
	if !ok {
		return Response{
			Success: false,
			Error:   "config is required",
		}
	}

	err := s.plugin.Initialize(config)
	if err != nil {
		return Response{
			Success: false,
			Error:   err.Error(),
		}
	}

	return Response{
		Success: true,
	}
}

// handleStart processes start requests
func (s *PluginServer) handleStart(req Request) Response {
	err := s.plugin.Start(s.ctx)
	if err != nil {
		return Response{
			Success: false,
			Error:   err.Error(),
		}
	}

	return Response{
		Success: true,
	}
}

// handleStop processes stop requests
func (s *PluginServer) handleStop(req Request) Response {
	err := s.plugin.Stop()
	if err != nil {
		return Response{
			Success: false,
			Error:   err.Error(),
		}
	}

	return Response{
		Success: true,
	}
}

// Shutdown gracefully shuts down the plugin server
func (s *PluginServer) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.listener != nil {
		s.listener.Close()
	}
}

// RunPluginServer is a helper function that plugins can use to start their server
func RunPluginServer(plugin PluginAPI) error {
	// Get port from command line argument
	port := 50051 // default port
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "--port=") {
				if p, err := strconv.Atoi(arg[7:]); err == nil {
					port = p
				}
			}
		}
	}

	server := NewPluginServer(plugin)
	return server.Serve(port)
}