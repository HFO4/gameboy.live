package driver

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"time"
)

type ControllerDriver interface {
	// Initialize controller status
	InitStatus(*byte)
	// Update input status for emulator, called by the emulator
	UpdateInput() bool
	// Update input status for the controller itself
	NewInput([]byte)
}

type TelnetController struct {
	inputStatus *byte
	Keymap      [8]KeyMap
}

type KeyMap struct {
	// Timestamp when the key was last pressed
	LastPress int64
	// Timestamp when the emulator last check the input status
	LastChecked int64
}

func (tel *TelnetController) InitStatus(statusPointer *byte) {
	tel.inputStatus = statusPointer
}

func (tel *TelnetController) UpdateInput() bool {
	var requestInterrupt bool
	var statusCopy byte
	statusCopy = *tel.inputStatus
	timeNow := time.Now().UnixNano() / int64(time.Millisecond)
	for key, offset := range tel.Keymap {
		/*
			If the key was pressed in 200ms,
			We consider player "hold" the key.
		*/
		if timeNow-offset.LastPress <= 200 {
			statusCopy = util.ClearBit(statusCopy, uint(key))
			if timeNow-offset.LastChecked >= 100 {
				requestInterrupt = true
			}
			offset.LastChecked = timeNow
		} else {
			statusCopy = util.SetBit(statusCopy, uint(key))
		}
	}
	*tel.inputStatus = statusCopy
	return requestInterrupt
}

func (tel *TelnetController) NewInput(data []byte) {

	/*
		Key map for the input key byte and the control
		bit in Gameboy register.
	*/
	keyDataMap := map[byte]int{
		65:  2,
		66:  3,
		68:  1,
		67:  0,
		10:  7,
		8:   6,
		122: 4,
		120: 5,
		0:   7,
	}

	key := data[len(data)-1]
	keyId := keyDataMap[key]

	timeNow := time.Now().UnixNano() / int64(time.Millisecond)
	tel.Keymap[keyId].LastPress = timeNow

}
