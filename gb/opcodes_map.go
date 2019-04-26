package gb

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
	byte(0x07): {
		Func:  (*Core).OP07,
		Clock: 4,
		OP:    "RLCA",
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
	byte(0x1E): {
		Func:  (*Core).OP1E,
		Clock: 8,
		OP:    "LD E,d8",
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
	byte(0x26): {
		Func:  (*Core).OP26,
		Clock: 8,
		OP:    "LD H,d8",
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
	byte(0x2D): {
		Func:  (*Core).OP2D,
		Clock: 4,
		OP:    "DEC L",
	},
	byte(0x2E): {
		Func:  (*Core).OP2E,
		Clock: 8,
		OP:    "LD L,d8",
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
	byte(0x40): {
		Func:  (*Core).OP40,
		Clock: 4,
		OP:    "LD B,B",
	},
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
	byte(0x54): {
		Func:  (*Core).OP54,
		Clock: 4,
		OP:    "LD D,H",
	},
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
	byte(0x5D): {
		Func:  (*Core).OP5D,
		Clock: 4,
		OP:    "LD E,L",
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
	byte(0x62): {
		Func:  (*Core).OP62,
		Clock: 4,
		OP:    "LD H,D",
	},
	byte(0x67): {
		Func:  (*Core).OP67,
		Clock: 4,
		OP:    "LD H,A",
	},
	byte(0x69): {
		Func:  (*Core).OP69,
		Clock: 4,
		OP:    "LD L,C",
	},
	byte(0x6B): {
		Func:  (*Core).OP6B,
		Clock: 4,
		OP:    "LD L,E",
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
	byte(0x7D): {
		Func:  (*Core).OP7D,
		Clock: 4,
		OP:    "LD A,L",
	},
	byte(0x7E): {
		Func:  (*Core).OP7E,
		Clock: 8,
		OP:    "LD A,(HL)",
	},
	//0x8_
	byte(0x80): {
		Func:  (*Core).OP80,
		Clock: 4,
		OP:    "ADD A,B",
	},
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
	byte(0x89): {
		Func:  (*Core).OP89,
		Clock: 4,
		OP:    "ADC A,C",
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
	byte(0xA8): {
		Func:  (*Core).OPA8,
		Clock: 4,
		OP:    "XOR B",
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
	byte(0xB8): {
		Func:  (*Core).OPB8,
		Clock: 4,
		OP:    "CP B",
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
	byte(0xC6): {
		Func:  (*Core).OPC6,
		Clock: 8,
		OP:    "ADD A,d8",
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
	byte(0xF6): {
		Func:  (*Core).OPF6,
		Clock: 8,
		OP:    "OR d8",
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
