package cpu

// The implementation of the increment instructions.
//
// This will change the flags of the Core it's run in, and will change the value
// at the pointer passed into it.
func (c *Core) inc_impl(where *byte) {
	*where++

	if (*where) == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if (*where)&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// The implementation of the decrement instructions.
//
// This will change the flags of the Core it's run in, and will change the value
// at the pointer passed into it.
func (c *Core) dec_impl(where *byte) {
	*where--

	if (*where) == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if (*where)&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// Increment Memory by One - Absolute
func (c *Core) INC____a(addr uint16) { c.PC += 3; c.inc_impl(&c.Memory[addr]) }

// Increment Memory by One - Absolute indexed with X
func (c *Core) INC___ax(addr uint16) { c.PC += 3; c.inc_impl(&c.Memory[addr+uint16(c.X)]) }

// Increment Memory by One - Zero Page
func (c *Core) INC__ZPg(zp byte) { c.PC += 2; c.inc_impl(&c.Memory[zp]) }

// Increment Memory by One - Zero Page indexed with X
func (c *Core) INC__ZPx(zp byte) { c.PC += 2; c.inc_impl(&c.Memory[(zp+c.X)&0xFF]) }

// Increment X by One - Implied
func (c *Core) INX____i() { c.PC += 1; c.inc_impl(&c.X) }

// Increment Y by One - Implied
func (c *Core) INY____i() { c.PC += 1; c.inc_impl(&c.Y) }

// Decrement Memory by One - Absolute
func (c *Core) DEC____a(addr uint16) { c.PC += 3; c.dec_impl(&c.Memory[addr]) }

// Decrement Memory by One - Absolute indexed with X
func (c *Core) DEC___ax(addr uint16) { c.PC += 3; c.dec_impl(&c.Memory[addr+uint16(c.X)]) }

// Decrement Memory by One - Zero Page
func (c *Core) DEC__ZPg(zp byte) { c.PC += 2; c.dec_impl(&c.Memory[zp]) }

// Decrement Memory by One - Zero Page indexed with X
func (c *Core) DEC__ZPx(zp byte) { c.PC += 2; c.dec_impl(&c.Memory[(zp+c.X)&0xFF]) }

// Decrement X by One - Implied
func (c *Core) DEX____i() { c.PC += 1; c.dec_impl(&c.X) }

// Decrement Y by One - Implied
func (c *Core) DEY____i() { c.PC += 1; c.dec_impl(&c.Y) }

// 65c02 Instructions/Implementations below this line

// Increment Accumulator by One - Implied
//
// CMOS 65c02
func (c *Core) INA____i() { c.PC += 1; c.inc_impl(&c.A) }

// Decrement Accumulator by One - Implied
//
// CMOS 65c02
func (c *Core) DEA____i() { c.PC += 1; c.dec_impl(&c.A) }
