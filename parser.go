package main

import (
	"errors"
	"fmt"
	"io"
)

// These reserved words have special meaning when they are the first word of the
// sentence. It's fine to include them as normal words inside a sentence.
const (
	WordDeclare   = "declare"
	WordIf        = "if"
	WordOtherwise = "otherwise"
	WordUnless    = "unless"
	WordUntil     = "until"
	WordWhile     = "while"
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

func (parser *Parser) consumeSpecificWord(expected string) (word string, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	word, err = parser.consumeWord()
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

func (parser *Parser) consumeSentenceWord(varMap map[string]*VariableDefinition) (_ interface{}, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	var token Token
	token, err = parser.consumeToken(TokenKindWord)
	if err == nil {
		if _, ok := varMap[token.Value]; ok {
			return VariableReference(token.Value), nil
		} else {
			return token.Value, nil
		}
	}

	token, err = parser.consumeToken(TokenKindText)
	if err == nil {
		return NewText(token.Value), nil
	}

	token, err = parser.consumeToken(TokenKindNumber)
	if err == nil {
		return NewNumber(token.Value), nil
	}

	return nil, fmt.Errorf("expected sentence word, but found something else")
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

func (parser *Parser) consumeSentence(varMap map[string]*VariableDefinition) (sentence *Sentence, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	sentence = new(Sentence)

	for !parser.isFinished() {
		word, err := parser.consumeSentenceWord(varMap)
		if err != nil {
			break
		}

		sentence.Words = append(sentence.Words, word)
	}

	return
}

func (parser *Parser) consumeSentenceCall(varMap map[string]*VariableDefinition) (sentence *Sentence, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	sentence, err = parser.consumeSentence(varMap)
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return
	}

	return
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

		// if/unless ...
		ifStmt, err := parser.consumeIf(function.VariableMap())
		if err == nil {
			function.AppendStatement(ifStmt)
			continue
		}

		// while/until ...
		whileStmt, err := parser.consumeWhile(function.VariableMap())
		if err == nil {
			function.AppendStatement(whileStmt)
			continue
		}

		// Normal sentence.
		var sentence *Sentence
		sentence, err = parser.consumeSentenceCall(function.VariableMap())
		if err == nil {
			function.AppendStatement(sentence)
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

	var sentence *Sentence
	sentence, err = parser.consumeSentence(nil)
	if err != nil {
		return
	}

	function = &Function{
		Definition: sentence,
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

		for i, word := range function.Definition.Words {
			if ty, ok := vars[word.(string)]; ok {
				function.Definition.Words[i] = VariableReference(word.(string))

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
	// TODO: Can these be made easier?
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

func (parser *Parser) consumeIf(varMap map[string]*VariableDefinition) (ifStmt *If, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	ifStmt = &If{}

	_, err = parser.consumeSpecificWord(WordIf)
	if err != nil {
		_, err = parser.consumeSpecificWord(WordUnless)
		if err != nil {
			return nil, errors.New("expected if or unless")
		}

		ifStmt.Unless = true
	}

	// TODO: If we hit and if, we must not allow it to process the line as a
	//  sentence.

	ifStmt.Condition, err = parser.consumeCondition(varMap)
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindComma)
	if err != nil {
		return
	}

	ifStmt.True, err = parser.consumeSentence(varMap)
	if err != nil {
		return
	}

	// Bail out if safely if there is no "otherwise".
	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err == nil {
		return
	}

	_, err = parser.consumeToken(TokenKindComma)
	if err != nil {
		return
	}

	_, err = parser.consumeSpecificWord(WordOtherwise)
	if err != nil {
		return
	}

	ifStmt.False, err = parser.consumeSentence(varMap)
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return
	}

	return
}

func (parser *Parser) consumeCondition(varMap map[string]*VariableDefinition) (condition *Condition, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	condition = &Condition{}

	condition.Left, err = parser.consumeSentenceWord(varMap)
	if err != nil {
		return nil, err
	}

	condition.Operator, err = parser.consumeOperator()
	if err != nil {
		return nil, err
	}

	condition.Right, err = parser.consumeSentenceWord(varMap)
	if err != nil {
		return nil, err
	}

	return condition, nil
}

func (parser *Parser) consumeOperator() (string, error) {
	operatorToken, err := parser.consumeToken(TokenKindOperator)
	if err != nil {
		return "", err
	}

	return operatorToken.Value, nil
}

func (parser *Parser) consumeWhile(varMap map[string]*VariableDefinition) (whileStmt *While, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	whileStmt = &While{}

	_, err = parser.consumeSpecificWord(WordWhile)
	if err != nil {
		_, err = parser.consumeSpecificWord(WordUntil)
		if err != nil {
			return nil, errors.New("expected while or until")
		}

		whileStmt.Until = true
	}

	// TODO: If we hit a "while", we must not allow it to process the line as a
	//  sentence.

	whileStmt.Condition, err = parser.consumeCondition(varMap)
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindComma)
	if err != nil {
		return
	}

	whileStmt.True, err = parser.consumeSentenceCall(varMap)
	if err != nil {
		return
	}

	return
}
