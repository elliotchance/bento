package main

import "math/big"

const (
	VariableTypeText   = "text"
	VariableTypeNumber = "number"
)

type VariableDefinition struct {
	Name string
	Type string

	// LocalScope is true if the variable was declared within the function.
	LocalScope bool
}

type VariableReference string

func NewText(s string) *string {
	return &s
}

func NewNumber(s string) *big.Rat {
	number, success := big.NewRat(0, 1).SetString(s)
	if !success {
		panic(s)
	}

	return number
}
