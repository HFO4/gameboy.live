package gb

import "github.com/HFO4/gbc-in-cloud/util"

var cbMap = [0x100](func()){}

var CBCycles = []int{
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 0
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 1
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 2
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 3
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2, // 4
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2, // 5
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2, // 6
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2, // 7
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 8
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // 9
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // A
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // B
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // C
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // D
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // E
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2, // F
} //0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f

func (core *Core) initCB() {

	var getters = [8]func() byte{
		func() byte { return core.CPU.Registers.B },
		func() byte { return core.CPU.Registers.C },
		func() byte { return core.CPU.Registers.D },
		func() byte { return core.CPU.Registers.E },
		func() byte { return byte(core.CPU.Registers.HL >> 8) },
		func() byte { return byte(core.CPU.Registers.HL & 0x00FF) },
		func() byte { return core.ReadMemory(core.CPU.Registers.HL) },
		func() byte { return core.CPU.Registers.A },
	}

	var setters = [8]func(byte){
		core.setB,
		core.setC,
		core.setD,
		core.setE,
		core.setH,
		core.setL,
		func(val byte) { core.WriteMemory(core.CPU.Registers.HL, val) },
		core.setA,
	}

	//Every 8 instructions is a group,which use different registers
	for i := 0; i < 8; i++ {

		registerID := i

		cbMap[0x00+i] = func() {
			core.RLC(getters[registerID], setters[registerID])
		}

		cbMap[0x08+i] = func() {
			core.RRC(getters[registerID], setters[registerID])
		}

		cbMap[0x20+i] = func() {
			core.SLA(getters[registerID], setters[registerID])
		}

		cbMap[0x30+i] = func() {
			core.SWAP(getters[registerID], setters[registerID])
		}

		/*
			RES commands
		*/
		cbMap[0x80+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 0))
		}
		cbMap[0x88+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 1))
		}
		cbMap[0x90+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 2))
		}
		cbMap[0x98+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 3))
		}
		cbMap[0xA0+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 4))
		}
		cbMap[0xA8+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 5))
		}
		cbMap[0xB0+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 6))
		}
		cbMap[0xB8+i] = func() {
			setters[registerID](util.ClearBit(getters[registerID](), 7))
		}

	}
}

func (core *Core) SLA(getter func() byte, setter func(byte)) {
	val := getter()
	carry := val >> 7
	res := (val << 1) & 0xFF
	setter(res)

	core.CPU.Flags.Zero = (res == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = (carry == 1)
	core.CPU.updateAFLow()
}

func (core *Core) RLC(getter func() byte, setter func(byte)) {
	val := getter()
	var carry byte
	var rot byte
	carry = val >> 7
	rot = (val<<1)&0xFF | carry
	setter(rot)
	core.CPU.Flags.Zero = (rot == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = (carry == 1)
	core.CPU.updateAFLow()
}

func (core *Core) SWAP(getter func() byte, setter func(byte)) {
	val := getter()
	res := val<<4&240 | val>>4
	setter(res)
	core.CPU.Flags.Zero = (res == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = false
	core.CPU.updateAFLow()
}

func (core *Core) RRC(getter func() byte, setter func(byte)) {
	val := getter()
	carry := val & 1
	rot := (val >> 1) | (carry << 7)
	setter(rot)
	core.CPU.Flags.Zero = (rot == 0)
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = false
	core.CPU.Flags.Carry = (carry == 1)
	core.CPU.updateAFLow()
}
