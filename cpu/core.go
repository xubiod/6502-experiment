package cpu

import (
	"fmt"
)

// A Core is the main data structure of the emulator. It holds its own memory,
// registers, and execution maps for instruction execution.
//
// This emulator is made as a **generic 6502 CPU emulator**, and as such as features
// that can be toggled with the `Features` field (see `CoreFeatureFlags`) to make
// the emulator act more like a specific 6502-compatible CPU.
//
// The emulator was designed like this for the ability to be able to run multiple
// emulations at once independently in coroutines if ever necessary.
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

	// What to do before executing instructions in `StepOnce()`.
	PreStep func(this *Core)

	// What to do after executing instructions in `StepOnce()`.
	PostStep func(this *Core)

	// The byte -> implementation map for instructions with no operands.
	execMapNil map[byte]func()

	// The byte -> implementation map for instructions with a byte as an operand.
	// Signed or not is instruction dependent.
	execMapByte map[byte]func(uint8)

	// The byte -> implementation map for instructions with an unsigned short (2
	// bytes) as an operand.
	execMapShort map[byte]func(uint16)

	// The byte -> implementation map for instructions with no operands for 65c02
	// instructions. Separated to make CMOS toggleable.
	execMapNilCMOS map[byte]func()

	// The byte -> implementation map for instructions with a byte as an operand
	// for 65c02 instructions. Signed or not is instruction dependent, and separated
	// to make CMOS toggleable.
	execMapByteCMOS map[byte]func(uint8)

	// The byte -> implementation map for instructions with two byte operands for
	// 65c02 instructions, which are for the branch if bit set/cleared instructions
	// for how unique their structure is. Separated to make CMOS toggleable.
	execMapBitBranchCMOS map[byte]func(uint8, uint8)

	// The byte -> implementation map for instructions with an unsigned short (2
	// bytes) as an operand for 65c02 instructions. Separated to make CMOS
	// toggleable.
	execMapShortCMOS map[byte]func(uint16)

	writingPointer uint16 // The pointer to writing to memory with `*Core.Write()`.
}

// A struct for a set of feature flags that can be changed to have the emulator
// "specialized" to a specific 6502-compatible CPU instead of a generic 6502.
type CoreFeatureFlags struct {
	// Is decimal mode implemented? Toggling this off does not change the behaviour
	// of enabling/disabling the flag itself, but if off ADC/SBC will ignore the flag.
	//
	// Setting this to false will act like a NES CPU, and defaults to `true`.
	DecimalModeImplemented bool

	// Is the ROR instruction broken from early revisions? The earliest revisions of the
	// 6502 effectively had the ROR instruction as an arithmetic shift left that did not
	// affect the carry flag instead. This is no longer an issue in 6502 chips that exist
	// today.
	//
	// The opcode still _did_ something on the 6502, but it was not documented until it was
	// properly implemented into the 6502.
	//
	// This is defaulted to `false`.
	RotateRightBug bool

	// Is the indirect jump bugged? On the NMOS 6502 family the indirect jump gets the
	// incorrect address when at a high-end page boundary (`xxFF`) which causes the address
	// to be formed from the bytes at `xx00` and `xxFF` instead of the expected behaviour,
	// which would be getting the other byte from the next page and not the same page.
	//
	// The CMOS derivatives fix this issue.
	//
	// This is defaulted to `true`.
	NMOSAbsoluteIndirectBug bool

	// On the NMOS line, the flags after arithmetic operations while in decimal mode are
	// **based on the binary result** instead of the decimal result. As a result, the flags
	// after a decimal mode operation are generally seen as meaningless.
	//
	// As an example, in binary mode `$99 + $01` the result is `$9A` with carry cleared. In
	// decimal mode, the result would be `$00` with the carry set, but the zero flag remains
	// **cleared as the binary result was not zero**.
	//
	// The CMOS derivatives fix this issue.
	//
	// This is defaulted to `true`.
	NMOSDecimalModeFlagBug bool

	// This is just a simple flag to make `StepOnce()` treat invalid instructions as NOPs
	// instead of doing nothing. The NOP byte lengths are different and affect the amount
	// that the program counter is affected by.
	//
	// This is similar to the behaviour of CMOS derivatives.
	//
	// Note that this affects the returns of `StepOnce()`
	//
	// This is defaulted to `false`.
	IncrementPCOnInvalidInstruction bool

	// This enables the recognition and execution of CMOS instructions, notably the 65c02.
	//
	// The implemented CMOS instructions do not cover WDC 65c02 only instructions (STP and
	// WAI) at the time of writing.
	//
	// This is defaulted to `false`.
	EnableCMOSInstructions bool

	// Prints to console whenever a `BRK` is encountered. See `*Core.CompleteDump()`
	// for more about what is printed to console.
	//
	// At the time of writing, coloured output is OS-dependent using control codes,
	// and running on the Windows operating system will lead to uncoloured output.
	// This is simply determined with a runtime check, can cannot currently be
	// toggled otherwise.
	//
	// This is defaulted to true.
	ConsoleOutOnBreak bool
}

