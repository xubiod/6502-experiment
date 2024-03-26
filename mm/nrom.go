package mm

import "xubiod/6502-experiment/cpu"

// https://www.nesdev.org/wiki/NROM
type MemMapperNROM128 struct {
	PrgRam  []byte
	PrgRom0 [0x4000]byte

	ChrRom0 [0x2000]byte
}

func (m *MemMapperNROM128) SwapCpu(on *cpu.Core) bool {
	n := 0
	n += copy(on.Memory[0x8000:0xC000], m.PrgRom0[:])
	n += copy(on.Memory[0xC000:], m.PrgRom0[:])
	return (n == 0x8000)
}

func (*MemMapperNROM128) StepCpu(along *cpu.Core) bool { return true }

// https://www.nesdev.org/wiki/NROM
type MemMapperNROM256 struct {
	PrgRam  []byte
	PrgRom0 [0x4000]byte
	PrgRom1 [0x4000]byte

	ChrRom0 [0x2000]byte
}

func (m *MemMapperNROM256) SwapCpu(on *cpu.Core) bool {
	n := 0
	n += copy(on.Memory[0x8000:0xC000], m.PrgRom0[:])
	n += copy(on.Memory[0xC000:], m.PrgRom1[:])
	return (n == 0x8000)
}

func (*MemMapperNROM256) StepCpu(along *cpu.Core) bool { return true }
