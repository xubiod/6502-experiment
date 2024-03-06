package cpu

// Store Accumulator to Memory - Absolute
func (c *Core) STA____a(addr uint16) { c.PC += 3; c.Memory[addr] = c.A }

// Store Accumulator to Memory - Absolute indexed with X
func (c *Core) STA___ax(addr uint16) { c.PC += 3; c.Memory[addr+uint16(c.X)] = c.A }

// Store Accumulator to Memory - Absolute indexed with Y
func (c *Core) STA___ay(addr uint16) { c.PC += 3; c.Memory[addr+uint16(c.Y)] = c.A }

// Store Accumulator to Memory - Zero Page
func (c *Core) STA__ZPg(zp byte) { c.PC += 2; c.Memory[zp] = c.A }

// Store Accumulator to Memory - Zero Page Indexed Indirect
func (c *Core) STA_IZPx(zp byte) { c.PC += 2; c.Memory[c.indirectZpX(zp)] = c.A }

// Store Accumulator to Memory - Zero Page indexed with X
func (c *Core) STA__ZPx(zp byte) { c.PC += 2; c.Memory[(zp+c.X)&0xFF] = c.A }

// Store Accumulator to Memory - Zero Page Indirect Indexed with Y
func (c *Core) STA_IZPy(zp byte) { c.PC += 2; c.Memory[c.indirectZpY(zp)] = c.A }

// Store X to Memory - Absolute
func (c *Core) STX____a(addr uint16) { c.PC += 3; c.Memory[addr] = c.X }

// Store X to Memory - Zero Page
func (c *Core) STX__ZPg(zp byte) { c.PC += 2; c.Memory[zp] = c.X }

// Store X to Memory - Zero Page indexed with Y
func (c *Core) STX__ZPy(zp byte) { c.PC += 2; c.Memory[zp+c.Y] = c.X }

// Store Y to Memory - Absolute
func (c *Core) STY____a(addr uint16) { c.PC += 3; c.Memory[addr] = c.Y }

// Store Y to Memory - Zero Page
func (c *Core) STY__ZPg(zp byte) { c.PC += 2; c.Memory[zp] = c.Y }

// Store Y to Memory - Zero Page indexed with X
func (c *Core) STY__ZPx(zp byte) { c.PC += 2; c.Memory[(zp+c.X)&0xFF] = c.Y }

// Store Zero to Memory - Absolute
//
// CMOS 65c02
func (c *Core) STZ____a(addr uint16) { c.PC += 3; c.Memory[addr] = 0 }

// Store Zero to Memory - Absolute indexed with X
//
// CMOS 65c02
func (c *Core) STZ___ax(addr uint16) { c.PC += 3; c.Memory[addr+uint16(c.X)] = 0 }

// Store Zero to Memory - Zero Page
//
// CMOS 65c02
func (c *Core) STZ__ZPg(zp byte) { c.PC += 2; c.Memory[zp] = 0 }

// Store Zero to Memory - Zero Page indexed with X
//
// CMOS 65c02
func (c *Core) STZ__ZPx(zp byte) { c.PC += 2; c.Memory[(zp+c.X)&0xFF] = 0 }

// Store Accumulator into Memory - Zero Page Indirect
//
// CMOS 65c02
func (c *Core) STA__IZP(zp byte) { c.PC += 2; c.Memory[c.indirectZp(zp)] = c.A }
