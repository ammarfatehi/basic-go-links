# Go Links - Personal URL Shortening Service

A lightweight, self-hosted "Go Link" service that allows you to create memorable shortcuts for long or complex URLs. Type `localhost:3001/gh` in your browser to be redirected to `https://github.com` (after setting up the shortcut).

## Features

- 🔗 **Simple URL Redirection**: Create memorable shortcuts like `gh` → `https://github.com`
- 🚀 **Fast Response**: < 50ms redirect times for instant navigation
- 💾 **Persistent Storage**: Links survive application restarts via JSON file storage
- 🐳 **Docker Ready**: One-command deployment with automatic restarts
- 🎨 **Clean Web UI**: Simple form to add and view all your shortcuts
- 🔒 **Privacy Focused**: Runs locally, your data stays with you

## Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Clone and start the service
git clone <your-repo-url>
cd basic-go-links
docker compose up -d

# Service will be available at http://go
```

### Option 2: Docker Run

```bash
# Build the image
docker build -t go-links .

# Run the container
docker run -d \
  --name go-links \
  -p 3001:3001 \
  -v $(pwd)/data:/app/data \
  --restart unless-stopped \
  go-links
```

### Option 3: Local Development

```bash
# Install Go 1.24+
go mod tidy
go run main.go

# Service will be available at http://go
```

## Usage

### 1. Access the Web Interface

Open http://go in your browser to see the management interface.

### 2. Add Your First Link

- **Shortcut**: `gh`
- **URL**: `https://github.com`
- Click "Add Link"

### 3. Use Your Shortcut

Type `go/gh` in your browser and you'll be redirected to GitHub!

### 4. Popular Shortcuts to Set Up

```
gh          → https://github.com
gm          → https://gmail.com
drive       → https://drive.google.com
cal         → https://calendar.google.com
docs        → https://docs.google.com
aws         → https://console.aws.amazon.com
localhost   → http://localhost:8080
```

## File Structure

```
basic-go-links/
├── main.go              # Main application code
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Easy deployment configuration
├── data/               # Volume-mounted directory
│   └── links.json      # Your links (auto-created)
├── go.mod              # Go module definition
└── README.md           # This file
```

## Data Storage

Your links are stored in `./data/links.json` and automatically persist across container restarts. The format is:

```json
[
  {
    "shortcut": "gh",
    "url": "https://github.com"
  },
  {
    "shortcut": "gm",
    "url": "https://gmail.com"
  }
]
```

## Advanced Usage

### Custom Port

Edit `docker-compose.yml` to change the port:

```yaml
ports:
  - "8080:3001" # Access via localhost:8080
```

### Multiple Instances

Run multiple instances for different purposes:

```bash
# Work links on port 3001
docker run -d --name work-links -p 3001:3001 -v $(pwd)/work-data:/app/data go-links

# Personal links on port 3002
docker run -d --name personal-links -p 3002:3001 -v $(pwd)/personal-data:/app/data go-links
```

### Backup Your Links

```bash
# Backup
cp data/links.json backup-$(date +%Y%m%d).json

# Restore
cp backup-20240112.json data/links.json
docker compose restart
```

## Browser Integration

For the ultimate experience, set up a bookmark with this JavaScript:

```javascript
javascript: (function () {
  var s = prompt("Shortcut:");
  if (s) {
    location.href = "http://localhost:3001/" + s;
  }
})();
```

Or create browser search engines:

- **Chrome**: Settings → Search engines → Add
- **Keyword**: `go`
- **URL**: `http://localhost:3001/%s`

Then type `go gh` in your address bar!

## Troubleshooting

### Service Not Starting

```bash
# Check logs
docker compose logs go-links

# Common issues:
# - Port 3001 already in use
# - Docker not running
# - Volume permission issues
```

### Links Not Persisting

```bash
# Check if data directory exists and is writable
ls -la data/
chmod 755 data/

# Recreate container with proper volume
docker compose down
docker compose up -d
```

### Performance Issues

```bash
# Check container resources
docker stats go-links

# The service should use minimal resources
# If not, check for infinite redirect loops
```

## Development

### Local Development

```bash
# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Or manual restart
go run main.go
```

### Building from Source

```bash
# Build binary
go build -o go-links .

# Build Docker image
docker build -t go-links .
```

## Security Note

This service is designed for personal, local use. It does not include authentication or HTTPS. Do not expose it to the public internet without additional security measures.

## Contributing

This is a personal project based on the PRD requirements. Feel free to fork and customize for your needs!

## License

MIT License - Feel free to use this for personal or commercial projects.
