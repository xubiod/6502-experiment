package cpu

// Push Accumulator to Stack - Implied
func (c *Core) PHA____i() {
	c.PC += 1

	c.Memory[0x0100+uint16(c.S)] = c.A
	c.S--
}

// Pull Accumulator from Stack - Implied
func (c *Core) PLA____i() {
	c.PC += 1

	c.S++
	c.A = c.Memory[0x0100+uint16(c.S)]

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

// Push Processor State to Stack - Implied
func (c *Core) PHP____i() {
	c.PC += 1

	c.Memory[0x0100+uint16(c.S)] = (c.Flags & ^FLAG_BREAK) & ^FLAG_UNUSED
	c.S--
}

// Pull Processor State from Stack - Implied
func (c *Core) PLP____i() {
	c.PC += 1

	c.S++
	c.Flags = c.Memory[0x0100+uint16(c.S)] & ^FLAG_UNUSED
}

// Push X to Stack - Implied
//
// CMOS 65c02
func (c *Core) PHX____i() {
	c.PC += 1

	c.Memory[0x0100+uint16(c.S)] = c.X
	c.S--
}

// Pull X from Stack - Implied
//
// CMOS 65c02
func (c *Core) PLX____i() {
	c.PC += 1

	c.S++
	c.X = c.Memory[0x0100+uint16(c.S)]

	if c.X == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.X&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// Push Y to Stack - Implied
//
// CMOS 65c02
func (c *Core) PHY____i() {
	c.PC += 1

	c.Memory[0x0100+uint16(c.S)] = c.Y
	c.S--
}

// Pull Y from Stack - Implied
//
// CMOS 65c02
func (c *Core) PLY____i() {
	c.PC += 1

	c.S++
	c.Y = c.Memory[0x0100+uint16(c.S)]

	if c.Y == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.Y&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}
