package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultFaceDetectorPath = "../models/yunet.onnx"
	defaultFaceEmbedderPath = "../models/mobilefacenet.onnx"
	defaultPort             = ":8080"
)

func main() {
	faceDetectorPath := os.Getenv("FACE_DETECTOR_MODEL")
	if faceDetectorPath == "" {
		faceDetectorPath = defaultFaceDetectorPath
	}

	faceEmbedderPath := os.Getenv("FACE_EMBEDDER_MODEL")
	if faceEmbedderPath == "" {
		faceEmbedderPath = defaultFaceEmbedderPath
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logx.Infof("loading face detector: %s", faceDetectorPath)
	var err error
	faceDetector, err = NewFaceDetector(faceDetectorPath)
	if err != nil {
		logx.Severef("failed to load face detector: %v", err)
	}
	defer faceDetector.Close()

	logx.Infof("loading face embedder: %s", faceEmbedderPath)
	faceEmbedder, err = NewFaceEmbedder(faceEmbedderPath)
	if err != nil {
		logx.Severef("failed to load face embedder: %v", err)
	}
	defer faceEmbedder.Close()

	logx.Info("models loaded successfully")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.POST("/face/process", handleFaceProcess)
	router.POST("/face/embed", handleFaceEmbed)
	router.GET("/health", handleHealth)
	router.GET("/metrics", handleFaceMetrics)

	logx.Infof("server starting on %s", port)
	if err := router.Run(port); err != nil {
		logx.Severef("server error: %v", err)
	}
}
