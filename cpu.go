package main

import (
	"encoding/binary"
	"errors"
	"slices"
)

type CPU struct {
	Memory   []byte         // memory is an array of bytes
	RegNames []string       // registerNames is an array of risc-v register names
	Regs     []uint32       // registers is an array of 32-bit words
	RegMap   map[string]int // registerMap is a map of register names to register numbers (0-31)
	PC       uint32         // program counter
}

func NewCPU() CPU {
	cpu := CPU{
		Memory:   make([]byte, 65536),
		RegNames: []string{"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2", "s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6"},
		Regs:     make([]uint32, 32),
		RegMap:   make(map[string]int),
		PC:       0,
	}

	// populate registerMap
	for i := 0; i < len(cpu.RegNames); i++ {
		cpu.RegMap[cpu.RegNames[i]] = i
	}

	return cpu
}

func (cpu *CPU) LoadProgram(program []byte) {
	copy(cpu.Memory, program)
}

// SetRegisterValue sets the value of a register
func (cpu *CPU) SetRegisterValue(register string, value uint32) error {
	if slices.Contains(cpu.RegNames, register) {
		cpu.Regs[cpu.RegMap[register]] = value
		return nil
	}
	return errors.New("register not found")
}

// GetRegisterValue gets the value of a register
func (cpu *CPU) GetRegisterValue(register string) (uint32, error) {
	if slices.Contains(cpu.RegNames, register) {
		return cpu.Regs[cpu.RegMap[register]], nil
	}
	return 0, errors.New("register not found")
}

func (cpu *CPU) FetchAndDecode() (instr uint32, err error) {
	// fetch instruction from memory
	// and convert it to a 32-bit word
	instr = binary.LittleEndian.Uint32(cpu.Memory[cpu.PC : cpu.PC+4]) // instr will now be ordered this way: [funct7][rs2][rs1][funct3][rd][opcode]

	cpu.PC += 4 // increment program counter by 4 bytes (32 bits) because each instruction is 4 bytes
	return instr, nil
}

func (cpu *CPU) Execute(instr uint32) error {
	// the bit masking (instr & 0x7F) extracts the opcode from the instruction
	// the bitwise AND op with 0x7F masks out all but the lowest 7 bits of the instruction, which is the opcode
	// 7 bits because `0x7F` is `0111 1111` - 7 ones and 1 zero
	// and remember an AND op is a binary operation that takes two operands and returns 1 if both are 1, otherwise 0
	// so our `instr` which is a 32-bit word (4 bytes) will be masked to only the lowest 7 bits
	// (according to the risc-v specs, the opcode takes up 7 bits)

	var (
		// risc-v instructions are ordered this way: [funct7][rs2][rs1][funct3][rd][opcode]...togther they form a 32-bit word
		// [31:25] funct7 | [24:20] rs2 | [19:15] rs1 | [14:12] funct3 | [11:7] rd | [6:0] opcode
		// 7 bits for funct7, 5 bits for rs2, 5 bits for rs1, 3 bits for funct3, 5 bits for rd, 7 bits for opcode
		// so we can extract the funct7, funct3, rs1, rs2, rd and opcode from the instruction by shifting and masking

		opcode = instr & 0x7F         // mask out all but the lowest 7 bits to get the opcode (as explained above)
		funct3 = (instr >> 12) & 0x7  // shift right by 12 bits and mask out all but the lowest 3 bits to get the funct3 (function code)
		funct7 = (instr >> 25) & 0x7F // shift right by 25 bits and mask out all but the lowest 7 bits to get the funct7 (function code)

		rs1 = (instr >> 15) & 0x1F // shift right by 15 bits and mask out all but the lowest 5 bits to get the rs1 (source register 1)
		rs2 = (instr >> 20) & 0x1F // shift right by 20 bits and mask out all but the lowest 5 bits to get the rs2 (source register 2)
		rd  = (instr >> 7) & 0x1F  // shift right by 7 bits and mask out all but the lowest 5 bits to get the rd (destination register)
	)

	// use a switch statement to determine which instruction to execute based on the opcode
	switch opcode {
	case 0x0:
	case 0x33: // R-type arithmetic: add, sub, etc (they share the same opcode but are differentiated by funct3 and funct7, which tell the operation variant)
		switch funct3 {
		case 0x0:
			switch funct7 {
			case 0x00:
				// add
				return cpu.executeAdd(instr, rs1, rs2, rd)
			case 0x20:
				// sub (unimplemented)
				return errors.New("sub is unimplemented yet")
			}
		}

	case ADDI:
		return cpu.executeAddi(instr)
	case LW:
		return cpu.executeLw(instr)
	case SW:
		return cpu.executeSw(instr)
	case BEQ:
		return cpu.executeBeq(instr)
	default:
		return errors.New("invalid instruction")
	}

	return errors.New("unimplemented instruction variant")
}

// ... follow a fetch-decode-execute cycle ...//
func (cpu *CPU) Run() {
	for {
		instr, err := cpu.FetchAndDecode()
		if err != nil {
			return
		}
		cpu.Execute(instr)
	}
}

// Instruction implementations

// ADD
func (cpu *CPU) executeAdd(instr uint32, rs1 uint32, rs2 uint32, rd uint32) error {

	// extract the rs1, rs2, and rd from the instruction
	// rs1 is bits [19:15]

	// rs2 is bits [24:20]

	// rd is bits [11:7]

	return nil
}

// ADDI (add immediate - adds a 12-bit immediate value to a register)
func (cpu *CPU) executeAddi(instr uint32) error {
	panic("unimplemented")
}

// LW (load word - loads a 32-bit value from memory into a register)
func (cpu *CPU) executeLw(instr uint32) error {
	panic("unimplemented")
}

// SW (store word - stores a 32-bit value from a register into memory)
func (cpu *CPU) executeSw(instr uint32) error {
	panic("unimplemented")
}

// BEQ (branch if equal - branches PC if two registers are equal)
func (cpu *CPU) executeBeq(instr uint32) error {
	panic("unimplemented")
}
