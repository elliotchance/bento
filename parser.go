package main

import (
	"errors"
	"fmt"
	"io"
)

// These reserved words have special meaning when they are the first word of the
// sentence. It's fine to include them as normal words inside a sentence.
const (
	WordDeclare = "declare"
)

type Parser struct {
	r               io.Reader
	tokens          []Token
	offset          int
	program         *Program
	currentFunction string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		r: r,
	}
}

func (parser *Parser) Parse() (*Program, error) {
	var err error
	parser.tokens, err = Tokenize(parser.r)
	if err != nil {
		return nil, err
	}

	parser.program = &Program{
		Variables: map[string]*Variable{},
		Functions: map[string]*Function{},
	}
	parser.currentFunction = "start"

	// We have to parse the whole file for known variables and functions first.
	parser.prepareAllVariablesAndFunctionDeclarations()

	// Now we can compile the program.
	for !parser.isFinished() {
		// declare ...
		_, _, err := parser.consumeDeclare()
		if err == nil {
			// Ignore this, it's already been handled by
			// prepareAllVariablesAndFunctionDeclarations.

			continue
		}

		// function declaration:
		declarationSyntax, _, err := parser.consumeFunctionDeclaration()
		if err == nil {
			parser.currentFunction = declarationSyntax

			continue
		}

		// sentence
		sentence, args, err := parser.consumeSentenceCall()
		if err == nil {
			syntax := sentence.Syntax()

			// Local function.
			sentence := parser.program.SentenceForSyntax(syntax, args)
			if sentence != nil {
				goto found
			}

			// System function.
			sentence = System.SentenceForSyntax(syntax, args)
			if sentence != nil {
				goto found
			}

			return nil, fmt.Errorf("cannot understand: %s", syntax)

		found:
			parser.program.Functions[parser.currentFunction].Sentences = append(
				parser.program.Functions[parser.currentFunction].Sentences,
				sentence)
			continue
		}

		return nil, fmt.Errorf("unexpected %s", parser.tokens[parser.offset])
	}

	return parser.program, nil
}

func (parser *Parser) consumeToken(kind string) (Token, error) {
	if parser.offset >= len(parser.tokens) {
		return Token{},
			errors.New("expected token, but the file ended unexpectedly")
	}

	if kind2 := parser.tokens[parser.offset].Kind; kind2 != kind {
		return Token{}, fmt.Errorf("expected %s, but got %s", kind, kind2)
	}

	parser.offset++

	return parser.tokens[parser.offset-1], nil
}

func (parser *Parser) consumeTokens(kinds []string) (tokens []Token) {
	for ; parser.offset < len(parser.tokens); parser.offset++ {
		token := parser.tokens[parser.offset]

		found := false
		for _, kind := range kinds {
			if kind == token.Kind {
				found = true
				break
			}
		}

		if !found {
			break
		}

		tokens = append(tokens, token)
	}

	return
}

func (parser *Parser) consumeSpecificWord(expected string) (string, error) {
	word, err := parser.consumeWord()
	if err != nil {
		return "", err
	}

	if word != expected {
		return "", fmt.Errorf(`expected "%s", but got "%s"`, expected, word)
	}

	return word, nil
}

func (parser *Parser) consumeWord() (string, error) {
	token, err := parser.consumeToken(TokenKindWord)
	if err != nil {
		return "", err
	}

	return token.Value, err
}

func (parser *Parser) consumeType() (string, error) {
	ty, err := parser.consumeWord()

	if ty != VariableTypeText {
		return "", fmt.Errorf("expected variable type, but got %s", ty)
	}

	return ty, err
}

// some-variable is text
func (parser *Parser) consumeVariableIsType() (name, ty string, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	name, err = parser.consumeWord()
	if err != nil {
		return "", "", err
	}

	_, err = parser.consumeSpecificWord("is")
	if err != nil {
		return "", "", err
	}

	ty, err = parser.consumeType()
	if err != nil {
		return "", "", err
	}

	return
}

