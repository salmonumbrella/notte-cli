# notte-cli

A command-line interface for the [notte.cc](https://notte.cc) browser automation platform. Control browser sessions, AI agents, and web scraping through a simple, resource-based CLI.

## Features

- Full access to the Notte API through intuitive commands
- Browser session management and automation
- AI agent orchestration and workflow execution
- Secure credential storage via system keyring
- Multiple output formats (text, JSON)
- Web scraping with structured data extraction

## Installation

### Via Go Install

```bash
go install github.com/salmonumbrella/notte-cli/cmd/notte@latest
```

### Build from Source

```bash
git clone https://github.com/salmonumbrella/notte-cli.git
cd notte-cli
make build    # Creates ./notte binary
make install  # Installs to $GOPATH/bin
```

## Authentication

notte-cli supports three authentication methods (checked in order):

1. **Environment variable** (recommended for CI/CD):
   ```bash
   export NOTTE_API_KEY="your-api-key"
   ```

2. **System keyring** (recommended for local development):
   ```bash
   notte auth login
   # Enter your API key when prompted
   ```

3. **Config file** (`~/.config/notte/config.json`):
   ```json
   {
     "api_key": "your-api-key",
     "api_url": "https://api.notte.cc"
   }
   ```

Check authentication status:
```bash
notte auth status
```

## Quick Start

### Start a browser session and visit a page

```bash
# Start a new session
notte sessions start

# Start with specific options
notte sessions start --headless --proxy-region us-east

# Execute commands in a session
notte session execute <session-id> '{"action": "goto", "url": "https://example.com"}'
```

### Scrape a webpage

```bash
# Quick scrape with structured output
notte scrape https://example.com --schema '{"title": "string", "price": "number"}'

# Get raw HTML
notte scrape-html https://example.com
```

### Run an AI agent workflow

```bash
# List available workflows
notte workflows list

# Run a workflow
notte workflow run <workflow-id> --input '{"url": "https://example.com"}'

# Watch workflow progress
notte workflow run <workflow-id> --watch
```

### Manage credentials securely

```bash
# Create a vault
notte vaults create --name "Production Credentials"

# Add credentials to vault
notte vault credentials <vault-id> add \
  --username "user@example.com" \
  --password "secret" \
  --url "https://app.example.com"
```

## Command Reference

### Authentication
- `notte auth login` - Store API key in system keyring
- `notte auth logout` - Remove API key from keyring
- `notte auth status` - Show authentication status and source

### Browser Sessions (Multiple)
- `notte sessions list` - List all active sessions
- `notte sessions start` - Start a new browser session

### Browser Session (Single)
- `notte session status <id>` - Get session status
- `notte session stop <id>` - Stop a session
- `notte session observe <id>` - Watch session in real-time
- `notte session execute <id>` - Execute browser actions
- `notte session scrape <id>` - Scrape current page
- `notte session cookies <id>` - Manage session cookies
- `notte session debug <id>` - Get debug information
- `notte session network <id>` - View network activity
- `notte session replay <id>` - Get session replay data
- `notte session offset <id>` - Adjust session timing

### AI Agents (Multiple)
- `notte agents list` - List all agents
- `notte agents start` - Start a new agent

### AI Agent (Single)
- `notte agent status <id>` - Get agent status
- `notte agent stop <id>` - Stop an agent
- `notte agent workflow-code <id>` - Get agent's workflow code
- `notte agent replay <id>` - Get agent execution replay

### Vaults (Multiple)
- `notte vaults list` - List all vaults
- `notte vaults create` - Create a new vault

### Vault (Single)
- `notte vault update <id>` - Update vault metadata
- `notte vault delete <id>` - Delete a vault
- `notte vault credentials <id>` - Manage vault credentials
- `notte vault card <id>` - Manage payment cards in vault

### Personas (Multiple)
- `notte personas list` - List all personas
- `notte personas create` - Create a new persona

### Persona (Single)
- `notte persona show <id>` - View persona details
- `notte persona delete <id>` - Delete a persona
- `notte persona emails <id>` - Manage persona email addresses
- `notte persona sms <id>` - Manage persona SMS numbers
- `notte persona phone <id>` - Manage persona phone numbers

### Workflows (Multiple)
- `notte workflows list` - List all workflows
- `notte workflows create` - Create a new workflow

### Workflow (Single)
- `notte workflow show <id>` - View workflow details
- `notte workflow update <id>` - Update workflow
- `notte workflow delete <id>` - Delete workflow
- `notte workflow fork <id>` - Fork workflow to new version
- `notte workflow run <id>` - Execute workflow
- `notte workflow runs <id>` - List workflow runs
- `notte workflow schedule <id>` - Schedule recurring execution
- `notte workflow unschedule <id>` - Remove schedule
- `notte workflow run-stop <run-id>` - Stop a running workflow
- `notte workflow run-metadata <run-id>` - Get run metadata

### Files
- `notte files list` - List uploaded files
- `notte files upload <path>` - Upload a file
- `notte files download <id>` - Download a file

### Web Scraping
- `notte scrape <url>` - Scrape with structured schema
- `notte scrape-html <url>` - Get raw HTML content

### Usage & Monitoring
- `notte usage` - View API usage statistics
- `notte usage logs` - View usage logs

### Utilities
- `notte health` - Check API health status
- `notte prompt-improve` - Improve a prompt with AI
- `notte prompt-nudge` - Get prompt suggestions
- `notte version` - Show CLI version

## Configuration

Configuration is stored at `~/.config/notte/config.json`:

```json
{
  "api_key": "your-api-key",
  "api_url": "https://api.notte.cc"
}
```

You can also override the API URL:
```bash
# In config file
{
  "api_url": "https://custom-api.example.com"
}

# Or via environment
export NOTTE_API_URL="https://custom-api.example.com"
```

## Output Formats

### Text Output (default)

Human-readable, colorized output:
```bash
notte sessions list
```

### JSON Output

Machine-readable output for scripting:
```bash
notte sessions list --output json
notte sessions list -o json
```

Disable colors in text output:
```bash
notte sessions list --no-color
```

## Global Flags

All commands support these flags:

- `-o, --output <format>` - Output format: `text` or `json` (default: `text`)
- `--no-color` - Disable colored output
- `-v, --verbose` - Enable verbose logging
- `--timeout <seconds>` - API request timeout (default: 30)

## Examples

### Automated web scraping

```bash
# Start session, scrape, and stop
SESSION_ID=$(notte sessions start --headless -o json | jq -r '.id')
notte session execute $SESSION_ID '{"action": "goto", "url": "https://news.ycombinator.com"}'
notte session scrape $SESSION_ID --schema '{"articles": [{"title": "string", "url": "string"}]}'
notte session stop $SESSION_ID
```

### Running workflows with input

```bash
# Create and run a workflow
WORKFLOW_ID=$(notte workflows create --name "Price Checker" --code "..." -o json | jq -r '.id')
notte workflow run $WORKFLOW_ID --input '{"product_url": "https://example.com/product"}' --watch
```

### Managing credentials

```bash
# Create vault and add credentials
VAULT_ID=$(notte vaults create --name "E-commerce Sites" -o json | jq -r '.id')
notte vault credentials $VAULT_ID add \
  --username "shop@example.com" \
  --password "$SECURE_PASSWORD" \
  --url "https://shop.example.com"
```

## License

MIT License - see LICENSE file for details.

## Resources

- [Notte API Documentation](https://notte.cc/docs)
- [GitHub Repository](https://github.com/salmonumbrella/notte-cli)
- [Report Issues](https://github.com/salmonumbrella/notte-cli/issues)
