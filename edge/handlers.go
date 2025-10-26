package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

var lastMetrics *Metrics

func handleDetect(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		logx.Errorf("no image in request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image provided"})
		return
	}

	logx.Infof("received image: %s, size: %d, content-type: %s",
		file.Filename, file.Size, file.Header.Get("Content-Type"))

	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}

	f, err := file.Open()
	if err != nil {
		logx.Errorf("failed to open file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	detections, latency, err := detector.Detect(f)
	if err != nil {
		logx.Errorf("detection failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics := computeMetrics(file.Filename, detections, latency)
	lastMetrics = metrics

	go logToMLflow(metrics)

	logx.Infof("detection complete: %d objects, %.2fms", len(detections), latency)

	c.JSON(http.StatusOK, gin.H{
		"detections": detections,
		"latency_ms": latency,
		"metrics":    metrics,
	})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleMetrics(c *gin.Context) {
	if lastMetrics == nil {
		c.JSON(http.StatusOK, gin.H{"metrics": nil, "message": "no detections yet"})
		return
	}
	c.JSON(http.StatusOK, lastMetrics)
}

func computeMetrics(imageName string, detections []Detection, latencyMs float64) *Metrics {
	metrics := &Metrics{
		ImageName:     imageName,
		NumDetections: len(detections),
		LatencyMs:     latencyMs,
		Detections:    detections,
	}

	if len(detections) > 0 {
		var sum float32
		minConf := float32(1.0)
		maxConf := float32(0.0)

		for _, d := range detections {
			sum += d.Confidence
			if d.Confidence < minConf {
				minConf = d.Confidence
			}
			if d.Confidence > maxConf {
				maxConf = d.Confidence
			}
		}

		metrics.AvgConfidence = sum / float32(len(detections))
		metrics.MinConfidence = minConf
		metrics.MaxConfidence = maxConf
	}

	return metrics
}

func logToMLflow(metrics *Metrics) {
	payload := map[string]interface{}{
		"image_name":    metrics.ImageName,
		"detections":    metrics.Detections,
		"latency_ms":    metrics.LatencyMs,
		"model_version": "yolov8n",
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://127.0.0.1:5001/log", "application/json", bytes.NewReader(body))
	if err != nil {
		logx.Errorf("failed to log to mlflow: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logx.Infof("mlflow logging returned status: %d", resp.StatusCode)
	}
}
