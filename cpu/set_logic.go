package cpu

// This is the implementation of the AND instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) and_impl(with byte) {
	c.A = c.A & with

	if c.A == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.A&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// This is the implementation of the ORA instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) ora_impl(with byte) {
	c.A = c.A | with

	if c.A == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.A&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// This is the implementation of the EOR instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) eor_impl(with byte) {
	c.A = c.A ^ with

	if c.A == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.A&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// Bitwise AND Accumulator with Memory - Absolute
func (c *Core) AND____a(addr uint16) { c.PC += 3; c.and_impl(c.Memory[addr]) }

// Bitwise AND Accumulator with Memory - Absolute indexed with X
func (c *Core) AND___ax(addr uint16) { c.PC += 3; c.and_impl(c.Memory[addr+uint16(c.X)]) }

// Bitwise AND Accumulator with Memory - Absolute indexed with Y
func (c *Core) AND___ay(addr uint16) { c.PC += 3; c.and_impl(c.Memory[addr+uint16(c.Y)]) }

// Bitwise AND Accumulator with Memory - Immediate
func (c *Core) AND__Imm(literal byte) { c.PC += 2; c.and_impl(literal) }

// Bitwise AND Accumulator with Memory - Zero Page
func (c *Core) AND__ZPg(zp byte) { c.PC += 2; c.and_impl(c.Memory[zp]) }

// Bitwise AND Accumulator with Memory - Zero Page Indexed Indirect
func (c *Core) AND_IZPx(zp byte) { c.PC += 2; c.and_impl(c.Memory[c.indirectZpX(zp)]) }

// Bitwise AND Accumulator with Memory - Zero Page indexed with X
func (c *Core) AND__ZPx(zp byte) { c.PC += 2; c.and_impl(c.Memory[zp+c.X]) }

// Bitwise AND Accumulator with Memory - Zero Page Indirect Indexed with Y
func (c *Core) AND_IZPy(zp byte) { c.PC += 2; c.and_impl(c.Memory[c.indirectZpY(zp)]) }

// Bitwise OR Accumulator with Memory - Absolute
func (c *Core) ORA____a(addr uint16) { c.PC += 3; c.ora_impl(c.Memory[addr]) }

// Bitwise OR Accumulator with Memory - Absolute indexed with X
func (c *Core) ORA___ax(addr uint16) { c.PC += 3; c.ora_impl(c.Memory[addr+uint16(c.X)]) }

// Bitwise OR Accumulator with Memory - Absolute indexed with Y
func (c *Core) ORA___ay(addr uint16) { c.PC += 3; c.ora_impl(c.Memory[addr+uint16(c.Y)]) }

// Bitwise OR Accumulator with Memory - Immediate
func (c *Core) ORA__Imm(literal byte) { c.PC += 2; c.ora_impl(literal) }

// Bitwise OR Accumulator with Memory - Zero Page
func (c *Core) ORA__ZPg(zp byte) { c.PC += 2; c.ora_impl(c.Memory[zp]) }

// Bitwise OR Accumulator with Memory - Zero Page Indexed Indirect
func (c *Core) ORA_IZPx(zp byte) { c.PC += 2; c.ora_impl(c.Memory[c.indirectZpX(zp)]) }

// Bitwise OR Accumulator with Memory - Zero Page indexed with X
func (c *Core) ORA__ZPx(zp byte) { c.PC += 2; c.ora_impl(c.Memory[zp+c.X]) }

// Bitwise OR Accumulator with Memory - Zero Page Indirect Indexed with Y
func (c *Core) ORA_IZPy(zp byte) { c.PC += 2; c.ora_impl(c.Memory[c.indirectZpY(zp)]) }

// Bitwise Exclusive OR Accumulator with Memory - Absolute
func (c *Core) EOR____a(addr uint16) { c.PC += 3; c.eor_impl(c.Memory[addr]) }

// Bitwise Exclusive OR Accumulator with Memory - Absolute indexed with X
func (c *Core) EOR___ax(addr uint16) { c.PC += 3; c.eor_impl(c.Memory[addr+uint16(c.X)]) }

// Bitwise Exclusive OR Accumulator with Memory - Absolute indexed with Y
func (c *Core) EOR___ay(addr uint16) { c.PC += 3; c.eor_impl(c.Memory[addr+uint16(c.Y)]) }

// Bitwise Exclusive OR Accumulator with Memory - Immediate
func (c *Core) EOR__Imm(literal byte) { c.PC += 2; c.eor_impl(literal) }

// Bitwise Exclusive OR Accumulator with Memory - Zero Page
func (c *Core) EOR__ZPg(zp byte) { c.PC += 2; c.eor_impl(c.Memory[zp]) }

// Bitwise Exclusive OR Accumulator with Memory - Zero Page Indexed Indirect
func (c *Core) EOR_IZPx(zp byte) { c.PC += 2; c.eor_impl(c.Memory[c.indirectZpX(zp)]) }

// Bitwise Exclusive OR Accumulator with Memory - Zero Page indexed with X
func (c *Core) EOR__ZPx(zp byte) { c.PC += 2; c.eor_impl(c.Memory[zp+c.X]) }

// Bitwise Exclusive OR Accumulator with Memory - Zero Page Indirect Indexed with Y
func (c *Core) EOR_IZPy(zp byte) { c.PC += 2; c.eor_impl(c.Memory[c.indirectZpY(zp)]) }
