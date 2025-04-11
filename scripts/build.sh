#!/bin/bash

# Build for all platforms
# Linux (arm64)
BUILD_DIR=build

# Create build directory if it doesn't exist
mkdir -p $BUILD_DIR

echo "Building for Linux (arm64)"
GOOS=linux GOARCH=arm64 go build -o $BUILD_DIR/secnex-reverse-proxy-linux-arm64 .

# Linux (amd64)
echo "Building for Linux (amd64)"
GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/secnex-reverse-proxy-linux-amd64 .

# Linux (arm)
echo "Building for Linux (arm)"
GOOS=linux GOARCH=arm go build -o $BUILD_DIR/secnex-reverse-proxy-linux-arm .

# MacOS (amd64)
echo "Building for MacOS (amd64)"
GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/secnex-reverse-proxy-darwin-amd64 .

# MacOS (arm64)
echo "Building for MacOS (arm64)"
GOOS=darwin GOARCH=arm64 go build -o $BUILD_DIR/secnex-reverse-proxy-darwin-arm64 .

# Windows (amd64)
echo "Building for Windows (amd64)"
GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/secnex-reverse-proxy-windows-amd64.exe .

# Windows (arm64)
echo "Building for Windows (arm64)"
GOOS=windows GOARCH=arm64 go build -o $BUILD_DIR/secnex-reverse-proxy-windows-arm64.exe .

# Windows (arm)
echo "Building for Windows (arm)"
GOOS=windows GOARCH=arm go build -o $BUILD_DIR/secnex-reverse-proxy-windows-arm.exe .

echo "Build completed successfully!"