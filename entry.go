// Giperion October - November 2018
// [EUREKA] 3.8 Beta
package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// Create my own type of window
type MandelbrotWindow struct {
	MainWindow
	outputViewport *walk.CustomWidget
	outBitmap      *walk.Bitmap
	cachedImage2   MandelbrotImage
	viewSettings   MandelbrotView
	cameraViewSettings []MandelbrotView
	LastSize 	   walk.Size
	fractalPrecision int
}

var movementFactor *big.Float = big.NewFloat(0.03)

func (mw *MandelbrotWindow) OnKeyDownHandler(key walk.Key) {
	// Movement. Up, down, left, right
	if key == walk.KeyW || key == walk.KeyUp || key == walk.KeyNumpad8 {
		var movementDelta big.Float
		movementDelta.Mul(movementFactor, &mw.viewSettings.cameraZoom)
		mw.viewSettings.cameraY.Sub(&mw.viewSettings.cameraY, &movementDelta)
	}
	if key == walk.KeyS || key == walk.KeyDown || key == walk.KeyNumpad2 {
		var movementDelta big.Float
		movementDelta.Mul(movementFactor, &mw.viewSettings.cameraZoom)
		mw.viewSettings.cameraY.Add(&mw.viewSettings.cameraY, &movementDelta)
	}
	if key == walk.KeyA || key == walk.KeyLeft || key == walk.KeyNumpad4 {
		var movementDelta big.Float
		movementDelta.Mul(movementFactor, &mw.viewSettings.cameraZoom)
		mw.viewSettings.cameraX.Sub(&mw.viewSettings.cameraX, &movementDelta)
	}
	if key == walk.KeyD || key == walk.KeyRight || key == walk.KeyNumpad6 {
		var movementDelta big.Float
		movementDelta.Mul(movementFactor, &mw.viewSettings.cameraZoom)
		mw.viewSettings.cameraX.Add(&mw.viewSettings.cameraX, &movementDelta)
	}

	// Scaling.
	if key == walk.KeyQ || key == walk.KeyNumpad1 {
		mw.viewSettings.cameraZoom.Mul(&mw.viewSettings.cameraZoom, big.NewFloat(1.1))
	}
	if key == walk.KeyE || key == walk.KeyNumpad3 {
		mw.viewSettings.cameraZoom.Mul(&mw.viewSettings.cameraZoom, big.NewFloat(0.9))
	}

	// Iteration changing
	if key == walk.KeyZ {
		mw.viewSettings.iteration += 5
	}
	if key == walk.KeyC {
		mw.viewSettings.iteration -= 5
	}

	// if you try to invalidate entire form - it will erase current frame, since root ui object has that flag
	// so we invalidate only our widget, which is instructed to not do that
	walk.App().ActiveForm().Children().At(0).Invalidate()
	//walk.App().ActiveForm().Invalidate()
}

const (
	StandartFractal = 0
	PreciseFractal = 1
)

func InlineTesting() {
	var viewParam MandelbrotView = NewMandelbrotView()

	_ = Mandelbrot(50, 45, viewParam)
	_ = Mandelbrot_BigNumber(50, 45, viewParam)
}

func SaveFractalViewState(cameraName string, viewSettings *MandelbrotView, outConfig *string) {
	*outConfig += "\n"
	*outConfig += "[" + cameraName + "]"
	viewSettings.SaveToString(outConfig)
	*outConfig += "\n"
}

