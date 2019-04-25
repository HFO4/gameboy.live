package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"log"
)

var OPCodeFunctionMap = [0x100]OPCodeUnit{
	//0x0_
	byte(0x00): {
		Func:  (*Core).OP00,
		Clock: 4,
		OP:    "NOP",
	},
	byte(0x01): {
		Func:  (*Core).OP01,
		Clock: 12,
		OP:    "LD BC,d16",
	},
	byte(0x03): {
		Func:  (*Core).OP03,
		Clock: 8,
		OP:    "INC BC",
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
	byte(0x09): {
		Func:  (*Core).OP09,
		Clock: 8,
		OP:    "ADD HL,BC",
	},
	byte(0x0B): {
		Func:  (*Core).OP0B,
		Clock: 8,
		OP:    "DEC BC",
	},
	byte(0x0A): {
		Func:  (*Core).OP0A,
		Clock: 8,
		OP:    "LD A,(BC)",
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
	byte(0x11): {
		Func:  (*Core).OP11,
		Clock: 12,
		OP:    "LD DE,d16",
	},
	byte(0x12): {
		Func:  (*Core).OP12,
		Clock: 8,
		OP:    "LD (DE),A",
	},
	byte(0x13): {
		Func:  (*Core).OP13,
		Clock: 8,
		OP:    "INC DE",
	},
	byte(0x16): {
		Func:  (*Core).OP16,
		Clock: 8,
		OP:    "LD D,d8",
	},
	byte(0x18): {
		Func:  (*Core).OP18,
		Clock: 12,
		OP:    "JR r8",
	},
	byte(0x19): {
		Func:  (*Core).OP19,
		Clock: 8,
		OP:    "ADD HL,DE",
	},
	byte(0x1A): {
		Func:  (*Core).OP1A,
		Clock: 8,
		OP:    "LD A,(DE)",
	},
	byte(0x1C): {
		Func:  (*Core).OP1C,
		Clock: 4,
		OP:    "INC E",
	},
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
	byte(0x22): {
		Func:  (*Core).OP22,
		Clock: 8,
		OP:    "LD (HL+),A",
	},
	byte(0x23): {
		Func:  (*Core).OP23,
		Clock: 8,
		OP:    "INC HL",
	},
	byte(0x28): {
		Func:  (*Core).OP28,
		Clock: 8,
		OP:    "JR Z,r8",
	},
	byte(0x2A): {
		Func:  (*Core).OP2A,
		Clock: 8,
		OP:    "LD A,(HL+)",
	},
	byte(0x2C): {
		Func:  (*Core).OP2C,
		Clock: 4,
		OP:    "INC L",
	},
	byte(0x2F): {
		Func:  (*Core).OP2F,
		Clock: 4,
		OP:    "CPL",
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
	byte(0x34): {
		Func:  (*Core).OP34,
		Clock: 12,
		OP:    "INC (HL)",
	},
	byte(0x36): {
		Func:  (*Core).OP36,
		Clock: 12,
		OP:    "LD (HL),d8",
	},
	byte(0x35): {
		Func:  (*Core).OP35,
		Clock: 12,
		OP:    "DEC (HL)",
	},
	byte(0x3A): {
		Func:  (*Core).OP3A,
		Clock: 8,
		OP:    "LD A,(HL-)",
	},
	byte(0x3C): {
		Func:  (*Core).OP3C,
		Clock: 4,
		OP:    "INC A",
	},
	byte(0x3D): {
		Func:  (*Core).OP3D,
		Clock: 4,
		OP:    "DEC A",
	},
	byte(0x3E): {
		Func:  (*Core).OP3E,
		Clock: 8,
		OP:    "LD A,d8",
	},
	//0x4_
	byte(0x46): {
		Func:  (*Core).OP46,
		Clock: 8,
		OP:    "LD B,(HL)",
	},
	byte(0x47): {
		Func:  (*Core).OP47,
		Clock: 4,
		OP:    "LD B,A",
	},
	byte(0x4E): {
		Func:  (*Core).OP4E,
		Clock: 8,
		OP:    "LD C,(HL)",
	},
	byte(0x4F): {
		Func:  (*Core).OP4F,
		Clock: 4,
		OP:    "LD C,A",
	},
	//0x5_
	byte(0x56): {
		Func:  (*Core).OP56,
		Clock: 8,
		OP:    "LD D,(HL)",
	},
	byte(0x57): {
		Func:  (*Core).OP57,
		Clock: 4,
		OP:    "LD D,A",
	},
	byte(0x5E): {
		Func:  (*Core).OP5E,
		Clock: 8,
		OP:    "LD E,(HL)",
	},
	byte(0x5F): {
		Func:  (*Core).OP5F,
		Clock: 4,
		OP:    "LD E,A",
	},
	//0X6_
	byte(0x60): {
		Func:  (*Core).OP60,
		Clock: 4,
		OP:    "LD H,B",
	},
	byte(0x69): {
		Func:  (*Core).OP69,
		Clock: 4,
		OP:    "LD L,C",
	},
	byte(0x6F): {
		Func:  (*Core).OP6F,
		Clock: 4,
		OP:    "LD L,A",
	},
	//0X7_
	byte(0x71): {
		Func:  (*Core).OP71,
		Clock: 8,
		OP:    "LD (HL),C",
	},
	byte(0x72): {
		Func:  (*Core).OP72,
		Clock: 8,
		OP:    "LD (HL),D",
	},
	byte(0x73): {
		Func:  (*Core).OP73,
		Clock: 8,
		OP:    "LD (HL),E",
	},
	byte(0x77): {
		Func:  (*Core).OP77,
		Clock: 8,
		OP:    "LD (HL),A",
	},
	byte(0x78): {
		Func:  (*Core).OP78,
		Clock: 4,
		OP:    "LD A,B",
	},
	byte(0x79): {
		Func:  (*Core).OP79,
		Clock: 4,
		OP:    "LD A,C",
	},
	byte(0x7A): {
		Func:  (*Core).OP7A,
		Clock: 4,
		OP:    "LD A,D",
	},
	byte(0x7B): {
		Func:  (*Core).OP7B,
		Clock: 4,
		OP:    "LD A,E",
	},
	byte(0x7C): {
		Func:  (*Core).OP7C,
		Clock: 4,
		OP:    "LD A,H",
	},
	byte(0x7E): {
		Func:  (*Core).OP7E,
		Clock: 8,
		OP:    "LD A,(HL)",
	},
	//0x8_
	byte(0x85): {
		Func:  (*Core).OP85,
		Clock: 4,
		OP:    "ADD A,L",
	},
	byte(0x87): {
		Func:  (*Core).OP87,
		Clock: 4,
		OP:    "ADD A,A",
	},
	//0xA_
	byte(0xAF): {
		Func:  (*Core).OPAF,
		Clock: 4,
		OP:    "XOR A",
	},
	byte(0xA0): {
		Func:  (*Core).OPA0,
		Clock: 4,
		OP:    "AND B",
	},
	byte(0xA1): {
		Func:  (*Core).OPA1,
		Clock: 4,
		OP:    "AND C",
	},
	byte(0xA2): {
		Func:  (*Core).OPA2,
		Clock: 4,
		OP:    "AND D",
	},
	byte(0xA3): {
		Func:  (*Core).OPA3,
		Clock: 4,
		OP:    "AND E",
	},
	byte(0xA4): {
		Func:  (*Core).OPA4,
		Clock: 4,
		OP:    "AND H",
	},
	byte(0xA5): {
		Func:  (*Core).OPA5,
		Clock: 4,
		OP:    "AND L",
	},
	byte(0xA6): {
		Func:  (*Core).OPA6,
		Clock: 4,
		OP:    "AND (HL)",
	},
	byte(0xA7): {
		Func:  (*Core).OPA7,
		Clock: 4,
		OP:    "AND A",
	},
	byte(0xA9): {
		Func:  (*Core).OPA9,
		Clock: 4,
		OP:    "XOR C",
	},
	//0xB_
	byte(0xB0): {
		Func:  (*Core).OPB0,
		Clock: 4,
		OP:    "OR B",
	},
	byte(0xB1): {
		Func:  (*Core).OPB1,
		Clock: 4,
		OP:    "OR C",
	},
	//0xC_
	byte(0xC0): {
		Func:  (*Core).OPC0,
		Clock: 8,
		OP:    "RET NZ",
	},
	byte(0xC1): {
		Func:  (*Core).OPC1,
		Clock: 12,
		OP:    "POP BC",
	},
	byte(0xC2): {
		Func:  (*Core).OPC2,
		Clock: 12,
		OP:    "JP NZ,a16",
	},
	byte(0xC3): {
		Func:  (*Core).OPC3,
		Clock: 16,
		OP:    "JP a16",
	},
	byte(0xC5): {
		Func:  (*Core).OPC5,
		Clock: 16,
		OP:    "PUSH BC",
	},
	byte(0xC8): {
		Func:  (*Core).OPC8,
		Clock: 8,
		OP:    "RET Z",
	},
	byte(0xC9): {
		Func:  (*Core).OPC9,
		Clock: 16,
		OP:    "RET",
	},
	byte(0xCA): {
		Func:  (*Core).OPCA,
		Clock: 12,
		OP:    "JP Z,a16",
	},
	byte(0xCB): {
		Func:  (*Core).OPCB,
		Clock: 4,
		OP:    "PERFIX CB",
	},
	byte(0xCD): {
		Func:  (*Core).OPCD,
		Clock: 24,
		OP:    "CALL a16",
	},
	//0xD_
	byte(0xD1): {
		Func:  (*Core).OPD1,
		Clock: 12,
		OP:    "POP DE",
	},
	byte(0xD5): {
		Func:  (*Core).OPD5,
		Clock: 16,
		OP:    "PUSH DE",
	},
	byte(0xD9): {
		Func:  (*Core).OPD9,
		Clock: 16,
		OP:    "RETI",
	},
	//0xE_
	byte(0xE0): {
		Func:  (*Core).OPE0,
		Clock: 12,
		OP:    "LDH (a8),A",
	},
	byte(0xE1): {
		Func:  (*Core).OPE1,
		Clock: 12,
		OP:    "POP HL",
	},
	byte(0xE2): {
		Func:  (*Core).OPE2,
		Clock: 8,
		OP:    "LD (C),A",
	},
	byte(0xE5): {
		Func:  (*Core).OPE5,
		Clock: 16,
		OP:    "PUSH HL",
	},
	byte(0xE6): {
		Func:  (*Core).OPE6,
		Clock: 8,
		OP:    "AND d8",
	},
	byte(0xE9): {
		Func:  (*Core).OPE9,
		Clock: 4,
		OP:    "JP (HL)",
	},
	byte(0xEA): {
		Func:  (*Core).OPEA,
		Clock: 16,
		OP:    "LD (a16),A",
	},
	byte(0xEF): {
		Func:  (*Core).OPEF,
		Clock: 16,
		OP:    "RST 28H",
	},
	//0xF_
	byte(0xF0): {
		Func:  (*Core).OPF0,
		Clock: 12,
		OP:    "LDH A,(a8)",
	},
	byte(0xF1): {
		Func:  (*Core).OPF1,
		Clock: 12,
		OP:    "POP AF",
	},
	byte(0xF3): {
		Func:  (*Core).OPF3,
		Clock: 4,
		OP:    "DI",
	},
	byte(0xF5): {
		Func:  (*Core).OPF5,
		Clock: 16,
		OP:    "PUSH AF",
	},
	byte(0xFA): {
		Func:  (*Core).OPFA,
		Clock: 16,
		OP:    "LD A,(a16)",
	},
	byte(0xFB): {
		Func:  (*Core).OPFB,
		Clock: 4,
		OP:    "EI",
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
	OP:0x71 LD (HL),C
*/
func (core *Core) OP71() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.C)
	return 0
}

/*
	OP:0x72 LD (HL),D
*/
func (core *Core) OP72() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.D)
	return 0
}

/*
	OP:0x73 LD (HL),E
*/
func (core *Core) OP73() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.E)
	return 0
}

