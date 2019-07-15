package main

const (
	VariableTypeText   = "text"
	VariableTypeNumber = "number"
)

type VariableDefinition struct {
	Name string
	Type string

	// LocalScope is true if the variable was declared within the function.
	LocalScope bool

	// Precision is the decimal places for "number" type.
	Precision int
}

type VariableReference string

func NewText(s string) *string {
	return &s
}
