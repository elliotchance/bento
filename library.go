package main

type Library struct {
	Sentences []*Sentence
}

func (lib *Library) SentenceForSyntax(syntax string, args []interface{}) *Sentence {
	for _, sentence := range lib.Sentences {
		if sentence.Syntax == syntax {
			return &Sentence{
				Syntax:  syntax,
				Handler: sentence.Handler,
				Args:    args,
			}
		}
	}

	return nil
}
