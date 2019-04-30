package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
)

/**
FF41 - STAT - LCDC Status (R/W)
  Bit 6 - LYC=LY Coincidence Interrupt (1=Enable) (Read/Write)
  Bit 5 - Mode 2 OAM Interrupt         (1=Enable) (Read/Write)
  Bit 4 - Mode 1 V-Blank Interrupt     (1=Enable) (Read/Write)
  Bit 3 - Mode 0 H-Blank Interrupt     (1=Enable) (Read/Write)
  Bit 2 - Coincidence Flag  (0:LYC<>LY, 1:LYC=LY) (Read Only)
  Bit 1-0 - Mode Flag       (Mode 0-3, see below) (Read Only)
			0: During H-Blank
			1: During V-Blank
			2: During Searching OAM-RAM
			3: During Transfering Data to LCD Driver

The two lower STAT bits show the current status of the LCD controller.
  Mode 0: The LCD controller is in the H-Blank period and
		  the CPU can access both the display RAM (8000h-9FFFh)
		  and OAM (FE00h-FE9Fh)

  Mode 1: The LCD contoller is in the V-Blank period (or the
		  display is disabled) and the CPU can access both the
		  display RAM (8000h-9FFFh) and OAM (FE00h-FE9Fh)

  Mode 2: The LCD controller is reading from OAM memory.
		  The CPU <cannot> access OAM memory (FE00h-FE9Fh)
		  during this period.

  Mode 3: The LCD controller is reading from both OAM and VRAM,
		  The CPU <cannot> access OAM and VRAM during this period.
		  CGB Mode: Cannot access Palette Data (FF69,FF6B) either.

The following are typical when the display is enabled:
  Mode 2  2_____2_____2_____2_____2_____2___________________2____
  Mode 3  _33____33____33____33____33____33__________________3___
  Mode 0  ___000___000___000___000___000___000________________000
  Mode 1  ____________________________________11111111111111_____

The Mode Flag goes through the values 0, 2, and 3 at a cycle of about 109uS. 0 is present about 48.6uS, 2 about 19uS, and 3 about 41uS. This is interrupted every 16.6ms by the VBlank (1). The mode flag stays set at 1 for about 1.08 ms.

Mode 0 is present between 201-207 clks, 2 about 77-83 clks, and 3 about 169-175 clks. A complete cycle through these states takes 456 clks. VBlank lasts 4560 clks. A complete screen refresh occurs every 70224 clks.)
*/
func (core *Core) SetLCDStatus() {
	status := core.ReadMemory(0xFF41)
	if !core.IsLCDEnabled() {
		// set the mode to 1 during lcd disabled and reset scanline
		core.Timer.ScanlineCounter = 456
		core.Memory.MainMemory[0xFF44] = 0
		status &= 252
		status = util.ClearBit(status, 0)
		status = util.ClearBit(status, 1)
		core.WriteMemory(0xFF41, status)
		return
	}

	currentLine := core.ReadMemory(0xFF44)
	currentMode := status & 0x3
	mode := byte(0)
	reqInt := false

	// in vblank so set mode to 1
	if currentLine >= 144 {
		mode = 1
		status = util.SetBit(status, 0)
		status = util.ClearBit(status, 1)
		reqInt = util.TestBit(status, 4)
	} else {
		mode2bounds := 456 - 80
		mode3bounds := mode2bounds - 172

		// mode 2
		if core.Timer.ScanlineCounter >= mode2bounds {
			mode = 2
			status = util.SetBit(status, 1)
			status = util.ClearBit(status, 0)
			reqInt = util.TestBit(status, 5)
		} else if core.Timer.ScanlineCounter >= mode3bounds {
			mode = 3
			// mode 3
			status = util.SetBit(status, 1)
			status = util.SetBit(status, 0)
		} else {
			// mode 0
			mode = 0
			status = util.ClearBit(status, 1)
			status = util.ClearBit(status, 0)
			reqInt = util.TestBit(status, 3)
		}
	}

	// just entered a new mode so request interupt
	if reqInt && (mode != currentMode) {
		core.RequestInterrupt(1)
	}

	// check the conincidence flag
	if currentLine == core.ReadMemory(0xFF45) {
		status = util.SetBit(status, 2)
		if util.TestBit(status, 6) {
			core.RequestInterrupt(1)
		}
	} else {
		status = util.ClearBit(status, 2)
	}
	core.WriteMemory(0xFF41, status)

}

/*
	Check if LCD is enabled
*/
func (core *Core) IsLCDEnabled() bool {
	return util.TestBit(core.ReadMemory(0xFF40), 7)
}

/*
	Scan line and draw line
*/
func (core *Core) UpdateGraphics(cycles int) {
	//Set the LCD status register
	core.SetLCDStatus()

	//A complete cycle through Scan Line states takes 456 clks.
	//We use a counter to mark this.
	if core.IsLCDEnabled() {
		core.Timer.ScanlineCounter -= cycles
	} else {
		return
	}

	if core.Timer.ScanlineCounter <= 0 {
		// time to move onto next scanline
		core.Memory.MainMemory[0xFF44]++

		currentLine := core.ReadMemory(0xFF44)

		//Reset the counter
		core.Timer.ScanlineCounter += 456

		// we have entered vertical blank period
		if currentLine == 144 {
			core.DrawScanLine()
			core.RequestInterrupt(0)
		} else if currentLine > 153 {
			// if gone past scanline 153 reset to 0
			core.Memory.MainMemory[0xFF44] = 0
		} else if currentLine < 144 {
			//log.Println(currentLine)
			core.DrawScanLine()
		}

	}

}
