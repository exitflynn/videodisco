package main

import (
	"fmt"
	"image"
	"sort"

	ort "github.com/yalue/onnxruntime_go"
)

type Face struct {
	X1        float32 `json:"x1"`
	Y1        float32 `json:"y1"`
	X2        float32 `json:"x2"`
	Y2        float32 `json:"y2"`
	Score     float32 `json:"score"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type FaceDetector struct {
	session       *ort.Session[float32]
	inputShape    ort.Shape
	outputShape   ort.Shape
	inputTensor   *ort.Tensor[float32]
	outputTensor  *ort.Tensor[float32]
}

func NewFaceDetector(modelPath string) (*FaceDetector, error) {
	inputShape := ort.NewShape(1, 3, 320, 320)
	outputShape := ort.NewShape(1, 5600, 14)

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
		[]string{"input"},
		[]string{"output"},
		[]*ort.Tensor[float32]{inputTensor},
		[]*ort.Tensor[float32]{outputTensor},
	)
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &FaceDetector{
		session:      session,
		inputShape:   inputShape,
		outputShape:  outputShape,
		inputTensor:  inputTensor,
		outputTensor: outputTensor,
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

	output := fd.outputTensor.GetData()
	faces := decodeFaces(output, confThreshold)
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

	for i := 0; i < 5600; i++ {
		offset := i * 14
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

		x1 = (x1 + 0.5) * 320
		y1 = (y1 + 0.5) * 320
		x2 = (x2 + 0.5) * 320
		y2 = (y2 + 0.5) * 320

		x1 = max(0, x1)
		y1 = max(0, y1)
		x2 = min(320, x2)
		y2 = min(320, y2)

		faces = append(faces, Face{
			X1:    x1,
			Y1:    y1,
			X2:    x2,
			Y2:    y2,
			Score: conf,
		})
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

	scale := min(float32(320)/float32(w), float32(320)/float32(h))
	newW := int(float32(w) * scale)
	newH := int(float32(h) * scale)

	padW := (320 - newW) / 2
	padH := (320 - newH) / 2

	data := make([]float32, 3*320*320)

	srcX := bounds.Min.X
	srcY := bounds.Min.Y

	for y := 0; y < 320; y++ {
		for x := 0; x < 320; x++ {
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

			data[0*320*320+y*320+x] = r
			data[1*320*320+y*320+x] = g
			data[2*320*320+y*320+x] = b
		}
	}

	return data, nil
}
