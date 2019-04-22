package gb

import "github.com/HFO4/gbc-in-cloud/util"

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
	byte(0x0E): {
		Func:  (*Core).OP0E,
		Clock: 8,
		OP:    "LD C,d8",
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
	//0x1_
	byte(0x1F): {
		Func:  (*Core).OP1F,
		Clock: 4,
		OP:    "RRA",
	},
	//0x3_
	byte(0x32): {
		Func:  (*Core).OP32,
		Clock: 8,
		OP:    "LD (HL-),A",
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
}

type OPCodeUnit struct {
	Func  func(*Core) int
	Clock int
	OP    string
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
	core.CPU.Flags.Zero = core.CPU.Registers.B == 0
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
