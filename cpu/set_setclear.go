package cpu

// Clear Carry Flag - Implied
func (c *Core) CLC____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_CARRY }

// Set Carry Flag - Implied
func (c *Core) SEC____i() { c.PC += 1; c.Flags = c.Flags | FLAG_CARRY }

// Clear Decimal Flag - Implied
func (c *Core) CLD____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_DECIMAL }

// Set Decimal Flag - Implied
func (c *Core) SED____i() { c.PC += 1; c.Flags = c.Flags | FLAG_DECIMAL }

// Clear Interrupt Disable Flag - Implied
func (c *Core) CLI____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_INTERRUPT_DISABLE }

// Set Interrupt Disable Flag - Implied
func (c *Core) SEI____i() { c.PC += 1; c.Flags = c.Flags | FLAG_INTERRUPT_DISABLE }

// Clear Overflow Flag - Implied
func (c *Core) CLV____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_OVERFLOW }

// 65c02 Instructions/Implementations below this line

// Set Memory Bit - Func Generator
//
// This is like this because the opcode determines the bit, and there's no need to do
// this over and over
//
// CMOS 65c02
func (c *Core) SMB_G(bit uint) func(zp byte) {
	if bit > 7 {
		panic("can only check bits from 0 to 7")
	}
	return func(zp byte) {
		c.Memory[zp] = c.Memory[zp] | (0b00000001 << bit)
	}
}

// Clear Memory Bit - Func Generator
//
// This is like this because the opcode determines the bit, and there's no need to do
// this over and over
//
// CMOS 65c02
func (c *Core) RMB_G(bit uint) func(zp byte) {
	if bit > 7 {
		panic("can only check bits from 0 to 7")
	}
	return func(zp byte) {
		c.Memory[zp] = c.Memory[zp] & ^(0b00000001 << bit)
	}
}