var defaultFeatures CoreFeatureFlags = CoreFeatureFlags{
	DecimalModeImplemented:          true,
	RotateRightBug:                  false,
	NMOSAbsoluteIndirectBug:         true,
	NMOSDecimalModeFlagBug:          true,
	IncrementPCOnInvalidInstruction: false,
	EnableCMOSInstructions:          false,
	ConsoleOutOnBreak:               true,
}

const (
	FLAG_CARRY             byte = 1 << iota // C - Set when the last operation resulted in an overflow.
	FLAG_ZERO                               // Z - Set when the last operation resulted in a zero.
	FLAG_INTERRUPT_DISABLE                  // I - When set, interrupts are disabled.
	FLAG_DECIMAL                            // D - When set, math operations are done with BCD. No other operation is affected by the status of this flag.
	FLAG_BREAK                              // B - Set when a software interrupt happens with `BRK`.
	FLAG_UNUSED                             // _ - This flag is not used by the 6502. It is always set on the 6502.
	FLAG_OVERFLOW                           // V - Set when the last operation resulted in a *signed overflow* if the numbers were interpreted as signed.
	FLAG_NEGATIVE                           // N - Set when the last operation resulted as a negative number as a bit 7 check.
)

const (
	HIGHLIGHT_SEGMENT = "\033[33m" // This is a control code for yellow text on default background
	HIGHLIGHT_CLEAR   = "\033[0m"  // This is a control code to clear the text/background colour to default
)

// Does the calculations for a zero-page indirect indexed with Y address to get
// the valid address.
func (c *Core) indirectZpY(zp byte) (addr uint16) {
	var lsb, msb byte
	lsb = c.Memory[zp]
	msb = c.Memory[(zp+1)&0xFF]

	addr = (uint16(msb) << 8 & uint16(lsb)) + uint16(c.Y)
	return
}

// Does the calculations for a zero-page indexed indirect to get the address.
func (c *Core) indirectZpX(zp byte) (addr uint16) {
	var lsb, msb byte
	lsb = c.Memory[(zp+c.X)&0xFF]
	msb = c.Memory[(zp+c.X+1)&0xFF]

	addr = uint16(msb) << 8 & uint16(lsb)
	return
}

// Does the calculations for a zero-page indirect to get the address.
//
// Only used by 65c02 instructions.
func (c *Core) indirectZp(zp byte) (addr uint16) {
	var lsb, msb byte
	lsb = c.Memory[zp]
	msb = c.Memory[(zp+1)&0xFF]

	addr = (uint16(msb) << 8 & uint16(lsb))
	return
}

// Creates and prepares a *Core.
func NewCore() (c *Core) {
	c = &Core{Features: defaultFeatures}
	c.prepare()
	return
}

