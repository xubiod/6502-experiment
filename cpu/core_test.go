package cpu

import (
	"strings"
	"testing"
	"xubiod/6502-experiment/assembler"
)

func TestExists(t *testing.T) {
	voids := map[byte]uint8{
		0x02: 0, 0x03: 0, 0x04: 0, 0x07: 0, 0x0B: 0, 0x0C: 0, 0x0F: 0,
		0x12: 0, 0x13: 0, 0x14: 0, 0x17: 0, 0x1A: 0, 0x1B: 0, 0x1C: 0, 0x1F: 0,
		0x22: 0, 0x23: 0, 0x27: 0, 0x2B: 0, 0x2F: 0,
		0x32: 0, 0x33: 0, 0x34: 0, 0x37: 0, 0x3A: 0, 0x3B: 0, 0x3C: 0, 0x3F: 0,
		0x42: 0, 0x43: 0, 0x44: 0, 0x47: 0, 0x4B: 0, 0x4F: 0,
		0x52: 0, 0x53: 0, 0x54: 0, 0x57: 0, 0x5A: 0, 0x5B: 0, 0x5C: 0, 0x5F: 0,
		0x62: 0, 0x63: 0, 0x64: 0, 0x67: 0, 0x6B: 0, 0x6F: 0,
		0x72: 0, 0x73: 0, 0x74: 0, 0x77: 0, 0x7A: 0, 0x7B: 0, 0x7C: 0, 0x7F: 0,
		0x80: 0, 0x82: 0, 0x83: 0, 0x87: 0, 0x89: 0, 0x8B: 0, 0x8F: 0,
		0x92: 0, 0x93: 0, 0x97: 0, 0x9B: 0, 0x9C: 0, 0x9E: 0, 0x9F: 0,
		0xA3: 0, 0xA7: 0, 0xAB: 0, 0xAF: 0,
		0xB2: 0, 0xB3: 0, 0xB7: 0, 0xBB: 0, 0xBF: 0,
		0xC2: 0, 0xC3: 0, 0xC7: 0, 0xCB: 0, 0xCF: 0,
		0xD2: 0, 0xD3: 0, 0xD4: 0, 0xD7: 0, 0xDA: 0, 0xDB: 0, 0xDC: 0, 0xDF: 0,
		0xE2: 0, 0xE3: 0, 0xE7: 0, 0xEB: 0, 0xEF: 0,
		0xF2: 0, 0xF3: 0, 0xF4: 0, 0xF7: 0, 0xFA: 0, 0xFB: 0, 0xFC: 0, 0xFF: 0,
	}

	c := NewCore()

	var i uint16
	var ok bool
	var finalOk bool

	failCount := 0

	for i < 256 {
		finalOk = false

		_, ok = voids[byte(i)]
		finalOk = finalOk || ok

		_, ok = c.execMapNil[byte(i)]
		finalOk = finalOk || ok

		_, ok = c.execMapByte[byte(i)]
		finalOk = finalOk || ok

		_, ok = c.execMapShort[byte(i)]
		finalOk = finalOk || ok

		if !finalOk {
			t.Errorf("opcode %0X: should exist but doesn't!!", i)
			failCount++
		}

		i++
	}

	t.Logf("\n%d/%d exist", 256-failCount, 256)
}

