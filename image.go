// Giperion November 2018
// [EUREKA] 3.8 Beta

package main

import (
	"image"
	"image/color"
)

type MandelbrotImage struct {
	pixels []color.RGBA
	imageRect image.Rectangle
}

func NewMandelbrotImage (width, height int) MandelbrotImage {
	return MandelbrotImage{
		pixels: make([]color.RGBA, width * height),
		imageRect: image.Rect(0, 0, width, height),
	}
}

func (theImage MandelbrotImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (theImage MandelbrotImage) Bounds() image.Rectangle {
	return theImage.imageRect
}

func (theImage MandelbrotImage) At(x,y int) color.Color {
	if x > theImage.imageRect.Max.X - 1 || y > theImage.imageRect.Max.Y - 1 {
		return color.White
	}

	return theImage.pixels[(y * theImage.imageRect.Max.X) + x]
}

func (theImage *MandelbrotImage) Set(theColor color.RGBA, x, y int) {
	theImage.pixels[(y * theImage.imageRect.Max.X) + x] = theColor
}

func (theImage *MandelbrotImage) Resize(width, height int){
	theImage.pixels = make([]color.RGBA, width * height)
	theImage.imageRect = image.Rect(0, 0, width, height)
}