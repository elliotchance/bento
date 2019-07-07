package main

import (
	"fmt"
	"io"
)

// Reserved words.
const (
	WordDeclare = "declare"
)

func Parse(r io.Reader) (*Program, error) {
	tokens, err := Tokenize(r)
	if err != nil {
		return nil, err
	}

	program := &Program{
		Variables: map[string]*Variable{},
	}
	syntax := ""
	var args []interface{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Kind {
		case TokenKindWord:
			if token.Value == WordDeclare {
				name, ty, newI, err := consumeDeclare(tokens, i)
				if err != nil {
					return nil, err
				}

				program.Variables[name] = &Variable{
					Type:  ty,
					Value: "",
				}
				i = newI
				continue
			}

			if _, ok := program.Variables[token.Value]; ok {
				syntax = appendSyntax(syntax, "?")
				args = append(args, VariableReference(token.Value))
			} else {
				syntax = appendSyntax(syntax, token.Value)
			}

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

func consumeDeclare(tokens []Token, offset int) (string, string, int, error) {
	name := tokens[offset+1].Value
	ty := tokens[offset+3].Value

	return name, ty, offset + 4, nil
}

func appendSyntax(syntax, word string) string {
	if syntax == "" {
		return word
	}

	return syntax + " " + word
}
