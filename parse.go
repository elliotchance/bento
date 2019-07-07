package main

import (
	"fmt"
	"io"
)

func Parse(r io.Reader) (*Program, error) {
	tokens, err := Tokenize(r)
	if err != nil {
		return nil, err
	}

	program := &Program{}
	syntax := ""
	var args []interface{}

	for _, token := range tokens {
		switch token.Kind {
		case TokenKindWord:
			syntax = appendSyntax(syntax, token.Value)

		case TokenKindText:
			syntax = appendSyntax(syntax, "?")
			args = append(args, token.Value)

		case TokenKindEndline:
			sentence := System.SentenceForSyntax(syntax, args)
			if sentence == nil {
				return nil, fmt.Errorf("cannot understand: %s", syntax)
			}

			program.Sentences = append(program.Sentences, sentence)
			syntax = ""
			args = nil
		}
	}

	return program, nil
}

func appendSyntax(syntax, word string) string {
	if syntax == "" {
		return word
	}

	return syntax + " " + word
}
