package cpu

func (c *Core) asl_impl(loc *byte) {
	var shouldCarry = *loc & 0x80

	*loc = (*loc << 1) & 0xFE

	if shouldCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

func (c *Core) lsr_impl(loc *byte) {
	var shouldCarry = *loc & 0x01

	*loc = (*loc >> 1) & 0x7F

	if shouldCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

func (c *Core) rol_impl(loc *byte) {
	var futureCarry = *loc & 0b10000000

	*loc = ((*loc << 1) & 0xFE) | (c.Flags & FLAG_CARRY)

	if futureCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

func (c *Core) ror_impl(loc *byte) {
	var futureCarry = *loc & 0b00000001

	*loc = ((*loc >> 1) & 0x7F) | ((c.Flags & FLAG_CARRY) << 7)

	if futureCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if *loc&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if *loc == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}
}

func (c *Core) ASL____a(addr uint16) { c.PC += 3; c.asl_impl(&c.Memory[addr]) }
func (c *Core) ASL___ax(addr uint16) { c.PC += 3; c.asl_impl(&c.Memory[addr+uint16(c.X)]) }
func (c *Core) ASL____A()            { c.PC += 1; c.asl_impl(&c.A) }
func (c *Core) ASL__ZPg(zp byte)     { c.PC += 2; c.asl_impl(&c.Memory[zp]) }
func (c *Core) ASL__ZPx(zp byte)     { c.PC += 2; c.asl_impl(&c.Memory[zp+c.X]) }

func (c *Core) LSR____a(addr uint16) { c.PC += 3; c.lsr_impl(&c.Memory[addr]) }
func (c *Core) LSR___ax(addr uint16) { c.PC += 3; c.lsr_impl(&c.Memory[addr+uint16(c.X)]) }
func (c *Core) LSR____A()            { c.PC += 1; c.lsr_impl(&c.A) }
func (c *Core) LSR__ZPg(zp byte)     { c.PC += 2; c.lsr_impl(&c.Memory[zp]) }
func (c *Core) LSR__ZPx(zp byte)     { c.PC += 2; c.lsr_impl(&c.Memory[zp+c.X]) }

func (c *Core) ROL____a(addr uint16) { c.PC += 3; c.rol_impl(&c.Memory[addr]) }
func (c *Core) ROL___ax(addr uint16) { c.PC += 3; c.rol_impl(&c.Memory[addr+uint16(c.X)]) }
func (c *Core) ROL____A()            { c.PC += 1; c.rol_impl(&c.A) }
func (c *Core) ROL__ZPg(zp byte)     { c.PC += 2; c.rol_impl(&c.Memory[zp]) }
func (c *Core) ROL__ZPx(zp byte)     { c.PC += 2; c.rol_impl(&c.Memory[zp+c.X]) }

func (c *Core) ROR____a(addr uint16) { c.PC += 3; c.ror_impl(&c.Memory[addr]) }
func (c *Core) ROR___ax(addr uint16) { c.PC += 3; c.ror_impl(&c.Memory[addr+uint16(c.X)]) }
func (c *Core) ROR____A()            { c.PC += 1; c.ror_impl(&c.A) }
func (c *Core) ROR__ZPg(zp byte)     { c.PC += 2; c.ror_impl(&c.Memory[zp]) }
func (c *Core) ROR__ZPx(zp byte)     { c.PC += 2; c.ror_impl(&c.Memory[zp+c.X]) }