// foo is text, bar is text
func (parser *Parser) consumeVariableIsTypeList() (list map[string]string, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	list = map[string]string{}

	for !parser.isFinished() {
		name, ty, err := parser.consumeVariableIsType()
		if err != nil {
			return nil, err
		}

		list[name] = ty

		_, err = parser.consumeToken(TokenKindComma)
		if err != nil {
			break
		}
	}

	return
}

func (parser *Parser) consumeSentenceWords() (Tokens, error) {
	tokens := parser.consumeTokens([]string{TokenKindWord, TokenKindText})
	if len(tokens) == 0 {
		return nil, fmt.Errorf("expected sentence, but found something else")
	}

	return tokens, nil
}

func (parser *Parser) consumeSentenceCall() (tokens Tokens, args []interface{}, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	for ; !parser.lineIsFinished(); parser.offset++ {
		token := parser.tokens[parser.offset]
		switch token.Kind {
		case TokenKindWord:
			// Local variable.
			if _, ok := parser.program.Functions[parser.currentFunction].Variables[token.Value]; ok {
				args = append(args, VariableReference(token.Value))
				token.Value = "?"

				goto done
			}

			// Global variable.
			if _, ok := parser.program.Variables[token.Value]; ok {
				args = append(args, VariableReference(token.Value))
				token.Value = "?"

				goto done
			}

		case TokenKindText:
			args = append(args, token.Value)
			token.Value = "?"

		default:
			break
		}

	done:
		tokens = append(tokens, token)
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return nil, nil, err
	}

	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("expected sentence call")
	}

	return
}

func (parser *Parser) consumeFunctionDeclaration() (syntax string, vars map[string]*Variable, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	var sentenceTokens Tokens
	sentenceTokens, err = parser.consumeSentenceWords()
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindOpenBracket)
	if err == nil {
		varsList, err := parser.consumeVariableIsTypeList()
		if err != nil {
			return "", nil, err
		}

		_, err = parser.consumeToken(TokenKindCloseBracket)
		if err != nil {
			return "", nil, err
		}

		vars = map[string]*Variable{}
		position := 0
		for i := 0; i < len(sentenceTokens); i++ {
			for name, ty := range varsList {
				if sentenceTokens[i].Value == name {
					sentenceTokens[i].Value = "?"

					if _, ok := vars[name]; !ok {
						vars[name] = &Variable{
							Type:     ty,
							Position: position,
						}
						position++
					}

					break
				}
			}
		}
	}

	_, err = parser.consumeToken(TokenKindColon)
	if err != nil {
		return "", nil, err
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return "", nil, err
	}

	return sentenceTokens.Syntax(), vars, nil
}

func (parser *Parser) prepareAllVariablesAndFunctionDeclarations() {
	parser.program.Functions["start"] = &Function{}

	for parser.offset = 0; parser.offset < len(parser.tokens); parser.offset++ {
		// declare ...
		name, ty, err := parser.consumeDeclare()
		if err == nil {
			parser.program.Variables[name] = &Variable{
				Type:  ty,
				Value: "",
			}

			continue
		}

		// function declaration:
		declarationSyntax, vars, err := parser.consumeFunctionDeclaration()
		if err == nil {
			parser.program.Functions[declarationSyntax] = &Function{
				Variables: vars,
			}
			continue
		}

		// It didn't find a match. We should fast-forward to the start of the
		// next line.
		parser.offset++
		for ; parser.offset < len(parser.tokens); parser.offset++ {
			if parser.tokens[parser.offset].Kind == TokenKindEndOfLine {
				break
			}
		}
	}

	parser.offset = 0
}

func (parser *Parser) isFinished() bool {
	return parser.tokens[parser.offset].Kind == TokenKindEndOfFile
}

func (parser *Parser) lineIsFinished() bool {
	return parser.tokens[parser.offset].Kind == TokenKindEndOfLine
}

func (parser *Parser) consumeDeclare() (name, ty string, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	_, err = parser.consumeSpecificWord(WordDeclare)
	if err != nil {
		return
	}

	name, ty, err = parser.consumeVariableIsType()
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return
	}

	return name, ty, nil
}
