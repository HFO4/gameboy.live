package gb

import (
	"fmt"
	"github.com/HFO4/gbc-in-cloud/util"
	"log"
)

type CPU struct {
	Registers Registers
	Flags     Flags
	Halt      bool
}

/*
	Registers
	  16bit Hi   Lo   Name/Function
	  AF    A    -    Accumulator & Flags
	  BC    B    C    BC
	  DE    D    E    DE
	  HL    H    L    HL
	  SP    -    -    Stack Pointer
	  PC    -    -    Program Counter/Pointer
*/
type Registers struct {
	A  byte
	B  byte
	C  byte
	D  byte
	E  byte
	F  byte
	HL uint16
	PC uint16
	SP uint16
}

/*
	The Flag Register (lower 8bit of AF register)
	  Bit  Name  Set Clr  Expl.
	  7    zf    Z   NZ   Zero Flag
	  6    n     -   -    Add/Sub-Flag (BCD)
	  5    h     -   -    Half Carry Flag (BCD)
	  4    cy    C   NC   Carry Flag
	  3-0  -     -   -    Not used (always zero)
	Contains the result from the recent instruction which has affected flags.
*/
type Flags struct {
	Zero      bool
	Sub       bool
	HalfCarry bool
	Carry     bool
	//	IME - Interrupt Master Enable Flag (Write Only)
	//  	0 - Disable all Interrupts
	//  	1 - Enable all Interrupts that are enabled in IE Register (FFFF)
	InterruptMaster bool

	PendingInterruptEnabled bool
}

func (core *Core) initCPU() {
	log.Println("[Core] Initialize CPU flags and registers")

	//Initialize flags with default value.
	core.CPU.Flags.Zero = true
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = true
	core.CPU.Flags.InterruptMaster = false

	/*
		Initialize register after BIOS
		AF=$01B0
		BC=$0013
		DE=$00D8
		HL=$014D
		Stack Pointer=$FFFE
	*/
	core.CPU.Registers.A = 0x01
	core.CPU.Registers.B = 0x00
	core.CPU.Registers.C = 0x13
	core.CPU.Registers.D = 0x00
	core.CPU.Registers.E = 0xD8
	core.CPU.Registers.F = 0xB0
	core.CPU.Registers.HL = 0x014D
	core.CPU.Registers.PC = 0x0100
	core.CPU.Registers.SP = 0xFFFE

}

/*
	Execute the next  OPCode and return used CPU clock
*/
func (core *Core) ExecuteNextOPCode() int {
	opcode := core.ReadMemory(core.CPU.Registers.PC)
	core.CPU.Registers.PC++
	return core.ExecuteOPCode(opcode)
}

/*
	Break execution and print register/dump memory information for debug purpose.
*/
func (core *Core) Break(code byte) {
	af := core.CPU.getAF()
	bc := core.CPU.getBC()
	de := core.CPU.getDE()
	hl := core.CPU.Registers.HL
	sp := core.CPU.Registers.SP
	pc := core.CPU.Registers.PC - 1
	lcdc := core.Memory.MainMemory[0xFF40]
	IF := core.Memory.MainMemory[0xFF0F]
	IE := core.Memory.MainMemory[0xFFFF]
	core.Memory.Dump("memory.dump")
	log.Printf("[Debug] \n\033[34m[OP:%s]\nAF:%04X  BC:%04X  DE:%04X  HL:%04X  SP:%04X\nPC:%04X  LCDC:%02X  IF:%02X    IE:%02X    IME:%t\nLCD:%X \033[0m", OPCodeFunctionMap[code].OP, af, bc, de, hl, sp, pc, lcdc, IF, IE, core.CPU.Flags.InterruptMaster, core.DebugControl)
}

/*
	Execute given OPCode and return used CPU clock
*/
func (core *Core) ExecuteOPCode(code byte) int {
	if OPCodeFunctionMap[code].Clock != 0 {
		if core.DebugControl == core.CPU.Registers.PC-1 && core.Debug {
			core.Break(code)
			_, err := fmt.Scanf("%X", &core.DebugControl)
			if err != nil {
				core.CPU.Registers.PC = 0
			}
		}
		var extCycles int
		extCycles = OPCodeFunctionMap[code].Func(core)
		return OPCodeFunctionMap[code].Clock + extCycles
	} else {
		if core.Debug {
			core.Break(code)
		}
		log.Fatalf("Unable to resolve OPCode:%X   PC:%X\n", code, core.CPU.Registers.PC-1)
		return 0
	}
}

