package assembler

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MemLocation6502 uint16

type BlockType int

const (
	B_TEXT BlockType = iota
	B_DATA
	B_REM
)

// An Assembler hold all of the data responsible for assembling assembly language
// to byte code for a 6502 processor.
//
// Syntax is available in a more human-readable format with `README.md`.
type Assembler struct {
	// A map for strings to locations in memory. Preprocessing fills this map
	// while parsing uses it appropriately.
	//
	// Parsing treats labels as "copy and paste", treating relative addressed
	// instructions appropriately being the only exception.
	Labels map[string]MemLocation6502

	// The current memory location during the assembly process.
	//
	// Preprocessing uses this to keep track of where labels should be within
	// memory.
	//
	// Parsing uses this with labels when they are used in relative instructions
	// like branches to make the value for branches correct.
	CurrentLocation MemLocation6502

	// The current line number being processed.
	//
	// Parsing uses this for error reporting.
	Line uint16

	// The current processing mode.
	//
	// Parsing uses this to know when to parse instructions, data blocks, or
	// completely ignore for remark blocks.
	processingMode BlockType
}

var (
	ErrInvalidInstruction     = errors.New("invalid instruction name")
	ErrInvalidAddressingMode  = errors.New("invalid instruction addressing mode")
	ErrInvalidBlockType       = errors.New("invalid block type (not a text, data, or remark)")
	ErrInvalidBlockLineLen    = errors.New("data line missing nibble")
	ErrLabelLocationIllogical = errors.New("branch cannot reach this label")

	ErrHCF = errors.New("halt and catch fire? so funny hehe haha")
)

// Combines the error, raw line, and current line number the assembler was parsing
// into one error to simplify debugging the program being assembled.
func (a *Assembler) appendLine(err error, rawLine string) error {
	return fmt.Errorf("%s (line %d)\n\t-> %d | %s", err, a.Line, a.Line, rawLine)
}

// Creates and sets up an Assembler for use.
func New() *Assembler {
	return &Assembler{CurrentLocation: 0x200, Labels: make(map[string]MemLocation6502), processingMode: B_TEXT}
}

const (
	INST_PATTERN string = `^([a-z]{3})` // Constant for what an instruction looks like. Used by all instruction regex patterns.
)

var (
	reLabel = regexp.MustCompile(`^[A-Za-z_]\w*:`) // Regex for a label declaration pattern.
	reBlock = regexp.MustCompile(`^\.\w+$`)        // Regex for a block pattern.

	reIZPgY     = regexp.MustCompile(INST_PATTERN + `\s+\(\$([0-9a-f]{2})\),y`) // Regex for an indirect zero page indirect indexed with Y instruction.
	reIZPgX     = regexp.MustCompile(INST_PATTERN + `\s+\(\$([0-9a-f]{2}),x\)`) // Regex for an indirect zero page indexed indirect instruction.
	reIAbs      = regexp.MustCompile(INST_PATTERN + `\s+\(\$([0-9a-f]{4})\)`)   // Regex for an indirect absolute instruction.
	reAbsY      = regexp.MustCompile(INST_PATTERN + `\s+\$([0-9a-f]{4}),y`)     // Regex for an absolute address indexed with Y instruction.
	reAbsX      = regexp.MustCompile(INST_PATTERN + `\s+\$([0-9a-f]{4}),x`)     // Regex for an absolute address indexed with X instruction.
	reZPgY      = regexp.MustCompile(INST_PATTERN + `\s+\$([0-9a-f]{2}),y`)     // Regex for a zero page indexed with Y instruction.
	reZPgX      = regexp.MustCompile(INST_PATTERN + `\s+\$([0-9a-f]{2}),x`)     // Regex for a zero page indexed with X instruction.
	reAbs       = regexp.MustCompile(INST_PATTERN + `\s+\$([0-9a-f]{4})`)       // Regex for an absolute address instruction.
	reOneByte   = regexp.MustCompile(INST_PATTERN + `\s+\$([0-9a-f]{2})`)       // Regex for a single byte operand instruction.
	reLiteral   = regexp.MustCompile(INST_PATTERN + `\s+#\$([0-9a-f]{2})`)      // Regex for an immediate operand instruction.
	reNoOperand = regexp.MustCompile(INST_PATTERN)                              // Regex for the basic instruction with no operands.

	reHCF = regexp.MustCompile(`^hcf`) // Regex for a string starting with "hcf"

	allWhitespace = regexp.MustCompile(`\s`)
)