func TestResetRoutine(t *testing.T) {
	c := NewCore()

	resetRoutine := []byte{
		0xa2, 0xff, // LDX __Imm(0xFF)
		0x9a, // TXS ____i()

		0x78, // SEI ____i()
		0x18, // CLC ____i()
		0xd8, // CLD ____i()
		0x58, // CLI ____i()
		0xb8, // CLV ____i()

		0xa9, 0x00, // LDA __Imm(0x00)
		0xa2, 0x00, // LDX __Imm(0x00)
		0xa0, 0x00, // LDY __Imm(0x00)
	}

	c.Write(resetRoutine)
	c.PC = 0x0200

	// purposefully poison registers

	c.A = 0xDE
	c.X = 0xAD
	c.Y = 0x24

	executing := true

	for c.Memory[c.PC] != 0x00 && executing {
		executing, _, _ = c.StepOnce()
	}

	if c.A != 0 {
		t.Errorf("accumulator not 0, was \"%1X\"", c.A)
	}

	if c.X != 0 {
		t.Errorf("X not 0, was \"%1X\"", c.X)
	}

	if c.Y != 0 {
		t.Errorf("Y not 0, was \"%1X\"", c.Y)
	}

	if c.S != 0xFF {
		t.Errorf("stack pointer not FF, was \"%1X\"", c.S)
	}

	if c.PC != 0x020E {
		t.Errorf("program counter not 20E, was \"%1X\"", c.PC)
	}
	t.Log("\n" + c.CompleteDump())
}

func TestArithmeticADC(t *testing.T) {
	c := NewCore()

	// numbers to test adding
	lefts_t := []byte{0x50, 0x50, 0x50, 0x50, 0xD0, 0xD0, 0xD0, 0xD0}
	right_t := []byte{0x10, 0x50, 0x90, 0xD0, 0x10, 0x50, 0x90, 0xD0}

	// should carry be set before adding?
	carry_set_t := []byte{0xEA, 0x38, 0xEA, 0x38, 0xEA, 0x38, 0xEA, 0x38} // 0xEA is a NOP, 0x38 is a SEC

	// expected carry and overflow flags
	carry_out_t := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x01} // 0x01 means set
	overflow__t := []byte{0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00} // 0x40 means set

	for idx, left := range lefts_t {
		right := right_t[idx]

		prg := []byte{
			carry_set_t[idx], // either a NOP or SEC

			0xa9, left,
			0x8d, 0x00, 0x00,
			0xa9, right,

			0x6d, 0x00, 0x00,

			0x00,
		}

		stdProcedure(c, prg)

		var carryOnAdd byte = 0
		if carry_set_t[idx] == 0x38 {
			carryOnAdd += 1
		}

		if c.A != left+right+carryOnAdd {
			t.Errorf("adc fail - (%1x + %1x) - answer\texpected %1x\tgot %1x", left, right, left+right+carryOnAdd, c.A)
		}

		if c.Flags&FLAG_CARRY != carry_out_t[idx] {
			t.Errorf("adc fail - (%1x + %1x) - carry\texpected %1x\tgot %1x", left, right, carry_out_t[idx], c.Flags&FLAG_CARRY)
		}

		if c.Flags&FLAG_OVERFLOW != overflow__t[idx] {
			t.Errorf("adc fail - (%1x + %1x) - V flag\texpected %1x\tgot %1x", left, right, overflow__t[idx], c.Flags&FLAG_OVERFLOW)
		}
		t.Log("\n" + c.CompleteDump())
	}
}

