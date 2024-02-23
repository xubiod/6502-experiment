package cpu

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

func (c *Core) TSX____i() {
	c.PC += 1

	c.X = c.S

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

func (c *Core) TXS____i() {
	c.PC += 1

	c.S = c.X

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
