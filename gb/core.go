package gb

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Core struct {
	Rom Cartridge
}

func (core *Core) Init(romPath string) {
	romData := core.readRomFile(romPath)
	fmt.Printf("%s", string(romData[0x134:0x143]))
}

//Read cartridge data from file
func (core *Core) readRomFile(romPath string) []byte {
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

	return bytes
}
