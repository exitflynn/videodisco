package main

import (
	"fmt"
	"image"
	"math"

	ort "github.com/yalue/onnxruntime_go"
)

type FaceEmbedder struct {
	session      *ort.Session[float32]
	inputShape   ort.Shape
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
}

func NewFaceEmbedder(modelPath string) (*FaceEmbedder, error) {
	inputShape := ort.NewShape(1, 3, 112, 112)
	outputShape := ort.NewShape(1, 128)

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

	return &FaceEmbedder{
		session:      session,
		inputShape:   inputShape,
		inputTensor:  inputTensor,
		outputTensor: outputTensor,
	}, nil
}

func (fe *FaceEmbedder) Embed(cropData []float32) ([]float32, error) {
	copy(fe.inputTensor.GetData(), cropData)

	err := fe.session.Run()
	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	emb := fe.outputTensor.GetData()

	normalized := make([]float32, 128)
	copy(normalized, emb)
	l2Normalize(normalized)

	return normalized, nil
}

func (fe *FaceEmbedder) Close() {
	if fe.session != nil {
		fe.session.Destroy()
	}
}

func l2Normalize(emb []float32) {
	var norm float32
	for _, v := range emb {
		norm += v * v
	}
	norm = float32(math.Sqrt(float64(norm)))

	if norm > 0 {
		for i := range emb {
			emb[i] /= norm
		}
	}
}

func cropFace(img image.Image, face Face) ([]float32, error) {
	bounds := img.Bounds()
	x1 := int(face.X1)
	y1 := int(face.Y1)
	x2 := int(face.X2)
	y2 := int(face.Y2)

	if x1 < 0 {
		x1 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	if x2 > bounds.Max.X {
		x2 = bounds.Max.X
	}
	if y2 > bounds.Max.Y {
		y2 = bounds.Max.Y
	}

	cropW := x2 - x1
	cropH := y2 - y1

	if cropW <= 0 || cropH <= 0 {
		return nil, fmt.Errorf("invalid crop size: %dx%d", cropW, cropH)
	}

	scale := float32(112) / float32(maxInt(cropW, cropH))
	newW := int(float32(cropW) * scale)
	newH := int(float32(cropH) * scale)

	padW := (112 - newW) / 2
	padH := (112 - newH) / 2

	data := make([]float32, 3*112*112)

	for y := 0; y < 112; y++ {
		for x := 0; x < 112; x++ {
			ox := x - padW
			oy := y - padH

			var r, g, b float32 = 0.5, 0.5, 0.5

			if ox >= 0 && ox < newW && oy >= 0 && oy < newH {
				srcX := x1 + int(float32(ox)/scale)
				srcY := y1 + int(float32(oy)/scale)

				if srcX >= bounds.Min.X && srcX < bounds.Max.X && srcY >= bounds.Min.Y && srcY < bounds.Max.Y {
					pix := img.At(srcX, srcY)
					cr, cg, cb, _ := pix.RGBA()
					r = float32(cr>>8)/255.0*2.0 - 1.0
					g = float32(cg>>8)/255.0*2.0 - 1.0
					b = float32(cb>>8)/255.0*2.0 - 1.0
				}
			}

			data[0*112*112+y*112+x] = r
			data[1*112*112+y*112+x] = g
			data[2*112*112+y*112+x] = b
		}
	}

	return data, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
