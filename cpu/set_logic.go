package cpu

func (c *Core) and_impl(with byte) {
	c.A = c.A & with

	if c.A == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.A&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) ora_impl(with byte) {
	c.A = c.A | with

	if c.A == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.A&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) eor_impl(with byte) {
	c.A = c.A ^ with

	if c.A == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if c.A&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) AND____a(addr uint16)  { c.PC += 3; c.and_impl(c.Memory[addr]) }
func (c *Core) AND___ax(addr uint16)  { c.PC += 3; c.and_impl(c.Memory[addr+uint16(c.X)]) }
func (c *Core) AND___ay(addr uint16)  { c.PC += 3; c.and_impl(c.Memory[addr+uint16(c.Y)]) }
func (c *Core) AND__Imm(literal byte) { c.PC += 2; c.and_impl(literal) }
func (c *Core) AND__ZPg(zp byte)      { c.PC += 2; c.and_impl(c.Memory[zp]) }
func (c *Core) AND_IZPx(zp byte)      { c.PC += 2; c.and_impl(c.Memory[c.indirectZpX(zp)]) }
func (c *Core) AND__ZPx(zp byte)      { c.PC += 2; c.and_impl(c.Memory[zp+c.X]) }
func (c *Core) AND_IZPy(zp byte)      { c.PC += 2; c.and_impl(c.Memory[c.indirectZpY(zp)]) }

func (c *Core) ORA____a(addr uint16)  { c.PC += 3; c.ora_impl(c.Memory[addr]) }
func (c *Core) ORA___ax(addr uint16)  { c.PC += 3; c.ora_impl(c.Memory[addr+uint16(c.X)]) }
func (c *Core) ORA___ay(addr uint16)  { c.PC += 3; c.ora_impl(c.Memory[addr+uint16(c.Y)]) }
func (c *Core) ORA__Imm(literal byte) { c.PC += 2; c.ora_impl(literal) }
func (c *Core) ORA__ZPg(zp byte)      { c.PC += 2; c.ora_impl(c.Memory[zp]) }
func (c *Core) ORA_IZPx(zp byte)      { c.PC += 2; c.ora_impl(c.Memory[c.indirectZpX(zp)]) }
func (c *Core) ORA__ZPx(zp byte)      { c.PC += 2; c.ora_impl(c.Memory[zp+c.X]) }
func (c *Core) ORA_IZPy(zp byte)      { c.PC += 2; c.ora_impl(c.Memory[c.indirectZpY(zp)]) }

func (c *Core) EOR____a(addr uint16)  { c.PC += 3; c.eor_impl(c.Memory[addr]) }
func (c *Core) EOR___ax(addr uint16)  { c.PC += 3; c.eor_impl(c.Memory[addr+uint16(c.X)]) }
func (c *Core) EOR___ay(addr uint16)  { c.PC += 3; c.eor_impl(c.Memory[addr+uint16(c.Y)]) }
func (c *Core) EOR__Imm(literal byte) { c.PC += 2; c.eor_impl(literal) }
func (c *Core) EOR__ZPg(zp byte)      { c.PC += 2; c.eor_impl(c.Memory[zp]) }
func (c *Core) EOR_IZPx(zp byte)      { c.PC += 2; c.eor_impl(c.Memory[c.indirectZpX(zp)]) }
func (c *Core) EOR__ZPx(zp byte)      { c.PC += 2; c.eor_impl(c.Memory[zp+c.X]) }
func (c *Core) EOR_IZPy(zp byte)      { c.PC += 2; c.eor_impl(c.Memory[c.indirectZpY(zp)]) }
