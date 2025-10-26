package main

import (
	"fmt"
	"io"
	"time"

	ort "github.com/yalue/onnxruntime_go"
	"gocv.io/x/gocv"
)

type Detection struct {
	Class      int       `json:"class"`
	Confidence float32   `json:"confidence"`
	BBox       []float32 `json:"bbox"`
}

type Detector struct {
	session     *ort.AdvancedSession
	inputShape  ort.Shape
	inputName   string
	outputNames []string
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
	defer inputTensor.Destroy()

	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create output tensor: %w", err)
	}
	defer outputTensor.Destroy()

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"images"},
		[]string{"output0"},
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensor},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Detector{
		session:     session,
		inputShape:  inputShape,
		inputName:   "images",
		outputNames: []string{"output0"},
	}, nil
}

func (d *Detector) Detect(imageReader io.Reader) ([]Detection, float64, error) {
	start := time.Now()

	img, err := gocv.IMDecode(imageReader, gocv.IMReadColor)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode image: %w", err)
	}
	defer img.Close()

	inputData, err := preprocessImage(img)
	if err != nil {
		return nil, 0, fmt.Errorf("preprocessing failed: %w", err)
	}

	inputTensor, err := ort.NewTensor(d.inputShape, inputData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create tensor: %w", err)
	}
	defer inputTensor.Destroy()

	outputTensor, err := ort.NewEmptyTensor[float32](ort.NewShape(1, 84, 8400))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create output tensor: %w", err)
	}
	defer outputTensor.Destroy()

	err = d.session.Run(
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensor},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("inference failed: %w", err)
	}

	output := outputTensor.GetData()
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
