#!/bin/bash

# Build for all platforms
# Linux (arm64)
echo "Building for Linux (arm64)"
GOOS=linux GOARCH=arm64 go build -o secnex-reverse-proxy main.go

# Linux (amd64)
echo "Building for Linux (amd64)"
GOOS=linux GOARCH=amd64 go build -o secnex-reverse-proxy main.go

# Linux (arm)
echo "Building for Linux (arm)"
GOOS=linux GOARCH=arm go build -o secnex-reverse-proxy main.go

# MacOS (amd64)
echo "Building for MacOS (amd64)"
GOOS=darwin GOARCH=amd64 go build -o secnex-reverse-proxy main.go

# MacOS (arm64)
echo "Building for MacOS (arm64)"
GOOS=darwin GOARCH=arm64 go build -o secnex-reverse-proxy main.go

# Windows (amd64)
echo "Building for Windows (amd64)"
GOOS=windows GOARCH=amd64 go build -o secnex-reverse-proxy.exe main.go

# Windows (arm64)
echo "Building for Windows (arm64)"
GOOS=windows GOARCH=arm64 go build -o secnex-reverse-proxy.exe main.go

# Windows (arm)
echo "Building for Windows (arm)"
GOOS=windows GOARCH=arm go build -o secnex-reverse-proxy.exe main.go

echo "Build completed successfully!"