# VideoDisco System Architecture

## High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Hybrid Edge-Cloud System                 │
└─────────────────────────────────────────────────────────────┘

┌──────────────────┐                    ┌──────────────────┐
│   Edge Layer     │     Network        │   Cloud Layer    │
│   (Local M1)     │◄──────────────────►│   (FastAPI)      │
│                  │   HTTP/REST        │                  │
│  ┌────────────┐  │                    │  ┌────────────┐  │
│  │ YOLOv8     │  │                    │  │ Classifier │  │
│  │ Detection  │  │                    │  │ / Refiner  │  │
│  │ (Go+ONNX)  │  │                    │  │ (Python)   │  │
│  └────────────┘  │                    │  └────────────┘  │
└──────────────────┘                    └──────────────────┘
        ▲                                        ▲
        │                                        │
        └────────────────┬─────────────────────┘
                         │
                    ┌────▼────┐
                    │  MLflow  │
                    │ Tracking │
                    └──────────┘
```

## Component Architecture

### Phase 0: Setup ✅
- Development environment
- Python 3.11+ with virtual environment
- Go programming environment
- All ML and deployment tools

### Phase 1: Edge Inference (Goal)
- **Go Microservice** (Port 8080)
  - Framework: Gin
  - Input: Image file (JPEG/PNG)
  - Processing:
    - Image preprocessing (resize, normalize)
    - ONNX Runtime inference
    - YOLOv8-nano detection
  - Output: JSON with bounding boxes, class labels, confidence

- **Image Preprocessing Pipeline**
  - Input: Raw image
  - Resize to 640x640 (YOLOv8-nano standard)
  - Normalize: (img / 255.0)
  - Convert to tensor (CHW format)

- **ONNX Model**
  - YOLOv8-nano exported to ONNX
  - M1 acceleration via CoreML providers
  - Input: (1, 3, 640, 640) tensor
  - Output: Detection results

### Phase 2: MLflow Integration
- MLflow tracking server
- Logging predictions per image
- Latency measurement
- Model version tracking

### Phase 3: Cloud Tier
- **Python FastAPI Service** (Port 8000)
  - Receives cropped objects from edge
  - Runs classification or refinement
  - Logs results to MLflow
  - Returns enriched labels

- **Edge-Cloud Communication**
  - Edge detects all objects
  - Filters by confidence threshold
  - Sends interesting objects to cloud
  - Receives refined classifications

### Phase 4: Kubernetes
- Containerization with Docker
- Deployment manifests
- Service discovery
- Scaling policies

### Phase 5: Production Features
- Web UI for image upload
- Retraining pipeline
- Automatic dataset collection
- Advanced monitoring

## Data Flow

### Phase 1: Edge-Only Flow

```
User Image
    │
    ▼
┌─────────────────────┐
│ Go Edge Service     │
│ (Gin HTTP Server)   │
└─────────────────────┘
    │
    ├─► Image Preprocessing
    │   • Resize: 640x640
    │   • Normalize: [0, 1]
    │   • Convert to tensor
    │
    ├─► ONNX Inference
    │   • Load model (onnx-runtime)
    │   • Run on CoreML (M1)
    │   • Get detections
    │
    └─► JSON Response
        {
          "detections": [
            {
              "class": "person",
              "confidence": 0.95,
              "bbox": [x1, y1, x2, y2]
            }
          ],
          "latency_ms": 42.3
        }
```

### Phase 2: With MLflow

```
Edge Service (same as Phase 1)
    │
    ▼
MLflow Tracker
    │
    ├─► Log predictions
    ├─► Log latency
    ├─► Log model version
    └─► Update metrics

MLflow Server
    │
    └─► Dashboard (accessible locally)
```

### Phase 3: Edge-Cloud

```
User Image
    │
    ▼
Edge Service
    │
    ├─► Detect all objects
    │
    ├─► Filter confidence > threshold
    │
    ├─► Crop objects
    │
    └─► Send to Cloud API
            │
            ▼
        Cloud Service (FastAPI)
            │
            ├─► Load classification model
            ├─► Run inference
            └─► Return refined labels
                    │
                    ▼
                Combine Results
                    │
                    └─► Log to MLflow
