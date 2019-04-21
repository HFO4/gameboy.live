package gb

import (
	"log"
	"time"
)

type Core struct {
	Cartridge Cartridge
	CPU       CPU
	Memory    Memory

	/*
		Clock and speed options
	*/
	//Frames per-second
	FPS int
	//CPU clock
	Clock int
	//in CBG mode, clock might change to twice as original
	SpeedMultiple int

	/*
		Development options
	*/
	//Debug mode
	Debug bool

	/*
		Timer
	*/
	Timer Timer
}

type Timer struct {
	TimerCounter    int
	DividerRegister int
}

func (core *Core) Init(romPath string) {

	core.SpeedMultiple = 0
	core.Timer.TimerCounter = 0
	core.Timer.DividerRegister = 0

	core.initRom(romPath)
	core.initMemory()
	core.initCPU()
}

func (core *Core) Run() {
	ticker := time.NewTicker(time.Second / time.Duration(core.FPS))
	for range ticker.C {
		core.Update()
	}
}

/*
	Render a frame.
*/
func (core *Core) Update() {
	cyclesThisUpdate := 0
	for cyclesThisUpdate < (core.Clock+core.SpeedMultiple*core.Clock)/core.FPS {
		cycles := 1
		cyclesThisUpdate += cycles
		core.UpdateTimers(cycles)
	}
	//log.Println("Render finish")
}

/*
	Check and update timers.
*/
func (core *Core) UpdateTimers(cycles int) {
	core.DoDividerRegister(cycles)

	if core.IsClockEnabled() {
		core.Timer.TimerCounter += cycles
		if core.Timer.TimerCounter >= core.GetClockFreqCount() {
			// reset m_TimerTracer to the correct value
			core.SetClockFreq()
			// timer about to overflow
			if core.ReadMemory(0xFF05) == 255 {
				core.WriteMemory(0xFF05, core.ReadMemory(0xFF06))
				//RequestInterupt(2)
			} else {
				core.WriteMemory(0xFF05, core.ReadMemory(0xFF05)+1)
			}
		}
	}
}

func (core *Core) DoDividerRegister(cycles int) {
	core.Timer.DividerRegister += cycles
	if core.Timer.DividerRegister >= 255 {
		core.Timer.DividerRegister = 0
		core.Memory.MainMemory[0xFF04]++
	}
}

func (core *Core) SetClockFreq() {
	core.Timer.TimerCounter = 0
}

func (core *Core) IsClockEnabled() bool {
	if core.ReadMemory(0xFF07)&0x04 == 0x04 {
		return true
	}
	return false
}

func (core *Core) GetClockFreq() byte {
	return core.ReadMemory(0xFF07) & 0x3
}

func (core *Core) GetClockFreqCount() int {
	switch core.GetClockFreq() {
	case 0:
		return 1024
	case 1:
		return 16
	case 2:
		return 64
	case 3:
		return 256
	default:
		return 1024
	}
}

/*
	Initialize Cartridge, load rom file and decode rom props
*/
func (core *Core) initRom(romPath string) {
	romData := core.readRomFile(romPath)

	/*
		0134-0143 - Title

		Title of the game in UPPER CASE ASCII. If it is less than 16 characters
		then the remaining bytes are filled with 00's. When inventing the CGB, Nintendo
		has reduced the length of this area to 15 characters, and some months later they
		had the fantastic idea to reduce it to 11 characters only. The new meaning of the
		ex-title bytes is described below.
	*/
	log.Printf("[Cartridge] Game title: %s\n", string(romData[0x134:0x143]))

	/*
		0143 - CGB Flag

		80h - Game supports CGB functions, but works on old gameboys also.
		C0h - Game works on CGB only (physically the same as 80h).
	*/
	isCGB := (romData[0x143] == 0x80 || romData[0x143] == 0xC0)
	log.Printf("[Cartridge] CGB mode: %t\n", isCGB)

	/*
		0147 - Cartridge Type

		Specifies which Memory Bank Controller (if any) is used in the cartridge,
		and if further external hardware exists in the cartridge.
	*/
	CartridgeType := romData[0x147]
	if _, ok := cartridgeTypeMap[CartridgeType]; !ok {
		log.Fatalf("[Cartridge] Unknown cartridge type: %x\n", CartridgeType)
	}
	log.Printf("[Cartridge] Cartridge type: %s\n", cartridgeTypeMap[CartridgeType])

	/*
		Init Cartridge struct according to cartridge type
	*/
	core.Cartridge = Cartridge{}
	switch CartridgeType {
	case 0x00, 0x08, 0x09, 0x0B, 0x0C, 0x0D:
		core.Cartridge.MBC = MBCRom{
			rom: romData,
		}
		core.Cartridge.Props = CartridgeProps{
			MBCType:   "rom",
			ROMLength: len(romData),
		}
	case 0x01, 0x02, 0x03:
		log.Println("mbc1")
	case 0x05, 0x06:
		log.Println("mbc2")
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		log.Println("mbc3")
	case 0x15, 0x16, 0x17:
		log.Println("mbc4")
	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		log.Println("mbc5")
	default:
		log.Fatal("[Cartridge] Unsupported MBC type")
	}

	/*
		Get rom ban number according to ROM Size byte (0148)
		Specifies the ROM Size of the cartridge. Typically calculated as "32KB shl N".
		  00h -  32KByte (no ROM banking) - here we set the bank number to 2
		  01h -  64KByte (4 banks)
		  02h - 128KByte (8 banks)
		  03h - 256KByte (16 banks)
		  04h - 512KByte (32 banks)
		  05h -   1MByte (64 banks)  - only 63 banks used by MBC1
		  06h -   2MByte (128 banks) - only 125 banks used by MBC1
		  07h -   4MByte (256 banks)
		  52h - 1.1MByte (72 banks)
		  53h - 1.2MByte (80 banks)
		  54h - 1.5MByte (96 banks)
	*/
	if _, ok := RomBankMap[romData[0x148]]; !ok {
		log.Fatalf("[Cartridge] Unknown ROM size byte : %x\n", romData[0x148])
	}
	core.Cartridge.Props.ROMBank = RomBankMap[romData[0x148]]
	log.Printf("[Cartridge] ROM bank number: %d (%dKBytes)\n", core.Cartridge.Props.ROMBank, core.Cartridge.Props.ROMBank*16)

	/*
		Get rom ban number according to ROM Size byte (0149)
		Specifies the size of the external RAM in the cartridge (if any).
		  00h - None		(0 banks of 8KBytes each)
		  01h - 2 KBytes	(1 banks of 8KBytes each)
		  02h - 8 KBytes	(1 banks of 8KBytes each)
		  03h - 32 KBytes 	(4 banks of 8KBytes each)
	*/
	if _, ok := RamBankMap[romData[0x149]]; !ok {
		log.Fatalf("[Cartridge] Unknown RAM size byte : %x\n", romData[0x149])
	}
	core.Cartridge.Props.RAMBank = RamBankMap[romData[0x148]]
	log.Printf("[Cartridge] RAM bank number: %d (%dKBytes)\n", core.Cartridge.Props.RAMBank, core.Cartridge.Props.RAMBank*8)
}
