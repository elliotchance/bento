// This file contains structures produced by the parser to describe the syntax
// of the program.

package main

import "strings"

const (
	OperatorEqual            = "="
	OperatorNotEqual         = "!="
	OperatorGreaterThan      = ">"
	OperatorGreaterThanEqual = ">="
	OperatorLessThan         = "<"
	OperatorLessThanEqual    = "<="
)

type Statement interface{}

// Program is the result of root-level AST after parsing the source code.
//
// The program may not be valid. It has to be compiled before it can be
// executed.
type Program struct {
	Functions map[string]*Function
}

func (program *Program) AppendFunction(fn *Function) {
	// TODO: Check for duplicate.
	program.Functions[fn.Definition.Syntax()] = fn
}

type Function struct {
	Definition *Sentence

	// Variables includes the arguments and locally declared variables.
	Variables []*VariableDefinition

	Statements []Statement

	IsQuestion bool
}

func (fn *Function) VariableMap() map[string]*VariableDefinition {
	m := make(map[string]*VariableDefinition)

	for _, variable := range fn.Variables {
		m[variable.Name] = variable
	}

	return m
}

func (fn *Function) AppendArgument(name, ty string) {
	fn.Variables = append(fn.Variables, &VariableDefinition{
		Name:       name,
		Type:       ty,
		LocalScope: false,
	})
}

func (fn *Function) AppendVariable(definition *VariableDefinition) {
	fn.Variables = append(fn.Variables, definition)
}

func (fn *Function) AppendStatement(statement Statement) {
	fn.Statements = append(fn.Statements, statement)
}

// Sentence is part of the AST. A sentence may not yet exist, or be valid.
type Sentence struct {
	Words []interface{}
}

// Syntax like "add ? to ?"
func (sentence *Sentence) Syntax() string {
	var words []string

	for _, word := range sentence.Words {
		if s, ok := word.(string); ok {
			words = append(words, s)
		} else {
			words = append(words, "?")
		}
	}

	return strings.Join(words, " ")
}

// Each of the values of the placeholders.
func (sentence *Sentence) Args() (args []interface{}) {
	for _, word := range sentence.Words {
		if _, ok := word.(string); !ok {
			args = append(args, word)
		}
	}

	return
}

type Condition struct {
	Left, Right interface{}
	Operator    string
}

type If struct {
	// Unless is true if "unless" was used instead of "if". This inverts the
	// logic.
	Unless bool

	// Either Condition or Question will be not-nil, never both.
	Condition *Condition
	Question  *Sentence

	// The blocks containing the true and false branches.
	True, False Statement
}

type While struct {
	// Until is true if "until" was used instead of "while". This inverts the
	// logic.
	Until bool

	// Either Condition or Question will be not-nil, never both.
	Condition *Condition
	Question  *Sentence

	// The blocks containing the true and false branches. This is a sentence
	// because it makes no sense to allow yes/no answers here.
	True *Sentence
}

type QuestionAnswer struct {
	Yes bool
}
