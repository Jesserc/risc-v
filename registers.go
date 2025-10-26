package main

const (
	ZERO = 0 // zero register
	RA   = 1 // return address
	SP   = 2 // stack pointer
	GP   = 3 // global pointer
	TP   = 4 // thread pointer
	T0   = 5 // temporary register (t0-t6)
	T1   = 6
	T2   = 7
	S0   = 8 // saved register (s0-s11)
	S1   = 9
	A0   = 10 // argument register (a0-a7)
	A1   = 11
	A2   = 12
	A3   = 13
	A4   = 14
	A5   = 15
	A6   = 16
	A7   = 17
	S2   = 18
	S3   = 19
	S4   = 20
	S5   = 21
	S6   = 22
	S7   = 23
	S8   = 24
	S9   = 25
	S10  = 26
	S11  = 27
	T3   = 28
	T4   = 29
	T5   = 30
	T6   = 31
)

/*
Notes:
t0-t6 are scratch registers and can be used for any purpose by the program
s0-s11 are saved registers and are used for local variables that must persist across function calls. The callee must restore them before returning (put back the original value)
a0-a7 are argument registers and are used for function arguments
x0 (zero register) is always 0
*/
