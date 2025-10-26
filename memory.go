package main

type memory []uint32

func createMemory(size uint32) memory {
	return make(memory, size)
}
