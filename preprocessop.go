package main

import (
	"image"
	"math"

	"github.com/nfnt/resize"
	"hawx.me/code/img/contrast"
	"hawx.me/code/img/greyscale"
	"hawx.me/code/img/sharpen"
)

type PreProcessOp interface {
	Apply(img image.Image) image.Image
}

type GreyScaleOp struct{}

func (op *GreyScaleOp) Apply(img image.Image) image.Image {
	return greyscale.Greyscale(img)
}

type ResizeOp struct {
	Scale float64
}

func (op *ResizeOp) Apply(img image.Image) image.Image {
	width := uint(img.Bounds().Max.X - img.Bounds().Min.X)
	if op.Scale != 1.0 {
		newWidth := uint(math.Round(float64(width) * op.Scale))
		if newWidth != width {
			img = resize.Resize(newWidth, 0, img, resize.Lanczos3)
		}
	}
	return img
}

type UnsharpenOp struct {
	Radius    int
	Sigma     float64
	Amount    float64
	Threshold float64
}

func (op *UnsharpenOp) Apply(img image.Image) image.Image {
	return sharpen.UnsharpMask(img, op.Radius, op.Sigma, op.Amount, op.Threshold)
}

type ConstrastOp struct {
	Factor float64
}

func (op *ConstrastOp) Apply(img image.Image) image.Image {
	return contrast.Linear(img, op.Factor)
}
