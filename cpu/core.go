package cpu

import (
	"errors"
)

type Core struct {
	// A 65,536 byte array to fully represent the memory of a 6502.
	//
	// `0x0000`-`0x00FF` is the zero page; `0x0100`-`0x01FF` is the stack; `0x0200`-`0xFFFF`
	// is the general memory of the chip.
	//
	// There is no special structure for the stack, it is **entirely managed
	// manually.**
	Memory [0x10000]byte

	A     byte   // A - accumulator
	X     byte   // X
	Y     byte   // Y
	PC    uint16 // PC - program counter
	S     uint8  // S - stack pointer; starts at `0x01FF` and grows down to `0x0100`
	Flags byte   // P - status, flags

	// Some instructions have different behaviours depending on what 6502-compatible
	// CPU they were based on, a quick example being the NES CPU not implementing
	// decimal mode functionality but keeping the flag itself.
	//
	// This struct has the options that can be changed to act more like a specific
	// CPU instead of a generic 6502.
	Features CoreFeatureFlags

	// The byte -> implementation map for instructions with no operands.
	execMapNil map[byte]func()

	// The byte -> implementation map for instructions with a byte as an operand.
	// Signed or not is instruction dependent.
	execMapByte map[byte]func(uint8)

	// The byte -> implementation map for instructions with an unsigned short (2
	// bytes) as an operand.
	execMapShort map[byte]func(uint16)

	writingPointer uint16 // The pointer to writing to memory with `*Core.Write()`.
}

// A struct for a set of feature flags that can be changed to have the emulator
// "specialized" to a specific 6502-compatible CPU instead of a generic 6502.
type CoreFeatureFlags struct {
	// Is decimal mode implemented? Toggling this off does not change the behaviour
	// of enabling/disabling the flag itself, but if off ADC/SBC will ignore the flag.
	//
	// Setting this to false will act like a NES CPU.
	DecimalModeImplemented bool
}

const (
	FLAG_CARRY             byte = 1 << iota // C - Set when the last operation resulted in an overflow.
	FLAG_ZERO                               // Z - Set when the last operation resulted in a zero.
	FLAG_INTERRUPT_DISABLE                  // I - When set, interrupts are disabled.
	FLAG_DECIMAL                            // D - When set, math operations are done with BCD.
	FLAG_BREAK                              // B - Set when a software interrupt happens with `BRK`.
	FLAG_UNUSED                             // _ - This flag is not used by the 6502.
	FLAG_OVERFLOW                           // V - Set when the last operation resulted in a *signed overflow*.
	FLAG_NEGATIVE                           // N - Set when the last operation resulted as a negative number.
)

// Does the calculations for a zero-page indirect indexed with Y address to get
// the valid address.
func (c *Core) indirectZpY(zp byte) (addr uint16) {
	var lsb, msb byte
	lsb = c.Memory[zp]
	msb = c.Memory[zp+1]

	addr = (uint16(msb) << 8 & uint16(lsb)) + uint16(c.Y)
	return
}

// Does the calculations for a zero-page indexed indirect to get the address.
func (c *Core) indirectZpX(zp byte) (addr uint16) {
	var lsb, msb byte
	lsb = c.Memory[zp+c.X]
	msb = c.Memory[zp+c.X+1]

	addr = uint16(msb) << 8 & uint16(lsb)
	return
}

// Creates and prepares a *Core.
func NewCore() (c *Core) {
	c = &Core{Features: CoreFeatureFlags{DecimalModeImplemented: true}}
	c.prepare()
	return
}

