package cpu

func (c *Core) adc_impl(middle byte) {
	if c.Features.DecimalModeImplemented && c.Flags&FLAG_DECIMAL > 0 {
		c.adc_impl_decimal(middle)
		return
	}

	var u1 = uint16(c.A)
	var u2 = uint16(middle)
	var result uint16 = u1 + u2 + uint16(c.Flags&FLAG_CARRY)

	c.A = byte(result & 0xFF)

	if result&0xFF != result {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if result&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if result == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	// the fucking overflow flag

	if (u1^result)&(u2^result)&0x80 != 0 {
		c.Flags = c.Flags | FLAG_OVERFLOW
	} else {
		c.Flags = c.Flags & ^FLAG_OVERFLOW
	}
}

func (c *Core) adc_impl_decimal(middle byte) {
	var u1 = uint16(c.A)
	var u2 = uint16(middle)

	var normResult uint16 = u1 + u2 + uint16(c.Flags&FLAG_CARRY)

	var lo = (u1 & 0x0F) + (u2 & 0x0F) + uint16(c.Flags&FLAG_CARRY)
	var loCarry = 0
	if lo > 0x09 {
		lo += 6
		loCarry = 1
		lo = lo & 0x0F
	}

	var hi = ((u1 & 0xF0) >> 4) + ((u2 & 0xF0) >> 4) + uint16(loCarry)
	var hiCarry = 0
	if hi > 0x09 {
		hi += 6
		hiCarry = 1
		hi = hi & 0x0F
	}

	var result uint16 = (hi << 4) | lo

	c.A = byte(result & 0xFF)

	if hiCarry > 0 {
		c.Flags = c.Flags | FLAG_CARRY
	} else {
		c.Flags = c.Flags & ^FLAG_CARRY
	}

	if normResult&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if normResult == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	// the fucking overflow flag

	if (u1^normResult)&(u2^normResult)&0x80 != 0 {
		c.Flags = c.Flags | FLAG_OVERFLOW
	} else {
		c.Flags = c.Flags & ^FLAG_OVERFLOW
	}
}

func (c *Core) sbc_impl(middle byte) {
	if c.Features.DecimalModeImplemented && c.Flags&FLAG_DECIMAL > 0 {
		c.sbc_impl_decimal(middle)
		return
	}

	var u1 = uint16(c.A)
	var u2 = uint16(middle)
	var result uint16 = (u1 - u2) + uint16(c.Flags&FLAG_CARRY)

	opl := c.A

	c.A = byte(result & 0xFF)

	if result&0xFF != result {
		c.Flags = c.Flags & ^FLAG_CARRY
	} else {
		c.Flags = c.Flags | FLAG_CARRY
	}

	if result&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if result == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	// the fucking overflow flag

	if (opl^middle)&(opl^c.A)&0x80 != 0 {
		c.Flags = c.Flags | FLAG_OVERFLOW
	} else {
		c.Flags = c.Flags & ^FLAG_OVERFLOW
	}
}

func (c *Core) sbc_impl_decimal(middle byte) {
	var u1 = uint16(c.A)
	var u2 = uint16(middle)

	var normResult uint16 = (u1 - u2) + uint16(1-c.Flags&FLAG_CARRY)

	opl := c.A

	var lo = (u1 & 0x0F) - (u2 & 0x0F) - uint16(1-c.Flags&FLAG_CARRY)
	var loBorrow = lo & 0b10000000 >> 7
	lo = lo & 0x0F
	if loBorrow > 0 {
		lo += 0xA
		lo = lo & 0x0F
	}

	var hi = ((u1 & 0xF0) >> 4) - ((u2 & 0xF0) >> 4) - uint16(loBorrow)
	var hiBorrow = hi & 0b10000000 >> 7
	hi = hi & 0x0F
	if hiBorrow > 0 {
		hi += 0xA
		hi = hi & 0x0F
	}

	var result uint16 = (hi << 4) | lo

	c.A = byte(result & 0xFF)

	if hiBorrow > 0 {
		c.Flags = c.Flags & ^FLAG_CARRY
	} else {
		c.Flags = c.Flags | FLAG_CARRY
	}

	if normResult&0b10000000 > 0 {
		c.Flags = c.Flags | FLAG_NEGATIVE
	} else {
		c.Flags = c.Flags & ^FLAG_NEGATIVE
	}

	if normResult == 0 {
		c.Flags = c.Flags | FLAG_ZERO
	} else {
		c.Flags = c.Flags & ^FLAG_ZERO
	}

	// the fucking overflow flag

	if (opl^middle)&(opl^byte(normResult&0xFF))&0x80 != 0 {
		c.Flags = c.Flags | FLAG_OVERFLOW
	} else {
		c.Flags = c.Flags & ^FLAG_OVERFLOW
	}
}

func (c *Core) ADC____a(addr uint16)  { c.PC += 3; c.adc_impl(c.Memory[addr]) }              // `6D` - ADC a
func (c *Core) ADC___ax(addr uint16)  { c.PC += 3; c.adc_impl(c.Memory[addr+uint16(c.X)]) }  // `7D` - ADC a, x
func (c *Core) ADC___ay(addr uint16)  { c.PC += 3; c.adc_impl(c.Memory[addr+uint16(c.Y)]) }  // `79` - ADC a, y
func (c *Core) ADC__Imm(literal byte) { c.PC += 2; c.adc_impl(literal) }                     // `69` - ADC #
func (c *Core) ADC__ZPg(zp byte)      { c.PC += 2; c.adc_impl(c.Memory[zp]) }                // `65` - ADC zp
func (c *Core) ADC_IZPx(zp byte)      { c.PC += 2; c.adc_impl(c.Memory[c.indirectZpX(zp)]) } // `61` - ADC (zp,x)
func (c *Core) ADC__ZPx(zp byte)      { c.PC += 2; c.adc_impl(c.Memory[zp+c.X]) }            // `75` - ADC zp, x
func (c *Core) ADC_IZPy(zp byte)      { c.PC += 2; c.adc_impl(c.Memory[c.indirectZpY(zp)]) } // `71` - ADC (zp), y

func (c *Core) SBC____a(addr uint16)  { c.PC += 3; c.sbc_impl(c.Memory[addr]) }              // `ED` - SDC a
func (c *Core) SBC___ax(addr uint16)  { c.PC += 3; c.sbc_impl(c.Memory[addr+uint16(c.X)]) }  // `FD` - SDC a, x
func (c *Core) SBC___ay(addr uint16)  { c.PC += 3; c.sbc_impl(c.Memory[addr+uint16(c.Y)]) }  // `F9` - SDC a, y
func (c *Core) SBC__Imm(literal byte) { c.PC += 2; c.sbc_impl(literal) }                     // `E9` - SDC #
func (c *Core) SBC__Zpg(zp byte)      { c.PC += 2; c.sbc_impl(c.Memory[zp]) }                // `E5` - SDC zp
func (c *Core) SBC_IZPx(zp byte)      { c.PC += 2; c.sbc_impl(c.Memory[c.indirectZpX(zp)]) } // `E1` - SDC (zp, x)
func (c *Core) SBC__ZPx(zp byte)      { c.PC += 2; c.sbc_impl(c.Memory[zp+c.X]) }            // `F5` - SDC zp, x
func (c *Core) SBC_IZPy(zp byte)      { c.PC += 2; c.sbc_impl(c.Memory[c.indirectZpY(zp)]) } // `F1` - SDC (zp), y
