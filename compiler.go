package main

type CompiledFunction struct {
	Variables         []interface{}
	Instructions      []Instruction
	InstructionOffset int
}

type CompiledProgram struct {
	Functions map[string]*CompiledFunction
}

type Compiler struct {
	program  *Program
	function *Function
	cf       *CompiledFunction
}

func NewCompiler(program *Program) *Compiler {
	return &Compiler{
		program: program,
	}
}

func (compiler *Compiler) Compile() *CompiledProgram {
	cp := &CompiledProgram{
		Functions: make(map[string]*CompiledFunction),
	}

	for _, compiler.function = range compiler.program.Functions {
		syntax := compiler.function.Definition.Syntax()
		compiler.compileFunction()
		cp.Functions[syntax] = compiler.cf
	}

	return cp
}

func (compiler *Compiler) compileFunction() {
	compiler.cf = &CompiledFunction{}

	// Make spaces for the arguments and locally declared variables. These
	// placeholders will be nil. The virtual machine will fill in the real
	// values at the time the function is invoked.
	for _, variable := range compiler.function.Variables {
		var value interface{}

		switch variable.Type {
		case VariableTypeText:
			value = NewText("")

		case VariableTypeNumber:
			value = NewNumber("0", variable.Precision)
		}

		compiler.cf.Variables = append(compiler.cf.Variables, value)
	}

	// All of other constants are appended into the end.
	// TODO: Change this switch into an interface.
	for _, statement := range compiler.function.Statements {
		switch stmt := statement.(type) {
		case *Sentence:
			compiler.cf.Instructions = append(compiler.cf.Instructions,
				compiler.compileSentence(stmt))

		case *If:
			compiler.cf.Instructions = append(compiler.cf.Instructions,
				compiler.compileIf(stmt)...)

		case *While:
			compiler.cf.Instructions = append(compiler.cf.Instructions,
				compiler.compileWhile(stmt)...)

		default:
			panic(stmt)
		}
	}
}

func (compiler *Compiler) resolveArg(arg interface{}) int {
	switch a := arg.(type) {
	case VariableReference:
		if a == BlackholeVariable {
			return blackholeVariableIndex
		}

		for i, arg2 := range compiler.function.Variables {
			if string(a) == arg2.Name {
				return i
			}
		}

		// TODO: handle bad variable name

	case *string, *Number:
		compiler.cf.Variables = append(compiler.cf.Variables, a)
		return len(compiler.cf.Variables) - 1
	}

	// Not possible
	return 0
}

func (compiler *Compiler) compileSentence(sentence *Sentence) Instruction {
	instruction := &CallInstruction{
		Call: sentence.Syntax(),
		Args: nil,
	}

	// TODO: Check the syntax exists in the system or file.

	for _, arg := range sentence.Args() {
		instruction.Args = append(instruction.Args, compiler.resolveArg(arg))
	}

	return instruction
}

func (compiler *Compiler) compileIf(ifStmt *If) []Instruction {
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

	jumpInstruction.Left = compiler.resolveArg(ifStmt.Condition.Left)
	jumpInstruction.Right = compiler.resolveArg(ifStmt.Condition.Right)

	if ifStmt.Unless {
		jumpInstruction.True, jumpInstruction.False =
			jumpInstruction.False, jumpInstruction.True
	}

	instructions := []Instruction{
		jumpInstruction,
		compiler.compileSentence(ifStmt.True),
	}

	if ifStmt.False != nil {
		instructions = append(instructions,
			// This prevents the True case above from also running the else
			// clause.
			&JumpInstruction{Forward: 2},

			compiler.compileSentence(ifStmt.False))
	}

	return instructions
}

func (compiler *Compiler) compileWhile(whileStmt *While) []Instruction {
	jumpInstruction := &ConditionJumpInstruction{
		Operator: whileStmt.Condition.Operator,
		True:     1,
		False:    3,
	}

	jumpInstruction.Left = compiler.resolveArg(whileStmt.Condition.Left)
	jumpInstruction.Right = compiler.resolveArg(whileStmt.Condition.Right)

	if whileStmt.Until {
		jumpInstruction.True, jumpInstruction.False =
			jumpInstruction.False, jumpInstruction.True
	}

	instructions := []Instruction{
		jumpInstruction,
		compiler.compileSentence(whileStmt.True),
		&JumpInstruction{Forward: -2},
	}

	return instructions
}
