package cpu

func (c *Core) JMP____a(addr uint16) {
	c.PC = addr
}

func (c *Core) JMP___Ia(addrIndirect uint16) {
	var lsb, msb, page, within byte
	var addrL, addrM uint16

	page = byte(addrIndirect & 0xFF00 >> 8)
	within = byte(addrIndirect & 0x00FF)

	addrL = (uint16(page) << 8) | uint16(within)
	addrM = (uint16(page) << 8) | ((uint16(within) + 1) & 0xFF)

	lsb = c.Memory[addrL]
	msb = c.Memory[addrM]

	c.PC = (uint16(msb) << 8) | uint16(lsb)
}

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

func (c *Core) RTS____i() {
	var high, low byte

	c.S++
	high = c.Memory[0x0100+uint16(c.S)]

	c.S++
	low = c.Memory[0x0100+uint16(c.S)]

	var addr = (uint16(high) << 8) | uint16(low) + 1

	c.PC = addr
}

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
