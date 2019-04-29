package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"log"
)

/*
	OP:0xff RST 38H
*/
func (core *Core) OPFF() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0038
	return 0
}

/*
	OP:0xF9 LD SP,HL
*/
func (core *Core) OPF9() int {
	core.CPU.Registers.SP = core.CPU.Registers.HL
	return 0
}

/*
	OP:0xF8 LD HL,SP+r8
*/
func (core *Core) OPF8() int {
	val1 := int32(core.CPU.Registers.SP)
	val2 := int32(int8(core.getParameter8()))
	result := val1 + val2
	core.CPU.Registers.HL = uint16(result)
	tempVal := val1 ^ val2 ^ result
	core.CPU.Flags.Zero = false
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((tempVal & 0x10) == 0x10)
	core.CPU.Flags.Carry = ((tempVal & 0x100) == 0x100)
	return 0
}

/*
	OP:0xf7 RST 30H
*/
func (core *Core) OPF7() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0030
	return 0
}

/*
	OP:0xF2 LD A,(C)
*/
func (core *Core) OPF2() int {
	val := core.ReadMemory(0xFF00 + uint16(core.CPU.Registers.C))
	core.CPU.Registers.A = val
	return 0
}

/*
	OP:0xE8 ADD SP,r8
*/
func (core *Core) OPE8() int {
	origin1 := core.CPU.Registers.SP
	origin2 := int8(core.getParameter8())
	res := uint16(int32(core.CPU.Registers.SP) + int32(origin2))
	tmpVal := origin1 ^ uint16(origin2) ^ res
	core.CPU.Registers.SP = res

	core.CPU.Flags.Zero = false
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (tmpVal & 0x10) == 0x10
	core.CPU.Flags.Carry = ((tmpVal & 0x100) == 0x100)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xE7 RST 20H
*/
func (core *Core) OPE7() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0020
	return 0
}

/*
	OP:0xDF RST 18H
*/
func (core *Core) OPDF() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0018
	return 0
}

/*
	OP:0xDE SBC A,d8
*/
func (core *Core) OPDE() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val2 := core.getParameter8()
	origin := core.CPU.Registers.A
	dirtySum := int16(core.CPU.Registers.A) - int16(val2) - carry
	total := byte(dirtySum)
	core.CPU.Registers.A = total

	core.CPU.Flags.Zero = total == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val2&0xF)-carry < 0
	core.CPU.Flags.Carry = dirtySum < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xDB CALL C,a16
*/
func (core *Core) OPDC() int {
	address := core.getParameter16()
	if core.CPU.Flags.Carry {
		core.StackPush(core.CPU.Registers.PC)
		core.CPU.Registers.PC = address
		return 12
	}
	return 0
}

/*
	OP:0xDA JP C,a16
*/
func (core *Core) OPDA() int {
	address := core.getParameter16()
	if core.CPU.Flags.Carry {
		core.CPU.Registers.PC = address
		return 4
	}
	return 0
}

/*
	OP:0xD7 RST RST 10H
*/
func (core *Core) OPD7() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0010
	return 0
}

/*
	OP:0xD4 CALL NC,a16
*/
func (core *Core) OPD4() int {
	address := core.getParameter16()
	if !core.CPU.Flags.Carry {
		core.StackPush(core.CPU.Registers.PC)
		core.CPU.Registers.PC = address
		return 12
	}
	return 0
}

/*
	OP:0xD2 JP NC,a16
*/
func (core *Core) OPD2() int {
	address := core.getParameter16()
	if !core.CPU.Flags.Carry {
		core.CPU.Registers.PC = address
		return 4
	}
	return 0
}

/*
	OP:0xCF RST 08H
*/
func (core *Core) OPCF() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0008
	return 0
}

/*
	OP:0xCC CALL Z,a16
*/
func (core *Core) OPCC() int {
	address := core.getParameter16()
	if core.CPU.Flags.Zero {
		core.StackPush(core.CPU.Registers.PC)
		core.CPU.Registers.PC = address
		return 12
	}
	return 0
}

/*
	OP:0xC7 RST 00H
*/
func (core *Core) OPC7() int {
	core.StackPush(core.CPU.Registers.PC)
	core.CPU.Registers.PC = 0x0000
	return 0
}

