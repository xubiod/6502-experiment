package mm

import "xubiod/6502-experiment/cpu"

// https://www.nesdev.org/wiki/NROM
type MemMapperNROM128 struct {
	PrgRam  []byte
	PrgRom0 [0x4000]byte
}

func (m *MemMapperNROM128) SwapCpu(on *cpu.Core) bool {
	copy(on.Memory[0x8000:0xC000], m.PrgRom0[:])
	copy(on.Memory[0xC000:], m.PrgRom0[:])
	return true
}

func (*MemMapperNROM128) StepCpu(along *cpu.Core) bool { return true }

// https://www.nesdev.org/wiki/NROM
type MemMapperNROM256 struct {
	PrgRam  []byte
	PrgRom0 [0x4000]byte
	PrgRom1 [0x4000]byte
}

func (m *MemMapperNROM256) SwapCpu(on *cpu.Core) bool {
	copy(on.Memory[0x8000:0xC000], m.PrgRom0[:])
	copy(on.Memory[0xC000:], m.PrgRom1[:])
	return true
}

func (*MemMapperNROM256) StepCpu(along *cpu.Core) bool { return true }