// Creates the decoding tables. Must be called before any execution unless writing
// your own execution loop.
func (c *Core) prepare() {

	// NMOS 6502

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

	// CMOS 65c02

	c.execMapNilCMOS = map[byte]func(){
		0x1A: c.INA____i,
		0x3A: c.DEA____i,
		0x5A: c.PHY____i,
		0x7A: c.PLY____i,
		0xDA: c.PHX____i,
		0xFA: c.PLX____i,
	}

	c.execMapByteCMOS = map[byte]func(uint8){
		0x04: c.TSB__ZPg, 0x07: c.RMB_G(0),
		0x12: c.ORA__IZP, 0x14: c.TRB__ZPg, 0x17: c.RMB_G(1),
		0x27: c.RMB_G(2),
		0x32: c.AND__IZP, 0x34: c.BIT__ZPx, 0x37: c.RMB_G(3),
		0x47: c.RMB_G(4),
		0x52: c.EOR__IZP, 0x57: c.RMB_G(5),
		0x64: c.STZ__ZPg, 0x67: c.RMB_G(6),
		0x72: c.ADC__IZP, 0x74: c.STZ__ZPx, 0x77: c.RMB_G(7),
		0x80: c.BRA__rel, 0x87: c.SMB_G(0), 0x89: c.BIT__Imm,
		0x92: c.STA__IZP, 0x97: c.SMB_G(1),
		0xA7: c.SMB_G(2),
		0xB2: c.LDA__IZP, 0xB7: c.SMB_G(3),
		0xC7: c.SMB_G(4),
		0xD2: c.CMP__IZP, 0xD7: c.SMB_G(5),
		0xE7: c.SMB_G(6),
		0xF2: c.SBC__IZP, 0xF7: c.SMB_G(7),
	}

	c.execMapBitBranchCMOS = map[byte]func(uint8, uint8){
		0x0F: c.BBR_G(0),
		0x1F: c.BBR_G(1),
		0x2F: c.BBR_G(2),
		0x3F: c.BBR_G(3),
		0x4F: c.BBR_G(4),
		0x5F: c.BBR_G(5),
		0x6F: c.BBR_G(6),
		0x7F: c.BBR_G(7),
		0x8F: c.BBS_G(0),
		0x9F: c.BBS_G(1),
		0xAF: c.BBS_G(2),
		0xBF: c.BBS_G(3),
		0xCF: c.BBS_G(4),
		0xDF: c.BBS_G(5),
		0xEF: c.BBS_G(6),
		0xFF: c.BBS_G(7),
	}

	c.execMapShortCMOS = map[byte]func(uint16){
		0x0C: c.TSB____a,
		0x1C: c.TRB____a,
		0x3C: c.BIT___ax,
		0x7C: c.JMP__Iax,
		0x9C: c.STZ____a, 0x9E: c.STZ___ax,
	}

	c.Flags = c.Flags | FLAG_UNUSED

	_ = c.SetWriterPtr(0x0000)
}

// Does a single step of execution. If at an invalid instruction, the program
// counter will not increment.
//
// Returns true if the instruction was valid.
func (c *Core) StepOnce() (valid bool) {
	var validNMOS, validCMOS bool

	if c.PreStep != nil {
		c.PreStep(c)
	}

	var fOk, gOk, hOk, iOk, jOk, kOk, lOk bool

	var f, i func(uint8)
	var g, j func(uint16)
	var h, k func()
	var l func(uint8, uint8)

	inst := c.Memory[c.PC]
	valid, validNMOS, validCMOS = true, true, true

	f, fOk = c.execMapByte[inst]
	g, gOk = c.execMapShort[inst]
	h, hOk = c.execMapNil[inst]

	if c.Features.EnableCMOSInstructions {
		i, iOk = c.execMapByteCMOS[inst]
		j, jOk = c.execMapShortCMOS[inst]
		k, kOk = c.execMapNilCMOS[inst]
		l, lOk = c.execMapBitBranchCMOS[inst]
	}

	switch {
	case fOk:
		f(c.Memory[c.PC+1])

	case gOk:
		hi, lo := c.Memory[c.PC+1], c.Memory[c.PC+2]
		v := (uint16(hi) << 8) | uint16(lo)
		g(v)

	case hOk:
		h()

	case iOk:
		validNMOS = false
		i(c.Memory[c.PC+1])

	case jOk:
		validNMOS = false
		hi, lo := c.Memory[c.PC+1], c.Memory[c.PC+2]
		v := (uint16(hi) << 8) | uint16(lo)
		j(v)

	case kOk:
		validNMOS = false
		k()

	case lOk:
		validNMOS = false
		l(c.Memory[c.PC+1], c.Memory[c.PC+2])

	default:
		validNMOS = false
		validCMOS = false

		if c.Features.IncrementPCOnInvalidInstruction {
			switch inst & 0x0F {
			case 0x03, 0x0B:
				c.PC += 1
				validCMOS = true
			case 0x02, 0x04:
				c.PC += 2
				validCMOS = true
			case 0x0C:
				c.PC += 3
				validCMOS = true
			}
		}
	}
	valid = validCMOS || validNMOS

	if c.PostStep != nil {
		c.PostStep(c)
	}

	return
}

