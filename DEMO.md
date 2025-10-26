# videodisco demo guide

## prerequisites

ensure onnxruntime is installed:
```bash
brew install onnxruntime
```

## step 1: export model

export yolov8-nano to onnx (run once):
```bash
cd /Users/akshittyagi/projects/videodisco
pip install ultralytics
python scripts/export_model.py
mv yolov8n.onnx models/
```

this downloads the pretrained model (~6MB) and exports it to onnx format.

## step 2: start the server

```bash
cd edge
./videodisco-edge
```

you should see:
```
server starting on :8080
```

## step 3: test the api

### health check

```bash
curl http://localhost:8080/health
```

response:
```json
{"status":"ok"}
```

### detect objects

download a test image:
```bash
curl -o test.jpg https://ultralytics.com/images/bus.jpg
```

run detection:
```bash
curl -X POST -F "image=@test.jpg" http://localhost:8080/detect
```

expected response:
```json
{
  "detections": [
    {
      "class": 5,
      "confidence": 0.89,
      "bbox": [50.2, 100.5, 400.3, 550.8]
    },
    {
      "class": 0,
      "confidence": 0.92,
      "bbox": [120.1, 200.3, 300.5, 450.2]
    }
  ],
  "latency_ms": 28.5
}
```

class ids map to coco dataset classes:
- 0: person
- 1: bicycle
- 2: car
- 3: motorcycle
- 5: bus
- etc.

## performance

typical latency on m1 with onnxruntime:
- preprocessing: 2-5ms
- inference: 15-25ms
- postprocessing: 1-2ms
- **total: 20-30ms per image**

## troubleshooting

### model not found
```
failed to load model: model file not found
```
solution: export the model first (step 1)

### onnxruntime library not found
```
onnxruntime init failed
```
solution: 
```bash
brew install onnxruntime
# if still fails, set library path
export DYLD_LIBRARY_PATH=/opt/homebrew/lib:$DYLD_LIBRARY_PATH
```

### port already in use
```
server error: listen tcp :8080: bind: address already in use
```
solution: use different port
```bash
PORT=:8081 ./videodisco-edge
```

## next steps

- integrate with your application via http api
- adjust confidence threshold in detector.go (currently 0.5)
- add nms (non-maximum suppression) for overlapping boxes
- deploy as systemd service or docker container

