package main

type Function struct {
	Variables map[string]*Variable
	Sentences []*Sentence
}

func (fn *Function) Run(program *Program) {
	program.PrintTrace(program.CurrentFunction + ":", nil)

	program.StackLevel++
	for _, sentence := range fn.Sentences {
		sentence.Run(program)
	}
	program.StackLevel--
}
