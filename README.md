# videodisco

edge inference service for object detection using yolov8-nano

## stack

- go with gin framework
- onnxruntime native bindings
- gocv for image processing
- yolov8-nano model

## setup

install onnxruntime:
```bash
brew install onnxruntime
```

export model (run python script once):
```bash
pip install ultralytics
python scripts/export_model.py
mv yolov8n.onnx models/
```

build and run:
```bash
cd edge
make build
make run
```

## usage

detect objects in image:
```bash
curl -X POST -F "image=@photo.jpg" http://localhost:8080/detect
```

health check:
```bash
curl http://localhost:8080/health
```

## performance

typical latency on m1: 20-30ms per image
