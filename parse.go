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
		Functions: map[string]*Function{},
	}

	// We have to parse the whole file for known variables and functions first.
	syntaxes, err := getAllFunctionSyntaxes(tokens)
	if err != nil {
		return nil, err
	}

	for _, syntax := range syntaxes {
		program.Functions[syntax] = &Function{}
	}

	syntax := ""
	currentFunction := "start"
	var args []interface{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Kind {
		case TokenKindColon:
			currentFunction = syntax
			syntax = ""
			i++ // skip Endline

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
			// Local function.
			sentence := program.SentenceForSyntax(syntax, args)
			if sentence != nil {
				goto found
			}

			// System function.
			sentence = System.SentenceForSyntax(syntax, args)
			if sentence == nil {
				return nil, fmt.Errorf("cannot understand: %s", syntax)
			}

		found:
			program.Functions[currentFunction].Sentences =
				append(program.Functions[currentFunction].Sentences, sentence)
			syntax = ""
			args = nil
		}
	}

	return program, nil
}

func getAllFunctionSyntaxes(tokens []Token) ([]string, error) {
	syntaxes := []string{"start"}
	variables := map[string]struct{}{}
	syntax := ""

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Kind {
		case TokenKindColon:
			syntaxes = append(syntaxes, syntax)
			syntax = ""
			i++ // skip Endline

		case TokenKindWord:
			if token.Value == WordDeclare {
				name, _, newI, err := consumeDeclare(tokens, i)
				if err != nil {
					return nil, err
				}

				variables[name] = struct{}{}
				i = newI
				continue
			}

			if _, ok := variables[token.Value]; ok {
				syntax = appendSyntax(syntax, "?")
			} else {
				syntax = appendSyntax(syntax, token.Value)
			}

		case TokenKindText:
			syntax = appendSyntax(syntax, "?")

		case TokenKindEndline:
			syntax = ""
		}
	}

	return syntaxes, nil
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
