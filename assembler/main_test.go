package assembler

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	asm := New()

	questions := []string{
		"BRK", "RTI", "RTS", "PHP", "CLC", "PLP", "SEC", "PHA",
		"CLI", "PLA", "SEI", "DEY", "TYA", "TAY", "CLV", "INY",
		"CLD", "INX", "SED", "TXA", "TXS", "TAX", "TSX", "DEX",
		"NOP",
	}

	answers := [][]byte{
		{0x00}, {0x40}, {0x60}, {0x08}, {0x18}, {0x28}, {0x38}, {0x48},
		{0x58}, {0x68}, {0x78}, {0x88}, {0x98}, {0xA8}, {0xB8}, {0xC8},
		{0xD8}, {0xE8}, {0xF8}, {0x8A}, {0x9A}, {0xAA}, {0xBA}, {0xCA},
		{0xEA},
	}

	for idx, q := range questions {
		out, err := asm.ParseLine("\t" + q)
		if err != nil {
			t.Fatalf("simple - \"%s\" deadass did not assemble (%s)", q, err)
		}
		if slices.Compare(out, answers[idx]) != 0 {
			t.Fatalf("simple - \"%s\" should turn into %2X\tnot %2X", q, answers[idx], out)
		}
	}
}

func TestAssembleResetRoutine(t *testing.T) {
	asm := New()

	question := `.TXT ; this is unneeded usually but as an example it's here
	LDX #$FF
	TXS
	
	SEI
	CLC
	CLD
	CLI
	CLV
	
	LDA #$00
	LDX #$00
	LDY #$00`

	answer := []byte{
		0xa2, 0xff,
		0x9a,

		0x78,
		0x18,
		0xd8,
		0x58,
		0xb8,

		0xa9, 0x00,
		0xa2, 0x00,
		0xa0, 0x00,
	}

	out, err := asm.PreprocessAndParse(question)
	if err != nil {
		t.Fatalf("assemble_reset_routine - deadass did not assemble:\n%s", err)
	}
	if slices.Compare(out, answer) != 0 {
		t.Fatalf("assemble_reset_routine - program failed to assemble correctly")
	}
}

func TestAssembleFail(t *testing.T) {
	asm := New()

	question := `;
	LDX #$FF
	TXS
	
	SEI
	CLC
	CLD
	CLI
	CLV
	WTF	; Not an instruction
	
	LDA #$00
	LDX #$00
	LDY #$00`

	out, err := asm.PreprocessAndParse(question)
	if err == nil {
		t.Fatalf("assemble_fail - should have failed - did not")
	}

	if len(out) > 0 {
		t.Fatalf("assemble_fail - errored but returned partially assembled code, should be empty")
	}

	fmt.Printf("assemble_fail - failed successfully, error below:\n\n%s", err)
}

func TestHCF(t *testing.T) {
	asm := New()

	question := `;
	HCF`

	out, err := asm.PreprocessAndParse(question)
	if err == nil {
		t.Fatalf("hcf - should have failed - did not")
	}

	if len(out) > 0 {
		t.Fatalf("hcf - errored but returned partially assembled code, should be empty")
	}

	fmt.Printf("hcf - failed successfully, error below:\n\n%s", err)
}

func TestDataBlock(t *testing.T) {
	asm := New()

	question := `.DAT
	000102030405060708090A0B0C0D0E0F
	`

	answer := []byte{
		0x00,
		0x01,
		0x02,
		0x03,
		0x04,
		0x05,
		0x06,
		0x07,
		0x08,
		0x09,
		0x0A,
		0x0B,
		0x0C,
		0x0D,
		0x0E,
		0x0F,
	}

	out, err := asm.PreprocessAndParse(question)
	if err != nil {
		t.Fatalf("data_block - deadass did not assemble:\n%s", err)
	}
	if slices.Compare(out, answer) != 0 {
		t.Fatalf("data_block - program failed to assemble correctly")
	}
}

func TestDataMissingNibble(t *testing.T) {
	asm := New()

	question := `.DAT
	00010203040506
	0708090A0B0C0D
	0E0			; Missing a nibble, should be caught
	`

	out, err := asm.PreprocessAndParse(question)
	if err == nil {
		t.Fatalf("missing_nibble - should have failed - did not")
	}

	if len(out) > 0 {
		t.Fatalf("missing_nibble - errored but returned partially assembled code, should be empty")
	}

	fmt.Printf("missing_nibble - failed successfully, error below:\n\n%s", err)
}

func TestAssembleFailWithLabel(t *testing.T) {
	asm := New()

	question := "START:\n"
	question = question + strings.Repeat("\tNOP\n", 130)
	question = question + "\tBEQ START\n"

	out, err := asm.PreprocessAndParse(question)
	if err == nil {
		t.Fatalf("assemble_fail_with_labels - should have failed - did not")
	}

	if len(out) > 0 {
		t.Fatalf("assemble_fail_with_labels - errored but returned partially assembled code, should be empty")
	}

	fmt.Printf("assemble_fail_with_labels - failed successfully, error below:\n\n%s", err)
}
