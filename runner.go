package experiment

import (
	"errors"
	"xubiod/6502-experiment/cpu"
	"xubiod/6502-experiment/mm"
)

// A Runner has a CPU Core (cpu.Core) and a memory mapper (an implementation of
// mm.MemMapper) and runs them together in such a way to ensure they are synced
// together.
type Runner struct {
	CPU       *cpu.Core     // The CPU core to use
	MemMapper *mm.MemMapper // The memory manager implementation to use
}

func New(cpu *cpu.Core, mm *mm.MemMapper) (*Runner, error) {
	if cpu == nil {
		return nil, errors.New("core cannot be nil")
	}
	return &Runner{CPU: cpu, MemMapper: mm}, nil
}

func (r *Runner) StepOnce() (valid bool) {
	if r.CPU != nil {
		valid, _, _ = r.CPU.StepOnce()
	}
	if valid && r.MemMapper != nil {
		valid = valid && (*r.MemMapper).StepCpu(r.CPU)
		valid = valid && (*r.MemMapper).SwapCpu(r.CPU)
	}
	return
}
