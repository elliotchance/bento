package main

type Function struct {
	Sentences []*Sentence
}

func (fn *Function) Run(program *Program) {
	for _, sentence := range fn.Sentences {
		sentence.Run(program)
	}
}
