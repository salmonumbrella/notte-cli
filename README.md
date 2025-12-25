# notte-cli

Command-line interface for the notte.cc browser automation platform - control browser sessions, AI agents, and web scraping through intuitive resource-based commands.

## Features

- **Browser session management** - headless or headed Chrome/Firefox sessions with full control
- **AI agent orchestration** - run and monitor AI-powered workflows
- **Web scraping** - structured data extraction with custom schemas
- **Secure credential storage** - vaults for managing authentication and payment information
- **Personas** - create and manage digital identities with email, phone, and SMS
- **Multiple output formats** - human-readable text or JSON for scripting
- **System keyring integration** - secure API key storage using OS-native keychains

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

## Quick Start

### 1. Authenticate

```bash
notte auth login
# Enter your notte.cc API key when prompted
```

### 2. Verify Authentication

```bash
notte auth status
```

### 3. Start a Browser Session

```bash
notte sessions start --headless
```

### 4. Quick Scrape a Page

```bash
notte scrape https://news.ycombinator.com --instructions "Extract the top stories"
```

## Configuration

### API Key Storage

notte-cli supports three authentication methods (checked in order of priority):

1. **Environment variable** (recommended for CI/CD):
   ```bash
   export NOTTE_API_KEY="your-api-key"
   ```

2. **System keyring** (recommended for local development):
   ```bash
   notte auth login
   # Stores API key securely in macOS Keychain, Linux Secret Service, or Windows Credential Manager
   ```

3. **Config file** (`~/.config/notte/config.json`):
   ```json
   {
     "api_key": "your-api-key",
     "api_url": "https://api.notte.cc"
   }
   ```

### Environment Variables

