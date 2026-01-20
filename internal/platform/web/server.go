package web

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/gorilla/websocket"
)

//go:embed static/*
var staticFiles embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Server represents the web server
type Server struct {
	port       string
	binaryPath string
	clients    map[*websocket.Conn]bool
	clientsMux sync.Mutex
}

// NewServer creates a new web server
func NewServer(port, binaryPath string) *Server {
	return &Server{
		port:       port,
		binaryPath: binaryPath,
		clients:    make(map[*websocket.Conn]bool),
	}
}

// CommandRequest represents a command execution request
type CommandRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Start starts the web server
func (s *Server) Start() error {
	// API endpoints (register first to take precedence)
	http.HandleFunc("/api/commands", s.handleGetCommands)
	http.HandleFunc("/api/execute", s.handleExecuteCommand)
	http.HandleFunc("/ws", s.handleWebSocket)

	// Serve static files from the static subdirectory
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to create sub filesystem: %w", err)
	}
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	addr := ":" + s.port
	log.Printf("🌐 Web UI available at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleGetCommands returns available commands
func (s *Server) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	commands := []map[string]interface{}{
		{
			"id":          "sync",
			"name":        "Sync Reference Entity",
			"description": "Synchronize all records from a Reference Entity",
			"command":     "sync",
			"args": []map[string]interface{}{
				{"name": "entity-name", "type": "text", "placeholder": "brands", "required": true},
			},
			"flags": []map[string]interface{}{
				{"name": "debug", "type": "checkbox", "label": "Debug mode"},
			},
		},
		{
			"id":          "sync-product",
			"name":        "Sync Product Hierarchy",
			"description": "Synchronize a complete product hierarchy",
			"command":     "sync-product",
			"args": []map[string]interface{}{
				{"name": "identifier", "type": "text", "placeholder": "COMMON-001", "required": true},
			},
			"flags": []map[string]interface{}{
				{"name": "debug", "type": "checkbox", "label": "Debug mode"},
			},
		},
		{
			"id":          "sync-attribute",
			"name":        "Sync Attribute",
			"description": "Synchronize a single attribute",
			"command":     "sync-attribute",
			"args": []map[string]interface{}{
				{"name": "code", "type": "text", "placeholder": "sku", "required": true},
			},
			"flags": []map[string]interface{}{
				{"name": "debug", "type": "checkbox", "label": "Debug mode"},
			},
		},
		{
			"id":          "sync-category",
			"name":        "Sync Category",
			"description": "Synchronize a single category",
			"command":     "sync-category",
			"args": []map[string]interface{}{
				{"name": "code", "type": "text", "placeholder": "master", "required": true},
			},
			"flags": []map[string]interface{}{
				{"name": "debug", "type": "checkbox", "label": "Debug mode"},
			},
		},
		{
			"id":          "sync-family",
			"name":        "Sync Family",
			"description": "Synchronize a single family",
			"command":     "sync-family",
			"args": []map[string]interface{}{
				{"name": "code", "type": "text", "placeholder": "clothing", "required": true},
			},
			"flags": []map[string]interface{}{
				{"name": "debug", "type": "checkbox", "label": "Debug mode"},
			},
		},
		{
			"id":          "sync-updated-products",
			"name":        "Sync Updated Products",
			"description": "Synchronize products updated since a specific date",
			"command":     "sync-updated-products",
			"args": []map[string]interface{}{
				{"name": "date", "type": "datetime-local", "placeholder": "2024-01-01T00:00:00", "required": true},
			},
			"flags": []map[string]interface{}{
				{"name": "debug", "type": "checkbox", "label": "Debug mode"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// handleExecuteCommand executes a command and streams output via WebSocket
func (s *Server) handleExecuteCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build command
	args := append([]string{req.Command}, req.Args...)
	log.Printf("Executing command: %s %v", s.binaryPath, args)

	cmd := exec.Command(s.binaryPath, args...)

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		http.Error(w, fmt.Sprintf("Error creating stdout pipe: %v", err), http.StatusInternalServerError)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error creating stderr pipe: %v", err)
		http.Error(w, fmt.Sprintf("Error creating stderr pipe: %v", err), http.StatusInternalServerError)
		return
	}

	// Start command
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v", err)
		http.Error(w, fmt.Sprintf("Error starting command: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Command started successfully")

	// Stream output to all connected WebSocket clients
	go s.streamOutput(stdout, "stdout")
	go s.streamOutput(stderr, "stderr")

	// Wait for command to finish
	go func() {
		err := cmd.Wait()
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
		}

		s.broadcast(map[string]interface{}{
			"type":     "exit",
			"exitCode": exitCode,
		})
	}()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Command started",
	})
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	s.clientsMux.Lock()
	s.clients[conn] = true
	s.clientsMux.Unlock()

	// Send welcome message
	conn.WriteJSON(map[string]interface{}{
		"type":    "connected",
		"message": "Connected to Akeneo Migrator",
	})

	// Keep connection alive
	go func() {
		defer func() {
			s.clientsMux.Lock()
			delete(s.clients, conn)
			s.clientsMux.Unlock()
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// streamOutput streams command output to WebSocket clients
func (s *Server) streamOutput(reader io.Reader, streamType string) {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			s.broadcast(map[string]interface{}{
				"type":   "output",
				"stream": streamType,
				"data":   string(buf[:n]),
			})
		}
		if err != nil {
			break
		}
	}
}

// broadcast sends a message to all connected WebSocket clients
func (s *Server) broadcast(message map[string]interface{}) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	for client := range s.clients {
		err := client.WriteJSON(message)
		if err != nil {
			client.Close()
			delete(s.clients, client)
		}
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	for client := range s.clients {
		client.Close()
	}

	return nil
}
