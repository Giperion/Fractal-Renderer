package main

import "testing"

func TestMandelbrot(t *testing.T) {
	var viewParam MandelbrotView = NewMandelbrotView()

	color := Mandelbrot(50, 45, viewParam)
	var ExpectedColor Color = Color{22, 22, 22}
	if color != ExpectedColor {
		t.Fail()
	}
}

func TestMandelbrot_BigNumber(t *testing.T) {
	var viewParam MandelbrotView = NewMandelbrotView()

	color := Mandelbrot_BigNumber(50, 45, viewParam)
	var ExpectedColor Color = Color{22, 22, 22}
	if color != ExpectedColor {
		t.Fail()
	}
}