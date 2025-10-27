package main

import (
	"encoding/binary"
	"errors"
	"slices"
)
const (
	testNum uint32 = 0x12345678
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
	instr = binary.LittleEndian.Uint32(cpu.Memory[cpu.PC : cpu.PC+4]) // this includes the opcode, rs1, rs2, and rd (opcode and operands)

	cpu.PC += 4 // increment program counter by 4 bytes (32 bits) because each instruction is 4 bytes
	return instr, nil
}

func (cpu *CPU) Execute(instr uint32) error {
	// use a switch statement to determine which instruction to execute based on the opcode
	// the bit masking (instr & 0x7F) extracts the opcode from the instruction
	switch instr & 0x7F {
	case ADD:
		return cpu.executeAdd(instr)
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
}

// ... follow a fetch-decode-execute cycle ...//
func (cpu *CPU) Run() {

}

// Instruction implementations

// ADD
func (cpu *CPU) executeAdd(instr uint32) error {
	panic("unimplemented")
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