/*
	OP:0x7A LD A,D
*/
func (core *Core) OP7A() int {
	core.CPU.Registers.A = core.CPU.Registers.D
	return 0
}

/*
	OP:0x7B LD A,E
*/
func (core *Core) OP7B() int {
	core.CPU.Registers.A = core.CPU.Registers.E
	return 0
}

/*
	OP:0x57 LD D,A
*/
func (core *Core) OP57() int {
	core.CPU.Registers.D = core.CPU.Registers.A
	return 0
}

/*
	OP:0x3A LD A,(HL-)
*/
func (core *Core) OP3A() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.Registers.HL)
	core.CPU.Registers.HL--
	return 0
}

/*
	OP:0xC2 JP NZ,a16
*/
func (core *Core) OPC2() int {
	address := core.getParameter16()
	if !core.CPU.Flags.Zero {
		core.CPU.Registers.PC = address
		return 4
	}
	return 0
}

/*
	OP:0x6F LD L,A
*/
func (core *Core) OP6F() int {
	core.setL(core.CPU.Registers.A)
	return 0
}

/*
	OP:0x85 ADD A,L
*/
func (core *Core) OP85() int {
	origin1 := core.CPU.Registers.A
	origin2 := byte(core.CPU.Registers.HL & 0xFF)
	res := int16(core.CPU.Registers.A) + int16(core.CPU.Registers.HL&0xFF)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = res == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x03 INC BC
*/
func (core *Core) OP03() int {
	core.CPU.setBC(core.CPU.getBC() + 1)
	return 0
}

/*
	OP:0x0A LD A,(BC)
*/
func (core *Core) OP0A() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.getBC())
	return 0
}

