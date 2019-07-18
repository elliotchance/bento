package main

const (
	VariableTypeBlackhole = "blackhole"
	VariableTypeText      = "text"
	VariableTypeNumber    = "number"
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

var BlackholeVariable = VariableReference("_")

var blackholeVariableIndex = -1

func NewText(s string) *string {
	return &s
}