func TestDecimalADC(t *testing.T) {
	c := NewCore()

	// numbers to test adding
	lefts_t := []byte{0x50, 0x50, 0x50, 0x50, 0xD0, 0xD0, 0xD0, 0xD0}
	right_t := []byte{0x10, 0x50, 0x90, 0xD0, 0x10, 0x50, 0x90, 0xD0}

	// should carry be set before adding?
	carry_set_t := []byte{0xEA, 0xEA, 0xEA, 0xEA, 0xEA, 0xEA, 0xEA, 0xEA} // 0xEA is a NOP, 0x38 is a SEC

	// results
	results_t := []byte{0x60, 0x00, 0x40, 0x80, 0x40, 0x80, 0xC0, 0x00} // verified with visual6502.org

	// expected flags
	carry_out_t := []byte{0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01} // 0x01 means set
	overflow__t := []byte{0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00} // 0x40 means set
	negative__t := []byte{0x00, 0x80, 0x80, 0x00, 0x80, 0x00, 0x00, 0x80} // 0x80 means set

	for idx, left := range lefts_t {
		right := right_t[idx]

		prg := []byte{
			0xf8,
			carry_set_t[idx], // either a NOP or SEC

			0xa9, left,
			0x8d, 0x00, 0x00,
			0xa9, right,

			0x6d, 0x00, 0x00,

			0x00,
		}

		stdProcedure(c, prg)

		if c.A != results_t[idx] {
			t.Errorf("adc decimal fail - (%1x + %1x) - answer\texpected %1x\tgot %1x", left, right, results_t[idx], c.A)
		}

		if c.Flags&FLAG_CARRY != carry_out_t[idx] {
			t.Errorf("adc decimal fail - (%1x + %1x) - carry\texpected %1x\tgot %1x", left, right, carry_out_t[idx], c.Flags&FLAG_CARRY)
		}

		if c.Flags&FLAG_OVERFLOW != overflow__t[idx] {
			t.Errorf("adc decimal fail - (%1x + %1x) - V flag\texpected %1x\tgot %1x", left, right, overflow__t[idx], c.Flags&FLAG_OVERFLOW)
		}

		if c.Flags&FLAG_NEGATIVE != negative__t[idx] {
			t.Errorf("adc decimal fail - (%1x + %1x) - N flag\texpected %1x\tgot %1x", left, right, negative__t[idx], c.Flags&FLAG_NEGATIVE)
		}
		t.Log("\n" + c.CompleteDump())
	}
}

func TestArithmeticSBC(t *testing.T) {
	c := NewCore()

	// numbers to test subtracting
	right_t := []byte{0x50, 0x50, 0x50, 0x50, 0xD0, 0xD0, 0xD0, 0xD0}
	lefts_t := []byte{0xF0, 0xB0, 0x70, 0x30, 0xF0, 0xB0, 0x70, 0x30}

	// should carry be set before subtracting?
	carry_set_t := []byte{0xEA, 0x38, 0xEA, 0x38, 0xEA, 0x38, 0xEA, 0x38} // 0xEA is a NOP, 0x38 is a SEC

	// expected carry and overflow flags
	carry_out_t := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x01} // 0x01 means set
	overflow__t := []byte{0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00} // 0x40 means set

	for idx, left := range lefts_t {
		right := right_t[idx]

		prg := []byte{
			carry_set_t[idx], // either a NOP or SEC

			0xa9, left,
			0x8d, 0x00, 0x00,
			0xa9, right,

			0xed, 0x00, 0x00,

			0x00,
		}

		stdProcedure(c, prg)

		var carryOnSub byte = 0
		if carry_set_t[idx] == 0x38 {
			carryOnSub += 1
		}

		var pre = (uint16(right) - uint16(left)) + uint16(carryOnSub)
		var r = byte(pre & 0xFF)

		if c.A != r {
			t.Errorf("sbc fail - (%1x - %1x) - answer\texpected %1x\tgot %1x", right, left, r, c.A)
		}

		if c.Flags&FLAG_CARRY != carry_out_t[idx] {
			t.Errorf("sbc fail - (%1x - %1x) - carry\texpected %1x\tgot %1x", right, left, carry_out_t[idx], c.Flags&FLAG_CARRY)
		}

		if c.Flags&FLAG_OVERFLOW != overflow__t[idx] {
			t.Errorf("sbc fail - (%1x - %1x) - V flag\texpected %1x\tgot %1x", right, left, overflow__t[idx], c.Flags&FLAG_OVERFLOW)
		}
		t.Log("\n" + c.CompleteDump())
	}
}

