package main

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"strings"
)

type Instruction interface{}

type ConditionJumpInstruction struct {
	Left, Right int
	Operator    string
	True, False int
}

type JumpInstruction struct {
	Forward int
}

type CallInstruction struct {
	Call string
	Args []int
}

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

func (vm *VirtualMachine) Run() error {
	// TODO: Check start exists.
	return vm.call("start", nil)
}

func (vm *VirtualMachine) call(syntax string, args []int) error {
	fn := vm.program.Functions[syntax]

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

	for fn.InstructionOffset < len(fn.Instructions) {
		instruction := fn.Instructions[fn.InstructionOffset]

		var move int
		var err error

		// TODO: This switch needs to be refactored into an interface.
		switch ins := instruction.(type) {
		case *CallInstruction:
			move, err = vm.callInstruction(ins)

		case *ConditionJumpInstruction:
			move, err = vm.conditionJumpInstruction(ins)

		case *JumpInstruction:
			move, err = vm.jumpInstruction(ins)

		default:
			panic(ins)
		}

		if err != nil {
			return err
		}

		fn.InstructionOffset += move
	}

	return nil
}

func (vm *VirtualMachine) conditionJumpInstruction(instruction *ConditionJumpInstruction) (int, error) {
	cmp := 0
	left := vm.GetArg(instruction.Left)
	right := vm.GetArg(instruction.Right)

	leftText, leftIsText := left.(*string)
	rightText, rightIsText := right.(*string)

	if leftIsText && rightIsText {
		cmp = strings.Compare(*leftText, *rightText)
		goto done
	} else {
		leftNumber, leftIsNumber := left.(*big.Rat)
		rightNumber, rightIsNumber := right.(*big.Rat)

		if leftIsNumber && rightIsNumber {
			cmp = leftNumber.Cmp(rightNumber)
			goto done
		}
	}

	return 0, fmt.Errorf("cannot compare: %s %s %s",
		vm.GetArgType(instruction.Left),
		instruction.Operator,
		vm.GetArgType(instruction.Right))

done:
	var result bool
	switch instruction.Operator {
	case OperatorEqual:
		result = cmp == 0

	case OperatorNotEqual:
		result = cmp != 0

	case OperatorGreaterThan:
		result = cmp > 0

	case OperatorGreaterThanEqual:
		result = cmp >= 0

	case OperatorLessThan:
		result = cmp < 0

	case OperatorLessThanEqual:
		result = cmp <= 0
	}

	if result {
		return instruction.True, nil
	}

	return instruction.False, nil
}

func (vm *VirtualMachine) jumpInstruction(instruction *JumpInstruction) (int, error) {
	return instruction.Forward, nil
}

func (vm *VirtualMachine) callInstruction(instruction *CallInstruction) (int, error) {
	// Check if it is a system call?
	if handler, ok := System[instruction.Call]; ok {
		handler(vm, instruction.Args)

		return 1, nil
	}

	// Otherwise we have to increase the stack.
	return 1, vm.call(instruction.Call, instruction.Args)
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

func (vm *VirtualMachine) GetArgType(index int) string {
	switch vm.GetArg(index).(type) {
	case *string:
		return VariableTypeText

	case *big.Rat:
		return VariableTypeNumber
	}

	return reflect.TypeOf(vm.GetArg(index)).String()
}
