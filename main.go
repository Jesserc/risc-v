package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	fmt.Println("RISC-V CPU Emulator\n")

	cpu := NewCPU()

	// our program will perform the following operations:
	// Load immediate, add, subtract, store to memory
	// (LUI, ADDI, ADD, SUB, SW)

	// these are already-encoded machine code instructions (not assembly)
	// each instruction is a 32-bit word with fields packed together
	//
	// example breakdown of 0x12345537 (lui a0, 0x12345):
	//   binary: 00010010001101000101_01010_0110111
	//           |--- imm[31:12] ---||-rd-||opcode|
	//   [31:12] imm    = 0x12345 (immediate value)
	//   [11:7]  rd     = 10 (a0 register)
	//   [6:0]   opcode = 0x37 (LUI instruction)
	//
	// how the hex is formed from the binary:
	//   imm << 12:    0x12345 << 12 = 0x12345000 // shift left by 12 bits
	//   rd << 7:      10 << 7        = 0x00000500 // shift left by 7 bits
	//   opcode:       0x37           = 0x00000037 // opcode is already in the correct position
	//   combined (OR):                 0x12345537 // OR the values together to get the final instruction
	//
	// we write them in big-endian hex for readability, then convert to
	// little-endian bytes before loading into memory (risc-v spec)
	instructions := []uint32{
		0x12345537, // lui  a0, 0x12345
		0x02A00593, // addi a1, zero, 42
		0x00B50633, // add  a2, a0, a1
		0x40B606B3, // sub  a3, a2, a1
		0x00C12023, // sw   a2, 0(sp)
	}

	fmt.Println("Loading program...")
	for i, instr := range instructions {
		fmt.Printf("[%d] 0x%08X\n", i, instr)
	}

	// convert to little-endian bytes and load (risc-v is little-endian)
	program := make([]byte, len(instructions)*4) // times 4 because each instruction is 4 bytes
	for i, instr := range instructions {
		binary.LittleEndian.PutUint32(program[i*4:], instr)
	}
	cpu.LoadProgram(program)

	fmt.Println("\nExecuting...\n")

	for i := range instructions {
		fmt.Printf("Step %d: PC=0x%04X\n", i+1, cpu.PC)

		instr, err := cpu.FetchAndDecode()
		if err != nil {
			fmt.Printf("Error fetching instruction: %v\n", err)
			return
		}

		fmt.Printf("  Instruction: 0x%08X\n", instr)

		err = cpu.Execute(instr)
		if err != nil {
			fmt.Printf("Error executing instruction: %v\n", err)
			return
		}

		fmt.Printf("  a0=%08X a1=%08X a2=%08X a3=%08X\n",
			cpu.Regs[A0], cpu.Regs[A1], cpu.Regs[A2], cpu.Regs[A3])
		fmt.Println()
	}

	fmt.Println("\nFinal state:")
	// we only use the a0-a3 (argument) registers in this program.
	// display the values in the a0-a3 registers in 4 bytes hex and decimal
	fmt.Printf("a0 = %08X (%d)\n", cpu.Regs[A0], cpu.Regs[A0])
	fmt.Printf("a1 = %08X (%d)\n", cpu.Regs[A1], cpu.Regs[A1])
	fmt.Printf("a2 = %08X (%d)\n", cpu.Regs[A2], cpu.Regs[A2])
	fmt.Printf("a3 = %08X (%d)\n", cpu.Regs[A3], cpu.Regs[A3])

	// verify memory write
	storedValue := binary.LittleEndian.Uint32(cpu.Memory[cpu.Regs[SP] : cpu.Regs[SP]+4])
	fmt.Printf("\nMemory[sp] = %08X\n", storedValue)
	if storedValue == cpu.Regs[A2] {
		fmt.Println("Memory write verified...")
	}
}
