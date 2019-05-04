package gb

import (
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/util"
	"log"
	"time"
)

type Core struct {
	Cartridge Cartridge
	CPU       CPU
	Memory    Memory
	Sound     Sound
	cbMap     [0x100](func())

	/*
	   +++++++++++++++++++++++
	   +     Serial Port     +
	   +++++++++++++++++++++++
	*/
	Serial         driver.ChannelIO
	SerialByte     byte
	InterruptCount int

	/*
	   +++++++++++++++++++++++
	   +        Joypad       +
	   +++++++++++++++++++++++
	*/
	Controller   driver.ControllerDriver
	JoypadStatus byte

	/*
	   +++++++++++++++++++++++
	   +  Screen pixel data  +
	   +++++++++++++++++++++++
	*/

	//Screen pixel data
	Screen     [160][144][3]uint8
	ScanLineBG [160]bool
	//Display driver
	DisplayDriver driver.DisplayDriver
	// Signal to tell display driver to draw
	DrawSignal chan bool

	/*
	  +++++++++++++++++++++++++++
	  + Clock and speed options +
	  +++++++++++++++++++++++++++
	*/
	//Frames per-second
	FPS int
	//CPU clock
	Clock int
	//in CBG mode, clock might change to twice as original
	SpeedMultiple int

	/*
	  ++++++++++++++++++++++++++
	  +  Development options   +
	  ++++++++++++++++++++++++++
	*/
	//Debug mode
	Debug bool
	//Commands num to be executed in DEBUG mode
	DebugControl uint16

	StepExe int

	/*
	  ++++++++++++++++++++++++++
	  +      Other options     +
	  ++++++++++++++++++++++++++
	*/
	ToggleSound bool
	/*
		Timer
	*/
	Timer     Timer
	Exit      bool
	GameTitle string
}

type Timer struct {
	TimerCounter    int
	DividerRegister int
	ScanlineCounter int
}

// Initialize emulator
func (core *Core) Init(romPath string) {
	core.SpeedMultiple = 0
	core.Timer.TimerCounter = 0
	core.Timer.DividerRegister = 0
	core.JoypadStatus = 0xFF
	core.SerialByte = 0xFF
	core.Serial.Receive = make(chan byte)

	core.initRom(romPath)
	core.initMemory()
	core.initCPU()
	core.initCB()
	core.Controller.InitStatus(&core.JoypadStatus)
	core.DisplayDriver.Init(&core.Screen, core.GameTitle)

	/*
		If debug mode is ON, we set the DebugControl to 0x0100,
		where the ROM code was firstly executed.
	*/
	if core.Debug {
		core.DebugControl = 0x0100
	}

	if core.ToggleSound {
		core.Sound.Init()
	}
}

// Start the emulation loop
func (core *Core) Run() {
	// Execution interval depends on the FPS
	ticker := time.NewTicker(time.Second / time.Duration(core.FPS))
	for range ticker.C {
		core.Update()
		// Check controller input interrupt
		if core.Controller.UpdateInput() {
			core.RequestInterrupt(4)
		}
		// Check exit signal
		if core.Exit {
			close(core.DrawSignal)
			return
		}
	}
}

/*
	Render a frame.
*/
func (core *Core) Update() {
	cyclesThisUpdate := 0

	/*
		Gameboy's CPU speed is 4.194304MHz, so in every update loop,
		we need to execute `Clock / FPS` cycles. Some Gameboy Color games might
		use double speed mode, under these, `SpeedMultiple` will be set to `1`.
	*/
	for cyclesThisUpdate < ((core.SpeedMultiple+1)*core.Clock)/core.FPS {
		cycles := 4

		/*
			Check whether CPU is halted, when this happen, only an interrupt
			can stop halting.
		*/
		if !core.CPU.Halt {
			cycles = core.ExecuteNextOPCode()
		}
		cyclesThisUpdate += cycles
		core.UpdateTimers(cycles)
		core.UpdateGraphics(cycles)
		cyclesThisUpdate += core.Interrupt()
		core.UpdateIO(cycles)

	}
	core.RenderScreen()
}

