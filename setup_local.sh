#!/bin/bash
#
# VideoDisco Phase 0 Local Setup Script
# Run this on your M1 MacBook to set up the environment
# Note: This script requires internet access and should be run in your local terminal
#

set -e

echo "🎬 VideoDisco Phase 0 Setup"
echo "================================"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_section() {
    echo -e "\n${YELLOW}→ $1${NC}"
}

# Check if on M1/M2 Mac
print_section "Checking system architecture..."
ARCH=$(uname -m)
if [ "$ARCH" != "arm64" ]; then
    echo -e "${RED}⚠ Warning: This setup is optimized for M1/M2 Macs (arm64).${NC}"
    echo "  Detected: $ARCH"
fi
print_status "Architecture: $ARCH"

# Check Homebrew
print_section "Checking Homebrew..."
if ! command -v brew &> /dev/null; then
    echo -e "${RED}✗ Homebrew not found. Please install from https://brew.sh${NC}"
    exit 1
fi
print_status "Homebrew installed"

# Create project directories
print_section "Creating project structure..."
mkdir -p edge cloud shared/onnx_utils shared/preprocessing tests docs models data
print_status "Project directories created"

# Check Go
print_section "Go Installation..."
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    brew install go
    print_status "Go installed"
else
    GO_VERSION=$(go version | awk '{print $3}')
    print_status "Go already installed: $GO_VERSION"
fi

# Check Python 3.11+
print_section "Python 3.11+ Installation..."
if command -v python3 &> /dev/null; then
    PY_VERSION=$(python3 --version | awk '{print $2}')
    PY_MAJOR=$(echo $PY_VERSION | cut -d. -f1)
    PY_MINOR=$(echo $PY_VERSION | cut -d. -f2)
    if [ "$PY_MAJOR" -eq 3 ] && [ "$PY_MINOR" -ge 11 ]; then
        print_status "Python $PY_VERSION already installed"
    else
        echo "Installing Python 3.11..."
        brew install python@3.11
        print_status "Python 3.11 installed"
    fi
else
    echo "Installing Python 3.11..."
    brew install python@3.11
    print_status "Python 3.11 installed"
fi

# Create virtual environment
print_section "Setting up Python virtual environment..."
if [ ! -d "venv" ]; then
    python3 -m venv venv
    print_status "Virtual environment created"
else
    print_status "Virtual environment already exists"
fi

# Activate virtual environment
source venv/bin/activate
print_status "Virtual environment activated"

# Upgrade pip
print_section "Upgrading pip, setuptools, wheel..."
pip install --upgrade pip setuptools wheel
print_status "pip upgraded"

# Install all Python dependencies
print_section "Installing Python packages..."
echo "  Installing: onnxruntime-silicon"
pip install onnxruntime

echo "  Installing: ultralytics"
pip install ultralytics

echo "  Installing: mlflow"
pip install mlflow

echo "  Installing: fastapi, uvicorn"
pip install fastapi uvicorn

echo "  Installing: pytorch-related & vision packages"
pip install opencv-python pillow

echo "  Installing: development tools"
pip install pytest black flake8

print_status "All Python packages installed"

# Check Docker
print_section "Checking Docker..."
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠ Docker not found.${NC}"
    echo "  Please install Docker Desktop from:"
    echo "  https://docs.docker.com/desktop/install/mac-install/"
else
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    print_status "Docker installed: $DOCKER_VERSION"
fi

# Set up Go modules
print_section "Setting up Go modules..."
if [ ! -f "edge/go.mod" ]; then
    cd edge
    go mod init videodisco/edge
    go get github.com/gin-gonic/gin
    cd ..
    print_status "Go modules initialized"
else
    print_status "Go modules already exist"
fi

print_section "✓ Phase 0 Setup Complete!"
echo ""
echo "🎯 Next steps:"
echo "  1. Verify setup: python tests/test_setup.py"
echo "  2. Test ONNX: python tests/test_onnx.py"
echo "  3. Read: docs/setup.md for detailed information"
echo "  4. Begin Phase 1: On-device edge inference"
echo ""
echo "📝 Virtual environment:"
echo "  To activate: source venv/bin/activate"
echo "  To deactivate: deactivate"
echo ""
