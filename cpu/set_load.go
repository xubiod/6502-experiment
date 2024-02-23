package cpu

func (c *Core) ld_impl(to *byte, store byte) {
	*to = store

	if *to == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if *to&0b1000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) LDA____a(addr uint16)  { c.PC += 3; c.ld_impl(&c.A, c.Memory[addr]) }              // `AD` - LDA A
func (c *Core) LDA___ax(addr uint16)  { c.PC += 3; c.ld_impl(&c.A, c.Memory[addr+uint16(c.X)]) }  // `BD` - LDA A, x
func (c *Core) LDA___ay(addr uint16)  { c.PC += 3; c.ld_impl(&c.A, c.Memory[addr+uint16(c.Y)]) }  // `B9` - LDA A, y
func (c *Core) LDA__Imm(literal byte) { c.PC += 2; c.ld_impl(&c.A, literal) }                     // `A9` - LDA #
func (c *Core) LDA__ZPg(zp byte)      { c.PC += 2; c.ld_impl(&c.A, c.Memory[zp]) }                // `A5` - LDA zp
func (c *Core) LDA_IZPx(zp byte)      { c.PC += 2; c.ld_impl(&c.A, c.Memory[c.indirectZpX(zp)]) } // `A1` - LDA (zp, x)
func (c *Core) LDA__ZPx(zp byte)      { c.PC += 2; c.ld_impl(&c.A, c.Memory[zp+c.X]) }            // `B5` - LDA zp, x
func (c *Core) LDA_IZPy(zp byte)      { c.PC += 2; c.ld_impl(&c.A, c.Memory[c.indirectZpY(zp)]) } // `B1` - LDA (zp), y

func (c *Core) LDX____a(addr uint16)  { c.PC += 3; c.ld_impl(&c.X, c.Memory[addr]) }             // `AE` - LDX A
func (c *Core) LDX___ay(addr uint16)  { c.PC += 3; c.ld_impl(&c.X, c.Memory[addr+uint16(c.Y)]) } // `BE` - LDX A, Y
func (c *Core) LDX__Imm(literal byte) { c.PC += 2; c.ld_impl(&c.X, literal) }                    // `A2` - LDX #
func (c *Core) LDX__ZPg(zp byte)      { c.PC += 2; c.ld_impl(&c.X, c.Memory[zp]) }               // `A6` - LDX zp
func (c *Core) LDX__ZPy(zp byte)      { c.PC += 2; c.ld_impl(&c.X, c.Memory[zp+c.Y]) }           // `B6` - LDX zp, y

func (c *Core) LDY____a(addr uint16)  { c.PC += 3; c.ld_impl(&c.Y, c.Memory[addr]) }             // `AC` - LDY A
func (c *Core) LDY___ax(addr uint16)  { c.PC += 3; c.ld_impl(&c.Y, c.Memory[addr+uint16(c.X)]) } // `BC` - LDY A, X
func (c *Core) LDY__Imm(literal byte) { c.PC += 2; c.ld_impl(&c.Y, literal) }                    // `A0` - LDY #
func (c *Core) LDY__ZPg(zp byte)      { c.PC += 2; c.ld_impl(&c.Y, c.Memory[zp]) }               // `A4` - LDY zp
func (c *Core) LDY__ZPx(zp byte)      { c.PC += 2; c.ld_impl(&c.Y, c.Memory[zp+c.X]) }           // `B4` - LDY zp, X
