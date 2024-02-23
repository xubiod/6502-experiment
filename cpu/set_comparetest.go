package cpu

func (c *Core) cmp_impl(what, with byte) {
	if what == with {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
		c.Flags = c.Flags | FLAG_ZERO
		c.Flags = c.Flags | FLAG_CARRY
	} else if what < with {
		c.Flags = c.Flags | FLAG_NEGATIVE
		c.Flags = c.Flags & ^FLAG_ZERO
		c.Flags = c.Flags & ^FLAG_CARRY
	} else if what > with {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
		c.Flags = c.Flags & ^FLAG_ZERO
		c.Flags = c.Flags | FLAG_CARRY
	}
}

func (c *Core) bit_impl(with byte) {
	var r = c.A & with

	if r == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if r&0b01000000 > 0 {
		c.Flags = c.Flags | FLAG_OVERFLOW
	} else {
		c.Flags = c.Flags & ^FLAG_OVERFLOW
	}

	if r&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) CMP____a(addr uint16)  { c.PC += 3; c.cmp_impl(c.A, c.Memory[addr]) }
func (c *Core) CMP___ax(addr uint16)  { c.PC += 3; c.cmp_impl(c.A, c.Memory[addr+uint16(c.X)]) }
func (c *Core) CMP___ay(addr uint16)  { c.PC += 3; c.cmp_impl(c.A, c.Memory[addr+uint16(c.Y)]) }
func (c *Core) CMP__Imm(literal byte) { c.PC += 2; c.cmp_impl(c.A, literal) }
func (c *Core) CMP__ZPg(zp byte)      { c.PC += 2; c.cmp_impl(c.A, c.Memory[zp]) }
func (c *Core) CMP_IZPx(zp byte)      { c.PC += 2; c.cmp_impl(c.A, c.Memory[c.indirectZpX(zp)]) }
func (c *Core) CMP__ZPx(zp byte)      { c.PC += 2; c.cmp_impl(c.A, c.Memory[zp+c.X]) }
func (c *Core) CMP_IZPy(zp byte)      { c.PC += 2; c.cmp_impl(c.A, c.Memory[c.indirectZpY(zp)]) }

func (c *Core) CPX____a(addr uint16)  { c.PC += 3; c.cmp_impl(c.X, c.Memory[addr]) }
func (c *Core) CPX__Imm(literal byte) { c.PC += 2; c.cmp_impl(c.X, literal) }
func (c *Core) CPX__ZPg(zp byte)      { c.PC += 2; c.cmp_impl(c.X, c.Memory[zp]) }

// func (c *Core) CPX_ax(addr uint16)  { c.cmp_impl(c.X, c.Memory[addr+uint16(c.X)]) }
// func (c *Core) CPX_ay(addr uint16)  { c.cmp_impl(c.X, c.Memory[addr+uint16(c.Y)]) }
// func (c *Core) CPX_IndirectZPx(zp byte) { c.cmp_impl(c.X, c.Memory[c.indirectZpX(zp)]) }
// func (c *Core) CPX_ZPx(zp byte)         { c.cmp_impl(c.X, c.Memory[zp+c.X]) }
// func (c *Core) CPX_IndirectZPy(zp byte) { c.cmp_impl(c.X, c.Memory[c.indirectZpY(zp)]) }

func (c *Core) CPY____a(addr uint16)  { c.PC += 3; c.cmp_impl(c.Y, c.Memory[addr]) }
func (c *Core) CPY__Imm(literal byte) { c.PC += 2; c.cmp_impl(c.Y, literal) }
func (c *Core) CPY__ZPg(zp byte)      { c.PC += 2; c.cmp_impl(c.Y, c.Memory[zp]) }

// func (c *Core) CPY_ax(addr uint16)  { c.cmp_impl(c.Y, c.Memory[addr+uint16(c.X)]) }
// func (c *Core) CPY_ay(addr uint16)  { c.cmp_impl(c.Y, c.Memory[addr+uint16(c.Y)]) }
// func (c *Core) CPY_IndirectZPx(zp byte) { c.cmp_impl(c.Y, c.Memory[c.indirectZpX(zp)]) }
// func (c *Core) CPY_ZPx(zp byte)         { c.cmp_impl(c.Y, c.Memory[zp+c.X]) }
// func (c *Core) CPY_IndirectZPy(zp byte) { c.cmp_impl(c.Y, c.Memory[c.indirectZpY(zp)]) }

func (c *Core) BIT____a(addr uint16) { c.PC += 3; c.bit_impl(c.Memory[addr]) }
func (c *Core) BIT__ZPg(zp byte)     { c.PC += 2; c.bit_impl(c.Memory[zp]) }

// func (c *Core) BIT_Im(literal byte) { c.bit_impl(literal) }