/*
	Get 16bit parameter after OpCode
*/
func (core *Core) getParameter16() uint16 {
	b1 := uint16(core.ReadMemory(core.CPU.Registers.PC))
	b2 := uint16(core.ReadMemory(core.CPU.Registers.PC + 1))
	core.CPU.Registers.PC += 2
	return b2<<8 | b1
}

/*
	Set value of register A
*/
func (core *Core) setA(val byte) {
	core.CPU.Registers.A = val
}

/*
	Set value of register B
*/
func (core *Core) setB(val byte) {
	core.CPU.Registers.B = val
}

/*
	Set value of register C
*/
func (core *Core) setC(val byte) {
	core.CPU.Registers.C = val
}

/*
	Set value of register D
*/
func (core *Core) setD(val byte) {
	core.CPU.Registers.D = val
}

/*
	Set value of register D
*/
func (core *Core) setE(val byte) {
	core.CPU.Registers.E = val
}

/*
	Set value of register H
*/
func (core *Core) setH(val byte) {
	core.CPU.Registers.HL &= 0x00FF
	core.CPU.Registers.HL |= ((uint16(val) << 8) & 0xFF00) // OR in the desired mask
}

/*
	Set value of register L
*/
func (core *Core) setL(val byte) {
	core.CPU.Registers.HL &= 0xFF00
	core.CPU.Registers.HL |= (uint16(val) & 0x00FF) // OR in the desired mask
}

/*
	Get 8bit parameter after opcode
*/
func (core *Core) getParameter8() byte {
	b := core.ReadMemory(core.CPU.Registers.PC)
	core.CPU.Registers.PC += 1
	return b
}

/*
	Get value of AF register
*/
func (cpu *CPU) getAF() uint16 {
	return uint16(cpu.Registers.A)<<8 | uint16(cpu.Registers.F)
}

/*
	Set value of AF register
*/
func (cpu *CPU) setAF(val uint16) {
	cpu.Registers.A = byte((val & 0xFF00) >> 8)
	cpu.Registers.F = byte(val & 0xFF)
}

/*
	Set value of BC register
*/
func (cpu *CPU) setBC(val uint16) {
	cpu.Registers.B = byte((val & 0xFF00) >> 8)
	cpu.Registers.C = byte(val & 0xFF)
}

/*
	Set value of DE register
*/
func (cpu *CPU) setDE(val uint16) {
	cpu.Registers.D = byte((val & 0xFF00) >> 8)
	cpu.Registers.E = byte(val & 0xFF)
}

/*
	Get value of BC register
*/
func (cpu *CPU) getBC() uint16 {
	return uint16(cpu.Registers.B)<<8 | uint16(cpu.Registers.C)
}

/*
	Get value of DE register
*/
func (cpu *CPU) getDE() uint16 {
	return uint16(cpu.Registers.D)<<8 | uint16(cpu.Registers.E)
}

/*
	Update Low 8bit of AF register
	TODO: maybe this operation is useless in non-debug mode.
*/
func (cpu *CPU) updateAFLow() {
	newAF := cpu.Registers.F
	if cpu.Flags.Zero {
		newAF = util.SetBit(newAF, 7)
	} else {
		newAF = util.ClearBit(newAF, 7)
	}

	if cpu.Flags.Sub {
		newAF = util.SetBit(newAF, 6)
	} else {
		newAF = util.ClearBit(newAF, 6)
	}

	if cpu.Flags.HalfCarry {
		newAF = util.SetBit(newAF, 5)
	} else {
		newAF = util.ClearBit(newAF, 5)
	}

	if cpu.Flags.Carry {
		newAF = util.SetBit(newAF, 4)
	} else {
		newAF = util.ClearBit(newAF, 4)
	}

	cpu.Registers.F = newAF

}

/*
	Compare two values and set flags
*/

func (cpu *CPU) Compare(val1 byte, val2 byte) {
	cpu.Flags.Zero = (val1 == val2)
	cpu.Flags.Carry = (val1 > val2)
	cpu.Flags.HalfCarry = ((val1 & 0x0f) > (val2 & 0x0f))
	cpu.Flags.Sub = true

	cpu.updateAFLow()

}
