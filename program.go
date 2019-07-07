package main

type Program struct {
	Variables map[string]*Variable
	Sentences []*Sentence
}

func (program *Program) Run() {
	for _, sentence := range program.Sentences {
		sentence.Run(program)
	}
}

func (program *Program) ValueOf(val interface{}) interface{} {
	switch v := val.(type) {
	case VariableReference:
		return program.Variables[string(v)].Value

	default:
		return val
	}
}
