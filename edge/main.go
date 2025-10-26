package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultModelPath = "../models/yolov8n.onnx"
	defaultPort      = ":8080"
)

var detector *Detector

func main() {
	modelPath := os.Getenv("MODEL_PATH")
	if modelPath == "" {
		modelPath = defaultModelPath
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logx.Infof("loading model: %s", modelPath)

	var err error
	detector, err = NewDetector(modelPath)
	if err != nil {
		logx.Severef("failed to load model: %v", err)
	}
	defer detector.Close()

	logx.Info("model loaded successfully")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/detect", handleDetect)
	router.GET("/health", handleHealth)
	router.GET("/metrics", handleMetrics)

	logx.Infof("server starting on %s", port)
	if err := router.Run(port); err != nil {
		logx.Severef("server error: %v", err)
	}
}
