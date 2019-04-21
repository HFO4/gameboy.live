package gb

import (
	"log"
)

/*
	General Memory Map
	  0000-3FFF   16KB ROM Bank 00     (in cartridge, fixed at bank 00)
	  4000-7FFF   16KB ROM Bank 01..NN (in cartridge, switchable bank number)
	  8000-9FFF   8KB Video RAM (VRAM) (switchable bank 0-1 in CGB Mode)
	  A000-BFFF   8KB External RAM     (in cartridge, switchable bank, if any)
	  C000-CFFF   4KB Work RAM Bank 0 (WRAM)
	  D000-DFFF   4KB Work RAM Bank 1 (WRAM)  (switchable bank 1-7 in CGB Mode)
	  E000-FDFF   Same as C000-DDFF (ECHO)    (typically not used)
	  FE00-FE9F   Sprite Attribute Table (OAM)
	  FEA0-FEFF   Not Usable
	  FF00-FF7F   I/O Ports
	  FF80-FFFE   High RAM (HRAM)
	  FFFF        Interrupt Enable Register
*/
type Memory struct {
	MainMemory     [0x10000]byte
	CurrentROMBank byte
	RAMBank        [0x8000]byte
	CurrentRAMBank byte
	EnableRAM      bool
}

func (core *Core) initMemory() {
	log.Println("[Core] Start to initialize memory...")

	log.Println("[Memory] Load first 32KByte of rom data into memory")
	//Load first 32KB of ROM into 0000-7FFF
	for i := 0x0000; i < core.Cartridge.Props.ROMLength && i < 0x8000; i++ {
		core.Memory.MainMemory[i] = core.Cartridge.MBC.ReadRom(uint16(i))
	}

	//Specify which ROM bank is currently loaded into internal memory address 0x4000-0x7FFF.
	// As ROM Bank 0 is fixed into memory address 0x0-0x3FFF this variable should never be 0,
	// it should be at least 1. We need to initialize this variable on emulator load to 1.
	core.Memory.CurrentROMBank = 1

	//Specify which RAM bank is currently loaded into internal memory address 0xA000-0xBFFF.
	core.Memory.CurrentRAMBank = 0

	//Specify other register mapped in main memory according to http://bgb.bircd.org/pandocs.htm#powerupsequence
	core.Memory.MainMemory[0xFF05] = 0x00
	core.Memory.MainMemory[0xFF06] = 0x00
	core.Memory.MainMemory[0xFF07] = 0x00
	core.Memory.MainMemory[0xFF10] = 0x80
	core.Memory.MainMemory[0xFF11] = 0xBF
	core.Memory.MainMemory[0xFF12] = 0xF3
	core.Memory.MainMemory[0xFF14] = 0xBF
	core.Memory.MainMemory[0xFF16] = 0x3F
	core.Memory.MainMemory[0xFF17] = 0x00
	core.Memory.MainMemory[0xFF19] = 0xBF
	core.Memory.MainMemory[0xFF1A] = 0x7F
	core.Memory.MainMemory[0xFF1B] = 0xFF
	core.Memory.MainMemory[0xFF1C] = 0x9F
	core.Memory.MainMemory[0xFF1E] = 0xBF
	core.Memory.MainMemory[0xFF20] = 0xFF
	core.Memory.MainMemory[0xFF21] = 0x00
	core.Memory.MainMemory[0xFF22] = 0x00
	core.Memory.MainMemory[0xFF23] = 0xBF
	core.Memory.MainMemory[0xFF24] = 0x77
	core.Memory.MainMemory[0xFF25] = 0xF3
	core.Memory.MainMemory[0xFF26] = 0xF1
	core.Memory.MainMemory[0xFF40] = 0x91
	core.Memory.MainMemory[0xFF42] = 0x00
	core.Memory.MainMemory[0xFF43] = 0x00
	core.Memory.MainMemory[0xFF45] = 0x00
	core.Memory.MainMemory[0xFF47] = 0xFC
	core.Memory.MainMemory[0xFF48] = 0xFF
	core.Memory.MainMemory[0xFF49] = 0xFF
	core.Memory.MainMemory[0xFF4A] = 0x00
	core.Memory.MainMemory[0xFF4B] = 0x00
	core.Memory.MainMemory[0xFFFF] = 0x00

}

func (core *Core) ReadMemory(address uint16) byte {
	// are we reading from the rom memory bank?
	if (address >= 0x4000) && (address <= 0x7FFF) {
		newAddress := address - 0x4000
		return core.Cartridge.MBC.ReadRom(newAddress + (uint16(core.Memory.CurrentROMBank) * 0x4000))
	} else if (address >= 0xA000) && (address <= 0xBFFF) {
		// are we reading from ram memory bank?
		newAddress := address - 0xA000
		return core.Memory.RAMBank[newAddress+(uint16(core.Memory.CurrentRAMBank)*0x4000)]
	}
	return core.Memory.MainMemory[address]
}

func (core *Core) WriteMemory(address uint16, data byte) {
	if address < 0x8000 {
		log.Println("todo banking")
		//core.HandleBanking(address,data) ;
	} else if (address >= 0xE000) && (address < 0xFE00) {
		// writing to ECHO ram also writes in RAM
		core.Memory.MainMemory[address] = data
		core.WriteMemory(address-0x2000, data)
	} else if (address >= 0xFEA0) && (address < 0xFEFF) {
		// this area is restricted
	} else if 0xFF04 == address {
		core.Memory.MainMemory[0xFF04] = 0
	} else if address == 0xFF07 {
		currentFreq := core.GetClockFreq()
		core.Memory.MainMemory[0xFF07] = data
		newFreq := core.GetClockFreq()
		if currentFreq != newFreq {
			core.SetClockFreq()
		}

	} else {
		core.Memory.MainMemory[address] = data
	}
	log.Printf("Write to %X,data:%X\n", address, data)
}
