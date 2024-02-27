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
