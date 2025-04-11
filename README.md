# Secnex Reverse Proxy

A simple reverse proxy with API management and certificate management.

## Features

- Reverse Proxy with HTTP and HTTPS support
- In-Memory configuration cache
- REST API for managing the proxy configuration
- Self-Signed Certificate Generation

## Prerequisites

- Go 1.21 or higher

## Installation

1. Clone the repository
2. Run `go mod download`
3. Build the project with `go build`

## Usage

### Starting the server

```bash
# Without SSL
./secnex-reverse-proxy

# Mit Self-Signed Zertifikat
USE_SELF_SIGNED=true ./secnex-reverse-proxy
```

### API-Endpunkte

The API is available on port 8081:

- `POST /config` - Add a new proxy configuration
- `GET /config?host=example.com` - Get configuration
- `DELETE /config?host=example.com` - Delete configuration
- `GET /health` - Health-Check

## Security

- The API has no authentication implemented
- Self-Signed Certificates are only intended for testing purposes
- For production environments, you should implement additional security measures