/*
	OP:0x60 LD H,B
*/
func (core *Core) OP60() int {
	core.setH(core.CPU.Registers.B)
	return 0
}

/*
	OP:0x69 LD L,C
*/
func (core *Core) OP69() int {
	core.setL(core.CPU.Registers.C)
	return 0
}

/*
	OP:0x46 LD B,(HL)
*/
func (core *Core) OP46() int {
	core.CPU.Registers.B = core.ReadMemory(core.CPU.Registers.HL)
	return 0
}

/*
	OP:0x4E LD C,(HL)
*/
func (core *Core) OP4E() int {
	core.CPU.Registers.C = core.ReadMemory(core.CPU.Registers.HL)
	return 0
}

/*
	OP:0x09 ADD HL,BC
*/
func (core *Core) OP09() int {
	originHL := core.CPU.Registers.HL
	originBC := core.CPU.getBC()
	res := int32(originBC) + int32(originHL)
	core.CPU.Registers.HL = uint16(res)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = int32(originHL&0xFFF) > (res & 0xFFF)
	core.CPU.Flags.Carry = res > 0xFFFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x77 LD (HL),A
*/
func (core *Core) OP77() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.A)
	return 0
}

/*
	OP:0x2C INC L
*/
func (core *Core) OP2C() int {
	origin := byte(core.CPU.Registers.HL & 0xFF)
	newVal := origin + 1
	core.setL(newVal)
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x35 DEC (HL)
*/
func (core *Core) OP35() int {
	origin := core.ReadMemory(core.CPU.Registers.HL)
	newVal := origin - 1
	core.WriteMemory(core.CPU.Registers.HL, newVal)
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x18 JR r8
*/
func (core *Core) OP18() int {
	address := int32(core.CPU.Registers.PC) + int32(int8(core.getParameter8()))
	core.CPU.Registers.PC = uint16(address)
	return 0
}

/*
	OP:0x7E LD A,(HL)
*/
func (core *Core) OP7E() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.Registers.HL)
	return 0
}

