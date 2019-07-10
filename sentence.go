package main

type SentenceHandler func(*Program, []interface{})

type Sentence struct {
	Syntax  string
	Handler SentenceHandler
	Args    []interface{}
}

func (sentence *Sentence) Run(program *Program) {
	program.PrintTrace(sentence.Syntax, sentence.Args)
	sentence.Handler(program, sentence.Args)
}