func main() {
	hWindow := MandelbrotWindow{
		MainWindow: MainWindow{
			Title:   "Fractal Viewer",
			MinSize: Size{600, 400},
			Layout:  VBox{},
		},
	}

	hWindow.MenuItems = []MenuItem {
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Action{
					Text:        "E&xit",
					OnTriggered: func() { walk.App().Exit(0) },
				},
			},
		}, /* Precision is not working as intended
		Menu{
			Text: "&Precision",
			Items: []MenuItem{
				Action{
					Text:        "64Bit (Fast)",
					OnTriggered: func() { hWindow.fractalPrecision = StandartFractal; walk.App().ActiveForm().Children().At(0).Invalidate() },
				},
				Action{
					Text:        "Unlimited (Slow)",
					OnTriggered: func() { hWindow.fractalPrecision = PreciseFractal; walk.App().ActiveForm().Children().At(0).Invalidate() },
				},
			},
		}, */
		Menu{
			Text: "&Camera",
			Items: []MenuItem{
				Menu {
					Text: "Save",
					Items: []MenuItem {
						Action{
							Text:        "Camera 1",
							OnTriggered: func() { hWindow.cameraViewSettings[0] = hWindow.viewSettings; walk.App().ActiveForm().Children().At(0).Invalidate() },
						},
						Action{
							Text:        "Camera 2",
							OnTriggered: func() { hWindow.cameraViewSettings[1] = hWindow.viewSettings; walk.App().ActiveForm().Children().At(0).Invalidate() },
						},
					},
				},
				Menu {
					Text: "Load",
					Items: []MenuItem {
						Action{
							Text:        "Camera 1",
							OnTriggered: func() { hWindow.viewSettings = hWindow.cameraViewSettings[0]; walk.App().ActiveForm().Children().At(0).Invalidate() },
						},
						Action{
							Text:        "Camera 2",
							OnTriggered: func() { hWindow.viewSettings = hWindow.cameraViewSettings[1]; walk.App().ActiveForm().Children().At(0).Invalidate() },
						},
					},
				},
				Action {
					Text: "Reset",
					// Save width and height of old view settings and apply that to newer one
					OnTriggered: func() {
						var width, height int;
						width = hWindow.viewSettings.width;
						height = hWindow.viewSettings.height;
						hWindow.viewSettings = NewMandelbrotView();
						hWindow.viewSettings.width = width; hWindow.viewSettings.height = height
						walk.App().ActiveForm().Children().At(0).Invalidate()
						},
				},
			},
		},
		Menu{
			Text: "&Help",
			Items: []MenuItem{
				Action{
					Text: "About",
					OnTriggered: func() {
						walk.MsgBox(nil, "About", "Mandelbrot renderer. DirectDraw Engine in Golang. ver. 1.0\r\nControls: \r\nWASD for movement, QE for zoom, ZC for increasing/decreasing details", walk.MsgBoxOK|walk.MsgBoxIconInformation)
					},
				},
			},
		},
	}

	hWindow.viewSettings = NewMandelbrotView()

	hWindow.Children = []Widget{
		CustomWidget{
			AssignTo:            &hWindow.outputViewport,
			ClearsBackground:    false,
			InvalidatesOnResize: true,
			Paint:               hWindow.OnPaint,
			PaintMode:           PaintNoErase,
		},
	}

	//InlineTesting()

	// Load config file
	configData, fileError := ioutil.ReadFile("settings.cfg")
	if fileError == nil {
		var configStr = string(configData[:len(configData)])
		var clearedConfigStr string = strings.Replace(configStr, " ", "", -1)
		var configLines []string = strings.Fields(clearedConfigStr)

		GetIntValueFromConfigLambda := func(strKey string, KeyValuePair []string) int {
			if len(KeyValuePair) > 1 {
				if strings.EqualFold(KeyValuePair[0], strKey) {
					value, parseError := strconv.ParseInt(KeyValuePair[1], 10, 32)
					if parseError == nil {
						return int(value)
					} else {
						walk.MsgBox(nil, "Config error", parseError.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
					}
				}
			}
			return 0
		}

		FlushSectionLambda := func (currentSection string, cameraViewSetting *MandelbrotView) {
			if len(currentSection) > 0 {
				cameraViewSetting.width = hWindow.viewSettings.width
				cameraViewSetting.height = hWindow.viewSettings.height
				hWindow.cameraViewSettings = append(hWindow.cameraViewSettings, *cameraViewSetting)
				*cameraViewSetting = NewMandelbrotView()
			}
		}

		var currentSection string
		var cameraViewSetting MandelbrotView
		for _, line := range configLines {
			if line[0] == '[' {
				// Flush current section, if we are in it
				FlushSectionLambda(currentSection, &cameraViewSetting)
				strLen := len(line)
				currentSection = line[1:strLen-1]
			}
			var KeyValuePair []string = strings.Split(line, "=")

			if len(currentSection) == 0 { // Global space
				var width int = GetIntValueFromConfigLambda("width", KeyValuePair)
				if width != 0 {
					hWindow.viewSettings.width = width
				}

				var height int = GetIntValueFromConfigLambda("height", KeyValuePair)
				if height != 0 {
					hWindow.viewSettings.height = height
				}
			} else { // Section space
				// Currently section name means camera name. So we should load camera settings
				// read the viewSettings
				switch KeyValuePair[0] {
				case "cameraX":
					pCameraX, _, parseError := big.ParseFloat(KeyValuePair[1], 10, 128, big.ToNearestEven)
					if parseError == nil {
						cameraViewSetting.cameraX = *pCameraX
					}
				case "cameraY":
					pCameraY, _, parseError := big.ParseFloat(KeyValuePair[1], 10, 128, big.ToNearestEven)
					if parseError == nil {
						cameraViewSetting.cameraY = *pCameraY
					}
				case "cameraZoom":
					pCameraZoom, _, parseError := big.ParseFloat(KeyValuePair[1], 10, 128, big.ToNearestEven)
					if parseError == nil {
						cameraViewSetting.cameraZoom = *pCameraZoom
					}
				case "iteration":
					cameraViewSetting.iteration = GetIntValueFromConfigLambda("iteration", KeyValuePair)
				}
			}
		}

		FlushSectionLambda(currentSection, &cameraViewSetting)
	}


	// Set the size and create renderbuffer
	hWindow.Size = Size{Width: hWindow.viewSettings.width, Height: hWindow.viewSettings.height}
	hWindow.cachedImage2 = NewMandelbrotImage(hWindow.Size.Width, hWindow.Size.Height)

	// Create a default camera settings, if that was not readed from config file
	if len (hWindow.cameraViewSettings) < 2 {
		hWindow.cameraViewSettings = nil
		hWindow.cameraViewSettings = make([]MandelbrotView, 0)
		for i := 0; i < 2; i++ {
			hWindow.cameraViewSettings = append(hWindow.cameraViewSettings, hWindow.viewSettings)
		}
	}

	// Set keyboard handler to accept user commands and run message loop
	hWindow.OnKeyPress = hWindow.OnKeyDownHandler
	hWindow.Run()

	// save current settings
	var settingContent string
	settingContent += "width = "
	settingContent += strconv.Itoa(hWindow.LastSize.Width)
	settingContent += "\n"
	settingContent += "height = "
	settingContent += strconv.Itoa(hWindow.LastSize.Height)
	settingContent += "\n"

	// Serialize all cameras
	for iter, value := range hWindow.cameraViewSettings {
		settingContent += "[" + "camera" + strconv.Itoa(iter) + "]"
		settingContent += "\n"
		value.SaveToString(&settingContent)
		settingContent += "\n"
	}

	ioutil.WriteFile("settings.cfg", []byte (settingContent), 0)
}