/*
	OP:0xCA JP Z,a16
*/
func (core *Core) OPCA() int {
	address := core.getParameter16()
	if core.CPU.Flags.Zero {
		core.CPU.Registers.PC = address
		return 4
	}
	return 0
}

/*
	OP:0x1C INC E
*/
func (core *Core) OP1C() int {
	origin := core.CPU.Registers.E
	newVal := origin + 1
	core.CPU.Registers.E = newVal
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x7C LD A,H
*/
func (core *Core) OP7C() int {
	core.CPU.Registers.A = byte(core.CPU.Registers.HL >> 8)
	return 0
}

/*
	OP:0x22 LD (HL+),A
*/
func (core *Core) OP22() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.A)
	core.CPU.Registers.HL++
	return 0
}

/*
	OP:0x1A LD A,(DE)
*/
func (core *Core) OP1A() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.getDE())
	return 0
}

/*
	OP:0x13 INC DE
*/
func (core *Core) OP13() int {
	core.CPU.setDE(core.CPU.getDE() + 1)
	return 0
}

/*
	OP:0x12 LD (DE),A
*/
func (core *Core) OP12() int {
	core.WriteMemory(core.CPU.getDE(), core.CPU.Registers.A)
	return 0
}

/*
	OP:0x11 LD DE,d16
*/
func (core *Core) OP11() int {
	core.CPU.setDE(core.getParameter16())
	return 0
}

/*
	OP:0xE9 JP (HL)
*/
func (core *Core) OPE9() int {
	core.CPU.Registers.PC = core.CPU.Registers.HL
	return 0
}

