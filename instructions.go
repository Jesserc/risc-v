package main

// We'll define just the add, sub, addi, sw, and lui instructions for now

const (
	// add and sub share the same opcode (they are both R-type instructions, but the funct3 field differentiates them)
	ADD = 0x33
	SUB = 0x33

	ADDI = 0x13
	SW   = 0x23
	LUI  = 0x37
)
