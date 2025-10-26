package main

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

var lastFaceMetrics *FaceMetrics

type FaceMetrics struct {
	ImageName string    `json:"image_name"`
	NumFaces  int       `json:"num_faces"`
	LatencyMs float64   `json:"latency_ms"`
	Faces     []Face    `json:"faces"`
	Kind      string    `json:"kind"`
}

var faceDetector *FaceDetector
var faceEmbedder *FaceEmbedder

func handleFaceProcess(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		logx.Errorf("no image in request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image provided"})
		return
	}

	logx.Infof("face process: %s, size: %d", file.Filename, file.Size)

	f, err := file.Open()
	if err != nil {
		logx.Errorf("failed to open file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	start := time.Now()

	img, _, err := image.Decode(f)
	if err != nil {
		logx.Errorf("failed to decode image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode image"})
		return
	}

	faces, err := faceDetector.Detect(img, 0.6)
	if err != nil {
		logx.Errorf("detection failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range faces {
		cropData, err := cropFace(img, faces[i])
		if err != nil {
			logx.Errorf("crop failed for face %d: %v", i, err)
			continue
		}

		emb, err := faceEmbedder.Embed(cropData)
		if err != nil {
			logx.Errorf("embed failed for face %d: %v", i, err)
			continue
		}

		faces[i].Embedding = emb
	}

	latency := time.Since(start).Seconds() * 1000
	metrics := &FaceMetrics{
		ImageName: file.Filename,
		NumFaces:  len(faces),
		LatencyMs: latency,
		Faces:     faces,
		Kind:      "face_process",
	}
	lastFaceMetrics = metrics

	go logFaceToMLflow(metrics)

	logx.Infof("face process complete: %d faces, %.2fms", len(faces), latency)

	c.JSON(http.StatusOK, gin.H{
		"faces":     faces,
		"latency_ms": latency,
	})
}

func handleFaceEmbed(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		logx.Errorf("no image in request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image provided"})
		return
	}

	logx.Infof("face embed: %s", file.Filename)

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	start := time.Now()

	imgData, err := readAllBytes(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		logx.Errorf("failed to decode image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode image"})
		return
	}

	faceBox := Face{X1: 0, Y1: 0, X2: float32(img.Bounds().Dx()), Y2: float32(img.Bounds().Dy())}
	cropData, err := cropFace(img, faceBox)
	if err != nil {
		logx.Errorf("crop failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "crop failed"})
		return
	}

	emb, err := faceEmbedder.Embed(cropData)
	if err != nil {
		logx.Errorf("embed failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	latency := time.Since(start).Seconds() * 1000
	metrics := &FaceMetrics{
		ImageName: file.Filename,
		NumFaces:  1,
		LatencyMs: latency,
		Faces:     []Face{{Embedding: emb}},
		Kind:      "face_embed",
	}
	lastFaceMetrics = metrics

	go logFaceToMLflow(metrics)

	logx.Infof("face embed complete: %.2fms", latency)

	c.JSON(http.StatusOK, gin.H{
		"embedding":  emb,
		"latency_ms": latency,
	})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleFaceMetrics(c *gin.Context) {
	if lastFaceMetrics == nil {
		c.JSON(http.StatusOK, gin.H{"metrics": nil})
		return
	}
	c.JSON(http.StatusOK, lastFaceMetrics)
}

func logFaceToMLflow(metrics *FaceMetrics) {
	payload := map[string]interface{}{
		"image_name":    metrics.ImageName,
		"num_faces":     metrics.NumFaces,
		"latency_ms":    metrics.LatencyMs,
		"kind":          metrics.Kind,
		"faces":         metrics.Faces,
		"model_version": "yunet+mobilefacenet",
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

func readAllBytes(f interface {
	Read(p []byte) (n int, err error)
}) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(f)
	return buf.Bytes(), err
}
