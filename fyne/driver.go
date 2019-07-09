package fyne

import (
	"fmt"
	"image"
	"log"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"

	"github.com/HFO4/gbc-in-cloud/util"
)

type LCD struct {
	pixels *[160][144][3]uint8
	screen *image.RGBA

	frame, output fyne.CanvasObject

	inputStatus *byte
	interrupt   bool
	title       string
}

func (lcd *LCD) Init(pixels *[160][144][3]uint8, title string) {
	lcd.pixels = pixels
	lcd.title = title
	log.Println("[Display] Initialize Fyne GUI display")
}

func (lcd *LCD) InitStatus(statusPointer *byte) {
	lcd.inputStatus = statusPointer
}

func (lcd *LCD) UpdateInput() bool {
	if lcd.interrupt {
		lcd.interrupt = false

		return true
	}

	return false
}

func (lcd *LCD) NewInput(b []byte) {
}

func (lcd *LCD) draw(w, h int) image.Image {
	i := 0
	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			r, g, b := lcd.pixels[x][y][0], lcd.pixels[x][y][1], lcd.pixels[x][y][2]

			if r == 0xFF && g == 0xFF && b == 0xFF {
				lcd.screen.Pix[i] = 0x9b
				lcd.screen.Pix[i+1] = 0xbc
				lcd.screen.Pix[i+2] = 0x0f
			} else if r == 0xCC && g == 0xCC && b == 0xCC {
				lcd.screen.Pix[i] = 0x8b
				lcd.screen.Pix[i+1] = 0xac
				lcd.screen.Pix[i+2] = 0x0f
			} else if r == 0x77 && g == 0x77 && b == 0x77 {
				lcd.screen.Pix[i] = 0x30
				lcd.screen.Pix[i+1] = 0x62
				lcd.screen.Pix[i+2] = 0x30
			} else {
				lcd.screen.Pix[i] = 0x0f
				lcd.screen.Pix[i+1] = 0x38
				lcd.screen.Pix[i+2] = 0x0f
			}
			lcd.screen.Pix[i+3] = 0xff

			i += 4
		}
	}

	return lcd.screen
}

// Mapping from keys to GB index.
// Reference :https://github.com/Humpheh/goboy/blob/master/pkg/gbio/iopixel/pixels.go
var keyMap = map[fyne.KeyName]byte{
	// A button
	fyne.KeyZ: 5,
	// B button
	fyne.KeyX: 4,
	// SELECT button
	fyne.KeyBackspace: 6,
	// START button
	fyne.KeyReturn: 7,
	// RIGHT button
	fyne.KeyRight: 0,
	// LEFT button
	fyne.KeyLeft: 1,
	// UP button
	fyne.KeyUp: 2,
	// DOWN button
	fyne.KeyDown: 3,
}

func (lcd *LCD) buttonDown(ev *fyne.KeyEvent) {

	var statusCopy byte
	statusCopy = *lcd.inputStatus
	if offset, ok := keyMap[ev.Name]; ok {
		statusCopy = util.ClearBit(statusCopy, uint(offset))
		lcd.interrupt = true
	}

	*lcd.inputStatus = statusCopy
}

func (lcd *LCD) buttonUp(ev *fyne.KeyEvent) {

	var statusCopy byte
	statusCopy = *lcd.inputStatus
	if offset, ok := keyMap[ev.Name]; ok {
		statusCopy = util.SetBit(statusCopy, uint(offset))
		lcd.interrupt = true
	}

	*lcd.inputStatus = statusCopy
}

func (lcd *LCD) MinSize([]fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(520, 400)
}

func (lcd *LCD) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	lcd.frame.Resize(size)

	xScale := float32(size.Width) / 520.0
	yScale := float32(size.Height) / 400.0

	lcd.output.Resize(fyne.NewSize(int(320*xScale), int(296*yScale)))
	lcd.output.Move(fyne.NewPos(int(100*xScale), int(54*yScale)))
}

func (lcd *LCD) Run(drawSignal chan bool, onQuit func()) {
	a := app.New()
	win := a.NewWindow(fmt.Sprintf("GameBoy - %s", lcd.title))

	lcd.screen = image.NewRGBA(image.Rect(0, 0, 160, 144))
	lcd.output = canvas.NewRaster(lcd.draw)
	go func() {
		for {
			// drawSignal was sent by the emulator
			<-drawSignal

			canvas.Refresh(lcd.output)
		}
	}()

	lcd.frame = canvas.NewImageFromResource(resourceFrameSvg)
	content := fyne.NewContainerWithLayout(lcd, lcd.output, lcd.frame)

	win.SetPadded(false)
	win.SetContent(content)
	win.Canvas().(desktop.Canvas).SetOnKeyDown(lcd.buttonDown)
	win.Canvas().(desktop.Canvas).SetOnKeyUp(lcd.buttonUp)
	win.SetOnClosed(func() {
		onQuit()
		a.Quit()
	})
	win.ShowAndRun()
}
