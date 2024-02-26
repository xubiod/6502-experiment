package cpu

// Branch on Carry Clear
func (c *Core) BCC__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_CARRY == 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Carry Set
func (c *Core) BCS__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_CARRY != 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Result Not Zero
func (c *Core) BNE__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_ZERO == 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Result Zero
func (c *Core) BEQ__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_ZERO != 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Result Plus
func (c *Core) BPL__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_NEGATIVE == 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Result Minus
func (c *Core) BMI__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_NEGATIVE != 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Overflow Clear
func (c *Core) BVC__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_OVERFLOW == 0 {
		c.PC += uint16(relative)
	}
}

// Branch on Overflow Set
func (c *Core) BVS__rel(raw uint8) {
	var relative int8 = int8(raw & 0x7F)
	if raw&0b10000000 > 0 {
		relative *= -1
	}

	c.PC += 2
	if c.Flags&FLAG_OVERFLOW != 0 {
		c.PC += uint16(relative)
	}
}
