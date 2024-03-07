package cpu

// This is the implementation of what the `CMP`, `CPX`, and `CPY` instructions
// do.
//
// This will change the flags in the Core it's run in.
func (c *Core) cmp_impl(what, with byte) {
	if what == with {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
		c.Flags = c.Flags | FLAG_ZERO
		c.Flags = c.Flags | FLAG_CARRY
	} else if what < with {
		c.Flags = c.Flags | FLAG_NEGATIVE
		c.Flags = c.Flags & ^FLAG_ZERO
		c.Flags = c.Flags & ^FLAG_CARRY
	} else if what > with {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
		c.Flags = c.Flags & ^FLAG_ZERO
		c.Flags = c.Flags | FLAG_CARRY
	}
}

// This is the implementation of what the `BIT` instruction does.
//
// This will change the flags in the Core it's run in.
func (c *Core) bit_impl(with byte) {
	var r = c.A & with

	if r == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if r&0b01000000 > 0 {
		c.Flags = c.Flags | FLAG_OVERFLOW
	} else {
		c.Flags = c.Flags & ^FLAG_OVERFLOW
	}

	if r&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// Compare Memory with Accumulator - Absolute
func (c *Core) CMP____a(addr uint16) { c.PC += 3; c.cmp_impl(c.A, c.Memory[addr]) }

// Compare Memory with Accumulator - Absolute indexed with X
func (c *Core) CMP___ax(addr uint16) { c.PC += 3; c.cmp_impl(c.A, c.Memory[addr+uint16(c.X)]) }

// Compare Memory with Accumulator - Absolute indexed with Y
func (c *Core) CMP___ay(addr uint16) { c.PC += 3; c.cmp_impl(c.A, c.Memory[addr+uint16(c.Y)]) }

// Compare Memory with Accumulator - Immediate
func (c *Core) CMP__Imm(literal byte) { c.PC += 2; c.cmp_impl(c.A, literal) }

// Compare Memory with Accumulator - Zero Page
func (c *Core) CMP__ZPg(zp byte) { c.PC += 2; c.cmp_impl(c.A, c.Memory[zp]) }

// Compare Memory with Accumulator - Zero Page Indexed Indirect
func (c *Core) CMP_IZPx(zp byte) { c.PC += 2; c.cmp_impl(c.A, c.Memory[c.indirectZpX(zp)]) }

// Compare Memory with Accumulator - Zero Page indexed with X
func (c *Core) CMP__ZPx(zp byte) { c.PC += 2; c.cmp_impl(c.A, c.Memory[(zp+c.X)&0xFF]) }

// Compare Memory with Accumulator - Zero Page Indirect Indexed with Y
func (c *Core) CMP_IZPy(zp byte) { c.PC += 2; c.cmp_impl(c.A, c.Memory[c.indirectZpY(zp)]) }

// Compare Memory with X - Absolute
func (c *Core) CPX____a(addr uint16) { c.PC += 3; c.cmp_impl(c.X, c.Memory[addr]) }

// Compare Memory with X - Immediate
func (c *Core) CPX__Imm(literal byte) { c.PC += 2; c.cmp_impl(c.X, literal) }

// Compare Memory with X - Zero Page
func (c *Core) CPX__ZPg(zp byte) { c.PC += 2; c.cmp_impl(c.X, c.Memory[zp]) }

// Compare Memory with Y - Absolute
func (c *Core) CPY____a(addr uint16) { c.PC += 3; c.cmp_impl(c.Y, c.Memory[addr]) }

// Compare Memory with Y - Immediate
func (c *Core) CPY__Imm(literal byte) { c.PC += 2; c.cmp_impl(c.Y, literal) }

// Compare Memory with Y - Zero Page
func (c *Core) CPY__ZPg(zp byte) { c.PC += 2; c.cmp_impl(c.Y, c.Memory[zp]) }

// Bit Test Memory with Accumulator - Absolute
func (c *Core) BIT____a(addr uint16) { c.PC += 3; c.bit_impl(c.Memory[addr]) }

// Bit Test Memory with Accumulator - Zero Page
func (c *Core) BIT__ZPg(zp byte) { c.PC += 2; c.bit_impl(c.Memory[zp]) }

// 65c02 Instructions/Implementations below this line

func (c *Core) trb_impl(loc uint16) {
	what := c.Memory[loc]
	var r = c.A & what

	working := c.A ^ 0xFF
	working = working & what

	c.Memory[loc] = working

	if r == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

func (c *Core) tsb_impl(loc uint16) {
	what := c.Memory[loc]
	var r = c.A & what

	working := c.A
	working = working | what

	c.Memory[loc] = working

	if r == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

// Bit Test Memory with Accumulator - Absolute Indexed with X
//
// CMOS 65c02
func (c *Core) BIT___ax(addr uint16) { c.bit_impl(c.Memory[addr+uint16(c.X)]) }

// Bit Test Memory with Accumulator - Zero Page Indexed with X
//
// CMOS 65c02
func (c *Core) BIT__ZPx(zp byte) { c.bit_impl(c.Memory[(zp+c.X)&0xFF]) }

// Bit Test Memory with Accumulator - Immediate
//
// CMOS 65c02
func (c *Core) BIT__Imm(literal byte) { c.bit_impl(literal) }

// Compare Memory with Accumulator - Zero Page Indirect
//
// CMOS 65c02
func (c *Core) CMP__IZP(zp byte) { c.PC += 2; c.cmp_impl(c.A, c.Memory[c.indirectZp(zp)]) }

// Test and Reset Bits - Absolute
//
// CMOS 65c02
func (c *Core) TRB____a(addr uint16) { c.PC += 3; c.trb_impl(addr) }

// Test and Reset Bits - Zero Page
//
// CMOS 65c02
func (c *Core) TRB__ZPg(zp byte) { c.PC += 2; c.trb_impl(uint16(zp)) }

// Test and Set Bits - Absolute
//
// CMOS 65c02
func (c *Core) TSB____a(addr uint16) { c.PC += 3; c.tsb_impl(addr) }

// Test and Set Bits - Zero Page
//
// CMOS 65c02
func (c *Core) TSB__ZPg(zp byte) { c.PC += 2; c.tsb_impl(uint16(zp)) }
