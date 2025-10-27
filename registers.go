package main

const (
	ZERO = iota // zero register
	RA   = iota // return address
	SP   = iota // stack pointer - set, but not used
	GP   = iota // global pointer - not used
	TP   = iota // thread pointer - not used
	T0   = iota // temporary register (t0-t6)
	T1   = iota
	T2   = iota
	S0   = iota // saved register (s0-s11)
	S1   = iota
	A0   = iota // argument register (a0-a7)
	A1   = iota
	A2   = iota
	A3   = iota
	A4   = iota
	A5   = iota
	A6   = iota
	A7   = iota
	S2   = iota
	S3   = iota
	S4   = iota
	S5   = iota
	S6   = iota
	S7   = iota
	S8   = iota
	S9   = iota
	S10  = iota
	S11  = iota
	T3   = iota
	T4   = iota
	T5   = iota
	T6   = iota
)

/*
Notes:
t0-t6 are scratch registers and can be used for any purpose by the program
s0-s11 are saved registers and are used for local variables that must persist across function calls. The callee must restore them before returning (put back the original value)
a0-a7 are argument registers and are used for function arguments
x0 (zero register) is always 0
*/
