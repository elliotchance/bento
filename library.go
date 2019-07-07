package main

type Library struct {
	Sentences map[string]*Sentence
}

func (lib *Library) SentenceForSyntax(syntax string, args []interface{}) *Sentence {
	if sentence, ok := lib.Sentences[syntax]; ok {
		return &Sentence{
			Handler: sentence.Handler,
			Args:    args,
		}
	}

	return nil
}