/*
	OP:0xD8 RET C
*/
func (core *Core) OPD8() int {
	if core.CPU.Flags.Carry {
		core.CPU.Registers.PC = core.StackPop()
		return 12
	}
	return 0
}

/*
	OP:0xD0 RET NC
*/
func (core *Core) OPD0() int {
	if !core.CPU.Flags.Carry {
		core.CPU.Registers.PC = core.StackPop()
		return 12
	}
	return 0
}

/*
	OP:0xCE ADC ADC A,d8
*/
func (core *Core) OPCE() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.getParameter8()
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xBF CP A
*/
func (core *Core) OPBF() int {
	core.CPU.Compare(core.CPU.Registers.A, core.CPU.Registers.A)
	return 0
}

/*
	OP:0xBD CP L
*/
func (core *Core) OPBD() int {
	core.CPU.Compare(byte(core.CPU.Registers.HL&0xFF), core.CPU.Registers.A)
	return 0
}

/*
	OP:0xBC CP H
*/
func (core *Core) OPBC() int {
	core.CPU.Compare(byte((core.CPU.Registers.HL&0xFF00)>>8), core.CPU.Registers.A)
	return 0
}

/*
	OP:0xBB CP E
*/
func (core *Core) OPBB() int {
	core.CPU.Compare(core.CPU.Registers.E, core.CPU.Registers.A)
	return 0
}

/*
	OP:0xBA CP D
*/
func (core *Core) OPBA() int {
	core.CPU.Compare(core.CPU.Registers.D, core.CPU.Registers.A)
	return 0
}

/*
	OP:0xB9 CP C
*/
func (core *Core) OPB9() int {
	core.CPU.Compare(core.CPU.Registers.C, core.CPU.Registers.A)
	return 0
}

