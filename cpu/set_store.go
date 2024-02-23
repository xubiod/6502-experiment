package cpu

func (c *Core) STA____a(addr uint16) { c.PC += 3; c.Memory[addr] = c.A }              // `8D` - STA a
func (c *Core) STA___ax(addr uint16) { c.PC += 3; c.Memory[addr+uint16(c.X)] = c.A }  // `9D` - STA A, x
func (c *Core) STA___ay(addr uint16) { c.PC += 3; c.Memory[addr+uint16(c.Y)] = c.A }  // `99` - STA A, y
func (c *Core) STA__ZPg(zp byte)     { c.PC += 2; c.Memory[zp] = c.A }                // `85` - STA zp
func (c *Core) STA_IZPx(zp byte)     { c.PC += 2; c.Memory[c.indirectZpX(zp)] = c.A } // `81` - STA (zp, x)
func (c *Core) STA__ZPx(zp byte)     { c.PC += 2; c.Memory[zp+c.X] = c.A }            // `95` - STA zp, x
func (c *Core) STA_IZPy(zp byte)     { c.PC += 2; c.Memory[c.indirectZpY(zp)] = c.A } // `91` - STA (zp), y

func (c *Core) STX____a(addr uint16) { c.PC += 3; c.Memory[addr] = c.X }   // `8E` - STX a
func (c *Core) STX__ZPg(zp byte)     { c.PC += 2; c.Memory[zp] = c.X }     // `86` - STX zp
func (c *Core) STX__ZPy(zp byte)     { c.PC += 2; c.Memory[zp+c.Y] = c.X } // `96` - STX zp, y

func (c *Core) STY____a(addr uint16) { c.PC += 3; c.Memory[addr] = c.Y }   // `8C` - STY a
func (c *Core) STY__ZPg(zp byte)     { c.PC += 2; c.Memory[zp] = c.Y }     // `84` - STY zp
func (c *Core) STY__ZPx(zp byte)     { c.PC += 2; c.Memory[zp+c.X] = c.Y } // `94` - STY zp, x
