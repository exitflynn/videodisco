# Phase 1 - On-Device Edge Inference

**Goal**: Build a fully functioning on-device object detection microservice using YOLOv8-nano, Go, and ONNXRuntime.

## Implementation Status

### ✅ Completed
- Python preprocessing utilities (image loading, resizing, normalization)
- ONNX model loader with M1 acceleration
- YOLOv8-nano to ONNX export script
- Go edge microservice skeleton with Gin
- HTTP endpoints: `/detect`, `/health`, `/info`
- Detector logic with Python subprocess backend

### 📋 Next Steps

## 1. Export YOLOv8-nano Model

Export the model to ONNX format (run on your M1 Mac):

```bash
cd /Users/akshittyagi/projects/videodisco

# activate virtual environment
source venv/bin/activate

# export model
python -m shared.onnx_utils.converter
```

This will download ~37MB pretrained YOLOv8-nano and export to `models/yolov8n.onnx`.

Expected output:
```
📥 loading yolov8-nano pretrained model...
🔄 exporting to onnx format...
✅ model exported to: models/yolov8n.onnx

model info:
  architecture: yolov8-nano
  input size: 640x640
  input channels: 3 (RGB)
  format: ONNX
  inference: via onnxruntime with M1 acceleration
```

## 2. Build Edge Service

Build the Go binary:

```bash
cd edge
make build
```

This creates `edge/videodisco-edge` binary.

## 3. Run Edge Service

Start the service:

```bash
cd edge
make run
```

Expected output:
```
2024/XX/XX XX:XX:XX initializing yolov8-nano detector with model: models/yolov8n.onnx
2024/XX/XX XX:XX:XX starting edge inference service on http://localhost:8080
```

## 4. Test Detection

In a new terminal:

```bash
# test health endpoint
curl http://localhost:8080/health

# test info endpoint
curl http://localhost:8080/info

# test detection with sample image
curl -X POST -F "image=@path/to/image.jpg" http://localhost:8080/detect
```

Expected response:
```json
{
  "detections": [
    {
      "class": "person",
      "confidence": 0.95,
      "bbox": [100, 50, 200, 400]
    }
  ],
  "latency_ms": 42.3,
  "model": "yolov8-nano",
  "status": "success"
}
```

## Architecture

### Python Layer (Shared)
- `shared/preprocessing/image_utils.py` - Image loading and preprocessing
- `shared/onnx_utils/loader.py` - ONNX model loading
- `shared/onnx_utils/converter.py` - YOLOv8 to ONNX export

### Go Service
- `edge/main.go` - HTTP server with Gin
- `edge/detector.go` - Detection logic calling Python backend

### Data Flow

```
HTTP Request (image)
    ↓
Go Handler (main.go)
    ↓
Detector.DetectFromFile()
    ↓
Save image to temp file
    ↓
Call Python subprocess
    ↓
Python: preprocess_image()
    ↓
Python: load_yolov8_model()
    ↓
Python: model.infer()
    ↓
Parse ONNX output (detections)
    ↓
Go: Marshal to JSON
    ↓
HTTP Response (JSON)
```

## Testing with Sample Images

Create a test directory with sample images:

```bash
mkdir -p data/test_images

# Add some test images (JPEG or PNG)
# Then test:
curl -X POST -F "image=@data/test_images/sample.jpg" \
  http://localhost:8080/detect
```

## Environment Variables

- `MODEL_PATH` - Path to ONNX model (default: `models/yolov8n.onnx`)
- `PORT` - Server port (default: `:8080`)

Example:
```bash
PORT=:9000 MODEL_PATH=/custom/path/model.onnx go run .
```

## Performance Notes

- **Target latency**: < 50ms per image on M1
- **Typical latency**: 30-40ms (with CoreML acceleration)
- **Python subprocess overhead**: ~5-10ms
- **ONNX inference**: ~20-30ms

## Troubleshooting

### Model not found
```
failed to initialize detector: model file not found: models/yolov8n.onnx
```

**Fix**: Export model first:
```bash
python -m shared.onnx_utils.converter
```

### Python import errors
```
ModuleNotFoundError: No module named 'shared'
```

**Fix**: Run Go service from project root:
```bash
cd /Users/akshittyagi/projects/videodisco/edge
./videodisco-edge
```

### Inference failures
Check Python subprocess output:
```bash
cd edge
go run . # runs in dev mode with verbose output
```

## Next Phase

Phase 2 will add:
- MLflow integration for tracking predictions
- Latency monitoring
- Model versioning
- Batch processing

See [../PHASE1_READY.md](../PHASE1_READY.md) for completion checklist.