/*
	OP:0x56 LD D,(HL)
*/
func (core *Core) OP56() int {
	core.CPU.Registers.D = core.ReadMemory(core.CPU.Registers.HL)
	return 0
}

/*
	OP:0x23 INC HL
*/
func (core *Core) OP23() int {
	core.CPU.Registers.HL += 1
	return 0
}

/*
	OP:0x5E LD E,(HL)
*/
func (core *Core) OP5E() int {
	core.CPU.Registers.E = core.ReadMemory(core.CPU.Registers.HL)
	return 0
}

/*
	OP:0x19 ADD HL,DE
*/
func (core *Core) OP19() int {
	originHL := core.CPU.Registers.HL
	originDE := core.CPU.getDE()
	res := int32(originDE) + int32(originHL)
	core.CPU.Registers.HL = uint16(res)

	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = int32(originHL&0xFFF) > (res & 0xFFF)
	core.CPU.Flags.Carry = res > 0xFFFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x16 LD D,d8
*/
func (core *Core) OP16() int {
	core.CPU.Registers.D = core.getParameter8()
	return 0
}

/*
	OP:0x5F LD E,A
*/
func (core *Core) OP5F() int {
	core.CPU.Registers.E = core.CPU.Registers.A
	return 0
}

/*
	OP:0x87 ADD A,A
*/
func (core *Core) OP87() int {
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) + int16(core.CPU.Registers.A)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = res == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = 2*(origin&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xEF RST 28H
*/
func (core *Core) OPEF() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0028
	return 0
}

/*
	OP:0x79 LD A,C
*/
func (core *Core) OP79() int {
	core.CPU.Registers.A = core.CPU.Registers.C
	return 0
}

/*
	OP:0xA9 XOR C
*/
func (core *Core) OPA9() int {
	core.CPU.Registers.A = core.CPU.Registers.C ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x4F LD C,A
*/
func (core *Core) OP4F() int {
	core.CPU.Registers.C = core.CPU.Registers.A
	return 0
}

/*
	OP:0xB0 OR B
*/
func (core *Core) OPB0() int {
	core.CPU.Registers.B = core.CPU.Registers.A | core.CPU.Registers.B
	core.CPU.Flags.Zero = (core.CPU.Registers.B == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x47 LD B,A
*/
func (core *Core) OP47() int {
	core.CPU.Registers.B = core.CPU.Registers.A
	return 0
}

/*
	OP:0xCB PREFIX CB
*/
func (core *Core) OPCB() int {
	nextIns := core.getParameter8()
	if cbMap[nextIns] != nil {
		cbMap[nextIns]()
		return CBCycles[nextIns] * 4
	} else {
		log.Fatalf("Undefined CB Opcode: %X \n", nextIns)
	}
	return 0
}

/*
	OP:0xE6 AND d8
*/
func (core *Core) OPE6() int {
	val := core.getParameter8()
	core.CPU.Registers.A = core.CPU.Registers.A & val
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.Flags.Sub = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x2F CPL
*/
func (core *Core) OP2F() int {
	core.CPU.Registers.A = 0XFF ^ core.CPU.Registers.A
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = true
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xD9 RETI
*/
func (core *Core) OPD9() int {
	core.CPU.Registers.PC = core.StackPop()
	core.CPU.Flags.InterruptMaster = true
	return 0
}

/*
	OP:0xF1 POP AF
*/
func (core *Core) OPF1() int {
	core.CPU.setAF(core.StackPop())
	core.CPU.Flags.Zero = util.TestBit(core.CPU.Registers.F, 7)
	core.CPU.Flags.Sub = util.TestBit(core.CPU.Registers.F, 6)
	core.CPU.Flags.HalfCarry = util.TestBit(core.CPU.Registers.F, 5)
	core.CPU.Flags.Carry = util.TestBit(core.CPU.Registers.F, 4)
	return 0
}

/*
	OP:0xC1 POP BC
*/
func (core *Core) OPC1() int {
	core.CPU.setBC(core.StackPop())
	return 0
}

/*
	OP:0xD1 POP DE
*/
func (core *Core) OPD1() int {
	core.CPU.setDE(core.StackPop())
	return 0
}

/*
	OP:0xE1 POP HL
*/
func (core *Core) OPE1() int {
	core.CPU.Registers.HL = core.StackPop()
	return 0
}

/*
	OP:0x3C INC A
*/
func (core *Core) OP3C() int {
	origin := core.CPU.Registers.A
	newVal := origin + 1
	core.CPU.Registers.A = newVal
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x34 INC (HL)
*/
func (core *Core) OP34() int {
	origin := core.ReadMemory(core.CPU.Registers.HL)
	newVal := origin + 1
	core.WriteMemory(core.CPU.Registers.HL, newVal)
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x3D DEC A
*/
func (core *Core) OP3D() int {
	origin := core.CPU.Registers.A
	core.CPU.Registers.A--
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xC8 RET Z
*/
func (core *Core) OPC8() int {
	if core.CPU.Flags.Zero {
		core.CPU.Registers.PC = core.StackPop()
		return 12
	}
	return 0
}

/*
	OP:0xC0 LD A,(a16)
*/
func (core *Core) OPFA() int {
	core.CPU.Registers.A = core.ReadMemory(core.getParameter16())
	return 0
}

/*
	OP:0xC0 RET NZ
*/
func (core *Core) OPC0() int {
	if !core.CPU.Flags.Zero {
		core.CPU.Registers.PC = core.StackPop()
		return 12
	}
	return 0
}

/*
	OP:0x28 JR Z,r8
*/
func (core *Core) OP28() int {
	address := int8(core.getParameter8())
	if core.CPU.Flags.Zero {
		core.CPU.Registers.PC = uint16(int32(core.CPU.Registers.PC) + int32(address))
		return 4
	}
	return 0
}

/*
	OP:0xA6 AND (HL)
*/
func (core *Core) OPA6() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.Registers.HL) & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA5 AND L
*/
func (core *Core) OPA5() int {
	core.CPU.Registers.A = byte(core.CPU.Registers.HL&0xFF) & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA4 AND H
*/
func (core *Core) OPA4() int {
	core.CPU.Registers.A = byte(core.CPU.Registers.HL>>8) & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA3 AND E
*/
func (core *Core) OPA3() int {
	core.CPU.Registers.A = core.CPU.Registers.E & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA2 AND D
*/
func (core *Core) OPA2() int {
	core.CPU.Registers.A = core.CPU.Registers.D & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA1 AND C
*/
func (core *Core) OPA1() int {
	core.CPU.Registers.A = core.CPU.Registers.C & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA0 AND B
*/
func (core *Core) OPA0() int {
	core.CPU.Registers.A = core.CPU.Registers.B & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xA7 AND A
*/
func (core *Core) OPA7() int {
	core.CPU.Registers.A = core.CPU.Registers.A & core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xE5 PUSH HL
*/
func (core *Core) OPE5() int {
	core.StackPush(core.CPU.Registers.HL)
	return 0
}

/*
	OP:0xD5 PUSH DE
*/
func (core *Core) OPD5() int {
	core.StackPush(core.CPU.getDE())
	return 0
}

/*
	OP:0xC5 PUSH BC
*/
func (core *Core) OPC5() int {
	core.StackPush(core.CPU.getBC())
	return 0
}

/*
	OP:0xF5 PUSH AF
*/
func (core *Core) OPF5() int {
	core.StackPush(core.CPU.getAF())
	return 0
}

/*
	OP:0xFB EI
*/
func (core *Core) OPFB() int {
	core.CPU.Flags.PendingInterruptEnabled = true
	return 0
}

/*
	OP:0xC9 RET
*/
func (core *Core) OPC9() int {
	core.CPU.Registers.PC = core.StackPop()
	return 0
}

/*
	OP:0xB1 OR C
*/
func (core *Core) OPB1() int {
	core.CPU.Registers.A = core.CPU.Registers.A | core.CPU.Registers.C
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x78 LD A,B
*/
func (core *Core) OP78() int {
	core.CPU.Registers.A = core.CPU.Registers.B
	return 0
}

/*
	OP:0x0B DEC BC
*/
func (core *Core) OP0B() int {
	core.CPU.setBC(core.CPU.getBC() - 1)
	return 0
}

/*
	OP:0x01 LD BC,d16
*/
func (core *Core) OP01() int {
	core.CPU.setBC(core.getParameter16())
	return 0
}

/*
	OP:0xCD CALL a16
*/
func (core *Core) OPCD() int {
	nextAddress := core.getParameter16()
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = nextAddress
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
	core.CPU.Compare(core.getParameter8(), core.CPU.Registers.A)
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
