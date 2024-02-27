package cpu

func branchVal(i uint8) (o uint16) {
	m := int8(i & 0x7F)
	if i&0b10000000 > 0 {
		m *= -1
	}
	return uint16(m)
}

func (c *Core) BCC__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_CARRY == 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BCS__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_CARRY != 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BNE__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_ZERO == 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BEQ__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_ZERO != 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BPL__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_NEGATIVE == 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BMI__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_NEGATIVE != 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BVC__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_OVERFLOW == 0 {
		c.PC += branchVal(raw)
	}
}

func (c *Core) BVS__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_OVERFLOW != 0 {
		c.PC += branchVal(raw)
	}
}
