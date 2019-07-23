package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// TODO: Prevent a variable from being redefined by the same name in a function.

// TODO: You cannot declare a variable with the same name as one of the function
//  parameters.

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

	token := parser.tokens[parser.offset]

	parser.offset++

	// Consume a conditional multiline.
	if parser.offset+1 < len(parser.tokens) &&
		parser.tokens[parser.offset].Kind == TokenKindEllipsis {
		if kind2 := parser.tokens[parser.offset+1].Kind; kind2 != TokenKindEndOfLine {
			return Token{},
				fmt.Errorf("expected %s, but got %s", TokenKindEndOfLine, kind2)
		}

		parser.offset += 2
	}

	return token, nil
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

func (parser *Parser) consumeSpecificWord(expected ...string) (word string, err error) {
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

	for _, allowed := range expected {
		if word == allowed {
			return word, nil
		}
	}

	return "", fmt.Errorf(`expected one of "%v", but got "%s"`, expected, word)
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
		if _, ok := varMap[token.Value]; ok || token.Value == "_" {
			return VariableReference(token.Value), nil
		}

		return token.Value, nil
	}

	token, err = parser.consumeToken(TokenKindText)
	if err == nil {
		return NewText(token.Value), nil
	}

	token, err = parser.consumeToken(TokenKindNumber)
	if err == nil {
		return NewNumber(token.Value, UnlimitedPrecision), nil
	}

	return nil, fmt.Errorf("expected sentence word, but found something else")
}

func (parser *Parser) consumeInteger() (value int, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	var token Token
	token, err = parser.consumeToken(TokenKindNumber)
	if err != nil {
		return
	}

	value, err = strconv.Atoi(token.Value)
	if err != nil {
		return
	}

	return
}

func (parser *Parser) consumeNumberType() (precision int, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	_, err = parser.consumeSpecificWord(VariableTypeNumber)
	if err != nil {
		return
	}

	_, err = parser.consumeSpecificWord("with")
	if err != nil {
		// That's OK, we can safely bail out here.
		precision = DefaultNumericPrecision
		err = nil
		return
	}

	precision, err = parser.consumeInteger()
	if err != nil {
		return
	}

	_, err = parser.consumeSpecificWord("decimal")
	if err != nil {
		return
	}

	_, err = parser.consumeSpecificWord("places")
	if err != nil {
		// Even through "places" or "place" is allowed for any number of decimal
		// places, it reads better to say "1 decimal place".
		_, err = parser.consumeSpecificWord("place")
		if err != nil {
			return
		}

		err = nil
	}

	return
}

func (parser *Parser) consumeType() (ty string, precision int, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	// The "a" and "an" are optional so that it is easier to read in some cases:
	//
	//   is text
	//   is a number
	//   is an order receipt

	_, err = parser.consumeSpecificWord("a")
	if err == nil {
		goto consumeType
	}

	_, err = parser.consumeSpecificWord("an")
	if err == nil {
		goto consumeType
	}

consumeType:
	ty, err = parser.consumeSpecificWord(VariableTypeText)
	if err == nil {
		return ty, 0, nil
	}

	precision, err = parser.consumeNumberType()
	if err == nil {
		return "number", precision, err
	}

	return "", 0, fmt.Errorf("expected variable type, but got %s", ty)
}

// Examples:
//
//   some-variable is text
//   some-variable is number
//   some-variable is number with 2 decimal places
//
func (parser *Parser) consumeVariableIsType() (definition *VariableDefinition, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	definition = new(VariableDefinition)

	definition.Name, err = parser.consumeWord()
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeSpecificWord("is")
	if err != nil {
		return nil, err
	}

	definition.Type, definition.Precision, err = parser.consumeType()
	if err != nil {
		return nil, err
	}

	return
}

// foo is text, bar is text
func (parser *Parser) consumeVariableIsTypeList() (list map[string]*VariableDefinition, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	list = map[string]*VariableDefinition{}

	for !parser.isFinished() {
		definition, err := parser.consumeVariableIsType()
		if err != nil {
			return nil, err
		}

		list[definition.Name] = definition

		err = parser.consumeComma()
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

func (parser *Parser) consumeQuestionAnswer() (answer *QuestionAnswer, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	word, err := parser.consumeSpecificWord("yes", "no")
	if err != nil {
		return nil, err
	}

	return &QuestionAnswer{
		Yes: word == "yes",
	}, nil
}

func (parser *Parser) consumeQuestionAnswerCall() (answer *QuestionAnswer, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	answer, err = parser.consumeQuestionAnswer()
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return nil, err
	}

	return answer, nil
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

func (parser *Parser) consumeSentenceOrAnswer(varMap map[string]*VariableDefinition) (_ Statement, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	answer, err := parser.consumeQuestionAnswer()
	if err == nil {
		return answer, nil
	}

	return parser.consumeSentence(varMap)
}

func (parser *Parser) consumeSentenceCallOrAnswerCall(varMap map[string]*VariableDefinition) (_ Statement, err error) {
	originalOffset := parser.offset
	defer func() {
		if err != nil {
			parser.offset = originalOffset
		}
	}()

	answer, err := parser.consumeQuestionAnswerCall()
	if err == nil {
		return answer, nil
	}

	return parser.consumeSentenceCall(varMap)
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
		definition, err := parser.consumeDeclare()
		if err == nil {
			definition.LocalScope = true
			function.AppendVariable(definition)
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

		// TODO: yes/no cannot be used outside of questions
		sentenceOrAnswer, err :=
			parser.consumeSentenceCallOrAnswerCall(function.VariableMap())
		if err == nil {
			function.AppendStatement(sentenceOrAnswer)
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
				function.AppendArgument(word.(string), ty.Type)
			}
		}
	}

	_, err = parser.consumeToken(TokenKindQuestion)
	if err == nil {
		function.IsQuestion = true
	} else {
		_, err = parser.consumeToken(TokenKindColon)
		if err != nil {
			return nil, err
		}
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

func (parser *Parser) consumeDeclare() (definition *VariableDefinition, err error) {
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

	definition, err = parser.consumeVariableIsType()
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return
	}

	return
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
		// It must be a question instead of a condition.
		ifStmt.Question, err = parser.consumeSentence(varMap)

		if err != nil {
			return
		}
	}

	err = parser.consumeComma()
	if err != nil {
		return
	}

	ifStmt.True, err = parser.consumeSentenceOrAnswer(varMap)
	if err != nil {
		return
	}

	// Bail out if safely if there is no "otherwise".
	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err == nil {
		return
	}

	err = parser.consumeComma()
	if err != nil {
		return
	}

	_, err = parser.consumeSpecificWord(WordOtherwise)
	if err != nil {
		return
	}

	ifStmt.False, err = parser.consumeSentenceOrAnswer(varMap)
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

func (parser *Parser) consumeComma() error {
	_, err := parser.consumeToken(TokenKindComma)
	if err != nil {
		return err
	}

	// Ignore the error as the new line is optional.
	_, _ = parser.consumeToken(TokenKindEndOfLine)

	return nil
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
		// It must be a question instead of a condition.
		whileStmt.Question, err = parser.consumeSentence(varMap)

		if err != nil {
			return
		}
	}

	err = parser.consumeComma()
	if err != nil {
		return
	}

	whileStmt.True, err = parser.consumeSentence(varMap)
	if err != nil {
		return
	}

	_, err = parser.consumeToken(TokenKindEndOfLine)
	if err != nil {
		return
	}

	return
}