/*
	OP:0xB7 OR A
*/
func (core *Core) OPB7() int {
	core.CPU.Registers.A = core.CPU.Registers.A | core.CPU.Registers.A
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xB6 OR (HL)
*/
func (core *Core) OPB6() int {
	core.CPU.Registers.A = core.CPU.Registers.A | core.ReadMemory(core.CPU.Registers.HL)
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xB5 OR L
*/
func (core *Core) OPB5() int {
	core.CPU.Registers.A = core.CPU.Registers.A | byte(core.CPU.Registers.HL&0xFF)
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xB4 OR H
*/
func (core *Core) OPB4() int {
	core.CPU.Registers.A = core.CPU.Registers.A | byte((core.CPU.Registers.HL&0xFF00)>>8)
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xB3 OR E
*/
func (core *Core) OPB3() int {
	core.CPU.Registers.A = core.CPU.Registers.A | core.CPU.Registers.E
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xB2 OR D
*/
func (core *Core) OPB2() int {
	core.CPU.Registers.A = core.CPU.Registers.A | core.CPU.Registers.D
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xAE XOR (HL)
*/
func (core *Core) OPAE() int {
	core.CPU.Registers.A = core.ReadMemory(core.CPU.Registers.HL) ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xAD XOR L
*/
func (core *Core) OPAD() int {
	core.CPU.Registers.A = byte(core.CPU.Registers.HL&0xFF) ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xAC XOR H
*/
func (core *Core) OPAC() int {
	core.CPU.Registers.A = byte((core.CPU.Registers.HL&0xFF00)>>8) ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xAB XOR E
*/
func (core *Core) OPAB() int {
	core.CPU.Registers.A = core.CPU.Registers.E ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xAA XOR F
*/
func (core *Core) OPAA() int {
	core.CPU.Registers.A = core.CPU.Registers.D ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x9F SBC A,A
*/
func (core *Core) OP9F() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := core.CPU.Registers.A
	origin := val
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x9E SBC A,(HL)
*/
func (core *Core) OP9E() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := core.ReadMemory(core.CPU.Registers.HL)
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x9D SBC A,L
*/
func (core *Core) OP9D() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := byte(core.CPU.Registers.HL & 0xFF)
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x9C SBC A,H
*/
func (core *Core) OP9C() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x9B SBC A,E
*/
func (core *Core) OP9B() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := core.CPU.Registers.E
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x9A SBC A,D
*/
func (core *Core) OP9A() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := core.CPU.Registers.D
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x99 SBC A,C
*/
func (core *Core) OP99() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := core.CPU.Registers.C
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x98 SBC A,B
*/
func (core *Core) OP98() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	val := core.CPU.Registers.B
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val) - carry
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF)-carry < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x97 SUB A
*/
func (core *Core) OP97() int {
	val := core.CPU.Registers.A
	origin := val
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x95 SUB L
*/
func (core *Core) OP95() int {
	val := byte(core.CPU.Registers.HL & 0xFF)
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x94 SUB H
*/
func (core *Core) OP94() int {
	val := byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x93 SUB E
*/
func (core *Core) OP93() int {
	val := core.CPU.Registers.E
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x92 SUB D
*/
func (core *Core) OP92() int {
	val := core.CPU.Registers.D
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x91 SUB C
*/
func (core *Core) OP91() int {
	val := core.CPU.Registers.C
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x90 SUB B
*/
func (core *Core) OP90() int {
	val := core.CPU.Registers.B
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x0F RRCA
*/
func (core *Core) OP0F() int {
	value := core.CPU.Registers.A
	core.CPU.Registers.A = byte(value>>1) | byte((value&1)<<7)
	core.CPU.Flags.Zero = false
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = core.CPU.Registers.A > 0x7F
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x3F CCF
*/
func (core *Core) OP3F() int {
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = !core.CPU.Flags.Carry
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x8F ADC A,A
*/
func (core *Core) OP8F() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x8E ADC A,(HL)
*/
func (core *Core) OP8E() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.ReadMemory(core.CPU.Registers.HL)
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x8D ADC A,L
*/
func (core *Core) OP8D() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := byte(core.CPU.Registers.HL & 0xFF)
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x8C ADC A,H
*/
func (core *Core) OP8C() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x8B ADC A,E
*/
func (core *Core) OP8B() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.E
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x8A ADC A,F
*/
func (core *Core) OP8A() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.D
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x89 ADC A,C
*/
func (core *Core) OP89() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.C
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x88 ADC A,B
*/
func (core *Core) OP88() int {
	carry := int16(0)
	if core.CPU.Flags.Carry {
		carry = 1
	}
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.B
	res := int16(core.CPU.Registers.A) + int16(origin2) + carry

	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF)+byte(carry) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x86 ADD A,(HL)
*/
func (core *Core) OP86() int {
	origin1 := core.CPU.Registers.A
	origin2 := core.ReadMemory(core.CPU.Registers.HL)
	res := int16(core.CPU.Registers.A) + int16(origin2)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x84 ADD A,H
*/
func (core *Core) OP84() int {
	origin1 := core.CPU.Registers.A
	origin2 := byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	res := int16(core.CPU.Registers.A) + int16(origin2)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x83 ADD A,E
*/
func (core *Core) OP83() int {
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.E
	res := int16(core.CPU.Registers.A) + int16(origin2)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x81 ADD A,C
*/
func (core *Core) OP81() int {
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.C
	res := int16(core.CPU.Registers.A) + int16(origin2)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x76 HALT
*/
func (core *Core) OP76() int {
	core.CPU.Halt = true
	return 0
}

/*
	OP:0x75 LD (HL),L
*/
func (core *Core) OP75() int {
	core.WriteMemory(core.CPU.Registers.HL, byte(core.CPU.Registers.HL&0xFF))
	return 0
}

/*
	OP:0x74 LD (HL),H
*/
func (core *Core) OP74() int {
	core.WriteMemory(core.CPU.Registers.HL, byte((core.CPU.Registers.HL&0xFF00)>>8))
	return 0
}

/*
	OP:0x70 LD (HL),B
*/
func (core *Core) OP70() int {
	core.WriteMemory(core.CPU.Registers.HL, core.CPU.Registers.B)
	return 0
}

/*
	OP:0x6E LD L,(HL)
*/
func (core *Core) OP6E() int {
	core.setL(core.ReadMemory(core.CPU.Registers.HL))
	return 0
}

/*
	OP:0x6D LD L,L
*/
func (core *Core) OP6D() int {
	core.setL(byte(core.CPU.Registers.HL & 0xFF))
	return 0
}

/*
	OP:0x6C LD L,H
*/
func (core *Core) OP6C() int {
	core.setL(byte((core.CPU.Registers.HL & 0xFF00) >> 8))
	return 0
}

/*
	OP:0x6A LD L,D
*/
func (core *Core) OP6A() int {
	core.setL(core.CPU.Registers.D)
	return 0
}

/*
	OP:0x68 LD L,B
*/
func (core *Core) OP68() int {
	core.setL(core.CPU.Registers.B)
	return 0
}

/*
	OP:0x66 LD H,(HL)
*/
func (core *Core) OP66() int {
	core.setH(core.ReadMemory(core.CPU.Registers.HL))
	return 0
}

/*
	OP:0x65 LD H,L
*/
func (core *Core) OP65() int {
	core.setH(byte(core.CPU.Registers.HL & 0xFF))
	return 0
}

/*
	OP:0x64 LD H,H
*/
func (core *Core) OP64() int {
	core.setH(byte((core.CPU.Registers.HL & 0xFF00) >> 8))
	return 0
}

/*
	OP:0x63 LD H,e
*/
func (core *Core) OP63() int {
	core.setH(core.CPU.Registers.E)
	return 0
}

/*
	OP:0x61 LD H,C
*/
func (core *Core) OP61() int {
	core.setH(core.CPU.Registers.C)
	return 0
}

/*
	OP:0x5C LD E,H
*/
func (core *Core) OP5C() int {
	core.CPU.Registers.E = byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	return 0
}

/*
	OP:0x5B LD E,E
*/
func (core *Core) OP5B() int {
	core.CPU.Registers.E = core.CPU.Registers.E
	return 0
}

/*
	OP:0x5A LD E,D
*/
func (core *Core) OP5A() int {
	core.CPU.Registers.E = core.CPU.Registers.D
	return 0
}

/*
	OP:0x59 LD E,C
*/
func (core *Core) OP59() int {
	core.CPU.Registers.E = core.CPU.Registers.C
	return 0
}

/*
	OP:0x58 LD E,B
*/
func (core *Core) OP58() int {
	core.CPU.Registers.E = core.CPU.Registers.B
	return 0
}

/*
	OP:0x55 LD D,L
*/
func (core *Core) OP55() int {
	core.CPU.Registers.D = byte((core.CPU.Registers.HL & 0xFF))
	return 0
}

/*
	OP:0x53 LD D,E
*/
func (core *Core) OP53() int {
	core.CPU.Registers.D = core.CPU.Registers.E
	return 0
}

/*
	OP:0x52 LD D,D
*/
func (core *Core) OP52() int {
	core.CPU.Registers.D = core.CPU.Registers.D
	return 0
}

/*
	OP:0x51 LD D,C
*/
func (core *Core) OP51() int {
	core.CPU.Registers.D = core.CPU.Registers.C
	return 0
}

/*
	OP:0x50 LD D,B
*/
func (core *Core) OP50() int {
	core.CPU.Registers.D = core.CPU.Registers.B
	return 0
}

/*
	OP:0x4D LD C,L
*/
func (core *Core) OP4D() int {
	core.CPU.Registers.C = byte((core.CPU.Registers.HL & 0xFF))
	return 0
}

/*
	OP:0x4C LD C,H
*/
func (core *Core) OP4C() int {
	core.CPU.Registers.C = byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	return 0
}

/*
	OP:0x4B LD C,E
*/
func (core *Core) OP4B() int {
	core.CPU.Registers.C = core.CPU.Registers.E
	return 0
}

/*
	OP:0x4A LD C,D
*/
func (core *Core) OP4A() int {
	core.CPU.Registers.C = core.CPU.Registers.D
	return 0
}

/*
	OP:0x49 LD C,C
*/
func (core *Core) OP49() int {
	core.CPU.Registers.C = core.CPU.Registers.C
	return 0
}

/*
	OP:0x48 LD C,B
*/
func (core *Core) OP48() int {
	core.CPU.Registers.C = core.CPU.Registers.B
	return 0
}

/*
	OP:0x45 LD B,L
*/
func (core *Core) OP45() int {
	core.CPU.Registers.B = byte((core.CPU.Registers.HL & 0xFF))
	return 0
}

/*
	OP:0x44 LD B,H
*/
func (core *Core) OP44() int {
	core.CPU.Registers.B = byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	return 0
}

/*
	OP:0x43 LD B,E
*/
func (core *Core) OP43() int {
	core.CPU.Registers.B = core.CPU.Registers.E
	return 0
}

/*
	OP:0x42 LD B,D
*/
func (core *Core) OP42() int {
	core.CPU.Registers.B = core.CPU.Registers.D
	return 0
}

/*
	OP:0x41 LD B,C
*/
func (core *Core) OP41() int {
	core.CPU.Registers.B = core.CPU.Registers.C
	return 0
}

/*
	OP:0x3B DEC SP
*/
func (core *Core) OP3B() int {
	core.CPU.Registers.SP = core.CPU.Registers.SP - 1
	return 0
}

/*
	OP:0x39 ADD HL,SP
*/
func (core *Core) OP39() int {
	originHL := core.CPU.Registers.HL
	originSP := core.CPU.Registers.SP
	res := int32(originSP) + int32(originHL)
	core.CPU.Registers.HL = uint16(res)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = int32(originHL&0xFFF) > (res & 0xFFF)
	core.CPU.Flags.Carry = res > 0xFFFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x37 SCF
*/
func (core *Core) OP37() int {
	core.CPU.Flags.Carry = true
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x33 INC SP
*/
func (core *Core) OP33() int {
	core.CPU.Registers.SP = core.CPU.Registers.SP + 1
	return 0
}

/*
	OP:0x29 ADD HL,HL
*/
func (core *Core) OP29() int {
	originHL := core.CPU.Registers.HL
	res := int32(originHL) + int32(originHL)
	core.CPU.Registers.HL = uint16(res)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = int32(originHL&0xFFF) > (res & 0xFFF)
	core.CPU.Flags.Carry = res > 0xFFFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x27 DAA
*/
func (core *Core) OP27() int {
	if !core.CPU.Flags.Sub {

		if core.CPU.Flags.Carry || core.CPU.Registers.A > 0x99 {
			core.CPU.Registers.A = core.CPU.Registers.A + 0x60
			core.CPU.Flags.Carry = true
		}
		if core.CPU.Flags.HalfCarry || core.CPU.Registers.A&0xF > 0x9 {
			core.CPU.Registers.A = core.CPU.Registers.A + 0x06
			core.CPU.Flags.HalfCarry = false
		}
	} else if core.CPU.Flags.Carry && core.CPU.Flags.HalfCarry {
		core.CPU.Registers.A = core.CPU.Registers.A + 0x9A
		core.CPU.Flags.HalfCarry = false
	} else if core.CPU.Flags.Carry {
		core.CPU.Registers.A = core.CPU.Registers.A + 0xA0
	} else if core.CPU.Flags.HalfCarry {
		core.CPU.Registers.A += 0xFA
		core.CPU.Flags.HalfCarry = false
	}
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x25 DEC H
*/
func (core *Core) OP25() int {
	origin := byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	core.setH(origin - 1)
	core.CPU.Flags.Zero = ((origin - 1) == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x17 RLA
*/
func (core *Core) OP17() int {
	carryFlag := byte(0)
	if core.CPU.Flags.Carry {
		carryFlag = 1
	}
	core.CPU.Flags.Carry = ((core.CPU.Registers.A & 0x80) == 0x80)
	core.CPU.Registers.A = ((core.CPU.Registers.A << 1) & 0xFF) | carryFlag
	core.CPU.Flags.Zero = false
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x15 DEC D
*/
func (core *Core) OP15() int {
	origin := core.CPU.Registers.D
	core.CPU.Registers.D = origin - 1
	core.CPU.Flags.Zero = ((origin - 1) == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x10 STOP 0
*/
func (core *Core) OP10() int {
	//TODO STOP
	return 0
}

/*
	OP:0x08 LD (a16),SP
*/
func (core *Core) OP08() int {
	address := core.getParameter16()
	core.WriteMemory(address, byte(core.CPU.Registers.SP&0xff))
	core.WriteMemory(address+1, byte((core.CPU.Registers.SP&0xff00)>>8))
	return 0
}

/*
	OP:0x04 INC B
*/
func (core *Core) OP04() int {
	origin := core.CPU.Registers.B
	newVal := origin + 1
	core.CPU.Registers.B = newVal
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x02 LD (BC),A
*/
func (core *Core) OP02() int {
	core.WriteMemory(core.CPU.getBC(), core.CPU.Registers.A)
	return 0
}

/*
	OP:0x24 INC H
*/
func (core *Core) OP24() int {
	origin := byte((core.CPU.Registers.HL & 0xFF00) >> 8)
	newVal := origin + 1
	core.setH(newVal)
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xC4 CALL NZ,a16
*/
func (core *Core) OPC4() int {
	val := core.getParameter16()
	if !core.CPU.Flags.Zero {
		core.StackPush(core.CPU.Registers.PC)
		core.CPU.Registers.PC = val
		return 14
	}
	return 0
}

/*
	OP:0x14 INC D
*/
func (core *Core) OP14() int {
	origin := core.CPU.Registers.D
	newVal := origin + 1
	core.CPU.Registers.D = newVal
	core.CPU.Flags.Zero = (newVal == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = ((origin&0xF)+(1&0xF) > 0xF)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x30 JR NC,r8
*/
func (core *Core) OP30() int {
	address := int8(core.getParameter8())
	if !core.CPU.Flags.Carry {
		core.CPU.Registers.PC = uint16(int32(core.CPU.Registers.PC) + int32(address))
		return 4
	}
	return 0
}

/*
	OP:0x82 ADD A,D
*/
func (core *Core) OP82() int {
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.D
	res := int16(core.CPU.Registers.A) + int16(core.CPU.Registers.D)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x2B DEC HL
*/
func (core *Core) OP2B() int {
	core.CPU.Registers.HL = core.CPU.Registers.HL - 1
	return 0
}

/*
	OP:0x1B DEC DE
*/
func (core *Core) OP1B() int {
	core.CPU.setDE(core.CPU.getDE() - 1)
	return 0
}

/*
	OP:0x38 JR C,r8
*/
func (core *Core) OP38() int {
	address := int8(core.getParameter8())
	if core.CPU.Flags.Carry {
		core.CPU.Registers.PC = uint16(int32(core.CPU.Registers.PC) + int32(address))
		return 4
	}
	return 0
}

/*
	OP:0x96 SUB (HL)
*/
func (core *Core) OP96() int {
	val := core.ReadMemory(core.CPU.Registers.HL)
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)

	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xBE CP (HL)
*/
func (core *Core) OPBE() int {
	val := core.ReadMemory(core.CPU.Registers.HL)
	core.CPU.Compare(val, core.CPU.Registers.A)
	return 0
}

/*
	OP:0x1D DEC E
*/
func (core *Core) OP1D() int {
	origin := core.CPU.Registers.E
	core.CPU.Registers.E = origin - 1
	core.CPU.Flags.Zero = ((origin - 1) == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xD6 SUB d8
*/
func (core *Core) OPD6() int {
	val := core.getParameter8()
	origin := core.CPU.Registers.A
	res := int16(core.CPU.Registers.A) - int16(val)
	final := byte(res)
	core.CPU.Registers.A = final

	core.CPU.Flags.Zero = final == 0
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = int16(origin&0x0F)-int16(val&0xF) < 0
	core.CPU.Flags.Carry = res < 0
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xEE XOR d8
*/
func (core *Core) OPEE() int {
	val := core.getParameter8()
	core.CPU.Registers.A = core.CPU.Registers.A ^ val
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x2E LD L,d8
*/
func (core *Core) OP2E() int {
	core.setL(core.getParameter8())
	return 0
}

/*
	OP:0xA8 XOR B
*/
func (core *Core) OPA8() int {
	core.CPU.Registers.A = core.CPU.Registers.B ^ core.CPU.Registers.A
	core.CPU.Flags.Zero = core.CPU.Registers.A == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0xB8 CP B
*/
func (core *Core) OPB8() int {
	core.CPU.Compare(core.CPU.Registers.B, core.CPU.Registers.A)
	return 0
}

/*
	OP:0x80 ADD A,B
*/
func (core *Core) OP80() int {
	origin1 := core.CPU.Registers.A
	origin2 := core.CPU.Registers.B
	res := int16(core.CPU.Registers.A) + int16(origin2)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x07 RLCA
*/
func (core *Core) OP07() int {
	origin := core.CPU.Registers.A

	core.CPU.Registers.A = byte(core.CPU.Registers.A<<1) | (core.CPU.Registers.A >> 7)
	core.CPU.Flags.Zero = false
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = origin > 0x7F
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x26 LD H,d8
*/
func (core *Core) OP26() int {
	core.setH(core.getParameter8())
	return 0
}

/*
	OP:0x1E LD E,d8
*/
func (core *Core) OP1E() int {
	core.CPU.Registers.E = core.getParameter8()
	return 0
}

/*
	OP:0x40 LD B,B
*/
func (core *Core) OP40() int {
	core.CPU.Registers.B = core.CPU.Registers.B
	return 0
}

/*
	OP:0x62 LD H,D
*/
func (core *Core) OP62() int {
	core.setH(core.CPU.Registers.D)
	return 0
}

/*
	OP:0x6B LD L,E
*/
func (core *Core) OP6B() int {
	core.setL(core.CPU.Registers.E)
	return 0
}

/*
	OP:0xF6 OR d8
*/
func (core *Core) OPF6() int {
	val := core.getParameter8()
	core.CPU.Registers.A = core.CPU.Registers.A | val
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x54 LD D,H
*/
func (core *Core) OP54() int {
	core.CPU.Registers.D = byte(core.CPU.Registers.HL >> 8)
	return 0
}

/*
	OP:0x5D LD E,L
*/
func (core *Core) OP5D() int {
	core.CPU.Registers.E = byte(core.CPU.Registers.HL & 0xFF)
	return 0
}

/*
	OP:0xC6 ADD A,d8
*/
func (core *Core) OPC6() int {
	origin1 := core.CPU.Registers.A
	origin2 := core.getParameter8()
	res := int16(core.CPU.Registers.A) + int16(origin2)
	core.CPU.Registers.A = byte(res)
	core.CPU.Flags.Zero = byte(res) == 0
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = (origin1&0xF)+(origin2&0xF) > 0xF
	core.CPU.Flags.Carry = res > 0xFF
	core.CPU.updateAFLow()
	return 0
}

/*
	OP:0x7D LD A,L
*/
func (core *Core) OP7D() int {
	core.CPU.Registers.A = byte(core.CPU.Registers.HL & 0xFF)
	return 0
}

/*
	OP:0x67 LD H,A
*/
func (core *Core) OP67() int {
	core.setH(core.CPU.Registers.A)
	return 0
}

/*
	OP:0x2D DEC L
*/
func (core *Core) OP2D() int {
	origin := byte(core.CPU.Registers.HL & 0xFF)
	core.setL(origin - 1)
	core.CPU.Flags.Zero = ((origin - 1) == 0)
	core.CPU.Flags.Sub = true
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
	core.CPU.updateAFLow()
	return 0
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
	core.CPU.Flags.Zero = byte(res) == 0
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
	core.CPU.Flags.HalfCarry = (origin&0x0F == 0)
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
	OP:0x7F LD A,A
*/
func (core *Core) OP7F() int {
	core.CPU.Registers.A = core.CPU.Registers.A
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
	core.CPU.Registers.HL = 1 + core.CPU.Registers.HL
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
	core.CPU.Flags.Zero = byte(res) == 0
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
	core.CPU.Registers.A = core.CPU.Registers.A | core.CPU.Registers.B
	core.CPU.Flags.Zero = (core.CPU.Registers.A == 0)
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
	if core.cbMap[nextIns] != nil {
		core.cbMap[nextIns]()
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

	core.CPU.setAF(core.StackPop() & 0xFFF0)

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
	core.CPU.updateAFLow()
	core.StackPush(core.CPU.getAF())
	return 0
}

/*
	OP:0xFB EI
*/
func (core *Core) OPFB() int {
	core.CPU.Flags.PendingInterruptEnabled = true
	//log.Println(core.CPU.Registers.PC)
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
	core.CPU.Flags.InterruptMaster = false
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
*/
func (core *Core) OP1F() int {
	value := core.CPU.Registers.A
	var carry byte
	if core.CPU.Flags.Carry {
		carry = 0x80
	}
	result := byte(value>>1) | carry
	core.CPU.Registers.A = result
	core.CPU.Flags.Zero = false
	core.CPU.Flags.Carry = (1 & value) == 1
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Sub = false
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
