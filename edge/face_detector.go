package main

import (
	"fmt"
	"image"
	"sort"
	"strings"

	ort "github.com/yalue/onnxruntime_go"
)

type Face struct {
	X1        float32   `json:"x1"`
	Y1        float32   `json:"y1"`
	X2        float32   `json:"x2"`
	Y2        float32   `json:"y2"`
	Score     float32   `json:"score"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type FaceDetector struct {
	session       *ort.Session[float32]
	inputShape    ort.Shape
	inputTensor   *ort.Tensor[float32]
	outputTensors []*ort.Tensor[float32]
}

func NewFaceDetector(modelPath string) (*FaceDetector, error) {
	inputShape := ort.NewShape(1, 3, 640, 640)

	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}

	outputTensors := make([]*ort.Tensor[float32], 9)
	outputNames := []string{"cls_8", "cls_16", "cls_32", "obj_8", "obj_16", "obj_32", "bbox_8", "bbox_16", "bbox_32"}

	for i, name := range outputNames {
		var outputShape ort.Shape
		if strings.Contains(name, "_8") {
			if strings.Contains(name, "cls") || strings.Contains(name, "obj") {
				outputShape = ort.NewShape(1, 6400, 1)
			} else if strings.Contains(name, "bbox") {
				outputShape = ort.NewShape(1, 6400, 4)
			} else {
				outputShape = ort.NewShape(1, 6400, 10)
			}
		} else if strings.Contains(name, "_16") {
			if strings.Contains(name, "cls") || strings.Contains(name, "obj") {
				outputShape = ort.NewShape(1, 1600, 1)
			} else if strings.Contains(name, "bbox") {
				outputShape = ort.NewShape(1, 1600, 4)
			} else {
				outputShape = ort.NewShape(1, 1600, 10)
			}
		} else if strings.Contains(name, "_32") {
			if strings.Contains(name, "cls") || strings.Contains(name, "obj") {
				outputShape = ort.NewShape(1, 400, 1)
			} else if strings.Contains(name, "bbox") {
				outputShape = ort.NewShape(1, 400, 4)
			} else {
				outputShape = ort.NewShape(1, 400, 10)
			}
		}

		tensor, err := ort.NewEmptyTensor[float32](outputShape)
		if err != nil {
			for j := 0; j < i; j++ {
				outputTensors[j].Destroy()
			}
			inputTensor.Destroy()
			return nil, fmt.Errorf("failed to create output tensor %d: %w", i, err)
		}
		outputTensors[i] = tensor
	}

	session, err := ort.NewSession[float32](
		modelPath,
		[]string{"input"},
		outputNames,
		[]*ort.Tensor[float32]{inputTensor},
		outputTensors,
	)
	if err != nil {
		inputTensor.Destroy()
		for _, tensor := range outputTensors {
			tensor.Destroy()
		}
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &FaceDetector{
		session:       session,
		inputShape:    inputShape,
		inputTensor:   inputTensor,
		outputTensors: outputTensors,
	}, nil
}

func (fd *FaceDetector) Detect(img image.Image, confThreshold float32) ([]Face, error) {
	inputData, err := preprocessFaceDetector(img)
	if err != nil {
		return nil, err
	}

	copy(fd.inputTensor.GetData(), inputData)

	err = fd.session.Run()
	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	faces := []Face{}
	for _, outputTensor := range fd.outputTensors {
		output := outputTensor.GetData()
		decoded := decodeFaces(output, confThreshold)
		faces = append(faces, decoded...)
	}
	faces = nms(faces, 0.4)

	return faces, nil
}

func (fd *FaceDetector) Close() {
	if fd.session != nil {
		fd.session.Destroy()
	}
}

func decodeFaces(output []float32, confThreshold float32) []Face {
	faces := []Face{}

	stride := 14
	numDetections := len(output) / stride

	for i := 0; i < numDetections; i++ {
		offset := i * stride
		if offset+4 > len(output) {
			break
		}

		x := output[offset]
		y := output[offset+1]
		w := output[offset+2]
		h := output[offset+3]
		conf := output[offset+4]

		if conf < confThreshold {
			continue
		}

		x1 := x - w/2
		y1 := y - h/2
		x2 := x + w/2
		y2 := y + h/2

		x1 = max(0, min(640, x1))
		y1 = max(0, min(640, y1))
		x2 = max(0, min(640, x2))
		y2 = max(0, min(640, y2))

		if x2-x1 > 10 && y2-y1 > 10 {
			faces = append(faces, Face{
				X1:    x1,
				Y1:    y1,
				X2:    x2,
				Y2:    y2,
				Score: conf,
			})
		}
	}

	return faces
}

func nms(faces []Face, iouThresh float32) []Face {
	if len(faces) == 0 {
		return faces
	}

	sort.Slice(faces, func(i, j int) bool {
		return faces[i].Score > faces[j].Score
	})

	kept := []Face{}
	for i, face := range faces {
		keep := true
		for _, kf := range kept {
			iou := computeIoU(face, kf)
			if iou > iouThresh {
				keep = false
				break
			}
		}
		if keep {
			kept = append(kept, faces[i])
		}
	}

	return kept
}

func computeIoU(a, b Face) float32 {
	x1 := max(a.X1, b.X1)
	y1 := max(a.Y1, b.Y1)
	x2 := min(a.X2, b.X2)
	y2 := min(a.Y2, b.Y2)

	if x2 < x1 || y2 < y1 {
		return 0
	}

	inter := (x2 - x1) * (y2 - y1)
	areaA := (a.X2 - a.X1) * (a.Y2 - a.Y1)
	areaB := (b.X2 - b.X1) * (b.Y2 - b.Y1)
	union := areaA + areaB - inter

	return inter / union
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func preprocessFaceDetector(img image.Image) ([]float32, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	scale := min(float32(640)/float32(w), float32(640)/float32(h))
	newW := int(float32(w) * scale)
	newH := int(float32(h) * scale)

	padW := (640 - newW) / 2
	padH := (640 - newH) / 2

	data := make([]float32, 3*640*640)

	srcX := bounds.Min.X
	srcY := bounds.Min.Y

	for y := 0; y < 640; y++ {
		for x := 0; x < 640; x++ {
			ox := x - padW
			oy := y - padH

			var r, g, b float32
			if ox >= 0 && ox < newW && oy >= 0 && oy < newH {
				srcPx := img.At(srcX+int(float32(ox)/scale), srcY+int(float32(oy)/scale))
				cr, cg, cb, _ := srcPx.RGBA()
				r = float32(cr>>8) / 255.0
				g = float32(cg>>8) / 255.0
				b = float32(cb>>8) / 255.0
			}

			data[0*640*640+y*640+x] = r
			data[1*640*640+y*640+x] = g
			data[2*640*640+y*640+x] = b
		}
	}

	return data, nil
}