var (
	// Instructions that have an zero page address as an operand for an indirect
	// indexed with Y value (1 byte).
	TB_IZPgY = map[string]byte{
		"ora": 0x11,
		"and": 0x31,
		"eor": 0x51,
		"adc": 0x71,
		"sta": 0x91,
		"lda": 0xB1,
		"cmp": 0xD1,
		"sbc": 0xF1,
	}

	// Instructions that have an zero page address as an operand for an indirect
	// value indexed with X (1 byte).
	TB_IZPgX = map[string]byte{
		"ora": 0x01,
		"and": 0x21,
		"eor": 0x41,
		"adc": 0x61,
		"sta": 0x81,
		"lda": 0xA1,
		"cmp": 0xC1,
		"sbc": 0xE1,
	}

	// Instructions that have an absolute address as an operand for an indirect
	// value (2 bytes).
	TB_IAbs = map[string]byte{
		"jmp": 0x6c,
	}

	// Instructions that have an absolute address as an operand for indexed with
	// Y (2 bytes).
	TB_AbsY = map[string]byte{
		"ora": 0x09,
		"and": 0x29,
		"eor": 0x49,
		"adc": 0x69,
		"sta": 0x89,
		"lda": 0xA9,
		"cmp": 0xC9,
		"sbc": 0xE9,

		"ldx": 0xBE,
	}

	// Instructions that have an absolute address as an operand for indexed with
	// X (2 bytes).
	TB_AbsX = map[string]byte{
		"ora": 0x1D, "asl": 0x1E,
		"and": 0x3D, "rol": 0x3E,
		"eor": 0x5D, "lsr": 0x5E,
		"adc": 0x7D, "ror": 0x7E,
		"sta": 0x9D,
		"ldy": 0xBC, "lda": 0xBD,
		"cmp": 0xDD, "dec": 0xDE,
		"sbc": 0xFD, "inc": 0xFE,
	}

	// Instructions that have a zero page address as an operand for indexing with
	// Y (1 byte).
	TB_ZPgY = map[string]byte{
		"stx": 0x96,
		"ldx": 0xB6,
	}

	// Instructions that have a zero page address as an operand for indexing with
	// X (1 byte).
	TB_ZPgX = map[string]byte{
		"ora": 0x15, "asl": 0x16,
		"and": 0x35, "rol": 0x36,
		"eor": 0x55, "lsr": 0x56,
		"adc": 0x75, "ror": 0x76,
		"sty": 0x94, "sta": 0x95,
		"ldy": 0xB4, "lda": 0xB5,
		"cmp": 0xD5, "dec": 0xD6,
		"sbc": 0xF5, "inc": 0xF6,
	}

	// Instructions that have an absolute address as an operand (2 bytes).
	TB_Abs = map[string]byte{
		"ora": 0x0D, "asl": 0x0E,
		"jsr": 0x20, "bit": 0x2C, "and": 0x2D, "rol": 0x2E,
		"jmp": 0x4C, "eor": 0x4D, "lsr": 0x4E,
		"adc": 0x6D, "ror": 0x6E,
		"sty": 0x8C, "sta": 0x8D, "stx": 0x8E,
		"ldy": 0xAC, "lda": 0xAD, "ldx": 0xAE,
		"cpy": 0xCC, "cmp": 0xCD, "dec": 0xCE,
		"cpx": 0xEC, "sbc": 0xED, "inc": 0xEE,
	}

	// Instructions that have a signed byte/relative jump as an operand (1 byte).
	TB_Relative = map[string]byte{
		"bpl": 0x10,
		"bmi": 0x30,
		"bvc": 0x50,
		"bvs": 0x70,
		"bcc": 0x90,
		"bcs": 0xB0,
		"bne": 0xD0,
		"beq": 0xF0,
	}

	// Instructions that have an immediate as an operand (1 byte).
	TB_Literal = map[string]byte{
		"ldy": 0xA0,
		"cpy": 0xC0,
		"cpx": 0xE0, "ldx": 0xA2, "cmp": 0xC9,
		"sbc": 0xE9, "lda": 0xA9, "ora": 0x09,
		"and": 0x29,
		"eor": 0x49,
		"adc": 0x69,
	}

	// Instructions that have a zero page address as an operand (1 byte).
	TB_Zp = map[string]byte{
		"ora": 0x05, "asl": 0x06,
		"bit": 0x24, "and": 0x25, "rol": 0x26,
		"eor": 0x45, "lsr": 0x46,
		"adc": 0x65, "ror": 0x66,
		"sty": 0x84, "sta": 0x85, "stx": 0x86,
		"ldy": 0xA4, "lda": 0xA5, "ldx": 0xA6,
		"cpy": 0xC4, "cmp": 0xC5, "dec": 0xC6,
		"cpx": 0xE4, "sbc": 0xE5, "inc": 0xE6,
	}

	// Instructions that have no operands.
	TB_NoOperand = map[string]byte{
		"brk": 0x00, "php": 0x08, "asl": 0x0A,
		"clc": 0x18,
		"plp": 0x28, "rol": 0x2A,
		"sec": 0x38,
		"rti": 0x40, "pha": 0x48, "lsr": 0x4A,
		"cli": 0x58,
		"rts": 0x60, "pla": 0x68, "ror": 0x6A,
		"sei": 0x78,
		"dey": 0x88, "txa": 0x8A,
		"tya": 0x98, "txs": 0x9A,
		"tay": 0xA8, "tax": 0xAA,
		"clv": 0xB8, "tsx": 0xBA,
		"iny": 0xC8, "dex": 0xCA,
		"cld": 0xD8,
		"inx": 0xE8, "nop": 0xEA,
		"sed": 0xF8,
	}
)

