# videodisco
Smol toy project to play around with doing ONNX inference in Go.

# VideoDisco: Distributed Face Detection & Clustering
> *I see, I learn*

This project is mostly just for toying around with onnxruntime's go interface. Going with an easy to use model, I've explored MobileFaceNet and YOLOV8 which makes the final binary good for real-time face detection and embedding generation. These embeddings can then be sent over to another backend elsewhere to perform operations like clustering without the images themselves leaving the device.

## Running Inference

The ONNX models are already included in the `models/` directory, you can run the Edge Service by:

```bash
cd edge/
go build -o videodisco
./videodisco
```

Or using the Makefile:
```bash
cd edge/
make build
make run
```

The service will start on `http://localhost:8080`

### Step 4: Test the Inference

**Option A: Face Detection & Embedding (Combined)**

```bash
curl -X POST http://localhost:8080/face/process \
  -H "Content-Type: application/json" \
  -d '{
    "image": "your-base64-encoded-image-here"
  }'
```

**Option B: Just Get Embeddings**

```bash
curl -X POST http://localhost:8080/face/embed \
  -H "Content-Type: application/json" \
  -d '{
    "image": "your-base64-encoded-image-here"
  }'
```

**Option C: Get Service Metrics**

```bash
curl http://localhost:8080/metrics
```

**Option D: Health Check**

```bash
curl http://localhost:8080/health
```

### Example Response

```json
{
  "faces": [
    {
      "id": "face_1",
      "bbox": {
        "x": 100,
        "y": 120,
        "width": 80,
        "height": 100
      },
      "confidence": 0.95,
      "embedding": [0.123, 0.456, ..., 0.789]
    }
  ],
  "latency_ms": 45,
  "image_id": "img_001"
}
```

---

## Docker Deployment

### Build Edge Service Docker Image

```bash
cd edge/
docker build -t videodisco-edge:latest .
docker run -p 8080:8080 videodisco-edge:latest
```

