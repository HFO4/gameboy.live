package driver

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"time"
)

type ControllerDriver interface {
	InitStatus(*byte)
	UpdateInput() bool
	NewInput([]byte)
}

type TelnetController struct {
	inputStatus *byte

	Keymap [8]KeyMap
}

type KeyMap struct {
	LastPress   int64
	LaskChecked int64
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
		if timeNow-offset.LastPress <= 200 {
			statusCopy = util.ClearBit(statusCopy, uint(key))
			if timeNow-offset.LaskChecked >= 100 {
				requestInterrupt = true
			}
			offset.LaskChecked = timeNow
		} else {
			statusCopy = util.SetBit(statusCopy, uint(key))
		}
	}
	*tel.inputStatus = statusCopy
	return requestInterrupt
}

func (tel *TelnetController) NewInput(data []byte) {

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