```

## Technology Stack Details

### Edge Layer
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.21+ | Performance, concurrency |
| Framework | Gin | Latest | HTTP routing, middleware |
| ML Runtime | ONNXRuntime | 1.16+ | Cross-platform inference |
| Optimization | CoreML | Native | M1 hardware acceleration |
| Model Format | ONNX | 1.0+ | Portable model format |

### Model Layer
| Component | Technology | Purpose |
|-----------|-----------|---------|
| Detection | YOLOv8-nano | Real-time object detection |
| Export | Ultralytics | Model export to ONNX |
| Inference | onnxruntime | Hardware-accelerated inference |

### Cloud Layer
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Python | 3.11+ | ML ecosystem access |
| Framework | FastAPI | 0.104+ | Modern async API |
| Server | Uvicorn | 0.24+ | ASGI application server |

### MLOps Layer
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Tracking | MLflow | 2.0+ | Experiment tracking |
| Artifacts | Local filesystem | - | Model & data storage |

### Deployment
| Component | Technology | Purpose |
|-----------|-----------|---------|
| Containerization | Docker | Image packaging |
| Orchestration | Kubernetes/Minikube | Container orchestration |
| Registry | Docker Hub/Local | Image distribution |

## Key Design Decisions

### Why Go for Edge?
- ✅ Fast, compiled language
- ✅ Excellent concurrency (goroutines)
- ✅ Small memory footprint
- ✅ Native ONNXRuntime bindings
- ✅ Easy deployment as single binary

### Why ONNX?
- ✅ Portable across frameworks
- ✅ Hardware acceleration support
- ✅ Optimized inference runtime
- ✅ Industry standard format

### Why Python for Cloud?
- ✅ Rich ML ecosystem
- ✅ FastAPI for async performance
- ✅ Easy to integrate new models
- ✅ Strong data science tools

### Why MLflow?
- ✅ Lightweight, installable locally
- ✅ Tracks experiments without cloud
- ✅ Extensible to cloud later
- ✅ Industry standard

## Directory Structure

```
videodisco/
├── edge/                          # Go microservice
│   ├── main.go                    # Entry point
│   ├── handlers/
│   │   └── detect.go             # /detect endpoint
│   ├── models/
│   │   └── onnx_loader.go        # ONNX model loading
│   ├── preprocessing/
│   │   ├── image.go              # Image loading & preprocessing
│   │   └── tensor.go             # Tensor operations
│   ├── go.mod
│   └── go.sum
│
├── cloud/                         # Python FastAPI service
│   ├── main.py                    # Entry point, app initialization
│   ├── models/
│   │   └── classifier.py         # Classification model
│   ├── api/
│   │   └── routes.py             # API endpoints
│   ├── utils/
│   │   └── mlflow_tracker.py    # MLflow integration
│   └── requirements.txt
│
├── shared/                        # Shared utilities
│   ├── onnx_utils/
│   │   └── converter.py          # YOLOv8 to ONNX export
│   └── preprocessing/
│       └── image_utils.py        # Common preprocessing
│
├── models/                        # Model storage
│   ├── yolov8n.onnx             # YOLOv8-nano ONNX model
│   └── README.md                 # Model info
│
├── data/                          # Test data
│   ├── sample_images/
│   └── coco_classes.txt          # Class labels
│
├── tests/                         # Testing
│   ├── test_setup.py             # Environment verification
│   ├── test_onnx.py              # ONNX runtime test
│   ├── test_edge_api.py          # Edge service tests
│   └── test_cloud_api.py         # Cloud service tests
│
├── docs/                          # Documentation
│   ├── setup.md                  # Setup instructions
│   ├── development.md            # Development guide
│   ├── architecture.md           # This file
│   └── deployment.md             # Deployment guide
│
├── k8s/                           # Kubernetes manifests (Phase 4)
│   ├── edge-deployment.yaml
│   ├── cloud-deployment.yaml
│   ├── mlflow-deployment.yaml
│   └── services.yaml
│
├── docker/                        # Docker files
│   ├── Dockerfile.edge
│   ├── Dockerfile.cloud
│   └── docker-compose.yml
│
├── venv/                          # Python virtual environment
├── requirements.txt               # Python dependencies
├── setup.sh                       # Setup script
├── setup_local.sh                 # Local setup script
├── README.md                      # Project README
└── roadmap.md                     # Project roadmap

```

## Performance Considerations

### M1 Optimization
- ONNXRuntime automatically uses CoreML provider
- Inference typically 2-5x faster than CPU-only
- Memory usage optimized for mobile-class hardware

### Inference Latency
- YOLOv8-nano target: < 50ms per image
- With M1 CoreML: typically 30-40ms
- Network latency (edge→cloud): 1-10ms

### Throughput
- Edge: ~20-30 images/sec on M1
- Cloud: Scales with model complexity

## Security Considerations

- Local-only in Phase 1 (no external exposure)
- MLflow dashboard: Local access only
- HTTPS recommended for Phase 3+ (cloud tier)
- Model signing for reproducibility

## Future Extensions

- GPU support (NVIDIA/AMD)
- Distributed inference
- Model quantization
- Advanced monitoring
- A/B testing framework
- Automated retraining

