package main

type Program struct {
	Variables map[string]*Variable
	Functions map[string]*Function
}

func (program *Program) Run() {
	program.Functions["start"].Run(program)
}

func (program *Program) ValueOf(val interface{}) interface{} {
	switch v := val.(type) {
	case VariableReference:
		return program.Variables[string(v)].Value

	default:
		return val
	}
}

func (program *Program) SentenceForSyntax(syntax string, args []interface{}) *Sentence {
	if _, ok := program.Functions[syntax]; ok {
		return &Sentence{
			Handler: func(program *Program, _ []interface{}) {
				program.Functions[syntax].Run(program)
			},
			Args: args,
		}
	}

	return nil
}