// Main render routine. Called by walk.
func (this *MandelbrotWindow) OnPaint(canvas *walk.Canvas, updateRectangle walk.Rectangle) error {
	//Create a new bitmap if it wasn't created yet
	var theError error

	var currentSize walk.Size = updateRectangle.Size()
	var previousSize walk.Size = walk.Size{Width: int (this.viewSettings.width), Height: int (this.viewSettings.height)}
	if previousSize != currentSize {
		// If resizing was done, we should recreate bitmap and cached image
		this.viewSettings.width = currentSize.Width
		this.viewSettings.height = currentSize.Height
		for _, value := range this.cameraViewSettings {
			value.width = currentSize.Width
			value.height = currentSize.Height
		}
		this.cachedImage2.Resize(currentSize.Width, currentSize.Height)
	}

	this.LastSize = walk.App().ActiveForm().Size()

	this.Render()

	if this.outBitmap != nil {
		this.outBitmap.Dispose()
	}

	// Flush image to output format
	this.outBitmap, theError = walk.NewBitmapFromImage(this.cachedImage2)
	if theError != nil {
		fmt.Println(theError)
		return theError
	}

	nullPosition := walk.Point{X: 0, Y: 0}
	canvas.DrawImage(this.outBitmap, nullPosition)

	return nil
}

func DrawMandelbrotInLines(targetImage *MandelbrotImage, startY int, endY int, viewSetting *MandelbrotView, precise bool, finishingSignal chan int) {
	for y := startY; y < endY; y++ {
		for x := 0; x < viewSetting.width; x++ {
			var pixelColor color.RGBA
			if precise {
				pixelColor = Mandelbrot_BigNumber(x, y, *viewSetting)
			} else {
				pixelColor = Mandelbrot(x, y, *viewSetting)
			}
			targetImage.Set(pixelColor, x, y)
		}
	}

	finishingSignal <- endY // send the last pixel line number
}

func (theWindow *MandelbrotWindow) Render() {
	var StartRenderingTimer time.Time = time.Now()
	var LastRenderY chan int = make(chan int)

	var ImageSize walk.Size = walk.Size{Width: theWindow.viewSettings.width, Height: theWindow.viewSettings.height}

	//var NumberOfThreadStarted int = runtime.NumCPU() * 2
	var NumberOfThreadStarted int = ImageSize.Height
	var YOffset int = ImageSize.Height / NumberOfThreadStarted

	for threadId := 0; threadId < NumberOfThreadStarted; threadId++ {
		var StartY int = threadId * YOffset
		var EndY int = (threadId + 1) * YOffset
		go DrawMandelbrotInLines(&theWindow.cachedImage2, StartY, EndY, &theWindow.viewSettings, theWindow.fractalPrecision == 1, LastRenderY)
	}

	for threadId := 0; threadId < NumberOfThreadStarted; threadId++ {
		_ = <-LastRenderY
	}

	var RenderingDuration time.Duration = time.Since(StartRenderingTimer)
	var RenderTime float64 = float64(RenderingDuration.Nanoseconds()) / 1000000.0
	fmt.Printf("'%f' ms \r\n", RenderTime)
}