// Moves the writer pointer of the Core.
func (c *Core) SetWriterPtr(value uint16) (err error) {
	// if value < 0x0200 {
	// 	err = errors.New("writing pointer must be set in general purpose memory")
	// 	return
	// }
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

// Returns the processor state for printing to console, or any other human-readable
// logging format.
//
// The program counter, stack pointer, accumulator, X, and Y registers are printed
// as hexadecimal. The flags are presented as appropriate letters, with uppercase
// meaning set.
//
// The general flag letters are standard with other 6502 emulators:
//
//	Negative, oVerflow, Break, Decimal, Interrupt disable, Zero, Carry
func (c *Core) StateDump() (out string) {
	out = fmt.Sprintf("PC: %04x | S: %02x | A: %02x | X: %02x | Y: %02x | Fl: ",
		c.PC, c.S, c.A, c.X, c.Y)
	for idx, chr := range "nv-bdizc" {
		realRune := chr
		if c.Flags<<idx&0b10000000 > 0 && chr != '-' {
			realRune -= 32
		}
		out += string(realRune)
	}
	return
}

// A general wrapper for dumping segments of memory within the Core.
//
// If `start` is larger than `end`, they are swapped.
//
// The starting address is rounded down to the lowest high nibble (`A4` -> `A0`,
// `FF` -> `F0`) and at minimum this will return the contents of the next 16 bytes
// within the CPU's memory.
//
// If the highlight address is within the range, it is surrounded with square
// brackets. If `highlightColoured` is true, the location will be coloured yellow
// using control codes. This does not work out of the box on Windows.
func (c *Core) MemoryDump(start, end, highlight uint16, highlightColoured bool) (out string) {
	if start > end {
		start, end = end, start
	}

	var point uint16 = start & 0xFFF0
	var i uint16
	var width uint16 = 16 // 8 is a good smaller width

	for ; point < end; point += width {
		out += fmt.Sprintf("0x%04X |", point)
		for i = 0; i < width; i++ {
			if point+i == highlight {
				if highlightColoured {
					out += HIGHLIGHT_SEGMENT
				}

				out += fmt.Sprintf("[%02x]", c.Memory[point+i])

				if highlightColoured {
					out += HIGHLIGHT_CLEAR
				}
			} else {
				if point+i-1 == highlight && ((point+i-1)&0xF0 == (point+i)&0xF0) {
					out += fmt.Sprintf("%02x", c.Memory[point+i])
				} else {
					out += fmt.Sprintf(" %02x", c.Memory[point+i])
				}
			}
		}
		out += "\n"
	}
	return out
}

// Returns the stack dump for printing to console, or any other human-readable
// logging format. If coloured is true, it adds characters to colour the output
// for terminals, which will not be supported by Windows out of the box.
//
// The output is a general memory output, starting at the high nibble of the stack
// pointer (for example, if the stack pointer was `B4`, the output starts at
// address`0x01B0`) and continues to the end of the stack at address `0x01FF`.
//
// See `*Core.MemoryDump` for detailed output documentation.
func (c *Core) StackDump(coloured bool) (out string) {
	out = "Full Stack:\n"
	out += c.MemoryDump(0x0100+uint16(c.S), 0x01FF, 0x0100+uint16(c.S), coloured)

	return out
}

// Returns the memory dump for printing to console, or any other human-readable
// logging format. If coloured is true, it adds characters to colour the output
// for terminals, which will not be supported by Windows out of the box.
//
// The output is a general memory output, starting at the high nibble of the
// program counter - 48 bytes, continuing until the program counter + 16 bytes.
//
// See `*Core.MemoryDump` for detailed output documentation.
func (c *Core) ProgramCounterDump(coloured bool) (out string) {
	out = "Around PC:\n"
	out += c.MemoryDump(uint16(max(int32(c.PC)-0x31, 0)), c.PC+0x11, c.PC, coloured)
	return out
}

// Returns a combination of all the dump methods for a Core as one string.
//
// See `*Core.StateDump`, `*Core.StackDump`, and `*Core.ProgramCounterDump` for
// a complete documentation; in short the processor state is outputted, followed
// by a dump of the stack starting at the stack pointer.
func (c *Core) CompleteDump(coloured bool) string {
	return c.StateDump() + "\n" + c.StackDump(coloured) + "\n" + c.ProgramCounterDump(coloured)
}
