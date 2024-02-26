package cpu

func (c *Core) inc_impl(where *byte) {
	*where++

	r := *where

	if r == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if r&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) dec_impl(where *byte) {
	*where--

	r := *where

	if r == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	if r&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}
}

func (c *Core) INC____a(addr uint16) { c.PC += 3; c.inc_impl(&c.Memory[addr]) }             // `EE` - INC a
func (c *Core) INC___ax(addr uint16) { c.PC += 3; c.inc_impl(&c.Memory[addr+uint16(c.X)]) } // `FE` - INC a, x
func (c *Core) INC__ZPg(zp byte)     { c.PC += 2; c.inc_impl(&c.Memory[zp]) }               // `E6` - INC zp
func (c *Core) INC__ZPx(zp byte)     { c.PC += 2; c.inc_impl(&c.Memory[zp+c.X]) }           // `F6` - INC zp, x
func (c *Core) INX____i()            { c.PC += 1; c.inc_impl(&c.X) }                        // `E8` - INX
func (c *Core) INY____i()            { c.PC += 1; c.inc_impl(&c.Y) }                        // `C8` - INY

func (c *Core) DEC____a(addr uint16) { c.PC += 3; c.dec_impl(&c.Memory[addr]) }             // `CE` - DEC a
func (c *Core) DEC___ax(addr uint16) { c.PC += 3; c.dec_impl(&c.Memory[addr+uint16(c.X)]) } // `DE` - DEC a, x
func (c *Core) DEC__ZPg(zp byte)     { c.PC += 2; c.dec_impl(&c.Memory[zp]) }               // `C6` - DEC zp
func (c *Core) DEC__ZPx(zp byte)     { c.PC += 2; c.dec_impl(&c.Memory[zp+c.X]) }           // `D6` - DEC zp, x
func (c *Core) DEX____i()            { c.PC += 1; c.dec_impl(&c.X) }                        // `CA` - DEX
func (c *Core) DEY____i()            { c.PC += 1; c.dec_impl(&c.Y) }                        // `88` - DEY