// Creates the decoding tables. Must be called before any execution unless writing
// your own execution loop.
func (c *Core) prepare() {
	c.execMapNil = map[byte]func(){
		0x00: c.BRK____i, 0x08: c.PHP____i, 0x0A: c.ASL____A,
		0x18: c.CLC____i,
		0x28: c.PLP____i, 0x2A: c.ROL____A,
		0x38: c.SEC____i,
		0x40: c.RTI____i, 0x48: c.PHA____i, 0x4A: c.LSR____A,
		0x58: c.CLI____i,
		0x60: c.RTS____i, 0x68: c.PLA____i, 0x6A: c.ROR____A,
		0x78: c.SEI____i,
		0x88: c.DEY____i, 0x8A: c.TXA____i,
		0x98: c.TYA____i, 0x9A: c.TXS____i,
		0xA8: c.TAY____i, 0xAA: c.TAX____i,
		0xB8: c.CLV____i, 0xBA: c.TSX____i,
		0xC8: c.INY____i, 0xCA: c.DEX____i,
		0xD8: c.CLD____i,
		0xE8: c.INX____i, 0xEA: c.NOP____i,
		0xF8: c.SED____i,
	}

	c.execMapByte = map[byte]func(uint8){
		0x01: c.ORA_IZPx, 0x05: c.ORA__ZPg, 0x06: c.ASL__ZPg, 0x09: c.ORA__Imm,
		0x10: c.BPL__rel, 0x11: c.ORA_IZPy, 0x15: c.ORA__ZPx, 0x16: c.ASL__ZPx,
		0x21: c.AND_IZPx, 0x24: c.BIT__ZPg, 0x25: c.AND__ZPg, 0x26: c.ROL__ZPg, 0x29: c.AND__Imm,
		0x30: c.BMI__rel, 0x31: c.AND_IZPy, 0x35: c.AND__ZPx, 0x36: c.ROL__ZPx,
		0x41: c.EOR_IZPx, 0x45: c.EOR__ZPg, 0x46: c.LSR__ZPg, 0x49: c.EOR__Imm,
		0x50: c.BVC__rel, 0x51: c.EOR_IZPy, 0x55: c.EOR__ZPx, 0x56: c.LSR__ZPx,
		0x61: c.ADC_IZPx, 0x65: c.ADC__ZPg, 0x66: c.ROR__ZPg, 0x69: c.ADC__Imm,
		0x70: c.BVS__rel, 0x71: c.ADC_IZPy, 0x75: c.ADC__ZPx, 0x76: c.ROR__ZPx,
		0x81: c.STA_IZPx, 0x84: c.STY__ZPg, 0x85: c.STA__ZPg, 0x86: c.STX__ZPg,
		0x90: c.BCC__rel, 0x91: c.STA_IZPy, 0x94: c.STA__ZPx, 0x95: c.STA__ZPx, 0x96: c.STX__ZPy,
		0xA0: c.LDY__Imm, 0xA1: c.LDA_IZPx, 0xA2: c.LDX__Imm, 0xA4: c.LDY__ZPg, 0xA5: c.LDA__ZPg, 0xA6: c.LDX__ZPg, 0xA9: c.LDA__Imm,
		0xB0: c.BCS__rel, 0xB1: c.LDA_IZPy, 0xB4: c.LDY__ZPx, 0xB5: c.LDA__ZPx, 0xB6: c.LDX__ZPy,
		0xC0: c.CPY__Imm, 0xC1: c.CMP_IZPx, 0xC4: c.CPY__ZPg, 0xC5: c.CMP__ZPg, 0xC6: c.DEC__ZPg, 0xC9: c.CMP__Imm,
		0xD0: c.BNE__rel, 0xD1: c.CMP_IZPy, 0xD5: c.CMP__ZPx, 0xD6: c.DEC__ZPx,
		0xE0: c.CPX__Imm, 0xE1: c.SBC_IZPx, 0xE4: c.CPX__ZPg, 0xE5: c.SBC__Zpg, 0xE6: c.INC__ZPg, 0xE9: c.SBC__Imm,
		0xF0: c.BEQ__rel, 0xF1: c.SBC_IZPy, 0xF5: c.SBC__ZPx, 0xF6: c.INC__ZPx,
	}

	c.execMapShort = map[byte]func(uint16){
		0x0D: c.ORA____a, 0x0E: c.ASL____a,
		0x19: c.ORA___ay, 0x1D: c.ORA___ax, 0x1E: c.ASL___ax,
		0x20: c.JSR____a, 0x2C: c.BIT____a, 0x2D: c.AND____a, 0x2E: c.ROL____a,
		0x39: c.AND___ay, 0x3D: c.AND___ax, 0x3E: c.ROL___ax,
		0x4C: c.JMP____a, 0x4D: c.EOR____a, 0x4E: c.LSR____a,
		0x59: c.EOR___ay, 0x5D: c.EOR___ax, 0x5E: c.LSR___ax,
		0x6C: c.JMP___Ia, 0x6D: c.ADC____a, 0x6E: c.ROR____a,
		0x79: c.ADC___ay, 0x7D: c.ADC___ax, 0x7E: c.ROR___ax,
		0x8C: c.STY____a, 0x8D: c.STA____a, 0x8E: c.STX____a,
		0x99: c.STA___ay, 0x9D: c.STA___ax,
		0xAC: c.LDY____a, 0xAD: c.LDA____a, 0xAE: c.LDX____a,
		0xB9: c.LDA___ay, 0xBC: c.LDY___ax, 0xBD: c.LDA___ax, 0xBE: c.LDX___ay,
		0xCC: c.CPY____a, 0xCD: c.CMP____a, 0xCE: c.DEC____a,
		0xD9: c.CMP___ay, 0xDD: c.CMP___ax, 0xDE: c.DEC___ax,
		0xEC: c.CPX____a, 0xED: c.SBC____a, 0xEE: c.INC____a,
		0xF9: c.SBC___ay, 0xFD: c.SBC___ax, 0xFE: c.INC___ax,
	}

	_ = c.SetWriterPtr(0x0200)
}

// Does a single step of execution. If at an invalid instruction, the program
// counter will not increment.
//
// Returns true if execution was successful, false if not. This can be used to
// easily make an execution loop depending on the results of this method.
func (c *Core) StepOnce() (valid bool) {
	inst := c.Memory[c.PC]
	valid = true

	f, fOk := c.execMapByte[inst]
	g, gOk := c.execMapShort[inst]
	h, hOk := c.execMapNil[inst]

	switch {
	case fOk:
		f(c.Memory[c.PC+1])

	case gOk:
		hi, lo := c.Memory[c.PC+1], c.Memory[c.PC+2]
		v := (uint16(hi) << 8) | uint16(lo)
		g(v)

	case hOk:
		h()

	default:
		valid = false
	}
	return
}

// Moves the writer pointer of the Core. As of writing this pointer cannot be set
// to before `x0200` as the first 512 bytes of memory are the zero page and stack.
func (c *Core) SetWriterPtr(value uint16) (err error) {
	if value < 0x0200 {
		err = errors.New("writing pointer must be set in general purpose memory")
		return
	}
	c.writingPointer = value
	return
}

// Writes the contents of the byte slice to general memory, always stopping at
// the end of general memory (`0xFFFF`); will return the amount of bytes written.
//
// This uses the `*Core.writingPointer` which can be moved with `*Core.SetWriterPtr`.
func (c *Core) Write(what []byte) (n int) {
	limit := 0xFFFF - c.writingPointer
	n = min(int(limit), len(what))

	for i := 0; i < n; i++ {
		c.Memory[c.writingPointer] = what[i]
		c.writingPointer++
	}
	return
}