// Does the preprocessing pass on the given line.
//
// Preprocessing does label discovery and transverses a line while appropriately
// accommodating how much memory instructions would take when converted to byte
// code.
//
// This fills up the `*Assembler.Labels` for the parsing pass.
func (a *Assembler) PreprocessLine(line string) {
	line, _, _ = strings.Cut(line, ";")

	if len(strings.TrimSpace(line)) == 0 {
		return
	}

	if reLabel.MatchString(line) {
		line = strings.TrimSpace(line)
		line = strings.Trim(line, ":")
		a.Labels[line] = a.CurrentLocation
		return
	}

	isRel := false
	for op := range TB_Relative {
		if strings.Contains(line, op) {
			isRel = true
			break
		}
	}

	for label := range a.Labels {
		if strings.Contains(line, label) {
			if isRel {
				line = strings.ReplaceAll(line, label, "$DE")
			} else {
				line = strings.ReplaceAll(line, label, "$DEAD")
			}
		}
	}

	line = strings.TrimSpace(strings.ToLower(line))

	switch {
	case reIZPgY.MatchString(line), reIZPgX.MatchString(line), reZPgY.MatchString(line),
		reZPgX.MatchString(line), reOneByte.MatchString(line), reLiteral.MatchString(line):
		a.CurrentLocation += 2

	case reIAbs.MatchString(line), reAbsY.MatchString(line), reAbsX.MatchString(line),
		reAbs.MatchString(line):
		a.CurrentLocation += 3

	case reNoOperand.MatchString(line):
		a.CurrentLocation++

	default:
	}
}

