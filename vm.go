package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Instruction interface{}

type ConditionJumpInstruction struct {
	Left, Right int
	Operator    string
	True, False int
}

type QuestionJumpInstruction struct {
	True, False int
}

type JumpInstruction struct {
	Forward int
}

type CallInstruction struct {
	Call string
	Args []int
}

type QuestionAnswerInstruction struct {
	Yes bool
}

type VirtualMachine struct {
	program     *CompiledProgram
	memory      []interface{}
	stackOffset []int
	out         io.Writer
	answer      bool
	backends    []*Backend
}

func NewVirtualMachine(program *CompiledProgram) *VirtualMachine {
	return &VirtualMachine{
		program: program,
		out:     os.Stdout,
	}
}

func (vm *VirtualMachine) Run() error {
	vm.stackOffset = []int{0}

	// TODO: Check start exists.
	return vm.call("start", nil)
}

func (vm *VirtualMachine) call(syntax string, args []int) error {
	fn := vm.program.Functions[syntax]

	if fn == nil {
		// Maybe it belongs to a backend?
		for _, arg := range args {
			// TODO: It is ambiguous if a sentence contains more than one
			//  backend.
			if backend, ok := vm.GetArg(arg).(*Backend); ok {
				// TODO: Make sure syntax exists
				var realArgs []string
				for _, realArg := range args {
					realArgs = append(realArgs, fmt.Sprintf("%v", vm.GetArg(realArg)))
				}
				result, err := backend.send(&BackendRequest{
					Sentence: syntax,
					Args:     realArgs,
				})
				if err != nil {
					panic(err)
				}

				for key, value := range result.Set {
					index, err := strconv.Atoi(key[1:])
					if err != nil {
						panic(err)
					}

					vm.SetArg(index, NewText(value))
				}

				return nil
			}
		}

		return fmt.Errorf("no such function: %s", syntax)
	}

	// Start backends.
	// TODO: Backends are not closed.
	for _, variable := range fn.Variables {
		if backend, ok := variable.(*Backend); ok {
			err := backend.Start()
			if err != nil {
				return err
			}
		}
	}

	fn.InstructionOffset = 0

	// Expand the memory to accommodate the variables (arguments and constants
	// used in the function).
	// TODO: Refactor this in a much more efficient way.
	offset := vm.stackOffset[len(vm.stackOffset)-1]
	for i, v := range fn.Variables {
		for offset+i >= len(vm.memory) {
			vm.memory = append(vm.memory, nil)
		}
		vm.memory[offset+i] = v
	}

	// Load in the arguments.
	for i, arg := range args {
		to := vm.stackOffset[len(vm.stackOffset)-1] + i
		from := vm.stackOffset[len(vm.stackOffset)-2] + arg
		vm.memory[to] = vm.memory[from]
	}

	vm.stackOffset = append(vm.stackOffset,
		vm.stackOffset[len(vm.stackOffset)-1]+len(fn.Variables))

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

		case *QuestionJumpInstruction:
			move, err = vm.questionJumpInstruction(ins)

		case *QuestionAnswerInstruction:
			move, err = vm.questionAnswerInstruction(ins)

		default:
			panic(ins)
		}

		if err != nil {
			return err
		}

		fn.InstructionOffset += move
	}

	vm.stackOffset = vm.stackOffset[:len(vm.stackOffset)-1]

	return nil
}

func (vm *VirtualMachine) questionAnswerInstruction(instruction *QuestionAnswerInstruction) (int, error) {
	vm.answer = instruction.Yes

	return 1, nil
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
		leftNumber, leftIsNumber := left.(*Number)
		rightNumber, rightIsNumber := right.(*Number)

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

func (vm *VirtualMachine) questionJumpInstruction(instruction *QuestionJumpInstruction) (int, error) {
	if vm.answer {
		return instruction.True, nil
	}

	return instruction.False, nil
}

func (vm *VirtualMachine) jumpInstruction(instruction *JumpInstruction) (int, error) {
	return instruction.Forward, nil
}

func (vm *VirtualMachine) callInstruction(instruction *CallInstruction) (int, error) {
	// We technically only need to do this when calling a question.
	vm.answer = false

	// Check if it is a system call?
	if handler, ok := System[instruction.Call]; ok {
		handler(vm, instruction.Args)

		return 1, nil
	}

	// Otherwise we have to increase the stack.
	return 1, vm.call(instruction.Call, instruction.Args)
}

func (vm *VirtualMachine) GetArg(index int) interface{} {
	if index == blackholeVariableIndex {
		return nil
	}

	return vm.memory[vm.previousOffset()+index]
}

func (vm *VirtualMachine) previousOffset() int {
	return vm.stackOffset[len(vm.stackOffset)-2]
}

func (vm *VirtualMachine) SetArg(index int, value interface{}) {
	if index == blackholeVariableIndex {
		return
	}

	vm.memory[vm.previousOffset()+index] = value
}

func (vm *VirtualMachine) GetNumber(index int) *Number {
	if index == blackholeVariableIndex {
		return NewNumber("0", DefaultNumericPrecision)
	}

	return vm.memory[vm.previousOffset()+index].(*Number)
}

func (vm *VirtualMachine) GetText(index int) *string {
	if index == blackholeVariableIndex {
		return NewText("")
	}

	return vm.memory[vm.previousOffset()+index].(*string)
}

func (vm *VirtualMachine) GetArgType(index int) string {
	switch vm.GetArg(index).(type) {
	case nil:
		return VariableTypeBlackhole

	case *string:
		return VariableTypeText

	case *Number:
		return VariableTypeNumber
	}

	return reflect.TypeOf(vm.GetArg(index)).String()
}
