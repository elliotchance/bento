package main

type SentenceHandler func(*Program, []interface{})

type Sentence struct {
	Handler SentenceHandler
	Args    []interface{}
}

func (sentence *Sentence) Run(program *Program) {
	sentence.Handler(program, sentence.Args)
}
