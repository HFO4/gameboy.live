package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
)

var OPCodeFunctionMap = map[byte]OPCodeUnit{
	//0x0_
	byte(0x00): {
		Func:  (*Core).OP00,
		Clock: 4,
		OP:    "NOP",
	},
	byte(0x05): {
		Func:  (*Core).OP05,
		Clock: 4,
		OP:    "DEC B",
	},
	byte(0x06): {
		Func:  (*Core).OP06,
		Clock: 8,
		OP:    "LD B,d8",
	},
	byte(0x0C): {
		Func:  (*Core).OP0C,
		Clock: 4,
		OP:    "INC C",
	},
	byte(0x0D): {
		Func:  (*Core).OP0D,
		Clock: 4,
		OP:    "DEC C",
	},
	byte(0x0E): {
		Func:  (*Core).OP0E,
		Clock: 8,
		OP:    "LD C,d8",
	},
	//0x1_
	byte(0x1F): {
		Func:  (*Core).OP1F,
		Clock: 4,
		OP:    "RRA",
	},
	//0x2_
	byte(0x20): {
		Func:  (*Core).OP20,
		Clock: 8,
		OP:    "JR NZ,r8",
	},
	byte(0x21): {
		Func:  (*Core).OP21,
		Clock: 12,
		OP:    "LD HL,d16",
	},
	byte(0x2A): {
		Func:  (*Core).OP2A,
		Clock: 8,
		OP:    "LD A,(HL+)",
	},
	//0x3_
	byte(0x31): {
		Func:  (*Core).OP31,
		Clock: 12,
		OP:    "LD SP,d16",
	},
	byte(0x32): {
		Func:  (*Core).OP32,
		Clock: 8,
		OP:    "LD (HL-),A",
	},
	byte(0x36): {
		Func:  (*Core).OP36,
		Clock: 12,
		OP:    "LD (HL),d8",
	},
	byte(0x3E): {
		Func:  (*Core).OP3E,
		Clock: 8,
		OP:    "LD A,d8",
	},
	//0xA_
	byte(0xAF): {
		Func:  (*Core).OPAF,
		Clock: 4,
		OP:    "XOR A",
	},
	//0xC_
	byte(0xC3): {
		Func:  (*Core).OPC3,
		Clock: 16,
		OP:    "JP a16",
	},
	byte(0xCD): {
		Func:  (*Core).OPCD,
		Clock: 24,
		OP:    "CALL a16",
	},
	//0xE_
	byte(0xE0): {
		Func:  (*Core).OPE0,
		Clock: 12,
		OP:    "LDH (a8),A",
	},
	byte(0xE2): {
		Func:  (*Core).OPE2,
		Clock: 8,
		OP:    "LD (C),A",
	},
	byte(0xEA): {
		Func:  (*Core).OPEA,
		Clock: 16,
		OP:    "LD (a16),A",
	},
	//0xF_
	byte(0xF0): {
		Func:  (*Core).OPF0,
		Clock: 12,
		OP:    "LDH A,(a8)",
	},
	byte(0xF3): {
		Func:  (*Core).OPF3,
		Clock: 4,
		OP:    "DI",
	},
	byte(0xFE): {
		Func:  (*Core).OPFE,
		Clock: 8,
		OP:    "CP d8",
	},
}

type OPCodeUnit struct {
	Func  func(*Core) int
	Clock int
	OP    string
}

/*
	OP:0xCD CALL a16
*/
func (core *Core) OPCD() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = core.getParameter16()
	return 0
}

/*
	OP:0x0C INC C
*/
func (core *Core) OP0C() int {
	origin := core.CPU.Registers.C
	core.CPU.Registers.C++
	core.CPU.Flags.Zero = (core.CPU.Registers.C == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin&0xF)+(1&0xF) > 0xF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xE2 LD (C),A
*/
func (core *Core) OPE2() int {
	core.WriteMemory(0xFF00+uint16(core.CPU.Registers.C), core.CPU.Registers.A)
	return 0
}

/*
	OP:0x2A LD A,(HL+)
*/
func (core *Core) OP2A() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.Registers.HL)
	core.CPU.Registers.HL++
	return 0
}

