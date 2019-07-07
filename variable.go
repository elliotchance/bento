package main

const (
	VariableTypeText = "text"
)

type VariableReference string

type Variable struct {
	Type  string
	Value interface{}
}
