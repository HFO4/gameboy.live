package gb

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/*
	00h  ROM ONLY                 13h  MBC3+RAM+BATTERY
	01h  MBC1                     15h  MBC4
	02h  MBC1+RAM                 16h  MBC4+RAM
	03h  MBC1+RAM+BATTERY         17h  MBC4+RAM+BATTERY
	05h  MBC2                     19h  MBC5
	06h  MBC2+BATTERY             1Ah  MBC5+RAM
	08h  ROM+RAM                  1Bh  MBC5+RAM+BATTERY
	09h  ROM+RAM+BATTERY          1Ch  MBC5+RUMBLE
	0Bh  MMM01                    1Dh  MBC5+RUMBLE+RAM
	0Ch  MMM01+RAM                1Eh  MBC5+RUMBLE+RAM+BATTERY
	0Dh  MMM01+RAM+BATTERY        FCh  POCKET CAMERA
	0Fh  MBC3+TIMER+BATTERY       FDh  BANDAI TAMA5
	10h  MBC3+TIMER+RAM+BATTERY   FEh  HuC3
	11h  MBC3                     FFh  HuC1+RAM+BATTERY
	12h  MBC3+RAM
*/
var cartridgeTypeMap = map[byte]string{
	byte(0x00): "ROM ONLY",
	byte(0x01): "MBC1",
	byte(0x02): "MBC1+RAM",
	byte(0x03): "MBC1+RAM+BATTERY",
	byte(0x05): "MBC2",
	byte(0x06): "MBC2+BATTERY",
	byte(0x08): "ROM+RAM",
	byte(0x09): "ROM+RAM+BATTERY",
	byte(0x0B): "MMM01",
	byte(0x0C): "MMM01+RAM",
	byte(0x0D): "MMM01+RAM+BATTERY",
	byte(0x0F): "MBC3+TIMER+BATTERY",
	byte(0x10): "MBC3+TIMER+RAM+BATTERY",
	byte(0x11): "MBC3",
	byte(0x12): "MBC3+RAM",
	byte(0x13): "MBC3+RAM+BATTERY",
	byte(0x15): "MBC4",
	byte(0x16): "MBC4+RAM",
	byte(0x17): "MBC4+RAM+BATTERY",
	byte(0x19): "MBC5",
	byte(0x1A): "MBC5+RAM",
	byte(0x1B): "MBC5+RAM+BATTERY",
	byte(0x1C): "MBC5+RUMBLE",
	byte(0x1D): "MBC5+RUMBLE+RAM",
	byte(0x1E): "MBC5+RUMBLE+RAM+BATTERY",
	byte(0xFC): "POCKET CAMERA",
	byte(0xFD): "BANDAI TAMA5",
	byte(0xFE): "HuC3",
	byte(0x1F): "HuC1+RAM+BATTERY",
}

/*
	ROM bank number is linked to the ROM Size byte (0148).
		1 bank = 16 KBytes
	0x00 means no bank required.
*/
var RomBankMap = map[byte]uint8{
	byte(0x00): 2,
	byte(0x01): 4,
	byte(0x02): 8,
	byte(0x03): 16,
	byte(0x04): 32,
	byte(0x05): 64,
	byte(0x06): 128,
	byte(0x52): 72,
	byte(0x53): 80,
	byte(0x54): 96,
}

/*
	RAM bank number is linked to the RAM Size byte (0149).
		1 bank = 8 KBytes
	0x00 means no bank required.
*/
var RamBankMap = map[byte]uint8{
	byte(0x00): 0,
	byte(0x01): 1,
	byte(0x02): 1,
	byte(0x03): 4,
}

type Cartridge struct {
	Props CartridgeProps
	MBC   MBC
}

/*
	Cartridge props
*/
type CartridgeProps struct {
	MBCType string
	ROMBank uint8
	RAMBank uint8
}

type MBC interface {
}

/*
	Single ROM without MBC
*/
type MBCRom struct {
	rom []byte
}

func InitCartridge() {
	fmt.Println("ss")
}

/*
	Read cartridge data from file
*/
func (core *Core) readRomFile(romPath string) []byte {
	log.Println("[Core] Loading rom file...")
	romFile, err := os.Open(romPath)
	defer romFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	stats, statsErr := romFile.Stat()
	if statsErr != nil {
		log.Fatal(statsErr)
	}
	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufReader := bufio.NewReader(romFile)
	_, err = bufReader.Read(bytes)

	log.Printf("[Core] %d Bytes rom loaded\n", size)
	return bytes
}
