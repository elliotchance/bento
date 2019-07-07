package main

type SentenceHandler func([]interface{})

type Sentence struct {
	Syntax  string
	Handler SentenceHandler
	Args    []interface{}
}

func (sentence *Sentence) Run() {
	sentence.Handler(sentence.Args)
}
