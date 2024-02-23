package cpu

func (c *Core) PHA____i() {
	c.PC += 1

	c.Memory[0x0100+uint16(c.S)] = c.A
	c.S--
}

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

func (c *Core) PHP____i() {
	c.PC += 1

	c.Memory[0x0100+uint16(c.S)] = c.Flags & ^FLAG_BREAK
	c.S--
}

func (c *Core) PLP____i() {
	c.PC += 1

	c.S++
	c.Flags = c.Memory[0x0100+uint16(c.S)]
}
