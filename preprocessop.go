package main

import (
	"image"
	"math"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"hawx.me/code/img/contrast"
	"hawx.me/code/img/greyscale"
	"hawx.me/code/img/sharpen"
)

type PreProcessOp interface {
	Apply(img image.Image) image.Image
}

type GreyScaleOp struct{}

func NewGreyScaleOp() *GreyScaleOp {
	return &GreyScaleOp{}
}
func (op *GreyScaleOp) Apply(img image.Image) image.Image {
	return greyscale.Greyscale(img)
}

type ResizeOp struct {
	Scale float64
}

func NewResizeOp(scale float64) *ResizeOp {
	return &ResizeOp{Scale: scale}
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

func NewUnsharpenOp(radius int, sigma float64, amount float64, threshold float64) *UnsharpenOp {
	return &UnsharpenOp{
		Radius:    radius,
		Sigma:     sigma,
		Amount:    amount,
		Threshold: threshold,
	}
}
func NewUnsharpenOpFromParams(params []float64) *UnsharpenOp {
	// Default params, no effect
	unsharpen := []float64{0, 1.0, 0.0, 1.0}
	if len(params) > 4 {
		params = params[:4]
	}
	copy(unsharpen, params)
	return NewUnsharpenOp(int(unsharpen[0]), unsharpen[1], unsharpen[2], unsharpen[3])
}
func NewUnsharpenOpFromString(paramStr string) *UnsharpenOp {
	paramStr = strings.Trim(paramStr, "'\"")
	unsharpenParams := strings.Split(paramStr, ",")
	// Default params, no effect
	unsharpen := []float64{0, 1.0, 0.0, 1.0}
	if len(unsharpenParams) > 4 {
		unsharpenParams = unsharpenParams[:4]
	}
	for index, param := range unsharpenParams {
		param = strings.TrimSpace(param)
		val, err := strconv.ParseFloat(param, 64)
		if err == nil {
			unsharpen[index] = val
		}
	}
	return NewUnsharpenOpFromParams(unsharpen)
}
func (op *UnsharpenOp) Apply(img image.Image) image.Image {
	return sharpen.UnsharpMask(img, op.Radius, op.Sigma, op.Amount, op.Threshold)
}

type ConstrastOp struct {
	Factor float64
}

func NewContrastOp(factor float64) *ConstrastOp {
	return &ConstrastOp{Factor: factor}
}
func (op *ConstrastOp) Apply(img image.Image) image.Image {
	return contrast.Linear(img, op.Factor)
}
