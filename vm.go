package main

import (
	"io"
	"math/big"
	"os"
)

type VirtualMachine struct {
	program        *CompiledProgram
	memory         []interface{}
	memoryOffset   int
	previousOffset int
	out            io.Writer
}

func NewVirtualMachine(program *CompiledProgram) *VirtualMachine {
	return &VirtualMachine{
		program: program,
		out:     os.Stdout,
	}
}

func (vm *VirtualMachine) Run() {
	// TODO: Check start exists.
	vm.call("start", nil)
}

func (vm *VirtualMachine) call(syntax string, args []int) {
	fn := vm.program.Functions[syntax]

	if fn == nil {
		// TODO: Remove me
		panic(syntax)
	}

	// Expand the memory to accommodate the call.
	vm.memory = append(vm.memory, fn.Variables...)

	// Load in the arguments.
	for i, arg := range args {
		to := vm.memoryOffset + i
		from := vm.previousOffset + arg
		vm.memory[to] = vm.memory[from]
	}

	vm.previousOffset = vm.memoryOffset
	vm.memoryOffset += len(fn.Variables)

	for _, instruction := range fn.Instructions {
		// Check if it is a system call?
		if handler, ok := System[instruction.Call]; ok {
			handler(vm, instruction.Args)
			continue
		}

		// Otherwise we have to increase the stack.
		vm.call(instruction.Call, instruction.Args)
	}
}

func (vm *VirtualMachine) GetArg(index int) interface{} {
	return vm.memory[vm.previousOffset+index]
}

func (vm *VirtualMachine) SetArg(index int, value interface{}) {
	vm.memory[vm.previousOffset+index] = value
}

func (vm *VirtualMachine) GetNumber(index int) *big.Rat {
	return vm.memory[vm.previousOffset+index].(*big.Rat)
}
