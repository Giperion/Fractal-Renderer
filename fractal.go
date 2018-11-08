// Giperion October - November 2018
// [EUREKA] 3.8 Beta
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/big"
	"strconv"
)

type MandelbrotView struct {
	cameraX    big.Float
	cameraY    big.Float
	cameraZoom big.Float
	width      int
	height     int
	iteration  int
}

func (this *MandelbrotView) SaveToString(outString *string){
	*outString += "cameraX = "
	*outString += this.cameraX.Text('e', int (this.cameraX.Prec()))
	*outString += "\n"

	*outString += "cameraY = "
	*outString += this.cameraY.Text('e', int (this.cameraY.Prec()))
	*outString += "\n"

	*outString += "cameraZoom = "
	*outString += this.cameraZoom.Text('e', int (this.cameraZoom.Prec()))
	*outString += "\n"

	/* I decided to not save window size
	*outString += "width = "
	*outString += strconv.Itoa(this.width)
	*outString += "\n"

	*outString += "height = "
	*outString += strconv.Itoa(this.height)
	*outString += "\n"
	*/

	*outString += "iteration = "
	*outString += strconv.Itoa(this.iteration)
	*outString += "\n"
}

func NewMandelbrotView () MandelbrotView {
	var viewParam MandelbrotView
	viewParam.cameraZoom.SetFloat64(1.0)
	viewParam.width = 600
	viewParam.height = 400
	viewParam.iteration = 100
	return viewParam
}


// Perform mandelbrot fractal for specified pixel (cameraX, cameraY) and settings (viewParams)
// The result is a RGBA pixel, where Alpha = 255
func Mandelbrot(x, y int, viewParams MandelbrotView) color.RGBA {
	const centerX float64 = -0.5
	const centerY float64 = 0.5

	var fltScale float64
	fltScale, _ = viewParams.cameraZoom.Float64()

	var MinimumResultX float64 = centerX - fltScale
	var MaximumResultX float64 = centerX + fltScale
	var MinimumResultY float64 = centerY - fltScale
	var MaximumResultY float64 = centerY + fltScale

	var PixelWidth float64 = (MaximumResultX - MinimumResultX) / float64(viewParams.width)
	var PixelHeight float64 = (MaximumResultY - MinimumResultY) / float64(viewParams.height)
	var Iterations int
	var MaxIterations int = viewParams.iteration

	var Zx float64
	var Zy float64
	var ZxX2 float64
	var ZyX2 float64

	const EscapeRadius float64 = 2.0
	const EscapeRadiusX2 float64 = EscapeRadius * EscapeRadius

	var fltUserX float64
	var fltUserY float64
	fltUserX, _ = viewParams.cameraX.Float64()
	fltUserY, _ = viewParams.cameraY.Float64()

	var ResultX float64 = (MinimumResultX + PixelWidth*float64(x)) + fltUserX
	var ResultY float64 = (MinimumResultY + PixelHeight*float64(y)) + fltUserY

	if math.Abs(ResultY) < PixelHeight/2 {
		ResultY = 0.0
	}

	for ; Iterations < MaxIterations && ((ZxX2 + ZyX2) < EscapeRadiusX2); Iterations++ {
		Zy = 2*Zx*Zy + ResultY
		Zx = ZxX2 - ZyX2 + ResultX
		ZxX2 = Zx * Zx
		ZyX2 = Zy * Zy
	}

	var rawValue float64 = float64 (Iterations % MaxIterations)

	var outColor2 color.RGBA = HSV2RGB(float64(Iterations % 361), 0.01, rawValue)
	return outColor2
}

func PrintBigNumber(num *big.Float) {
	theNum, _ := num.Float64()
	fmt.Println(theNum)
}

var centerX *big.Float = big.NewFloat(-0.5)
var centerY *big.Float = big.NewFloat( 0.5)
var EscapeRadius *big.Float = big.NewFloat(2.0)
var EscapeRadiusX2 *big.Float = big.NewFloat(2.0 * 2.0)

