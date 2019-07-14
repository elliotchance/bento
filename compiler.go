package main

import "math/big"

type CompiledFunction struct {
	Variables         []interface{}
	Instructions      []Instruction
	InstructionOffset int
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
	for _, statement := range function.Statements {
		switch stmt := statement.(type) {
		case *Sentence:
			cf.Instructions = append(cf.Instructions,
				compileSentence(cf, function, stmt))

		case *If:
			cf.Instructions = append(cf.Instructions,
				compileIf(cf, function, stmt)...)

		default:
			panic(stmt)
		}
	}

	return cf
}

func compileSentence(cf *CompiledFunction, function *Function, sentence *Sentence) Instruction {
	instruction := &CallInstruction{
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

	return instruction
}

func compileIf(cf *CompiledFunction, function *Function, ifStmt *If) []Instruction {
	jumpInstruction := &ConditionJumpInstruction{
		Operator: ifStmt.Condition.Operator,
		True:     1,
		False:    2,
	}

	if ifStmt.False != nil {
		// This is to compensate for the added JumpInstruction that has to be
		// added below.
		jumpInstruction.False++
	}

	// TODO: This is duplicate code from compileSentence
	switch a := ifStmt.Condition.Left.(type) {
	case VariableReference:
		for i, arg2 := range function.Variables {
			if string(a) == arg2.Name {
				jumpInstruction.Left = i
				break
			}
		}

	// TODO: handle bad variable name

	case *string, *big.Rat:
		jumpInstruction.Left = len(cf.Variables)
		cf.Variables = append(cf.Variables, a)

	default:
		// TODO: This shouldn't be possible, it can be removed when the
		//  compiler is stable.
		panic(a)
	}

	// TODO: This is duplicate code from compileSentence (and above).
	switch a := ifStmt.Condition.Right.(type) {
	case VariableReference:
		for i, arg2 := range function.Variables {
			if string(a) == arg2.Name {
				jumpInstruction.Right = i
				break
			}
		}

	// TODO: handle bad variable name

	case *string, *big.Rat:
		jumpInstruction.Right = len(cf.Variables)
		cf.Variables = append(cf.Variables, a)

	default:
		// TODO: This shouldn't be possible, it can be removed when the
		//  compiler is stable.
		panic(a)
	}

	instructions := []Instruction{
		jumpInstruction,
		compileSentence(cf, function, ifStmt.True),
	}

	if ifStmt.False != nil {
		instructions = append(instructions,
			// This prevents the True case above from also running the else
			// clause.
			&JumpInstruction{Forward: 2},

			compileSentence(cf, function, ifStmt.False))
	}

	return instructions
}
