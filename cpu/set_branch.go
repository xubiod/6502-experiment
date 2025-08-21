package cpu

// Takes a unsigned 8 bit integer and mangles it into a signed 8 bit integer.
//
// The result is converted into a unsigned short to add to the program counter
// by the caller.
func branchVal(i uint8) (o uint16) {
	m := int8(i & 0x7F)
	if i&0b10000000 > 0 {
		m *= -1
	}
	return uint16(m)
}

// Branch on Carry Clear - Relative
func (c *Core) BCC__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_CARRY == 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Carry Set - Relative
func (c *Core) BCS__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_CARRY > 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Result Not Zero - Relative
func (c *Core) BNE__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_ZERO == 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Result Zero - Relative
func (c *Core) BEQ__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_ZERO > 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Result Plus - Relative
func (c *Core) BPL__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_NEGATIVE == 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Result Minus - Relative
func (c *Core) BMI__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_NEGATIVE > 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Overflow Clear - Relative
func (c *Core) BVC__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_OVERFLOW == 0 {
		c.PC += branchVal(raw)
	}
}

// Branch on Overflow Set - Relative
func (c *Core) BVS__rel(raw uint8) {
	c.PC += 2
	if c.Flags&FLAG_OVERFLOW > 0 {
		c.PC += branchVal(raw)
	}
}

// 65c02 Instructions/Implementations below this line

// Branch Always - Relative
//
// CMOS 65c02
func (c *Core) BRA__rel(raw uint8) { c.PC += 2 + branchVal(raw) }

// Branch if Bit Cleared - Func Generator
//
// This is like this because the opcode determines the bit, and there's no need to do
// this over and over
//
// CMOS 65c02
func (c *Core) BBR_G(bit uint) func(zp byte, raw uint8) {
	if bit > 7 {
		panic("can only check bits from 0 to 7")
	}
	return func(zp byte, raw uint8) {
		c.PC += 3
		if (c.Memory[zp]>>bit)&0b00000001 == 0 {
			c.PC += branchVal(raw)
		}
	}
}

// Branch if Bit Set - Func Generator
//
// This is like this because the opcode determines the bit, and there's no need to do
// this over and over
//
// CMOS 65c02
func (c *Core) BBS_G(bit uint) func(zp byte, raw uint8) {
	if bit > 7 {
		panic("can only check bits from 0 to 7")
	}
	return func(zp byte, raw uint8) {
		c.PC += 3
		if (c.Memory[zp]>>bit)&0b00000001 > 0 {
			c.PC += branchVal(raw)
		}
	}
}
