#!/bin/bash
#
# install-deps.sh - Install all PairAdmin development dependencies
#
# Usage: ./scripts/install-deps.sh
#

set -e

echo "=== PairAdmin Dependency Installer ==="
echo ""

# Check for required tools
echo "Checking prerequisites..."
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "ERROR: npm is not installed. Please install Node.js and npm first."
    exit 1
fi

echo "  ✓ Go: $(go version)"
echo "  ✓ npm: $(npm --version)"
echo ""

# Install Wails CLI
echo "Installing Wails CLI..."
go install github.com/wailsapp/wails/v2/cmd/wails@latest
echo "  ✓ Wails CLI installed"
echo ""

# Download Go dependencies
echo "Downloading Go dependencies..."
go mod download
echo "  ✓ Go dependencies downloaded"
echo ""

# Install frontend dependencies
echo "Installing frontend dependencies..."
cd frontend
npm install
echo "  ✓ Frontend dependencies installed"
cd ..
echo ""

# Verify installations
echo "Verifying installations..."
echo "  Wails version: $(wails version)"
echo "  Frontend packages: $(npm list --depth=0 | wc -l) packages installed"
echo ""

echo "=== Installation Complete ==="
echo ""
echo "Next steps:"
echo "  - Run 'wails dev' to start development mode"
echo "  - Run 'wails build' to create production binary"
