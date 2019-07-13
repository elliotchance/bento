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
	r       io.Reader
	tokens  []Token
	offset  int
	program *Program
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
		Functions: map[string]*Function{},
	}

	// Now we can compile the program.
	for !parser.isFinished() {
		function, err := parser.consumeFunction()
		if err != nil {
			return nil, err
		}

		parser.program.AppendFunction(function)
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

	if ty != VariableTypeText && ty != VariableTypeNumber {
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

func (parser *Parser) consumeSentenceWords(varMap map[string]*VariableDefinition) ([]interface{}, error) {
	tokens := parser.consumeTokens([]string{
		TokenKindWord, TokenKindText, TokenKindNumber,
	})
	if len(tokens) == 0 {
		return nil, fmt.Errorf("expected sentence, but found something else")
	}

	var words []interface{}
	for _, token := range tokens {
		switch token.Kind {
		case TokenKindWord:
			if _, ok := varMap[token.Value]; ok {
				words = append(words, VariableReference(token.Value))
			} else {
				words = append(words, token.Value)
			}

		case TokenKindText:
			words = append(words, NewText(token.Value))

		case TokenKindNumber:
			words = append(words, NewNumber(token.Value))
		}
	}

	return words, nil
}

func (parser *Parser) consumeSentenceCall(varMap map[string]*VariableDefinition) (tokens []interface{}, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	words, err := parser.consumeSentenceWords(varMap)
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return nil, err
	}

	return words, nil
}

func (parser *Parser) consumeFunction() (function *Function, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	function, err = parser.consumeFunctionDeclaration()
	if err != nil {
		return nil, err
	}

	for !parser.isFinished() {
		// declare ...
		name, ty, err := parser.consumeDeclare()
		if err == nil {
			function.AppendDeclare(name, ty)
			continue
		}

		// Normal sentence.
		var words []interface{}
		words, err = parser.consumeSentenceCall(function.VariableMap())
		if err == nil {
			function.AppendSentence(words)
			continue
		}

		return function, nil
	}

	return function, nil
}

func (parser *Parser) consumeFunctionDeclaration() (function *Function, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	var words []interface{}
	words, err = parser.consumeSentenceWords(nil)
	if err != nil {
		return
	}

	function = &Function{
		Definition: &Sentence{
			Tokens: words,
		},
	}

	_, err = parser.consumeToken(TokenKindOpenBracket)
	if err == nil {
		vars, err := parser.consumeVariableIsTypeList()
		if err != nil {
			return nil, err
		}

		_, err = parser.consumeToken(TokenKindCloseBracket)
		if err != nil {
			return nil, err
		}

		for i, word := range function.Definition.Tokens {
			if ty, ok := vars[word.(string)]; ok {
				function.Definition.Tokens[i] = VariableReference(word.(string))

				// Note: It's important that we add the arguments in the order
				// that they appear rather than the order that they are defined.
				// Appending them in this loop will ensure that.
				function.AppendArgument(word.(string), ty)
			}
		}
	}

	_, err = parser.consumeToken(TokenKindColon)
	if err != nil {
		return nil, err
	}

	parser.consumeTokens([]string{TokenKindEndOfLine})

	return function, nil
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