func Mandelbrot_BigNumber(x, y int, viewParams MandelbrotView) color.RGBA {
	var UserWidth *big.Float = big.NewFloat(float64(viewParams.width))
	var UserHeight *big.Float = big.NewFloat(float64(viewParams.height))

	var MinimumResultX big.Float
	MinimumResultX.Sub(centerX, &viewParams.cameraZoom)
	var MaximumResultX big.Float
	MaximumResultX.Add(centerX, &viewParams.cameraZoom)
	var MinimumResultY big.Float
	MinimumResultY.Sub(centerY, &viewParams.cameraZoom)
	var MaximumResultY big.Float
	MaximumResultY.Add(centerY, &viewParams.cameraZoom)

	//var PixelWidth float64 = (MaximumResultX - MinimumResultX) / float64(viewParams.width)
	var PixelWidth big.Float
	PixelWidth.Sub(&MaximumResultX, &MinimumResultX)
	PixelWidth.Quo(&PixelWidth, UserWidth)
	//var PixelHeight float64 = (MaximumResultY - MinimumResultY) / float64(viewParams.height)
	var PixelHeight big.Float
	PixelHeight.Sub(&MaximumResultY, &MinimumResultY)
	PixelHeight.Quo(&PixelHeight, UserHeight)
	var Iterations int
	var MaxIterations int = viewParams.iteration

	var Zx big.Float
	var Zy big.Float
	var ZxX2 big.Float
	var ZyX2 big.Float

	//var ResultX float64 = (MinimumResultX + PixelWidth*float64(cameraX)) + viewParams.cameraX
	var ResultX big.Float
	ResultX.Mul(&PixelWidth, big.NewFloat(float64(x)))
	ResultX.Add(&MinimumResultX, &ResultX)
	ResultX.Add(&ResultX, &viewParams.cameraX)
	//var ResultY float64 = (MinimumResultY + PixelHeight*float64(cameraY)) + viewParams.cameraY
	var ResultY big.Float
	ResultY.Mul(&PixelHeight, big.NewFloat(float64(y)))
	ResultY.Add(&MinimumResultY, &ResultY)
	ResultY.Add(&ResultY, &viewParams.cameraY)

	//if math.Abs(ResultY) < PixelHeight/2 {
	//	ResultY = 0.0
	//}
	var AbsResultY big.Float
	AbsResultY.Abs(&ResultY)
	if AbsResultY.Cmp(big.NewFloat(0.0).Quo(&PixelHeight, EscapeRadius)) == -1{
		ResultY.SetFloat64(0.0)
	}

	//for ; Iterations < MaxIterations && ((ZxX2 + ZyX2) < EscapeRadiusX2); Iterations++ {
	//	Zy = 2*Zx*Zy + ResultY
	//	Zx = ZxX2 - ZyX2 + ResultX
	//	ZxX2 = Zx * Zx
	//	ZyX2 = Zy * Zy
	//}

	var TempFloat big.Float
	var TempFloat2 big.Float

	for ; Iterations < MaxIterations; Iterations++ {
		TempFloat.Add(&ZxX2, &ZyX2)
		if TempFloat.Cmp(EscapeRadiusX2) == 1 {
			break
		}

		TempFloat2.Set(&Zy)
		Zy.Mul(EscapeRadius, &Zx)
		Zy.Mul(&TempFloat2, &Zy)
		Zy.Add(&Zy, &ResultY)
		//PrintBigNumber(&Zy)

		Zx.Sub(&ZxX2, &ZyX2)
		Zx.Add(&Zx, &ResultX)
		//PrintBigNumber(&Zx)

		ZxX2.Mul(&Zx, &Zx)
		ZyX2.Mul(&Zy, &Zy)
	}

	var rawValue float64 = (1.0 / 80.0) * float64(Iterations)
	var finalValue byte = Lerp(0, 255, rawValue)

	var outColor color.RGBA = color.RGBA{finalValue, finalValue, finalValue, 255}
	return outColor
}

func Lerp(start, end int, value float64) byte {
	var StartValue float64 = float64(start + (end - start))
	var ScaledValue float64 = StartValue * value
	return byte(ScaledValue)
}

func HSV2RGB(hue, saturation, value float64) color.RGBA {
	var hh, p, q, t, ff float64
	var i int
	var out color.RGBA
	out.A = 255

	convertFloatToByte := func (value float64) uint8 {
		return uint8(255.0 * value)
	}

	if saturation <= 0.0 {       // < is bogus, just shuts up warnings
		out.R = convertFloatToByte(value)
		out.G = convertFloatToByte(value)
		out.B = convertFloatToByte(value)
		return out
	}
	hh = hue
	if hh >= 360.0 {
		hh = 0.0
	}
	hh /= 60.0
	i = int (hh)
	ff = hh - float64(i)
	p = value * (1.0 - saturation)
	q = value * (1.0 - (saturation * ff))
	t = value * (1.0 - (saturation * (1.0 - ff)))

	switch i {
	case 0:
		out.R = convertFloatToByte(value)
		out.G = convertFloatToByte(t)
		out.B = convertFloatToByte(p)
	case 1:
		out.R = convertFloatToByte(q)
		out.G = convertFloatToByte(value)
		out.B = convertFloatToByte(p)
	case 2:
		out.R = convertFloatToByte(p)
		out.G = convertFloatToByte(value)
		out.B = convertFloatToByte(t)
	case 3:
		out.R = convertFloatToByte(p)
		out.G = convertFloatToByte(q)
		out.B = convertFloatToByte(value)
	case 4:
		out.R = convertFloatToByte(t)
		out.G = convertFloatToByte(p)
		out.B = convertFloatToByte(value)
	case 5:
	default:
		out.R = convertFloatToByte(value)
		out.G = convertFloatToByte(p)
		out.B = convertFloatToByte(q)
		break;
	}
	return out;
}