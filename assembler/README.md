# assembler

- [assembler](#assembler)
  - [Process](#process)
  - [Design](#design)
  - [Schema](#schema)
    - [Comments](#comments)
    - [Labels](#labels)
    - [Blocks](#blocks)
      - [Remark](#remark)
      - [Text](#text)
      - [Data](#data)
    - [Instructions](#instructions)
      - [Addressing mode priority](#addressing-mode-priority)
  - [Errors](#errors)

## Process

There are 2 passes:

- Preprocessing
  - Discovers labels at their appropriate memory locations
- Assembling
  - Processes line by line
  - Labels are replaced appropriately
  - Bytecode is generated from text blocks
  - Bytes are generated straight from data blocks

This assembler starts as if line 1 is going to be placed into `$0200` which is
the start of the 6502's generic purpose memory.

## Design

**The general principle behind this assembler is no assuming.**

There are no specialized instructions for the assembler to use. The assembler
directly translates the assembly to bytecode, and makes **no assumptions** about
anything. All literals **must be hexadecimal** and be fully complete, with **no
missing nibbles or bytes**.

An absolute address defined as a single byte/two nibbles will be read as a zero
page as that is how zero page addressing is defined, even if absolute addressing
has a [higher priority](#addressing-mode-priority).

There is **no way** to define any number as a base other than 16. `#$10` is 16,
*unless* decimal mode is on in which case `ADC` and `SBC` will do math with binary
coded decimal on the processor. The exception to this would be, for instance, the
NES cpu which, while 6502-compatible, did not implement decimal mode.

This assembler was made in conjunction with a 6502 emulator and as such has been
made with that emulator in mind.

## Schema

See [SCHEMA](./SCHEMA) for a general idea.

### Comments

Comments are single line, starting with a semicolon. Once a comment starts, they
terminate at the end of the line.

```asm
; This is a comment.
```

Comments are ignored by the assembler.

### Labels

Labels are defined with no whitespace preceding them, terminated by a colon:

```asm
LABEL:   ; This is a label.
```

Labels are case-sensitive, word characters (`0-9`, `A-Z`, `a-z`, and `_`). A label
that is called `LABEL`, `label`, and `Label` are all considered as different
labels. Labels will terminate at non-word characters.

Labels **cannot** start as a number.

Labels will always be the address of the current memory location, or how many
bytes away it is if the instruction is relative. **Labels will *always* decay into
2 byte addresses unless in a branch instruction, which it will decay into a single
byte.**

If a branch is more than 129 bytes *ahead* or 126 bytes *behind* a label it is
using, the number will not be able to be represented as a byte. This will error
out the assembler with the current line. This is not the normal range for a signed
byte, but since the program counter is incremented during execution the effective
range is given here in the documentation.

This is an example with `START` being defined first at `$0200`:

```txt
branch cannot reach this label (line 132)
        -> 132 |        BEQ START
```

### Blocks

Blocks are defined with a preceding period, and are case-insensitive.

There are three types of blocks: text, data, and remarks. They are declared as
the following:

```asm
.TEXT
  NOP     ; Instructions

.DATA
00000000  ; Raw bytes

.REMARK
Anything
```

Text blocks can be defined as `.TEXT`, `.TXT`, or `.T`. Text blocks contain only
instructions.

Data blocks can be defined as `.DATA`, `.DAT`, or `.D`.  Data blocks only contain
hexadecimal data.

Remark blocks can be defined as `.REMARK`, `.REM` or `.R`. Remark blocks are
completely ignored.

**The assembler always starts in text processing for instructions.**

#### Remark

Remark blocks are completely ignored, and can be used as a block comment if
desired.

#### Text

See [Instructions](#instructions) for instruction parsing.

#### Data

Data blocks consist of hexadecimal numbers. It is case-insensitive, and all
whitespace is completely ignored.

The spacing between them is ignored as well:

```asm
0010203040506070        ; This is the same as
80 90 A0 B0 C0 D0 E0 F0 ; this.
```

Note that all data in a *line* must be complete bytes with two nibbles each:

```asm
00 00 ; Okay
00 0  ; Not okay

0000  ; Okay
000   ; Not okay
```

This will error and the assembler will stop with no complete data or bytecode.

### Instructions

Instructions must be preceded by whitespace, the only exception is line feeds
(`\n`). *The assembler treats all lines that aren't labels as case-insensitive.*

The following is an arithmetic shift left with no operands, so it will shift the
accumulator left:

```asm
  ASL
```

Bytes are read as hexadecimal numbers only, ranging from `00` to `FF`. Both nibbles
must be present, capitalization of `A` through `F` is irrelevant.

Addresses are read as hexadecimal numbers only, ranging from `0000` to `FFFF`.
All four nibbles must be present, capitalization of `A` through `F` is irrelevant.

If an instruction has operands, the following is how the assembler sees them:

| Assembly  | Addressing Mode                                   |
|-----------|---------------------------------------------------|
| `none`    | Accumulator/Implied (depends on instruction)      |
| `#$xx`    | Immediate, `xx` is a byte                         |
| `$xxxx`   | Absolute, `xxxx` is an address                    |
| `$xx`     | Zero page OR relative branch, `xx` is a byte      |
| `($xxxx)` | Absolute indirect, `xxxx` is an address           |
| `$xxxx,X` | Absolute indexed with X, `xxxx` is an address     |
| `$xxxx,Y` | Absolute indexed with Y, `xxxx` is an address     |
| `$xx,X`   | Zero page indexed with X, `xx` is a byte          |
| `$xx,Y`   | Zero page indexed with Y, `xx` is a byte          |
| `($xx,X)` | Zero page indexed indirect, `xx` is a byte        |
| `($xx),Y` | Zero page indirect indexed with Y, `xx` is a byte |

#### Addressing mode priority

While it should not be an issue, instructions are checked in the given order of
priority, starting at the top and going down:

1. Zero page indirect indexed with Y
2. Zero page indexed indirect
3. Absolute indirect
4. Absolute indexed with Y
5. Absolute indexed with X
6. Zero page indexed with Y
7. Zero page indexed with X
8. Absolute
9. Relative
10. Immediate
11. Zero page
12. Accumulator/implied (no operands)

The thought process was that the most specific addressing modes are checked before
getting more broad.

## Errors

The assembler will error out on invalid instructions, and will not output any
incomplete bytecode. The error should give out a line number for the line with
the invalid instruction, along with the raw line contents to assist with debugging.

Here is an example error from testing the assembler:

```txt
invalid instruction (line 10)
  -> 10 |         WTF     ; Not an instruction
```
