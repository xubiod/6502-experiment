package cpu

import (
	"fmt"
	"runtime"
)

// Break - Implied
func (c *Core) BRK____i() {
	c.PC += 2

	var high, low byte

	high = byte(c.PC & 0xFF00 >> 8)
	low = byte(c.PC & 0x00FF)

	c.Memory[0x0100+uint16(c.S)] = low
	c.S--

	c.Memory[0x0100+uint16(c.S)] = high
	c.S--

	c.Memory[0x0100+uint16(c.S)] = c.Flags
	c.S--

	c.Flags = c.Flags | FLAG_BREAK

	c.PC = (uint16(c.Memory[0xFFFF]) << 8) | uint16(c.Memory[0xFFFE])

	if c.Features.ConsoleOutOnBreak {
		fmt.Println("Break!\n\n" + c.CompleteDump(runtime.GOOS != "windows") + "\n")
	}
}

// No Operation - Implied
func (c *Core) NOP____i() {
	c.PC += 1
}
