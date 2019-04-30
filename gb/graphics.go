package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"log"
)

/*
	Draw the current scan line
*/
func (core *Core) DrawScanLine() {

	//	FF40 - LCDC - LCD Control (R/W)
	//  	Bit 7 - LCD Display Enable             (0=Off, 1=On)
	//  	Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	//  	Bit 5 - Window Display Enable          (0=Off, 1=On)
	//  	Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
	//  	Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
	//  	Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
	//  	Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
	//  	Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
	control := core.ReadMemory(0xFF40)

	//	LCDC.0 - 1) Monochrome Gameboy and SGB: BG Display
	//	When Bit 0 is cleared, the background becomes blank (white).
	//	Window and Sprites may still be displayed (if enabled in Bit 1 and/or Bit 5).
	if util.TestBit(control, 0) {
		core.RenderTiles()
	}

	if util.TestBit(control, 1) {
		core.RenderSprites()
	}
}

/*
	Render Sprites
*/
func (core *Core) RenderSprites() {
	use8x16 := false
	lcdControl := core.ReadMemory(0xFF40)

	if util.TestBit(lcdControl, 2) {
		use8x16 = true
	}
	for sprite := 0; sprite < 40; sprite++ {
		// sprite occupies 4 bytes in the sprite attributes table
		index := sprite * 4
		yPos := core.ReadMemory(0xFE00+uint16(index)) - 16
		xPos := core.ReadMemory(0xFE00+uint16(index)+1) - 8
		tileLocation := core.ReadMemory(uint16(0xFE00 + index + 2))
		attributes := core.ReadMemory(0xFE00 + uint16(index) + 3)

		yFlip := util.TestBit(attributes, 6)
		xFlip := util.TestBit(attributes, 5)
		priority := !util.TestBit(attributes, 7)
		scanline := core.ReadMemory(0xFF44)

		ysize := 8
		if use8x16 {
			ysize = 16
		}

		// does this sprite intercept with the scanline?
		if (scanline >= yPos) && (scanline < (yPos + byte(ysize))) {
			line := int(scanline - yPos)
			// read the sprite in backwards in the y axis
			if yFlip {
				line -= ysize
				line *= -1
			}
			line *= 2 // same as for tiles
			dataAddress := (uint16(int(tileLocation)*16 + line))
			data1 := core.ReadMemory(0x8000 + dataAddress)
			data2 := core.ReadMemory(0x8000 + dataAddress + 1)

			// its easier to read in from right to left as pixel 0 is
			// bit 7 in the colour data, pixel 1 is bit 6 etc...
			for tilePixel := 7; tilePixel >= 0; tilePixel-- {
				colourbit := tilePixel

				// read the sprite in backwards for the x axis
				if xFlip {
					colourbit -= 7
					colourbit *= -1
				}

				colourNum := util.GetVal(data2, uint(colourbit))
				colourNum <<= 1
				colourNum |= util.GetVal(data1, uint(colourbit))
				colourAddress := uint16(0xFF48)
				if util.TestBit(attributes, 4) {
					colourAddress = 0xFF49
				}
				// now we have the colour id get the actual
				// colour from palette 0xFF47
				colour := core.GetColour(colourNum, colourAddress)

				// white is transparent for sprites.
				if colourNum == 0 {
					continue
				}

				red := uint8(0)
				green := uint8(0)
				blue := uint8(0)

				switch colour {
				case 0:
					red = 255
					green = 255
					blue = 255
				case 1:
					red = 0xCC
					green = 0xCC
					blue = 0xCC
				case 2:
					red = 0x77
					green = 0x77
					blue = 0x77
				default:
					red = 0
					green = 0
					blue = 0
				}

				xPix := 0 - tilePixel
				xPix += 7

				pixel := int(xPos) + xPix

				// sanity check
				if (scanline < 0) || (scanline > 143) || (pixel < 0) || (pixel > 159) {
					continue
				}

				if core.ScanLineBG[pixel] || priority {
					core.Screen[pixel][scanline-1][0] = red
					core.Screen[pixel][scanline-1][1] = green
					core.Screen[pixel][scanline-1][2] = blue
				}

			}
		}

	}
}

var last uint16

