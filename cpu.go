package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"slices"
)

type CPU struct {
	Memory   []byte            // memory is an array of bytes
	RegNames []string          // registerNames is an array of risc-v register names
	Regs     [32]uint32        // registers is an array of 32-bit words (we use a fixed array to match the exact register count)
	RegMap   map[string]uint32 // registerMap is a map of register names to register numbers (0-31)
	PC       int               // program counter
}

func NewCPU() CPU {
	cpu := CPU{
		Memory:   make([]byte, 65536),
		RegNames: []string{"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2", "s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6"},
		RegMap:   make(map[string]uint32),
		PC:       0,
	}

	// populate registerMap
	for i := 0; i < len(cpu.RegNames); i++ {
		cpu.RegMap[cpu.RegNames[i]] = uint32(i)
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
	opcode := instr & 0x7F // mask out all but the lowest 7 bits to get the opcode (as explained above)

	// Decode and execute based on opcode
	switch opcode {
	case 0x0:
		// No-op
		return nil

	case 0x33: // R-type arithmetic: add, sub, etc (they share the same opcode but are differentiated by funct3 and funct7, which tell the operation variant)
		// R-type format: [funct7][rs2][rs1][funct3][rd][opcode]
		// [31:25] funct7 | [24:20] rs2 | [19:15] rs1 | [14:12] funct3 | [11:7] rd | [6:0] opcode
		// 7 bits for funct7, 5 bits for rs2, 5 bits for rs1, 3 bits for funct3, 5 bits for rd, 7 bits for opcode...togther they form a 32-bit word
		//
		// so we can extract the funct7, funct3, rs1, rs2, rd and opcode from the instruction by shifting and masking
		funct3 := (instr >> 12) & 0x7  // shift right by 12 bits and mask out all but the lowest 3 bits to get the funct3 (function code)
		funct7 := (instr >> 25) & 0x7F // shift right by 25 bits and mask out all but the lowest 7 bits to get the funct7 (function code)
		rs1 := (instr >> 15) & 0x1F    // shift right by 15 bits and mask out all but the lowest 5 bits to get the rs1 (source register 1)
		rs2 := (instr >> 20) & 0x1F    // shift right by 20 bits and mask out all but the lowest 5 bits to get the rs2 (source register 2)
		rd := (instr >> 7) & 0x1F      // shift right by 7 bits and mask out all but the lowest 5 bits to get the rd (destination register)

		switch funct3 {
		case 0x0:
			switch funct7 {
			case 0x00:
				return cpu.executeAdd(rs1, rs2, rd)
			case 0x20:
				return cpu.executeSub(rs1, rs2, rd)
			}
		}
		return errors.New("unimplemented R-type instruction variant")

	case ADDI:
		// I-type format: [imm[11:0]][rs1][funct3][rd][opcode]
		// [31:20] imm[11:0] | [19:15] rs1 | [14:12] funct3 | [11:7] rd | [6:0] opcode
		return cpu.executeAddi(instr)

	case SW:
		// S-type format: [imm[11:5]][rs2][rs1][funct3][imm[4:0]][opcode]
		// [31:25] imm[11:5] | [24:20] rs2 | [19:15] rs1 | [14:12] funct3 | [11:7] imm[4:0] | [6:0] opcode
		return cpu.executeSw(instr)

	case LUI:
		// U-type format: [imm[31:12]][rd][opcode]
		// [31:12] imm[31:12] | [11:7] rd | [6:0] opcode
		return cpu.executeLui(instr)

	default:
		return errors.New("invalid instruction")
	}
}

// ============================================================================
// Fetch-Decode-Execute Cycle
// ============================================================================
func (cpu *CPU) Run() {
	for {
		instr, err := cpu.FetchAndDecode()
		if err != nil {
			fmt.Println(err)
			return
		}
		cpu.Execute(instr)
	}
}

// ============================================================================
// Instruction implementations
// ============================================================================

// ADD
func (cpu *CPU) executeAdd(rs1 uint32, rs2 uint32, rd uint32) error {
	// add the value of rs2 to rs1 and store in rd
	val1 := cpu.Regs[rs1]
	val2 := cpu.Regs[rs2]

	// store the result in the destination register
	cpu.Regs[rd] = val1 + val2

	return nil
}

// SUB
func (cpu *CPU) executeSub(rs1 uint32, rs2 uint32, rd uint32) error {
	// subtract the value of rs2 from rs1 and store in rd
	val1 := cpu.Regs[rs1]
	val2 := cpu.Regs[rs2]

	cpu.Regs[rd] = val1 - val2

	return nil
}

// ADDI (add immediate - adds a 12-bit immediate value to a register)
func (cpu *CPU) executeAddi(instr uint32) error {
	panic("unimplemented")
}

// SW (store word - stores a 32-bit value from a register into memory)
func (cpu *CPU) executeSw(instr uint32) error {
	panic("unimplemented")
}

// LUI (load upper immediate - loads a 20-bit value into the upper 20 bits of a register)
func (cpu *CPU) executeLui(instr uint32) error {
	panic("unimplemented")
}

/*
Note:
we'll use regMap, GetRegisterValue and SetRegisterValue for testing purposes.
so we can easily populate values to regs and test them without going through lui/addi
*/
