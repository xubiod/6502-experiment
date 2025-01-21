# 6502-experiment

An experimental 6502 and 6502-compatible emulator that was made as a personal
challenge.

Tests are compared with [visual6502](http://visual6502.org/) results with the same
machine code as an attempt to be accurate.

## Segments

* [cpu](./cpu/) - The main part of the emulation. Throughly documented.
* [assembler](./assembler/) - A basic assembler, mainly for making tests easier.
* [mm](./mm/) - An incomplete part for memory managers. Are not implemented.