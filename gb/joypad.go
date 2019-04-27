package gb

import "github.com/HFO4/gbc-in-cloud/util"

func (core *Core) GetJoypadStatus() byte {
	res := core.Memory.MainMemory[0xFF00]
	// flip all the bits
	res ^= 0xFF

	// are we interested in the standard buttons?
	if !util.TestBit(res, 4) {
		topJoypad := core.JoypadStatus >> 4
		topJoypad |= 0xF0 // turn the top 4 bits on
		res &= topJoypad  // show what buttons are pressed
	} else if !util.TestBit(res, 5) {
		bottomJoypad := core.JoypadStatus & 0xF
		bottomJoypad |= 0xF0
		res &= bottomJoypad
	}
	return res
}