/*
	OP:0x31 LD SP,d16
*/
func (core *Core) OP31() int {
	core.CPU.Registers.SP = core.getParameter16()
	return 0
}

/*
	OP:0xEA LD (a16),A
*/
func (core *Core) OPEA() int {
	core.WriteMemory(core.getParameter16(), core.CPU.Registers.A)
	return 0
}

/*
	OP:0x36 LD (HL),d8
*/
func (core *Core) OP36() int {
	core.WriteMemory(core.CPU.Registers.HL, core.getParameter8())
	return 0
}

/*
	OP:0xFE CP d8
*/
func (core *Core) OPFE() int {
	core.CPU.Compare(core.CPU.Registers.A, core.getParameter8())
	return 0
}

/*
	OP:0xF0 LDH LDH A,(a8)
*/
func (core *Core) OPF0() int {
	core.CPU.Registers.A = core.ReadMemory(0xFF00 + uint16(core.getParameter8()))
	return 0
}

/*
	OP:0xE0 LDH (a8),A
*/
func (core *Core) OPE0() int {
	core.WriteMemory(0xFF00+uint16(core.getParameter8()), core.CPU.Registers.A)
	return 0
}

/*
	OP:0xF3 DI
*/
func (core *Core) OPF3() int {
	core.CPU.Flags.PendingInterruptDisabled = true
	return 0
}

/*
	OP:0x3E LD A,d8
*/
func (core *Core) OP3E() int {
	core.CPU.Registers.A = core.getParameter8()
	return 0
}

/*
	OP:0X0D DEC C
*/
func (core *Core) OP0D() int {
	origin := core.CPU.Registers.C
	core.CPU.Registers.C--
	core.CPU.Flags.Zero = (core.CPU.Registers.C == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0X1F RRA
	TODO BE TESTED
*/
func (core *Core) OP1F() int {
	isLSBSet := util.TestBit(core.CPU.Registers.A, 0)
	core.CPU.Registers.A = core.CPU.Registers.A >> 1
	core.CPU.Flags.Zero = false
	core.CPU.Flags.Carry = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Sub = false
	if isLSBSet {
		core.CPU.Flags.Carry = true
		core.CPU.Registers.A = util.SetBit(core.CPU.Registers.A, 7)
	}
	if core.CPU.Registers.A == 0 {
		core.CPU.Flags.Zero = true
	}
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0X20 JR NZ,r8
*/
func (core *Core) OP20() int {
	address := int8(core.getParameter8())
	if !core.CPU.Flags.Zero {
		core.CPU.Registers.PC = uint16(int32(core.CPU.Registers.PC) + int32(address))
		return 4
	}
	return 0
}

/*
	OP:0X05 DEC B
*/
func (core *Core) OP05() int {
	origin := core.CPU.Registers.B
	core.CPU.Registers.B--
	core.CPU.Flags.Zero = (core.CPU.Registers.B == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0X32 LD (HL-),A
*/
func (core *Core) OP32() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.A)
	core.CPU.Registers.HL--
	return 0
}

/*


/*
	OP:0X00 NOP
*/
func (core *Core) OP00() int {
	return 0
}

/*
	OP:0X06 LD B,d8
*/
func (core *Core) OP06() int {
	core.CPU.Registers.B = core.getParameter8()
	return 0
}

/*
	OP:0X0E LD C,d8
*/
func (core *Core) OP0E() int {
	core.CPU.Registers.C = core.getParameter8()
	return 0
}

/*
	OP:0X21 LD HL,d16
*/
func (core *Core) OP21() int {
	core.CPU.Registers.HL = core.getParameter16()
	return 0
}

/*
	OP:0XAF XOR A
*/
func (core *Core) OPAF() int {
	core.CPU.Registers.A = core.CPU.Registers.A ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0XC3 JP a16
*/
func (core *Core) OPC3() int {
	core.CPU.Registers.PC = core.getParameter16()
	return 0
}
