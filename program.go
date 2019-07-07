package main

type Program struct {
	Sentences []*Sentence
}

func (program *Program) Run() {
	for _, sentence := range program.Sentences {
		sentence.Run()
	}
}
