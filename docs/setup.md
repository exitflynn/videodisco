# Phase 0 Setup Guide

Complete step-by-step setup instructions for the VideoDisco project.

## Overview

Phase 0 prepares your M1 MacBook with all necessary tools and libraries for edge-cloud AI inference.

## Prerequisites

- M1/M2 MacBook (arm64 architecture)
- macOS 11.0 or later
- Homebrew installed (https://brew.sh)
- ~4GB free disk space
- Internet connection for downloads

## Step 1: Install Homebrew (if not already installed)

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Step 2: Install Go

Go is needed for the edge microservice written in Go with Gin framework.

```bash
brew install go
```

Verify:
```bash
go version
```

Should output something like: `go version go1.21.0 darwin/arm64`

## Step 3: Install Python 3.11+

Python is needed for ML tools, MLflow, and cloud services.

```bash
brew install python@3.11
```

Create an alias if needed:
```bash
echo 'alias python3.11="/usr/local/bin/python3.11"' >> ~/.zshrc
source ~/.zshrc
```

Verify:
```bash
python3 --version
```

Should output: `Python 3.11.x` or higher

## Step 4: Install Docker

Docker is needed for containerization and Kubernetes deployment (Phase 4).

Download Docker Desktop for Mac (Apple Silicon):
https://docs.docker.com/desktop/install/mac-install/

After installation, verify:
```bash
docker --version
```

## Step 5: Clone/Setup VideoDisco Project

Navigate to your project directory:
```bash
cd /Users/akshittyagi/projects/videodisco
```

## Step 6: Run Automated Setup Script

The setup script automates the rest of the installation:

```bash
chmod +x setup.sh
./setup.sh
```

This script will:
- Create Python virtual environment
- Install ONNXRuntime (Silicon optimized)
- Install YOLOv8 & ultralytics
- Install MLflow
- Install FastAPI and dependencies
- Set up Go modules
- Create project structure

## Step 7: Manual Installation (Alternative)

If you prefer manual installation:

### Create Virtual Environment

```bash
python3 -m venv venv
source venv/bin/activate
```

### Upgrade pip

```bash
pip install --upgrade pip setuptools wheel
```

### Install ONNXRuntime for M1

This is the Silicon-optimized version critical for performance:

```bash
pip install onnxruntime-silicon
```

### Install YOLOv8 & Ultralytics

```bash
pip install ultralytics
```

This includes PyTorch automatically configured for M1.

### Install MLflow

```bash
pip install mlflow
```

### Install Cloud Service Dependencies

```bash
pip install fastapi uvicorn pillow numpy opencv-python
```

### Install Development Tools

```bash
pip install pytest black flake8
```

## Step 8: Verify Installation

Test Python imports:

```bash
python tests/test_setup.py
```

Expected output:
```
Testing imports...
✓ onnxruntime 1.16.0
✓ ultralytics 8.0.26
✓ mlflow 2.9.0
✓ fastapi 0.104.1
✓ torch 2.0.1

✓ All Phase 0 requirements verified!
```

Test ONNX inference:

```bash
python tests/test_onnx.py
```

Expected output:
```
🎬 ONNX Runtime Hello World Test

Creating test tensor:
  Input shape: (1, 3)
  Input data: [[1. 2. 3.]]

✓ ONNXRuntime version: 1.16.0
✓ Available providers: ['CoreMLExecutionProvider', 'CPUExecutionProvider']
✓ CoreML acceleration available (M1 optimization enabled!)

✓ ONNX Runtime Hello World test PASSED!
```

## Step 9: Go Module Setup

Initialize Go module for edge service:

```bash
cd edge
go mod init videodisco/edge
go get github.com/gin-gonic/gin
cd ..
```

Verify:
```bash
cat edge/go.mod
```

## Virtual Environment Usage

Always activate the virtual environment before working:

```bash
# Activate
source venv/bin/activate

# Your prompt will show: (venv) $

# Deactivate when done
deactivate
```

## Directory Structure Created

After setup, you should have:

```
videodisco/
├── venv/                    # Python virtual environment
├── edge/                    # Go edge service
│   ├── go.mod
│   └── go.sum
├── cloud/                   # Python cloud service (Phase 3)
├── shared/                  # Shared utilities
│   ├── onnx_utils/
│   └── preprocessing/
├── tests/                   # Test suite
│   ├── test_setup.py
│   └── test_onnx.py
├── docs/                    # Documentation
├── models/                  # Model storage
├── data/                    # Test images/data
├── requirements.txt         # Python dependencies
├── setup.sh                 # This setup script
├── README.md                # Project README
└── roadmap.md              # Project roadmap
```

## Troubleshooting

### Issue: "Python 3.11 not found"

Solution:
```bash
brew install python@3.11
python3.11 -m venv venv
```

### Issue: "onnxruntime-silicon not found"

Make sure you're in the virtual environment:
```bash
source venv/bin/activate
pip install onnxruntime-silicon
```

### Issue: "Go command not found"

```bash
brew install go
# Add to ~/.zshrc or ~/.bash_profile if needed:
export PATH=$PATH:/usr/local/go/bin
```

### Issue: "Docker command not found"

Download Docker Desktop from:
https://docs.docker.com/desktop/install/mac-install/

Then restart terminal.

### Issue: torch/ultralytics fails to install

This is usually a network issue. Try:
```bash
pip install --upgrade pip setuptools wheel
pip install --no-cache-dir ultralytics
```

### M1 Performance Note

ONNXRuntime-silicon uses Apple's CoreML framework for hardware acceleration. If setup.py shows:
```
✓ Available providers: ['CoreMLExecutionProvider', 'CPUExecutionProvider']
```

You're getting M1 acceleration! 🎉

## Next Steps

Once Phase 0 is complete:

1. Read [Phase 1 Guide](../roadmap.md) for on-device edge inference
2. Export YOLOv8-nano to ONNX format
3. Build the Go edge microservice with Gin
4. Implement inference pipeline

## References

- [Go Installation](https://golang.org/doc/install)
- [Python Official](https://www.python.org/downloads/)
- [ONNXRuntime Documentation](https://onnxruntime.ai/)
- [Ultralytics YOLOv8](https://docs.ultralytics.com/)
- [MLflow Documentation](https://mlflow.org/docs/latest/index.html)
- [Gin Web Framework](https://gin-gonic.com/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
