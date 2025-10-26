package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	var err error
	detector, err = NewDetector(modelPath)
	if err != nil {
		log.Fatalf("failed to load model: %v", err)
	}
	defer detector.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/detect", handleDetect)
	router.GET("/health", handleHealth)

	log.Printf("server starting on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

