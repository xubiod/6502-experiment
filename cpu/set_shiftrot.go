package cpu

// This is the implementation of the ASL instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) asl_impl(loc *byte) {
	var shouldCarry = *loc & 0b10000000

	*loc = (*loc << 1) & 0xFE

	if shouldCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

// This is the implementation of the LSR instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) lsr_impl(loc *byte) {
	var shouldCarry = *loc & 0b00000001

	*loc = (*loc >> 1) & 0x7F

	if shouldCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

// This is the implementation of the ROL instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) rol_impl(loc *byte) {
	var futureCarry = *loc & 0b10000000

	*loc = ((*loc << 1) & 0xFE) | (c.Flags & FLAG_CARRY)

	if futureCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

// This is the implementation of the ROR instructions.
//
// This will change the flags of the Core it's run in.
func (c *Core) ror_impl(loc *byte) {
	var futureCarry = *loc & 0b00000001

	if !c.Features.RotateRightBug {
		*loc = ((*loc >> 1) & 0x7F) | ((c.Flags & FLAG_CARRY) << 7)
	} else {
		*loc = (*loc << 1) & 0xFE
		futureCarry = 0
	}

	if futureCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

// Arithmetic Shift Left - Absolute
func (c *Core) ASL____a(addr uint16) { c.PC += 3; c.asl_impl(&c.Memory[addr]) }

// Arithmetic Shift Left - Absolute indexed with X
func (c *Core) ASL___ax(addr uint16) { c.PC += 3; c.asl_impl(&c.Memory[addr+uint16(c.X)]) }

// Arithmetic Shift Left - Accumulator
func (c *Core) ASL____A() { c.PC += 1; c.asl_impl(&c.A) }

// Arithmetic Shift Left - Zero Page
func (c *Core) ASL__ZPg(zp byte) { c.PC += 2; c.asl_impl(&c.Memory[zp]) }

// Arithmetic Shift Left - Zero Page indexed with X
func (c *Core) ASL__ZPx(zp byte) { c.PC += 2; c.asl_impl(&c.Memory[zp+c.X]) }

// Logical Shift Right - Absolute
func (c *Core) LSR____a(addr uint16) { c.PC += 3; c.lsr_impl(&c.Memory[addr]) }

// Logical Shift Right - Absolute indexed with X
func (c *Core) LSR___ax(addr uint16) { c.PC += 3; c.lsr_impl(&c.Memory[addr+uint16(c.X)]) }

// Logical Shift Right - Accumulator
func (c *Core) LSR____A() { c.PC += 1; c.lsr_impl(&c.A) }

// Logical Shift Right - Zero Page
func (c *Core) LSR__ZPg(zp byte) { c.PC += 2; c.lsr_impl(&c.Memory[zp]) }

// Logical Shift Right - Zero Page indexed with X
func (c *Core) LSR__ZPx(zp byte) { c.PC += 2; c.lsr_impl(&c.Memory[zp+c.X]) }

// Rotate Bits Left - Absolute
func (c *Core) ROL____a(addr uint16) { c.PC += 3; c.rol_impl(&c.Memory[addr]) }

// Rotate Bits Left - Absolute indexed with X
func (c *Core) ROL___ax(addr uint16) { c.PC += 3; c.rol_impl(&c.Memory[addr+uint16(c.X)]) }

// Rotate Bits Left - Accumulator
func (c *Core) ROL____A() { c.PC += 1; c.rol_impl(&c.A) }

// Rotate Bits Left - Zero Page
func (c *Core) ROL__ZPg(zp byte) { c.PC += 2; c.rol_impl(&c.Memory[zp]) }

// Rotate Bits Left - Zero Page indexed with X
func (c *Core) ROL__ZPx(zp byte) { c.PC += 2; c.rol_impl(&c.Memory[zp+c.X]) }

// Rotate Bits Right - Absolute
func (c *Core) ROR____a(addr uint16) { c.PC += 3; c.ror_impl(&c.Memory[addr]) }

// Rotate Bits Right - Absolute indexed with X
func (c *Core) ROR___ax(addr uint16) { c.PC += 3; c.ror_impl(&c.Memory[addr+uint16(c.X)]) }

// Rotate Bits Right - Accumulator
func (c *Core) ROR____A() { c.PC += 1; c.ror_impl(&c.A) }

// Rotate Bits Right - Zero Page
func (c *Core) ROR__ZPg(zp byte) { c.PC += 2; c.ror_impl(&c.Memory[zp]) }

// Rotate Bits Right - Zero Page indexed with X
func (c *Core) ROR__ZPx(zp byte) { c.PC += 2; c.ror_impl(&c.Memory[zp+c.X]) }
