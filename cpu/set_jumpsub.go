package cpu

// Jump - Absolute
func (c *Core) JMP____a(addr uint16) {
	c.PC = addr
}

// Jump - Indirect Absolute
func (c *Core) JMP___Ia(addrIndirect uint16) {
	var lsb, msb, page, within byte
	var addrL, addrM uint16

	if c.Features.NMOSAbsoluteIndirectBug {
		page = byte(addrIndirect & 0xFF00 >> 8)
		within = byte(addrIndirect & 0x00FF)

		addrL = (uint16(page) << 8) | uint16(within)
		addrM = (uint16(page) << 8) | ((uint16(within) + 1) & 0xFF)

		lsb = c.Memory[addrL]
		msb = c.Memory[addrM]
	} else {
		lsb = c.Memory[addrIndirect]
		msb = c.Memory[addrIndirect+1]
	}

	c.PC = (uint16(msb) << 8) | uint16(lsb)

}

// Jump + Push Return Address to Stack - Absolute
func (c *Core) JSR____a(addr uint16) {
	var nextInstr = c.PC + 2

	var high, low byte

	high = byte(nextInstr & 0xFF00 >> 8)
	low = byte(nextInstr & 0x00FF)

	c.Memory[0x0100+uint16(c.S)] = high
	c.S--

	c.Memory[0x0100+uint16(c.S)] = low
	c.S--

	c.PC = addr
}

// Return from Subroutine - Implied
func (c *Core) RTS____i() {
	var high, low byte

	c.S++
	high = c.Memory[0x0100+uint16(c.S)]

	c.S++
	low = c.Memory[0x0100+uint16(c.S)]

	var addr = (uint16(high) << 8) | uint16(low) + 1

	c.PC = addr
}

// Return from Interrupt - Implied
func (c *Core) RTI____i() {
	var highPC, lowPC, flags byte

	c.S++
	flags = c.Memory[0x0100+uint16(c.S)]

	c.S++
	highPC = c.Memory[0x0100+uint16(c.S)]

	c.S++
	lowPC = c.Memory[0x0100+uint16(c.S)]

	var addr = (uint16(highPC) << 8) | uint16(lowPC) + 1

	c.Flags = flags

	c.PC = addr
}

// Jump - Indirect Absolute Indexed
func (c *Core) JMP__Iax(addrIndirect uint16) {
	panic("unimplemented")
	// var lsb, msb, page, within byte
	// var addrL, addrM uint16

	// if c.Features.NMOSAbsoluteIndirectBug {
	// 	page = byte(addrIndirect & 0xFF00 >> 8)
	// 	within = byte(addrIndirect & 0x00FF)

	// 	addrL = (uint16(page) << 8) | uint16(within)
	// 	addrM = (uint16(page) << 8) | ((uint16(within) + 1) & 0xFF)

	// 	lsb = c.Memory[addrL]
	// 	msb = c.Memory[addrM]
	// } else {
	// 	lsb = c.Memory[addrIndirect]
	// 	msb = c.Memory[addrIndirect+1]
	// }

	// c.PC = (uint16(msb) << 8) | uint16(lsb)
}
