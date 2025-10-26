package main

import (
	"fmt"
	"io"
	"time"

	ort "github.com/yalue/onnxruntime_go"
)

type Detection struct {
	Class      int       `json:"class"`
	Confidence float32   `json:"confidence"`
	BBox       []float32 `json:"bbox"`
}

type Metrics struct {
	ImageName        string        `json:"image_name"`
	NumDetections    int           `json:"num_detections"`
	LatencyMs        float64       `json:"latency_ms"`
	AvgConfidence    float32       `json:"avg_confidence"`
	MaxConfidence    float32       `json:"max_confidence"`
	MinConfidence    float32       `json:"min_confidence"`
	PreprocessTimeMs float64       `json:"preprocess_time_ms"`
	InferenceTimeMs  float64       `json:"inference_time_ms"`
	Detections       []Detection   `json:"detections"`
}

type Detector struct {
	session      *ort.Session[float32]
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	inputShape   ort.Shape
}

func NewDetector(modelPath string) (*Detector, error) {
	ort.SetSharedLibraryPath("/opt/homebrew/lib/libonnxruntime.dylib")

	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("onnxruntime init failed: %w", err)
	}

	inputShape := ort.NewShape(1, 3, 640, 640)
	outputShape := ort.NewShape(1, 84, 8400)

	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}

	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		inputTensor.Destroy()
		return nil, fmt.Errorf("failed to create output tensor: %w", err)
	}

	session, err := ort.NewSession[float32](
		modelPath,
		[]string{"images"},
		[]string{"output0"},
		[]*ort.Tensor[float32]{inputTensor},
		[]*ort.Tensor[float32]{outputTensor},
	)
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Detector{
		session:      session,
		inputTensor:  inputTensor,
		outputTensor: outputTensor,
		inputShape:   inputShape,
	}, nil
}

func (d *Detector) Detect(imageReader io.Reader) ([]Detection, float64, error) {
	start := time.Now()

	inputData, err := preprocessImage(imageReader)
	if err != nil {
		return nil, 0, fmt.Errorf("preprocessing failed: %w", err)
	}

	copy(d.inputTensor.GetData(), inputData)

	err = d.session.Run()
	if err != nil {
		return nil, 0, fmt.Errorf("inference failed: %w", err)
	}

	output := d.outputTensor.GetData()
	detections := postprocess(output, 0.5)

	latency := time.Since(start).Seconds() * 1000
	return detections, latency, nil
}

func (d *Detector) Close() {
	if d.session != nil {
		d.session.Destroy()
	}
	ort.DestroyEnvironment()
}
