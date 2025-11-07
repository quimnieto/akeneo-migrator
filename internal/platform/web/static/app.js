// WebSocket connection
let ws = null;
let currentCommand = null;
let isExecuting = false;

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    connectWebSocket();
    loadCommands();
    setupEventListeners();
});

// Connect to WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
        console.log('WebSocket connected');
        addOutput('Connected to server', 'info');
    };
    
    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };
    
    ws.onclose = () => {
        console.log('WebSocket disconnected');
        addOutput('Disconnected from server', 'error');
        // Reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

// Handle WebSocket messages
function handleWebSocketMessage(message) {
    switch (message.type) {
        case 'connected':
            console.log(message.message);
            break;
        case 'output':
            addOutput(message.data, message.stream);
            break;
        case 'exit':
            isExecuting = false;
            updateExecuteButton();
            if (message.exitCode === 0) {
                addOutput('\n✅ Command completed successfully', 'success');
            } else {
                addOutput(`\n❌ Command failed with exit code ${message.exitCode}`, 'error');
            }
            break;
    }
}

// Load available commands
async function loadCommands() {
    try {
        const response = await fetch('/api/commands');
        const commands = await response.json();
        renderCommands(commands);
    } catch (error) {
        console.error('Error loading commands:', error);
        addOutput('Error loading commands', 'error');
    }
}

// Render commands list
function renderCommands(commands) {
    const commandsList = document.getElementById('commands-list');
    commandsList.innerHTML = '';
    
    commands.forEach(command => {
        const item = document.createElement('div');
        item.className = 'command-item';
        item.dataset.commandId = command.id;
        item.innerHTML = `
            <h3>${command.name}</h3>
            <p>${command.description}</p>
        `;
        item.addEventListener('click', () => selectCommand(command));
        commandsList.appendChild(item);
    });
}

// Select a command
function selectCommand(command) {
    currentCommand = command;
    
    // Update active state
    document.querySelectorAll('.command-item').forEach(item => {
        item.classList.remove('active');
    });
    document.querySelector(`[data-command-id="${command.id}"]`).classList.add('active');
    
    // Render command form
    renderCommandForm(command);
}

// Render command form
function renderCommandForm(command) {
    const formContainer = document.getElementById('command-form');
    
    let html = `
        <h2>${command.name}</h2>
        <p class="description">${command.description}</p>
        <form id="execute-form">
    `;
    
    // Render arguments
    command.args.forEach((arg, index) => {
        html += `
            <div class="form-group">
                <label for="arg-${index}">${arg.name}${arg.required ? ' *' : ''}</label>
                <input 
                    type="${arg.type}" 
                    id="arg-${index}" 
                    name="${arg.name}"
                    placeholder="${arg.placeholder || ''}"
                    ${arg.required ? 'required' : ''}
                >
            </div>
        `;
    });
    
    // Render flags
    if (command.flags && command.flags.length > 0) {
        command.flags.forEach((flag, index) => {
            html += `
                <div class="form-group">
                    <label class="checkbox-label">
                        <input 
                            type="${flag.type}" 
                            id="flag-${index}" 
                            name="${flag.name}"
                        >
                        ${flag.label}
                    </label>
                </div>
            `;
        });
    }
    
    html += `
            <button type="submit" class="btn-execute" id="execute-btn">
                Execute Command
            </button>
        </form>
    `;
    
    formContainer.innerHTML = html;
    
    // Setup form submission
    document.getElementById('execute-form').addEventListener('submit', handleExecute);
}

// Handle command execution
async function handleExecute(event) {
    event.preventDefault();
    
    if (isExecuting) {
        return;
    }
    
    const formData = new FormData(event.target);
    const args = [];
    
    // Collect arguments
    currentCommand.args.forEach((arg, index) => {
        const value = formData.get(arg.name);
        if (value) {
            // For datetime-local, convert to ISO format
            if (arg.type === 'datetime-local') {
                args.push(value.replace('T', 'T') + ':00');
            } else {
                args.push(value);
            }
        }
    });
    
    // Collect flags
    if (currentCommand.flags) {
        currentCommand.flags.forEach((flag, index) => {
            const value = formData.get(flag.name);
            if (value === 'on') {
                args.push(`--${flag.name}`);
            }
        });
    }
    
    // Clear output
    document.getElementById('output').innerHTML = '';
    
    // Add command to output
    addOutput(`$ akeneo-migrator ${currentCommand.command} ${args.join(' ')}\n`, 'info');
    
    isExecuting = true;
    updateExecuteButton();
    
    try {
        const response = await fetch('/api/execute', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                command: currentCommand.command,
                args: args,
            }),
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
    } catch (error) {
        console.error('Error executing command:', error);
        addOutput(`Error: ${error.message}`, 'error');
        isExecuting = false;
        updateExecuteButton();
    }
}

// Add output to terminal
function addOutput(text, type = 'stdout') {
    const output = document.getElementById('output');
    const line = document.createElement('div');
    line.className = `output-line ${type}`;
    line.textContent = text;
    output.appendChild(line);
    
    // Auto-scroll to bottom
    output.scrollTop = output.scrollHeight;
}

// Update execute button state
function updateExecuteButton() {
    const btn = document.getElementById('execute-btn');
    if (btn) {
        btn.disabled = isExecuting;
        btn.textContent = isExecuting ? 'Executing...' : 'Execute Command';
    }
}

// Setup event listeners
function setupEventListeners() {
    document.getElementById('clear-output').addEventListener('click', () => {
        document.getElementById('output').innerHTML = '';
    });
}