- `NOTTE_API_KEY` - API key for authentication
- `NOTTE_API_URL` - Override API endpoint (default: https://api.notte.cc)

### Config File Location

`~/.config/notte/config.json`

## Security

### Credential Storage

API keys are stored securely in your system's keychain:
- **macOS**: Keychain Access
- **Linux**: Secret Service (GNOME Keyring, KWallet)
- **Windows**: Credential Manager

### Best Practices

- **Never pass API keys on the command line** - use `notte auth login` or environment variables
- **Use vaults for sensitive credentials** - store website passwords and payment cards in encrypted vaults
- **Rotate API keys regularly** - obtain new keys from notte.cc dashboard
- **Clear credentials when done** - use `notte auth logout` to remove stored API keys

## Commands

### Authentication

```bash
notte auth login                 # Store API key in system keychain
notte auth logout                # Remove API key from keychain
notte auth status                # Show authentication status and source
```

### Browser Sessions

```bash
# List and create sessions
notte sessions list              # List all active browser sessions
notte sessions start [flags]     # Start a new browser session

# Session management
notte session status --id <id>   # Get session status
notte session stop --id <id>     # Stop a browser session
notte session observe --id <id>  # Watch session in real-time
notte session execute --id <id>  # Execute browser actions (goto, click, type, etc.)
notte session scrape --id <id>   # Scrape content from current page

# Session data
notte session cookies --id <id>           # Get all cookies from session
notte session cookies-set --id <id>       # Set cookies from JSON file
notte session network --id <id>           # View network activity logs
notte session debug --id <id>             # Get debug information
notte session replay --id <id>            # Get session replay data
notte session offset --id <id>            # Get session timing offset
```

#### Session Start Options

```bash
notte sessions start \
  --browser chromium|chrome|firefox    # Browser type (default: chromium)
  --headless                           # Run in headless mode (default: true)
  --timeout <minutes>                  # Session timeout 1-15 minutes (default: 3)
  --user-agent <string>                # Custom user agent
  --viewport-width <pixels>            # Viewport width
  --viewport-height <pixels>           # Viewport height
  --proxies                            # Use default proxy rotation
  --solve-captchas                     # Automatically solve captchas
  --cdp-url <url>                      # CDP URL of remote session provider
```

### AI Agents

```bash
# List and create agents
notte agents list                      # List all AI agents
notte agents start                     # Start a new AI agent

# Agent management
notte agent status --id <id>           # Get agent status
notte agent stop --id <id>             # Stop an agent
notte agent workflow-code --id <id>    # Get agent's workflow code
notte agent replay --id <id>           # Get agent execution replay
```

### Workflows

```bash
# List and create workflows
notte workflows list                   # List all workflows
notte workflows create                 # Create a new workflow

# Workflow management
notte workflow show --id <id>          # View workflow details
notte workflow update --id <id>        # Update workflow configuration
notte workflow delete --id <id>        # Delete a workflow
notte workflow fork --id <id>          # Fork workflow to new version

# Workflow execution
notte workflow run --id <id>                      # Execute workflow
notte workflow runs --id <id>                     # List workflow runs
notte workflow run-stop --id <run-id>             # Stop a running workflow
notte workflow run-metadata --id <run-id>         # Get run metadata
notte workflow run-metadata-update --id <run-id>  # Update run metadata

# Workflow scheduling
notte workflow schedule --id <id>      # Schedule recurring execution (cron)
notte workflow unschedule --id <id>    # Remove schedule
```

### Vaults

```bash
# List and create vaults
notte vaults list                      # List all credential vaults
notte vaults create                    # Create a new vault

# Vault management
notte vault update --id <id>           # Update vault metadata
notte vault delete --id <id>           # Delete a vault

# Vault credentials
notte vault credentials list --id <id>           # List all credentials
notte vault credentials add --id <id>            # Add credentials
notte vault credentials get --id <id>            # Get credentials for URL
notte vault credentials delete --id <id>         # Delete credentials

# Vault payment cards
notte vault card --id <id>             # Manage payment cards in vault
```

### Personas

```bash
# List and create personas
notte personas list                    # List all personas
notte personas create                  # Create a new persona

# Persona management
notte persona show --id <id>           # View persona details
notte persona delete --id <id>         # Delete a persona

# Persona contact methods
notte persona emails --id <id>         # Manage email addresses
notte persona sms --id <id>            # Manage SMS numbers
notte persona phone --id <id>          # Manage phone numbers
```

### Files

```bash
notte files list                       # List uploaded files
notte files upload <path>              # Upload a file to notte.cc
notte files download <id>              # Download a file by ID
```

### Web Scraping

```bash
# Quick scraping without sessions
notte scrape <url> [flags]             # Scrape with structured extraction
notte scrape-html <url>                # Get raw HTML content

# Scraping options
notte scrape <url> \
  --instructions <text>                # Extraction instructions
  --only-main-content                  # Extract only main content area
```

### Usage & Monitoring

```bash
notte usage                            # View API usage statistics
notte usage logs                       # View detailed usage logs
```

### Utilities

```bash
notte health                           # Check API health status
notte prompt-improve                   # Improve a prompt with AI assistance
notte prompt-nudge                     # Get prompt optimization suggestions
notte version                          # Show CLI version information
```

## Output Formats

### Text Output (default)

Human-readable, colorized output:

```bash
$ notte sessions list
ID                                   STATUS    BROWSER     CREATED
ses_abc123def456                     ACTIVE    chromium    2024-01-15 10:30:00
ses_xyz789uvw012                     STOPPED   chrome      2024-01-15 09:15:00
```

### JSON Output

Machine-readable output for scripting and automation:

```bash
$ notte sessions list --output json
{
  "sessions": [
    {
      "id": "ses_abc123def456",
      "status": "ACTIVE",
      "browser": "chromium",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

Disable colors in text output:

```bash
notte sessions list --no-color
```

## Examples

### Automated Web Scraping Pipeline

```bash
# Start session, navigate, scrape, and cleanup
SESSION_ID=$(notte sessions start --headless -o json | jq -r '.id')

# Navigate to page
notte session execute --id $SESSION_ID << 'EOF'
{"action": "goto", "url": "https://news.ycombinator.com"}
EOF

# Extract structured data
notte session scrape --id $SESSION_ID \
  --instructions "Extract top 10 stories with title and URL"

# Cleanup
notte session stop --id $SESSION_ID
```

### Running a Workflow with Input

```bash
# List workflows to find ID
notte workflows list

# Run workflow with specific input
notte workflow run --id wfl_abc123 << 'EOF'
{
  "product_url": "https://example.com/product/12345",
  "notify_email": "alerts@example.com"
}
EOF
```

### Managing Credentials Securely

```bash
# Create a vault for production credentials
VAULT_ID=$(notte vaults create --name "Production Sites" -o json | jq -r '.id')

# Add website credentials
notte vault credentials add --id $VAULT_ID \
  --username "admin@example.com" \
  --password "$SECURE_PASSWORD" \
  --url "https://app.example.com"

# List stored credentials
notte vault credentials list --id $VAULT_ID
```

### Creating a Persona for Testing

```bash
# Create a test persona
PERSONA_ID=$(notte personas create \
  --name "Test User" \
  --date-of-birth "1990-01-01" \
  -o json | jq -r '.id')

# Add email to persona
notte persona emails --id $PERSONA_ID add --email "testuser@example.com"

# Add phone number
notte persona phone --id $PERSONA_ID add --number "+15555551234"
```

### Multi-Step Browser Automation

```bash
# Start browser with specific configuration
SESSION_ID=$(notte sessions start \
  --browser chrome \
  --viewport-width 1920 \
  --viewport-height 1080 \
  --solve-captchas \
  -o json | jq -r '.id')

# Execute multiple actions
notte session execute --id $SESSION_ID '{"action": "goto", "url": "https://example.com"}'
notte session execute --id $SESSION_ID '{"action": "click", "selector": "#login-button"}'
notte session execute --id $SESSION_ID '{"action": "type", "selector": "#username", "text": "user@example.com"}'

# Get current page state
notte session observe --id $SESSION_ID

# Stop when done
notte session stop --id $SESSION_ID
```

## Global Flags

All commands support these flags:

- `-o, --output <format>` - Output format: `text` or `json` (default: text)
- `--no-color` - Disable colored output
- `-v, --verbose` - Enable verbose logging
- `--timeout <seconds>` - API request timeout in seconds (default: 30)
- `-h, --help` - Show help for any command

## Shell Completions

Generate shell completions for your preferred shell:

### Bash

```bash
notte completion bash > /usr/local/etc/bash_completion.d/notte
# Or for Linux:
notte completion bash > /etc/bash_completion.d/notte
```

### Zsh

```zsh
notte completion zsh > "${fpath[1]}/_notte"
# Or add to .zshrc:
echo 'eval "$(notte completion zsh)"' >> ~/.zshrc
```

### Fish

```fish
notte completion fish > ~/.config/fish/completions/notte.fish
```

### PowerShell

```powershell
notte completion powershell | Out-String | Invoke-Expression
# Or add to profile:
notte completion powershell >> $PROFILE
```

## Development

After cloning, install git hooks:

```bash
make setup
```

This installs [lefthook](https://github.com/evilmartians/lefthook) pre-commit and pre-push hooks for linting and testing.

### Available Make Targets

```bash
make build      # Build the notte binary
make install    # Install to $GOPATH/bin
make test       # Run tests
make lint       # Run golangci-lint
make fmt        # Format code with goimports and gofumpt
make clean      # Remove built binaries
```

## License

MIT License - see LICENSE file for details.

## Links

- [Notte API Documentation](https://notte.cc/docs)
- [GitHub Repository](https://github.com/salmonumbrella/notte-cli)
- [Report Issues](https://github.com/salmonumbrella/notte-cli/issues)
