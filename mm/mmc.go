package mm

import "xubiod/6502-experiment/cpu"

type MemMapper interface {
	SwapCpu(on *cpu.Core) bool
	StepCpu(along *cpu.Core) bool
}