func (core *Core) UpdateIO(cycles int) {
	data, reqInt := core.Serial.FetchByte(cycles)
	if reqInt {
		ret := core.Memory.MainMemory[0xFF02]
		ret = util.ClearBit(ret, 7)
		core.Memory.MainMemory[0xFF02] = ret
		//core.Serial.SetChannelStatus(util.TestBit(ret,0),util.TestBit(ret,7))
		core.SerialByte = data
		core.RequestInterrupt(3)
	}
}

/*
	Check interrupt.
*/
func (core *Core) Interrupt() int {

	/*
		If `EI`(Enable Interrupt) instruction was executed, Interrupt Mater Flag will
		be enable in next execution cycle.
	*/
	if core.CPU.Flags.PendingInterruptEnabled {
		core.CPU.Flags.PendingInterruptEnabled = false
		core.CPU.Flags.InterruptMaster = true
		return 0
	}

	/*
		If the CPU is neither interrupted nor halted,
		stop interrupt checking and return.
	*/
	if !core.CPU.Flags.InterruptMaster && !core.CPU.Halt {
		return 0
	}

	//Check the Interrupt Master Enable Flag
	if core.CPU.Flags.InterruptMaster || core.CPU.Halt {
		/*
			FF0F - IF - Interrupt Flag (R/W)
			  Bit 0: V-Blank  Interrupt Request (INT 40h)  (1=Request)
			  Bit 1: LCD STAT Interrupt Request (INT 48h)  (1=Request)
			  Bit 2: Timer    Interrupt Request (INT 50h)  (1=Request)
			  Bit 3: Serial   Interrupt Request (INT 58h)  (1=Request)
			  Bit 4: Joypad   Interrupt Request (INT 60h)  (1=Request)

		*/
		req := core.ReadMemory(0xFF0F)
		/*
			FFFF - IE - Interrupt Enable (R/W)
			  Bit 0: V-Blank  Interrupt Enable  (INT 40h)  (1=Enable)
			  Bit 1: LCD STAT Interrupt Enable  (INT 48h)  (1=Enable)
			  Bit 2: Timer    Interrupt Enable  (INT 50h)  (1=Enable)
			  Bit 3: Serial   Interrupt Enable  (INT 58h)  (1=Enable)
			  Bit 4: Joypad   Interrupt Enable  (INT 60h)  (1=Enable)
		*/
		enabled := core.ReadMemory(0xFFFF)
		if req > 0 {
			/*

			 */
			for i := 0; i < 5; i++ {
				if util.TestBit(req, uint(i)) {
					// Check whether this interrupt request is enabled in IE.
					if util.TestBit(enabled, uint(i)) {
						core.DoInterrupt(i)
						return 20
					}
				}
			}
		}
	}
	return 0
}

/*
	Performing an interrupt
*/
func (core *Core) DoInterrupt(id int) {

	if !core.CPU.Flags.InterruptMaster && core.CPU.Halt {
		core.CPU.Halt = false
		return
	}

	// Turn off the Interrupt Master Enable Flag
	core.CPU.Flags.InterruptMaster = false
	core.CPU.Halt = false

	req := core.ReadMemory(0xFF0F)
	req = util.ClearBit(req, uint(id))
	core.WriteMemory(0xFF0F, req)
	// We must save the current execution address by pushing it onto the stack
	core.StackPush(core.CPU.Registers.PC)

	/*
		Set the PC to correspond interrupt process program:
			V-Blank: 0x40
			LCD: 0x48
			TIMER: 0x50
			JOYPAD: 0x60
			Serial: 0x58
	*/
	switch id {
	case 0:
		core.CPU.Registers.PC = 0x40
	case 1:
		core.CPU.Registers.PC = 0x48
	case 2:
		core.CPU.Registers.PC = 0x50
	case 3:
		core.CPU.Registers.PC = 0x58
	case 4:
		core.CPU.Registers.PC = 0x60
	default:
		log.Fatalf("Unknown Interrupt: %d", id)
	}
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
				core.RequestInterrupt(2)
			} else {
				core.WriteMemory(0xFF05, core.ReadMemory(0xFF05)+1)
			}
		}
	}
}

