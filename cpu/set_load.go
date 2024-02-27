package cpu

func (c *Core) ld_impl(to *byte, store byte) {
	*to = store

	if *to == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if *to&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

// Load Memory into Accumulator - Absolute
func (c *Core) LDA____a(addr uint16) { c.PC += 3; c.ld_impl(&c.A, c.Memory[addr]) }

// Load Memory into Accumulator - Absolute indexed with X
func (c *Core) LDA___ax(addr uint16) { c.PC += 3; c.ld_impl(&c.A, c.Memory[addr+uint16(c.X)]) }

// Load Memory into Accumulator - Absolute indexed with Y
func (c *Core) LDA___ay(addr uint16) { c.PC += 3; c.ld_impl(&c.A, c.Memory[addr+uint16(c.Y)]) }

// Load Memory into Accumulator - Immediate
func (c *Core) LDA__Imm(literal byte) { c.PC += 2; c.ld_impl(&c.A, literal) }

// Load Memory into Accumulator - Zero Page
func (c *Core) LDA__ZPg(zp byte) { c.PC += 2; c.ld_impl(&c.A, c.Memory[zp]) }

// Load Memory into Accumulator - Zero Page Indexed Indirect
func (c *Core) LDA_IZPx(zp byte) { c.PC += 2; c.ld_impl(&c.A, c.Memory[c.indirectZpX(zp)]) }

// Load Memory into Accumulator - Zero Page indexed with X
func (c *Core) LDA__ZPx(zp byte) { c.PC += 2; c.ld_impl(&c.A, c.Memory[zp+c.X]) }

// Load Memory into Accumulator - Zero Page Indirect Indexed with Y
func (c *Core) LDA_IZPy(zp byte) { c.PC += 2; c.ld_impl(&c.A, c.Memory[c.indirectZpY(zp)]) }

// Load Memory into X - Absolute
func (c *Core) LDX____a(addr uint16) { c.PC += 3; c.ld_impl(&c.X, c.Memory[addr]) }

// Load Memory into X - Absolute indexed with Y
func (c *Core) LDX___ay(addr uint16) { c.PC += 3; c.ld_impl(&c.X, c.Memory[addr+uint16(c.Y)]) }

// Load Memory into X - Immediate
func (c *Core) LDX__Imm(literal byte) { c.PC += 2; c.ld_impl(&c.X, literal) }

// Load Memory into X - Zero Page
func (c *Core) LDX__ZPg(zp byte) { c.PC += 2; c.ld_impl(&c.X, c.Memory[zp]) }

// Load Memory into X - Zero Page indexed with Y
func (c *Core) LDX__ZPy(zp byte) { c.PC += 2; c.ld_impl(&c.X, c.Memory[zp+c.Y]) }

// Load Memory into Y - Absolute
func (c *Core) LDY____a(addr uint16) { c.PC += 3; c.ld_impl(&c.Y, c.Memory[addr]) }

// Load Memory into Y - Absolute indexed with X
func (c *Core) LDY___ax(addr uint16) { c.PC += 3; c.ld_impl(&c.Y, c.Memory[addr+uint16(c.X)]) }

// Load Memory into Y - Immediate
func (c *Core) LDY__Imm(literal byte) { c.PC += 2; c.ld_impl(&c.Y, literal) }

// Load Memory into Y - Zero Page
func (c *Core) LDY__ZPg(zp byte) { c.PC += 2; c.ld_impl(&c.Y, c.Memory[zp]) }

// Load Memory into Y - Zero Page indexed with X
func (c *Core) LDY__ZPx(zp byte) { c.PC += 2; c.ld_impl(&c.Y, c.Memory[zp+c.X]) }
