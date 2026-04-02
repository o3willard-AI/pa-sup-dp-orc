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

# Install system dependencies (Linux only)
if [[ "$(uname)" == "Linux" ]]; then
    echo "Detected Linux, installing system packages..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y \
            libatspi2.0-dev \
            libx11-dev \
            libxrandr-dev \
            libxinerama-dev \
            libxcursor-dev \
            libxi-dev \
            libxtst-dev \
            libgtk-3-dev \
            webkit2gtk-4.1-dev \
            pkg-config
        echo "  ✓ Linux system packages installed"
        
        # Ensure webkit2gtk-4.0.pc symlink exists for Wails
        if [ ! -f /usr/lib/x86_64-linux-gnu/pkgconfig/webkit2gtk-4.0.pc ] && [ -f /usr/lib/x86_64-linux-gnu/pkgconfig/webkit2gtk-4.1.pc ]; then
            sudo ln -sf /usr/lib/x86_64-linux-gnu/pkgconfig/webkit2gtk-4.1.pc /usr/lib/x86_64-linux-gnu/pkgconfig/webkit2gtk-4.0.pc
            echo "  ✓ Created webkit2gtk-4.0.pc symlink"
        fi
    else
        echo "WARNING: apt-get not found, system dependencies may be missing"
    fi
fi

echo ""

# Install Wails CLI
echo "Installing Wails CLI..."
go install github.com/wailsapp/wails/v2/cmd/wails@latest
if ! command -v wails &> /dev/null; then
    echo "ERROR: Wails CLI installation failed"
    exit 1
fi
echo "  ✓ Wails CLI installed"
echo ""

# Download Go dependencies
echo "Downloading Go dependencies..."
go mod download
echo "  ✓ Go dependencies downloaded"
echo ""

# Install frontend dependencies
echo "Installing frontend dependencies..."
if [ ! -d "frontend" ]; then
    echo "ERROR: frontend directory missing"
    exit 1
fi
cd frontend
npm install
echo "  ✓ Frontend dependencies installed"
cd ..

# Install root npm dependencies (for electron-builder fallback)
if [ -f "package.json" ]; then
    echo "Installing root npm dependencies..."
    npm install --only=development
    echo "  ✓ Root npm dependencies installed"
fi

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
