package cpu

func makeSigned(i uint8) (o int8) {
	o = int8(i & 0x7F)
	if i&0b10000000 > 0 {
		o *= -1
	}
	return
}

// Branch on Carry Clear
func (c *Core) BCC__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_CARRY == 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Carry Set
func (c *Core) BCS__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_CARRY != 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Result Not Zero
func (c *Core) BNE__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_ZERO == 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Result Zero
func (c *Core) BEQ__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_ZERO != 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Result Plus
func (c *Core) BPL__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_NEGATIVE == 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Result Minus
func (c *Core) BMI__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_NEGATIVE != 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Overflow Clear
func (c *Core) BVC__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_OVERFLOW == 0 {
		c.PC += uint16(makeSigned(raw))
	}
}

// Branch on Overflow Set
func (c *Core) BVS__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_OVERFLOW != 0 {
		c.PC += uint16(makeSigned(raw))
	}
}