func TestDecimalSBC(t *testing.T) {
	c := NewCore()

	// numbers to test subtracting
	lefts_t := []byte{0xF0, 0xB0, 0x70, 0x30, 0xF0, 0xB0, 0x70, 0x30}
	right_t := []byte{0x50, 0x50, 0x50, 0x50, 0xD0, 0xD0, 0xD0, 0xD0}

	// should carry be set before adding?
	carry_set_t := []byte{0x38, 0x38, 0x38, 0x38, 0x38, 0x38, 0x38, 0x38} // 0xEA is a NOP, 0x38 is a SEC

	// results
	results_t := []byte{0x00, 0x40, 0x80, 0x20, 0x80, 0x20, 0x60, 0xa0} // verified with visual6502.org

	// expected flags
	carry_out_t := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x01} // 0x01 means set
	overflow__t := []byte{0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00} // 0x40 means set
	negative__t := []byte{0x00, 0x80, 0x80, 0x00, 0x80, 0x00, 0x00, 0x80} // 0x80 means set

	for idx, left := range lefts_t {
		right := right_t[idx]

		prg := []byte{
			0xf8,
			carry_set_t[idx], // either a NOP or SEC

			0xa9, left,
			0x8d, 0x00, 0x00,
			0xa9, right,

			0xed, 0x00, 0x00,

			0x00,
		}

		stdProcedure(c, prg)

		if c.A != results_t[idx] {
			t.Errorf("sbc decimal fail - (M: %1x; A: %1x) - answer\texpected %1x\tgot %1x", left, right, results_t[idx], c.A)
		}

		if c.Flags&FLAG_CARRY != carry_out_t[idx] {
			t.Errorf("sbc decimal fail - (M: %1x; A: %1x) - carry\texpected %1x\tgot %1x", left, right, carry_out_t[idx], c.Flags&FLAG_CARRY)
		}

		if c.Flags&FLAG_OVERFLOW != overflow__t[idx] {
			t.Errorf("sbc decimal fail - (M: %1x; A: %1x) - V flag\texpected %1x\tgot %1x", left, right, overflow__t[idx], c.Flags&FLAG_OVERFLOW)
		}

		if c.Flags&FLAG_NEGATIVE != negative__t[idx] {
			t.Errorf("sbc decimal fail - (M: %1x; A: %1x) - N flag\texpected %1x\tgot %1x", left, right, negative__t[idx], c.Flags&FLAG_NEGATIVE)
		}
		t.Log("\n" + c.CompleteDump())
	}
}

func TestGeneralStackOps(t *testing.T) {
	var prg []byte
	loopCount := 32

	c := NewCore()
	asm := assembler.New()

	prg, _ = asm.Parse("LDX #$00\n" + strings.Repeat("TXA\nPHA\nINX\n", loopCount))

	stdProcedure(c, prg)

	var atStack byte

	for number := range loopCount {
		atStack = c.Memory[0x01FF-uint16(number)]
		if atStack != byte(number) {
			t.Errorf("general stack fail - expected %02x\tgot %02x", number, atStack)
		}
	}

	t.Log("\n" + c.CompleteDump())
}

// Writes reset procedure followed by the given program. Goes into a standard execution
// loop afterwards; breaking on `BRK` (`0x00`) or an invalid instruction. Does no checks
// itself, do that after calling this.
func stdProcedure(c *Core, program []byte) {
	c.prepare()

	resetRoutine := []byte{
		0xa2, 0xff, // LDX __Imm(0xFF)
		0x9a, // TXS ____i()

		0x78, // SEI ____i()
		0x18, // CLC ____i()
		0xd8, // CLD ____i()
		0x58, // CLI ____i()
		0xb8, // CLV ____i()

		0xa9, 0x00, // LDA __Imm(0x00)
		0xa2, 0x00, // LDX __Imm(0x00)
		0xa0, 0x00, // LDY __Imm(0x00)
	}

	c.Write(resetRoutine)
	c.Write(program)

	c.PC = 0x0200
	var exe bool = true

	for c.Memory[c.PC] != 0x00 && exe {
		exe, _, _ = c.StepOnce()
	}
}