// Resets the state for the parsing pass after the preprocessing pass finishes.
func (a *Assembler) PreprocessFinish() {
	a.CurrentLocation = 0x200
}

// Preprocesses a string like it was a file, breaking on newlines (`\n`). Calls
// `PreprocessLine` on these lines.
//
// After all lines are preprocessed, `PreprocessFinish` is called.
func (a *Assembler) Preprocess(prg string) {
	for _, line := range strings.Split(prg, "\n") {
		a.PreprocessLine(line)
	}
	a.PreprocessFinish()
}

// Does the parsing pass on the given line.
//
// Parsing uses regular expression patterns that are compiled with `*regexp.MustCompile`
// for efficiency and more of a guarantee that patterns are valid.
//
// Syntax is elaborated in the `README.md` file, and should be trusted as what
// the assembler sees as valid in a more human-readable way.
func (a *Assembler) ParseLine(line string) (out []byte, err error) {
	line, _, _ = strings.Cut(line, ";")

	if len(strings.TrimSpace(line)) == 0 {
		return
	}

	if reBlock.MatchString(strings.TrimSpace(line)) {
		line = strings.TrimSpace(strings.ToLower(line))

		switch line {
		case ".text", ".txt", ".t":
			a.processingMode = B_TEXT
		case ".data", ".dat", ".d":
			a.processingMode = B_DATA
		case ".remark", ".rem", ".r":
			a.processingMode = B_REM
		default:
			err = ErrInvalidBlockType
		}

		return
	}

	if reLabel.MatchString(line) {
		return
	}

	switch a.processingMode {
	case B_REM:
		return

	case B_TEXT:
		isRel := false
		for op := range TB_Relative {
			if strings.Contains(strings.ToLower(line), op) {
				isRel = true
				break
			}
		}

		for label, labelTo := range a.Labels {
			if strings.Contains(line, label) {
				if isRel {
					var diff = int32(a.CurrentLocation) - int32(labelTo)

					if diff > 127 || diff < -128 {
						err = ErrLabelLocationIllogical
						return
					}

					pos := byte(diff)
					line = strings.ReplaceAll(line, label, fmt.Sprintf("$%2X", pos))
				} else {
					line = strings.ReplaceAll(line, label, fmt.Sprintf("$%4X", labelTo))
				}
			}
		}

		var subs []string

		line = strings.TrimSpace(strings.ToLower(line))

		switch {
		case reIZPgY.MatchString(line):
			subs = reIZPgY.FindStringSubmatch(line)
			err = operationByte(subs, &TB_IZPgY, &out, &a.CurrentLocation)

		case reIZPgX.MatchString(line):
			subs = reIZPgX.FindStringSubmatch(line)
			err = operationByte(subs, &TB_IZPgX, &out, &a.CurrentLocation)

		case reIAbs.MatchString(line):
			subs = reIAbs.FindStringSubmatch(line)
			err = operationShort(subs, &TB_IAbs, &out, &a.CurrentLocation)

		case reAbsY.MatchString(line):
			subs = reAbsY.FindStringSubmatch(line)
			err = operationShort(subs, &TB_AbsY, &out, &a.CurrentLocation)

		case reAbsX.MatchString(line):
			subs = reAbsX.FindStringSubmatch(line)
			err = operationShort(subs, &TB_AbsX, &out, &a.CurrentLocation)

		case reZPgY.MatchString(line):
			subs = reZPgY.FindStringSubmatch(line)
			err = operationByte(subs, &TB_ZPgY, &out, &a.CurrentLocation)

		case reZPgX.MatchString(line):
			subs = reZPgX.FindStringSubmatch(line)
			err = operationByte(subs, &TB_ZPgX, &out, &a.CurrentLocation)

		case reAbs.MatchString(line):
			subs = reAbs.FindStringSubmatch(line)
			err = operationShort(subs, &TB_Abs, &out, &a.CurrentLocation)

		case reOneByte.MatchString(line):
			subs = reOneByte.FindStringSubmatch(line)
			if isRel {
				err = operationByte(subs, &TB_Relative, &out, &a.CurrentLocation)
			} else {
				err = operationByte(subs, &TB_Zp, &out, &a.CurrentLocation)
			}

		case reLiteral.MatchString(line):
			subs = reLiteral.FindStringSubmatch(line)
			err = operationByte(subs, &TB_Literal, &out, &a.CurrentLocation)

		case reHCF.MatchString(line):
			err = ErrHCF
			return

		case reNoOperand.MatchString(line):
			subs = reNoOperand.FindStringSubmatch(line)
			op, ok := TB_NoOperand[subs[0]]
			if !ok {
				err = ErrInvalidAddressingMode
				return
			}
			out = []byte{op}
			a.CurrentLocation++

		default:
			err = ErrInvalidInstruction
		}

	case B_DATA:
		line = allWhitespace.ReplaceAllString(line, "")

		if len(line)%2 != 0 {
			err = ErrInvalidBlockLineLen
			return
		}

		var convInter uint64
		for i := 0; i < len(line); i += 2 {
			convInter, _ = strconv.ParseUint(line[i:i+2], 16, 8)
			out = append(out, byte(convInter&0xFF))
		}
	}

	return
}