/*
	Request an Interrupt.
*/
func (core *Core) RequestInterrupt(id int) {
	//Read the present Interrupt Flag
	req := core.ReadMemory(0xFF0F)
	req = util.SetBit(req, uint(id))
	core.WriteMemory(0xFF0F, req)
}

/*
	Update divider register.
	This register is incremented at rate of 16384Hz (~16779Hz on SGB).
	In CGB Double Speed Mode it is incremented twice as fast, ie. at 32768Hz.
*/
func (core *Core) DoDividerRegister(cycles int) {
	core.Timer.DividerRegister += cycles
	if core.Timer.DividerRegister >= 255 {
		core.Timer.DividerRegister = 0
		core.Memory.MainMemory[0xFF04]++
	}
}

/*
	Reset clock frequency.
*/
func (core *Core) SetClockFreq() {
	core.Timer.TimerCounter = 0
}

/*
	Check whether clock is enabled.
*/
func (core *Core) IsClockEnabled() bool {
	if core.ReadMemory(0xFF07)&0x04 == 0x04 {
		return true
	}
	return false
}

/*
	Get clock frequency sign specified in TAC register.
*/
func (core *Core) GetClockFreq() byte {
	return core.ReadMemory(0xFF07) & 0x3
}

/*
	Get clock frequency sign according to clock frequency sign in TAC register.
	FF07 - TAC - Timer Control (R/W)
	  Bit 2    - Timer Stop  (0=Stop, 1=Start)
	  Bits 1-0 - Input Clock Select
             00:   4096 Hz    (~4194 Hz SGB)
             01: 262144 Hz  (~268400 Hz SGB)
             10:  65536 Hz   (~67110 Hz SGB)
             11:  16384 Hz   (~16780 Hz SGB)

*/
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
	core.GameTitle = string(romData[0x134:0x143])
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
	switch CartridgeType {
	case 0x00, 0x08, 0x09, 0x0B, 0x0C, 0x0D:
		core.Cartridge.MBC = &MBCRom{
			rom: romData,
			//Specify which ROM bank is currently loaded into internal memory address 0x4000-0x7FFF.
			//As ROM Bank 0 is fixed into memory address 0x0-0x3FFF this variable should never be 0,
			//it should be at least 1. We need to initialize this variable on emulator load to 1.
			CurrentROMBank: 1,
			//Specify which RAM bank is currently loaded into internal memory address 0xA000-0xBFFF.
			CurrentRAMBank: 1,
		}
		core.Cartridge.Props = &CartridgeProps{
			MBCType:   "rom",
			ROMLength: len(romData),
		}
	case 0x01, 0x02, 0x03:
		MBC := &MBC1{
			rom:            romData,
			CurrentROMBank: 1,
			CurrentRAMBank: 0,
			RAMBank:        make([]byte, 0x8000),
		}
		core.Cartridge.MBC = MBC
		core.Cartridge.Props = &CartridgeProps{
			MBCType:   "MBC1",
			ROMLength: len(romData),
		}
	case 0x05, 0x06:
		MBC := &MBC2{
			rom:            romData,
			CurrentROMBank: 1,
			CurrentRAMBank: 0,
			RAMBank:        make([]byte, 0x8000),
		}
		core.Cartridge.MBC = MBC
		core.Cartridge.Props = &CartridgeProps{
			MBCType:   "MBC2",
			ROMLength: len(romData),
		}
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		MBC := &MBC3{
			rom:            romData,
			CurrentROMBank: 1,
			CurrentRAMBank: 0,
		}
		core.Cartridge.MBC = MBC
		core.Cartridge.Props = &CartridgeProps{
			MBCType:   "MBC3",
			ROMLength: len(romData),
		}
	//case 0x15, 0x16, 0x17:
	//	log.Println("mbc4")
	//case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
	//	log.Println("mbc5")
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
