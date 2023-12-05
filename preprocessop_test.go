package main

import (
	"math"
	"testing"
)

func withinTolerance(a, b, e float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	if b == 0 {
		return d < e
	}
	return (d / math.Abs(b)) < e
}

func TestNewResizeOp(t *testing.T) {
	resizeOp := NewResizeOp(1.6)
	if !withinTolerance(resizeOp.Scale, 1.6, 1e-12) {
		t.Errorf("resizeOp.Scale = %f; want 1.6", resizeOp.Scale)
	}
}

func TestNewUnsharpenOp(t *testing.T) {
	unsharpenOp := NewUnsharpenOp(4, 1.0, 0.5, 0.06)
	if unsharpenOp.Radius != 4 {
		t.Errorf("unsharpenOp.Radius = %d; want 4", unsharpenOp.Radius)
	}
	if !withinTolerance(unsharpenOp.Sigma, 1.0, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 1.0", unsharpenOp.Sigma)
	}
	if !withinTolerance(unsharpenOp.Amount, 0.5, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 0.5", unsharpenOp.Amount)
	}
	if !withinTolerance(unsharpenOp.Threshold, 0.06, 1e-12) {
		t.Errorf("unsharpenOp.Threshold = %f; want 0.06", unsharpenOp.Threshold)
	}
}

func TestNewUnsharpenOpFromParams(t *testing.T) {
	unsharpenOp := NewUnsharpenOpFromParams([]float64{4, 1.0, 0.5, 0.06})
	if unsharpenOp.Radius != 4 {
		t.Errorf("unsharpenOp.Radius = %d; want 4", unsharpenOp.Radius)
	}
	if !withinTolerance(unsharpenOp.Sigma, 1.0, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 1.0", unsharpenOp.Sigma)
	}
	if !withinTolerance(unsharpenOp.Amount, 0.5, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 0.5", unsharpenOp.Amount)
	}
	if !withinTolerance(unsharpenOp.Threshold, 0.06, 1e-12) {
		t.Errorf("unsharpenOp.Threshold = %f; want 0.06", unsharpenOp.Threshold)
	}
}

func TestNewUnsharpenOpFromParamsTooFewValues(t *testing.T) {
	// Passing fewer than four values, the two last should be the default 0 and 1
	// which result in unsharpen having no effect
	unsharpenOp := NewUnsharpenOpFromParams([]float64{4, 1.0})
	if unsharpenOp.Radius != 4 {
		t.Errorf("unsharpenOp.Radius = %d; want 4", unsharpenOp.Radius)
	}
	if !withinTolerance(unsharpenOp.Sigma, 1.0, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 1.0", unsharpenOp.Sigma)
	}
	if !withinTolerance(unsharpenOp.Amount, 0.0, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 0.0", unsharpenOp.Amount)
	}
	if !withinTolerance(unsharpenOp.Threshold, 1.0, 1e-12) {
		t.Errorf("unsharpenOp.Threshold = %f; want 1.0", unsharpenOp.Threshold)
	}
}

func TestNewUnsharpenOpFromString(t *testing.T) {
	unsharpenOp := NewUnsharpenOpFromString("3,1.2,1.0,0.05")
	if unsharpenOp.Radius != 3 {
		t.Errorf("unsharpenOp.Radius = %d; want 3", unsharpenOp.Radius)
	}
	if !withinTolerance(unsharpenOp.Sigma, 1.2, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 1.2", unsharpenOp.Sigma)
	}
	if !withinTolerance(unsharpenOp.Amount, 1.0, 1e-12) {
		t.Errorf("unsharpenOp.Amount = %f; want 1.0", unsharpenOp.Amount)
	}
	if !withinTolerance(unsharpenOp.Threshold, 0.05, 1e-12) {
		t.Errorf("unsharpenOp.Threshold = %f; want 0.05", unsharpenOp.Threshold)
	}
}

func TestNewContrastOp(t *testing.T) {
	contrastOp := NewContrastOp(1.5)
	if !withinTolerance(contrastOp.Factor, 1.5, 1e-12) {
		t.Errorf("contrastOp.Factor = %f; want 1.5", contrastOp.Factor)
	}
}
