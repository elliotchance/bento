package main

const (
	VariableTypeText = "text"
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
