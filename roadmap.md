Hybrid Edge-Cloud Photo AI MVP Roadmap
Phase 0 – Setup & Basic Knowledge

Goal: Get the environment ready and understand the tools.

Tasks:

Install Go, Gin, Docker on M1 MacBook

Install Python 3.11+

Install ONNXRuntime for M1: onnxruntime-silicon

Install YOLOv8 + ultralytics (pip install ultralytics)

Familiarize yourself with MLflow basics (tracking experiments)

Deliverable:

Ability to run Python scripts for ONNX models locally

Run onnxruntime hello-world tensor inference

Phase 1 – On-device Edge Inference (YOLOv8-nano)

Goal: Build a working on-device object detection microservice.

Tasks:

Export YOLOv8-nano pre-trained or fine-tuned model to ONNX.

Build Go microservice with Gin:

POST endpoint /detect that accepts an image (JPEG/PNG).

Preprocess image: resize → normalize → convert to tensor.

Run ONNXRuntime inference.

Return bounding boxes + class labels in JSON.

Test locally on several images.

Optional Enhancements:

Pretty-print bounding boxes on images with Go imaging library.

Measure and log inference latency.

Deliverable:

Fully functioning edge inference microservice

CLI or web-based demo: upload a photo → get objects detected with bounding boxes

Can show this to anyone and it works without cloud

Why impressive:

Edge inference is production-relevant

Go + ONNXRuntime shows low-level system knowledge

Phase 2 – Edge Model Enhancements & MLflow Logging

Goal: Make the edge pipeline robust and trackable.

Tasks:

Integrate MLflow tracking in Python wrapper or Go microservice:

Log each image’s predictions (JSON + inference latency)

Track model version used

Add small preprocessing options: resize, normalize, optional augmentations

Optionally simulate edge-device constraints:

CPU-only inference

Artificial latency or batch processing

Deliverable:

Edge service logs predictions + metrics in MLflow

Dashboard shows inference time and sample outputs

Still fully functional on MacBook

Why impressive:

Shows awareness of MLOps & tracking, even without cloud

Phase 3 – Lightweight Cloud Model

Goal: Demonstrate tiered inference with cloud-side processing.

Tasks:

Choose a simple “heavier” cloud task:

Fine-grained classification of objects detected on edge

Optional image style transfer or enhancement

Build Python FastAPI microservice:

Accept cropped object images from edge

Run inference using larger pre-trained ONNX model

Return refined labels or enhanced image

Edge Go service sends selected crops to cloud API

Log cloud predictions + latency to MLflow

Deliverable:

Edge → Cloud inference pipeline works end-to-end

JSON outputs show “edge detection + cloud refinement”

MLflow dashboard tracks both layers

Why impressive:

Shows real edge-cloud orchestration

You demonstrate multi-tier inference and logging

Phase 4 – Optional Kubernetes Deployment

Goal: Show production readiness / cloud orchestration.

Tasks:

Containerize both edge and cloud services with Docker

Run cloud service in Minikube or AWS Free Tier EC2:

Use Kubernetes Deployment + Service

Optional: simulate multiple edge clients sending requests

Optional: include MLflow server container

Deliverable:

Deployed edge-cloud system running locally or on small cloud node

Can scale cloud service pods (simulate)

Dashboard tracks metrics in MLflow

Why impressive:

Demonstrates Kubernetes + containerized ML pipelines

Shows end-to-end production mindset

Phase 5 – Stretch Goals

Goal: Add bells and whistles for extra impact.

Ideas:

Automatic selection of “interesting objects” to send to cloud (thresholding)

Simulated edge device latency metrics

Optional retraining pipeline: cloud logs → dataset → fine-tune YOLOv8-nano

Optional web UI to upload images, view edge+cloud results, latency, and metrics

Deliverable:

Full hybrid edge-cloud AI pipeline

Demonstrates MLOps, deployment, tiered inference, monitoring

Works end-to-end on your MacBook + optional cloud

✅ Key Advantages of This Incremental MVP Approach

Phase 1 alone is impressive: Edge detection + Go + ONNXRuntime.

Phase 2 adds professional polish: MLflow tracking & latency measurement.

Phase 3 shows system design skills: tiered inference + cloud integration.

Phase 4 shows production-readiness: Docker + Kubernetes orchestration.

Phase 5 demonstrates foresight: retraining, monitoring, and automation.

Even if you only implement Phases 1–2, you already have a strong portfolio demo. Phases 3–4 make it full-stack production-grade, while Phase 5 gives extra “wow factor.”