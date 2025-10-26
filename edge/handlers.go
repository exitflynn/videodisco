package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

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

	logx.Infof("detection complete: %d objects, %.2fms", len(detections), latency)

	c.JSON(http.StatusOK, gin.H{
		"detections": detections,
		"latency_ms": latency,
	})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
