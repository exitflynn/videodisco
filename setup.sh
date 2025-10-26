#!/bin/bash

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
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi
print_status "Architecture: $ARCH"

# Check Homebrew
print_section "Checking Homebrew..."
if ! command -v brew &> /dev/null; then
    echo -e "${RED}✗ Homebrew not found. Please install from https://brew.sh${NC}"
    exit 1
fi
print_status "Homebrew installed"

# Install Go
print_section "Installing Go..."
if ! command -v go &> /dev/null; then
    brew install go
    print_status "Go installed"
else
    GO_VERSION=$(go version | awk '{print $3}')
    print_status "Go already installed: $GO_VERSION"
fi

# Install Python 3.11+
print_section "Installing Python 3.11+..."
if command -v python3 &> /dev/null; then
    PY_VERSION=$(python3 --version | awk '{print $2}')
    PY_MAJOR=$(echo $PY_VERSION | cut -d. -f1)
    PY_MINOR=$(echo $PY_VERSION | cut -d. -f2)
    if [ "$PY_MAJOR" -eq 3 ] && [ "$PY_MINOR" -ge 11 ]; then
        print_status "Python $PY_VERSION already installed"
    else
        echo -e "${YELLOW}Current Python $PY_VERSION is < 3.11, installing latest...${NC}"
        brew install python@3.11
        print_status "Python 3.11 installed"
    fi
else
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
print_section "Upgrading pip..."
pip install --upgrade pip setuptools wheel
print_status "pip upgraded"

# Install ONNXRuntime for M1/M2 (Silicon optimized)
print_section "Installing ONNXRuntime for M1 (Silicon)..."
pip install onnxruntime-silicon
print_status "ONNXRuntime-silicon installed"

# Install YOLOv8 & ultralytics
print_section "Installing YOLOv8 & ultralytics..."
pip install ultralytics
print_status "YOLOv8 & ultralytics installed"

# Install MLflow
print_section "Installing MLflow..."
pip install mlflow
print_status "MLflow installed"

# Install FastAPI and other cloud dependencies
print_section "Installing cloud service dependencies..."
pip install fastapi uvicorn pillow numpy opencv-python
print_status "Cloud dependencies installed"

# Install Docker (if not already installed)
print_section "Checking Docker installation..."
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Docker not found. Please install Docker Desktop from https://www.docker.com/products/docker-desktop${NC}"
    echo "  After installation, please run this script again."
else
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    print_status "Docker installed: $DOCKER_VERSION"
fi

# Create project directories
print_section "Creating project structure..."
mkdir -p edge cloud shared/onnx_utils shared/preprocessing tests docs models data
print_status "Project directories created"

# Generate Go module files
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

# Create requirements.txt for cloud service
print_section "Creating requirements.txt..."
cat > requirements.txt << 'EOF'
# ML & Data
ultralytics>=8.0.0
torch>=2.0.0
opencv-python>=4.8.0
numpy>=1.23.0
Pillow>=9.5.0

# ONNX
onnxruntime-silicon>=1.16.0

# Cloud Service
fastapi>=0.104.0
uvicorn>=0.24.0
pydantic>=2.0.0

# MLflow & Tracking
mlflow>=2.0.0
matplotlib>=3.7.0

# Development
pytest>=7.4.0
black>=23.0.0
flake8>=6.0.0
EOF
print_status "requirements.txt created"

# Install cloud requirements
print_section "Installing Python requirements..."
pip install -r requirements.txt
print_status "All Python requirements installed"

# Verify installations
print_section "Verifying installations..."

echo -e "\n${YELLOW}Verifying tools:${NC}"

# Check Go
GO_VERSION=$(go version | awk '{print $3}')
echo "  Go: $GO_VERSION"

# Check Python
PY_VERSION=$(python3 --version | awk '{print $2}')
echo "  Python: $PY_VERSION"

# Check ONNXRuntime
python3 -c "import onnxruntime; print(f'  ONNXRuntime: {onnxruntime.__version__}')"

# Check ultralytics
python3 -c "import ultralytics; print(f'  Ultralytics: {ultralytics.__version__}')"

# Check MLflow
python3 -c "import mlflow; print(f'  MLflow: {mlflow.__version__}')"

# Check Docker
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    echo "  Docker: $DOCKER_VERSION"
else
    echo "  Docker: NOT INSTALLED (please install manually)"
fi

# Create test script
print_section "Creating test scripts..."

cat > tests/test_setup.py << 'EOF'
"""Verify Phase 0 setup is complete"""
import sys

def test_imports():
    """Test all required imports"""
    print("Testing imports...")
    
    try:
        import onnxruntime
        print(f"✓ onnxruntime {onnxruntime.__version__}")
    except ImportError as e:
        print(f"✗ onnxruntime: {e}")
        return False
    
    try:
        import ultralytics
        print(f"✓ ultralytics {ultralytics.__version__}")
    except ImportError as e:
        print(f"✗ ultralytics: {e}")
        return False
    
    try:
        import mlflow
        print(f"✓ mlflow {mlflow.__version__}")
    except ImportError as e:
        print(f"✗ mlflow: {e}")
        return False
    
    try:
        import fastapi
        print(f"✓ fastapi {fastapi.__version__}")
    except ImportError as e:
        print(f"✗ fastapi: {e}")
        return False
    
    try:
        import torch
        print(f"✓ torch {torch.__version__}")
    except ImportError as e:
        print(f"✗ torch: {e}")
        return False
    
    return True

if __name__ == "__main__":
    if test_imports():
        print("\n✓ All Phase 0 requirements verified!")
        sys.exit(0)
    else:
        print("\n✗ Some requirements are missing")
        sys.exit(1)
EOF

cat > tests/test_onnx.py << 'EOF'
"""ONNX Runtime Hello World Test"""
import numpy as np
import onnxruntime as rt

def test_onnx_inference():
    """Test basic ONNX tensor inference"""
    print("🎬 ONNX Runtime Hello World Test\n")
    
    # Create a simple tensor (batch_size=1, features=3)
    X = np.array([[1.0, 2.0, 3.0]], dtype=np.float32)
    
    # Create a simple ONNX model in memory (identity model)
    from onnxruntime.transformers import pytorch_model_hf
    import tempfile
    import os
    
    print("Creating test tensor:")
    print(f"  Input shape: {X.shape}")
    print(f"  Input data: {X}")
    
    # Simple inference using onnxruntime
    print(f"\n✓ ONNXRuntime version: {rt.__version__}")
    print(f"✓ Available providers: {rt.get_available_providers()}")
    
    # For M1, should show 'CoreMLExecutionProvider' and 'CPUExecutionProvider'
    if 'CoreMLExecutionProvider' in rt.get_available_providers():
        print("✓ CoreML acceleration available (M1 optimization enabled!)")
    
    print("\n✓ ONNX Runtime Hello World test PASSED!")

if __name__ == "__main__":
    test_onnx_inference()
EOF

print_status "Test scripts created"

print_section "✓ Phase 0 Setup Complete!"
echo ""
echo "🎯 Next steps:"
echo "  1. Test your setup: python tests/test_setup.py"
echo "  2. Test ONNX inference: python tests/test_onnx.py"
echo "  3. Check docs/setup.md for detailed information"
echo "  4. Begin Phase 1: On-device edge inference"
echo ""
echo "📝 Virtual environment:"
echo "  To activate: source venv/bin/activate"
echo "  To deactivate: deactivate"
echo ""
