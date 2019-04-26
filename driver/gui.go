package driver

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"log"
	"os"
)

type LCD struct {
	pixels *[160][144][3]uint8
	window *pixelgl.Window

	pixelMap *pixel.PictureData
}

func (lcd *LCD) Init(pixels *[160][144][3]uint8) {
	lcd.pixels = pixels
	log.Println("[Display] Initialize GUI display")
	lcd.pixelMap = pixel.MakePictureData(pixel.R(0, 0, 160, 144))

}

func (lcd *LCD) Run(drawSignal chan bool) {
	cfg := pixelgl.WindowConfig{
		Title:  "TETRIS [FPS:60] [CLOCK:4194304]",
		Bounds: pixel.R(0, 0, 160*3, 144*3),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	lcd.window = win

	go func() {
		for !win.Closed() {
		}
		os.Exit(1)
	}()

	for {
		<-drawSignal
		for y := 0; y < 144; y++ {
			for x := 0; x < 160; x++ {
				colour := color.RGBA{R: lcd.pixels[x][y][0], G: lcd.pixels[x][y][1], B: lcd.pixels[x][y][2], A: 0xFF}
				lcd.pixelMap.Pix[(143-y)*160+x] = colour
			}
		}

		graph := pixel.NewSprite(pixel.Picture(lcd.pixelMap), pixel.R(0, 0, 160, 144))
		mat := pixel.IM
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(3, 3))
		graph.Draw(lcd.window, mat)
		win.Update()
	}

}
