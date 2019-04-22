package gb

import "log"

type CPU struct {
	Registers Registers
	Flags     Flags
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
	Conatins the result from the recent instruction which has affected flags.
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
	core.CPU.Registers.HL = 0x014D
	core.CPU.Registers.PC = 0x0100
	core.CPU.Registers.SP = 0xFFFE

}
