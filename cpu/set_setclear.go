package cpu

func (c *Core) CLC____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_CARRY }
func (c *Core) SEC____i() { c.PC += 1; c.Flags = c.Flags | FLAG_CARRY }
func (c *Core) CLD____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_DECIMAL }
func (c *Core) SED____i() { c.PC += 1; c.Flags = c.Flags | FLAG_DECIMAL }
func (c *Core) CLI____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_INTERRUPT_DISABLE }
func (c *Core) SEI____i() { c.PC += 1; c.Flags = c.Flags | FLAG_INTERRUPT_DISABLE }
func (c *Core) CLV____i() { c.PC += 1; c.Flags = c.Flags & ^FLAG_OVERFLOW }
