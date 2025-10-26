package main

import (
	"gocv.io/x/gocv"
	"image"
)

func preprocessImage(img gocv.Mat) ([]float32, error) {
	resized := gocv.NewMat()
	defer resized.Close()

	size := image.Point{X: 640, Y: 640}
	gocv.Resize(img, &resized, size, 0, 0, gocv.InterpolationLinear)

	gocv.CvtColor(resized, &resized, gocv.ColorBGRToRGB)

	normalized := gocv.NewMat()
	defer normalized.Close()

	resized.ConvertTo(&normalized, gocv.MatTypeCV32F)
	normalized.DivideFloat(255.0)

	data := make([]float32, 3*640*640)
	for c := 0; c < 3; c++ {
		for y := 0; y < 640; y++ {
			for x := 0; x < 640; x++ {
				val := normalized.GetVecfAt(y, x)
				data[c*640*640+y*640+x] = val[c]
			}
		}
	}

	return data, nil
}

func postprocess(output []float32, threshold float32) []Detection {
	detections := []Detection{}

	for i := 0; i < 8400; i++ {
		maxConf := float32(0)
		maxClass := 0

		for c := 0; c < 80; c++ {
			conf := output[84*i+4+c]
			if conf > maxConf {
				maxConf = conf
				maxClass = c
			}
		}

		if maxConf < threshold {
			continue
		}

		x := output[84*i+0]
		y := output[84*i+1]
		w := output[84*i+2]
		h := output[84*i+3]

		x1 := x - w/2
		y1 := y - h/2
		x2 := x + w/2
		y2 := y + h/2

		detections = append(detections, Detection{
			Class:      maxClass,
			Confidence: maxConf,
			BBox:       []float32{x1, y1, x2, y2},
		})
	}

	return detections
}

