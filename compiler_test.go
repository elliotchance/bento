package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var compileTests = map[string]struct {
	program  *Program
	expected *CompiledProgram
}{
	"Display": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", NewText("hello"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("hello"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
	},
	"DisplayVariable": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name: "name",
							Type: "text",
						},
					},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", VariableReference("name"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText(""),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
	},
	"DisplayVariable2": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name:       "name",
							Type:       "text",
							LocalScope: true,
						},
					},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", NewText("hi"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText(""), NewText("hi"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{1},
						},
					},
				},
			},
		},
	},
	"Display2": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name:       "name",
							Type:       "text",
							LocalScope: true,
						},
					},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", NewText("hi"),
							},
						},
						{
							Tokens: []interface{}{
								"set", NewText("foo"), "to", VariableReference("name"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText(""), NewText("hi"), NewText("foo"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{1},
						},
						{
							Call: "set ? to ?",
							Args: []int{2, 0},
						},
					},
				},
			},
		},
	},
	"CallFunctionWithoutArguments": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"print",
							},
						},
					},
				},
				"print": {
					Definition: &Sentence{Tokens: []interface{}{"print"}},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", NewText("hi"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Instructions: []Instruction{
						{
							Call: "print",
						},
					},
				},
				"print": {
					Variables: []interface{}{
						NewText("hi"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
	},
	"CallFunctionWithArguments": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"print", NewText("foo"),
							},
						},
					},
				},
				"print ?": {
					Definition: &Sentence{Tokens: []interface{}{
						"print", VariableReference("message"),
					}},
					Variables: []*VariableDefinition{
						{
							Name: "message",
							Type: "text",
						},
					},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", VariableReference("message"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"),
					},
					Instructions: []Instruction{
						{
							Call: "print ?",
							Args: []int{0},
						},
					},
				},
				"print ?": {
					Variables: []interface{}{
						NewText(""),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
	},
	"NumberVariable": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Tokens: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name: "num",
							Type: "number",
						},
					},
					Sentences: []*Sentence{
						{
							Tokens: []interface{}{
								"display", VariableReference("num"),
							},
						},
					},
				},
			},
		},
		expected: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewNumber("0"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
	},
}

func TestCompileProgram(t *testing.T) {
	for testName, test := range compileTests {
		t.Run(testName, func(t *testing.T) {
			cf := CompileProgram(test.program)

			diff := cmp.Diff(test.expected, cf,
				cmpopts.IgnoreTypes((func([]interface{}))(nil)),
				cmpopts.AcyclicTransformer("NumberToString",
					func(number *big.Rat) string {
						return number.FloatString(6)
					}))

			assert.Empty(t, diff)
		})
	}
}
