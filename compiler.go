package main

import "math/big"

type Instruction struct {
	Call string
	Args []int
}

type CompiledFunction struct {
	Variables    []interface{}
	Instructions []Instruction
}

type CompiledProgram struct {
	Functions map[string]*CompiledFunction
}

func CompileProgram(program *Program) *CompiledProgram {
	cp := &CompiledProgram{
		Functions: make(map[string]*CompiledFunction),
	}

	for _, function := range program.Functions {
		syntax := function.Definition.Syntax()
		cp.Functions[syntax] = CompileFunction(function)
	}

	return cp
}

func CompileFunction(function *Function) *CompiledFunction {
	cf := &CompiledFunction{}

	// Make spaces for the arguments and locally declared variables. These
	// placeholders will be nil. The virtual machine will fill in the real
	// values at the time the function is invoked.
	for _, variable := range function.Variables {
		var value interface{}

		switch variable.Type {
		case VariableTypeText:
			value = NewText("")

		case VariableTypeNumber:
			value = NewNumber("0")
		}

		cf.Variables = append(cf.Variables, value)
	}

	// All of other constants are appended into the end.
	for _, sentence := range function.Sentences {
		instruction := Instruction{
			Call: sentence.Syntax(),
			Args: nil,
		}

		// TODO: Check the syntax exists in the system or file.

		for _, arg := range sentence.Args() {
			switch a := arg.(type) {
			case VariableReference:
				for i, arg2 := range function.Variables {
					if string(a) == arg2.Name {
						instruction.Args = append(instruction.Args, i)
						break
					}
				}

				// TODO: handle bad variable name

			case *string, *big.Rat:
				instruction.Args = append(instruction.Args, len(cf.Variables))
				cf.Variables = append(cf.Variables, a)

			default:
				// TODO: This shouldn't be possible, it can be removed when the
				//  compiler is stable.
				panic(arg)
			}
		}

		cf.Instructions = append(cf.Instructions, instruction)
	}

	return cf
}
