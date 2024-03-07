package cpu

// Transfer Accumulator to X - Implied
func (c *Core) TAX____i() {
	c.PC += 1

	c.X = c.A

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

// Transfer X to Accumulator - Implied
func (c *Core) TXA____i() {
	c.PC += 1

	c.A = c.X

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

// Transfer Accumulator to Y - Implied
func (c *Core) TAY____i() {
	c.PC += 1

	c.Y = c.A

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

// Transfer Y to Accumulator - Implied
func (c *Core) TYA____i() {
	c.PC += 1

	c.A = c.Y

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

// Transfer Stack Pointer to X - Implied
func (c *Core) TSX____i() {
	c.PC += 1

	c.X = c.S

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

// Transfer X to Stack Pointer - Implied
func (c *Core) TXS____i() {
	c.PC += 1

	c.S = c.X
}
