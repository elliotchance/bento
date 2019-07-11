// This file contains structures produced by the parser to describe the syntax
// of the program.

package main

import "strings"

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

	// Variables inclues the arguments and locally declared variables.
	Variables []*VariableDefinition

	Sentences []*Sentence
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

func (fn *Function) AppendDeclare(name, ty string) {
	fn.Variables = append(fn.Variables, &VariableDefinition{
		Name:       name,
		Type:       ty,
		LocalScope: true,
	})
}

func (fn *Function) AppendSentence(tokens []interface{}) {
	fn.Sentences = append(fn.Sentences, &Sentence{Tokens: tokens})
}

// Sentence is part of the AST. A sentence may not yet exist, or be valid.
type Sentence struct {
	Tokens []interface{}
}

// Syntax like "add ? to ?"
func (sentence *Sentence) Syntax() string {
	var words []string

	for _, word := range sentence.Tokens {
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
	for _, word := range sentence.Tokens {
		if _, ok := word.(string); !ok {
			args = append(args, word)
		}
	}

	return
}
