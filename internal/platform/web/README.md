# Web UI

## Overview

Browser-based interface for executing Akeneo Migrator commands with real-time output streaming.

## Features

- ✅ Execute all available commands from your browser
- ✅ Real-time output streaming via WebSockets
- ✅ Clean, modern interface
- ✅ No external dependencies (embedded static files)

## Installation

```bash
# Install WebSocket dependency
go get github.com/gorilla/websocket
go mod tidy

# Build the project
make build
```

## Usage

### Start the web server

```bash
# Default port (3000)
./bin/akeneo-migrator web

# Custom port
./bin/akeneo-migrator web --port 8080
```

### Access the UI

Open your browser at: **http://localhost:3000**

## Architecture

```
Browser (WebSocket) ←→ Go Server ←→ CLI Commands
```

### Components

- **Server** (`server.go`): HTTP server with WebSocket support
- **Frontend** (`static/`): HTML + CSS + JavaScript
- **WebSocket**: Real-time bidirectional communication

### Flow

1. User selects a command in the UI
2. Frontend sends HTTP POST to `/api/execute`
3. Server spawns the CLI command as a subprocess
4. Output is streamed via WebSocket to all connected clients
5. Exit code is sent when command completes

## API Endpoints

### GET /api/commands

Returns list of available commands with their parameters.

**Response:**
```json
[
  {
    "id": "sync-product",
    "name": "Sync Product Hierarchy",
    "description": "...",
    "command": "sync-product",
    "args": [...],
    "flags": [...]
  }
]
```

### POST /api/execute

Executes a command.

**Request:**
```json
{
  "command": "sync-product",
  "args": ["COMMON-001", "--debug"]
}
```

**Response:**
```json
{
  "status": "Command started"
}
```

### WebSocket /ws

Real-time output streaming.

**Messages:**
```json
// Connected
{"type": "connected", "message": "..."}

// Output
{"type": "output", "stream": "stdout", "data": "..."}

// Exit
{"type": "exit", "exitCode": 0}
```

## Development

### File Structure

```
internal/web/
├── server.go           # HTTP server + WebSocket handler
├── static/
│   ├── index.html     # Main UI
│   ├── style.css      # Styles
│   └── app.js         # Frontend logic
└── README.md
```

### Adding New Commands

Commands are automatically discovered from the CLI. No changes needed in the web UI.

## Troubleshooting

### Error 500 when executing commands

Check the server logs for the actual error. Common issues:

1. **Binary not found**: Make sure the binary path is correct
2. **Permission denied**: Ensure the binary is executable
3. **Missing config**: Commands need proper Akeneo configuration

### WebSocket disconnects

The WebSocket will automatically reconnect after 3 seconds.

### Commands not appearing

Refresh the page or check browser console for errors.

## Security Notes

⚠️ **This is a development tool**. Do not expose to the internet without:

- Authentication
- HTTPS
- CORS restrictions
- Input validation
- Rate limiting

## Future Enhancements

- [ ] Command history
- [ ] Multiple concurrent executions
- [ ] Save/load command presets
- [ ] Export output to file
- [ ] Authentication
- [ ] Command scheduling