// Parses a string like it was a file, breaking on newlines (`\n`). Calls `ParseLine`
// on these lines.
//
// If `ParseLine` errors, the returned byte slice is emptied and the line that
// errored is appended to the error before returning it back. This is done for
// debugging simplicity.
func (a *Assembler) Parse(prg string) (out []byte, err error) {
	a.Line = 1
	var working []byte
	for _, line := range strings.Split(prg, "\n") {
		working, err = a.ParseLine(line)
		if err != nil {
			out = []byte{}
			err = a.appendLine(err, line)
			return
		}
		out = append(out, working...)
		a.Line++
	}
	return
}

// Executes `Preprocess` followed by `Parse`. The returns in `Parse` are returned
// with no modification.
func (a *Assembler) PreprocessAndParse(prg string) (out []byte, err error) {
	a.Preprocess(prg)
	out, err = a.Parse(prg)
	return
}

// Generates the bytecode for an instruction with one short/two byte operand.
//
// The contents of `subs` are expected to be the results of a `*regexp.Regexp.FindStringSubmatch`
// with no modifications to it.
//
// This reads from `opTable` but does not modify it. The value at the `out` pointer
// **will be overwritten**, and the value at the `mp` pointer is incremented twice.
//
// Output is the opcode, the low byte of the operand, and the high byte of the
// operand, in that order.
func operationShort(subs []string, opTable *map[string]byte, out *[]byte, mp *MemLocation6502) (err error) {
	op, ok := (*opTable)[subs[1]]
	if !ok {
		err = ErrInvalidAddressingMode
		return
	}
	convInter, _ := strconv.ParseUint(subs[2], 16, 16)
	or1 := byte(convInter & 0xFF)
	or2 := byte((convInter >> 8) & 0xFF)
	*out = []byte{op, or1, or2}

	*mp += 3
	return
}

// Generates the bytecode for an instruction with one byte operand.
//
// The contents of `subs` are expected to be the results of a `*regexp.Regexp.FindStringSubmatch`
// with no modifications to it.
//
// This reads from `opTable` but does not modify it. The value at the `out` pointer
// **will be overwritten**, and the value at the `mp` pointer is incremented twice.
//
// Output is the opcode and the operand, in that order.
func operationByte(subs []string, opTable *map[string]byte, out *[]byte, mp *MemLocation6502) (err error) {
	op, ok := (*opTable)[subs[1]]
	if !ok {
		err = ErrInvalidAddressingMode
		return
	}
	convInter, _ := strconv.ParseUint(subs[2], 16, 8)
	or1 := byte(convInter & 0xFF)
	*out = []byte{op, or1}

	*mp += 2
	return
}
