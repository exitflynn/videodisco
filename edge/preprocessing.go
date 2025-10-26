package main

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func preprocessImage(reader io.Reader) ([]float32, error) {
	imgData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	resized := resizeImage(img, 640, 640)
	tensor := imageToTensor(resized)
	
	return tensor, nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcW / width
			srcY := y * srcH / height
			dst.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return dst
}

func imageToTensor(img image.Image) []float32 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	data := make([]float32, 3*width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()

			data[0*width*height+y*width+x] = float32(r>>8) / 255.0
			data[1*width*height+y*width+x] = float32(g>>8) / 255.0
			data[2*width*height+y*width+x] = float32(b>>8) / 255.0
		}
	}

	return data
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