/*
	Render Tiles for the current scan line
*/
func (core *Core) RenderTiles() {
	var tileData uint16 = 0
	var backgroundMemory uint16 = 0
	var unsig bool = true
	lcdControl := core.ReadMemory(0xFF40)

	//	FF42 - SCY - Scroll Y (R/W)
	//	FF43 - SCX - Scroll X (R/W)
	//		Specifies the position in the 256x256 pixels BG map (32x32 tiles)
	//		which is to be displayed at the upper/left LCD display position.
	//		Values in range from 0-255 may be used for X/Y each, the video
	//		controller automatically wraps back to the upper (left) position
	//		in BG map when drawing exceeds the lower (right) border of the BG
	//		map area.
	scrollY := core.ReadMemory(0xFF42)
	scrollX := core.ReadMemory(0xFF43)

	//	FF4A - WY - Window Y Position (R/W)
	//	FF4B - WX - Window X Position minus 7 (R/W)
	//		Specifies the upper/left positions of the Window area.
	//		(The window is an alternate background area which can be
	//		displayed above of the normal background. OBJs (sprites)
	//		may be still displayed above or behinf the window, just as
	//		for normal BG.)
	//
	//		The window becomes visible (if enabled) when positions are set
	//		in range WX=0..166, WY=0..143. A postion of WX=7, WY=0 locates
	//		the window at upper left, it is then completly covering normal
	//		background.
	windowY := core.ReadMemory(0xFF4A)
	windowX := core.ReadMemory(0xFF4B) - 7

	usingWindow := false

	// is the window enabled?
	if util.TestBit(lcdControl, 5) {

		// is the current scan line we're drawing
		// within the windows Y pos?,
		if windowY <= core.ReadMemory(0xFF44) {
			usingWindow = true
		}
	}

	// which tile data are we using?
	if util.TestBit(lcdControl, 4) {
		tileData = 0x8000
	} else {
		// IMPORTANT: This memory region uses signed
		// bytes as tile identifiers
		tileData = 0x8800
		unsig = false
	}

	// which background mem?
	if !usingWindow {
		if util.TestBit(lcdControl, 3) {
			backgroundMemory = 0x9C00
		} else {
			backgroundMemory = 0x9800
		}
	} else {
		if util.TestBit(lcdControl, 6) {
			backgroundMemory = 0x9C00
		} else {
			backgroundMemory = 0x9800
		}
	}

	// yPos is used to calculate which of 32 vertical tiles the
	// current scanline is drawing
	var yPos byte = 0
	if !usingWindow {
		yPos = scrollY + core.ReadMemory(0xFF44)
	} else {
		yPos = core.ReadMemory(0xFF44) - windowY
	}

	// which of the 8 vertical pixels of the current
	// tile is the scanline on?
	var tileRow = ((uint16(yPos / 8)) * 32)

	// time to start drawing the 160 horizontal pixels
	// for this scanline
	for pixel := byte(0); pixel < 160; pixel++ {

		xPos := byte(pixel) + scrollX

		// translate the current x pos to window space if necessary
		if usingWindow {
			if pixel >= windowX {
				xPos = pixel - windowX
			}
		}

		// which of the 32 horizontal tiles does this xPos fall within?
		tileCol := uint16(xPos / 8)
		var tileNum int16

		// get the tile identity number. Remember it can be signed
		// or unsigned
		tileAddress := backgroundMemory + tileRow + tileCol
		if unsig {
			tileNum = int16(core.ReadMemory(tileAddress))
		} else {
			tileNum = int16(int8(core.ReadMemory(tileAddress)))
		}

		// deduce where this tile identifier is in memory.
		tileLocation := tileData
		if unsig {
			tileLocation += uint16(tileNum * 16)
		} else {
			tileLocation = uint16(int32(tileLocation) + int32((tileNum+128)*16))
		}

		// find the correct vertical line we're on of the
		// tile to get the tile data
		//	from in memory
		line := yPos % 8
		// each vertical line takes up two bytes of memory
		line *= 2
		data1 := core.ReadMemory(tileLocation + uint16(line))
		data2 := core.ReadMemory(tileLocation + uint16(line) + 1)
		if last == 0x86F0 && (tileLocation+uint16(line) == 0x8000) {
			log.Printf("%X\n", tileNum)

		}
		last = tileLocation + uint16(line)

		// pixel 0 in the tile is it 7 of data 1 and data2.
		// Pixel 1 is bit 6 etc..
		var colourBit int = int(xPos % 8)
		colourBit -= 7
		colourBit *= -1

		// combine data 2 and data 1 to get the colour id for this pixel
		// in the tile
		colourNum := util.GetVal(data2, uint(colourBit))
		colourNum <<= 1
		colourNum |= util.GetVal(data1, uint(colourBit))

		// now we have the colour id get the actual
		// colour from palette 0xFF47
		colour := core.GetColour(colourNum, 0xFF47)

		red := uint8(0)
		green := uint8(0)
		blue := uint8(0)

		switch colour {
		case 0:
			red = 255
			green = 255
			blue = 255
		case 1:
			red = 0xCC
			green = 0xCC
			blue = 0xCC
		case 2:
			red = 0x77
			green = 0x77
			blue = 0x77
		default:
			red = 0
			green = 0
			blue = 0
		}
		finally := int(core.ReadMemory(0xFF44))
		// safety check to make sure what im about
		// to set is int the 160x144 bounds
		if (finally < 0) || (finally > 143) || (pixel < 0) || (pixel > 159) {
			continue
		}

		// Store whether the background is white
		if colour == 0 {
			core.ScanLineBG[pixel] = true
		} else {
			core.ScanLineBG[pixel] = false
		}

		core.Screen[pixel][finally-1][0] = red
		core.Screen[pixel][finally-1][1] = green
		core.Screen[pixel][finally-1][2] = blue
	}

}

/*
	Get colour id via colour palette and colour Num
		0 - WHITE
		1 - LIGHT_GRAY
		2 - DARK_GRAY
		3 - BLACK
	TODO: GCB mode
*/
func (core *Core) GetColour(colourNum byte, address uint16) int {
	res := 0
	palette := core.ReadMemory(address)

	// which bits of the colour palette does the colour id map to?
	var hi uint
	var lo uint
	switch colourNum {
	case 0:
		hi = 1
		lo = 0
	case 1:
		hi = 3
		lo = 2
	case 2:
		hi = 5
		lo = 4
	case 3:
		hi = 7
		lo = 6
	default:
		hi = 1
		lo = 0
	}

	// use the palette to get the colour
	colour := byte(0)
	colour = util.GetVal(palette, hi) << 1
	colour |= util.GetVal(palette, lo)

	switch colour {
	case 0:
		res = 0
	case 1:
		res = 1
	case 2:
		res = 2
	case 3:
		res = 3
	default:
		res = 0
	}

	return res
}

func (core *Core) RenderScreen() {
	core.DrawSignal <- true
}
