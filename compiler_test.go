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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name: "name",
							Type: "text",
						},
					},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name:       "name",
							Type:       "text",
							LocalScope: true,
						},
					},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name:       "name",
							Type:       "text",
							LocalScope: true,
						},
					},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
								"display", NewText("hi"),
							},
						},
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
							Call: "display ?",
							Args: []int{1},
						},
						&CallInstruction{
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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
								"print",
							},
						},
					},
				},
				"print": {
					Definition: &Sentence{Words: []interface{}{"print"}},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
							Call: "print",
						},
					},
				},
				"print": {
					Variables: []interface{}{
						NewText("hi"),
					},
					Instructions: []Instruction{
						&CallInstruction{
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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
								"print", NewText("foo"),
							},
						},
					},
				},
				"print ?": {
					Definition: &Sentence{Words: []interface{}{
						"print", VariableReference("message"),
					}},
					Variables: []*VariableDefinition{
						{
							Name: "message",
							Type: "text",
						},
					},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
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
						&CallInstruction{
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
					Definition: &Sentence{Words: []interface{}{"start"}},
					Variables: []*VariableDefinition{
						{
							Name: "num",
							Type: "number",
						},
					},
					Statements: []Statement{
						&Sentence{
							Words: []interface{}{
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
						&CallInstruction{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
	},
	"InlineIf": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&If{
							Condition: &Condition{
								Left:     NewText("foo"),
								Right:    NewText("bar"),
								Operator: OperatorEqual,
							},
							True: &Sentence{
								Words: []interface{}{
									"display", NewText("match!"),
								},
							},
						},
						&Sentence{
							Words: []interface{}{
								"display", NewText("done"),
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
						NewText("foo"), NewText("bar"), NewText("match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     1,
							False:    2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
	},
	"InlineIfElse": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&If{
							Condition: &Condition{
								Left:     NewText("foo"),
								Right:    NewText("bar"),
								Operator: OperatorNotEqual,
							},
							True: &Sentence{
								Words: []interface{}{
									"display", NewText("match!"),
								},
							},
							False: &Sentence{
								Words: []interface{}{
									"display", NewText("no match!"),
								},
							},
						},
						&Sentence{
							Words: []interface{}{
								"display", NewText("done"),
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
						NewText("foo"), NewText("bar"), NewText("match!"), NewText("no match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorNotEqual,
							True:     1,
							False:    3,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&JumpInstruction{
							Forward: 2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{4},
						},
					},
				},
			},
		},
	},
	"InlineUnless": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&If{
							Unless: true,
							Condition: &Condition{
								Left:     NewText("foo"),
								Right:    NewText("bar"),
								Operator: OperatorEqual,
							},
							True: &Sentence{
								Words: []interface{}{
									"display", NewText("match!"),
								},
							},
						},
						&Sentence{
							Words: []interface{}{
								"display", NewText("done"),
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
						NewText("foo"), NewText("bar"), NewText("match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     2,
							False:    1,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
	},
	"InlineUnlessElse": {
		program: &Program{
			Functions: map[string]*Function{
				"start": {
					Definition: &Sentence{Words: []interface{}{"start"}},
					Statements: []Statement{
						&If{
							Unless: true,
							Condition: &Condition{
								Left:     NewText("foo"),
								Right:    NewText("bar"),
								Operator: OperatorNotEqual,
							},
							True: &Sentence{
								Words: []interface{}{
									"display", NewText("match!"),
								},
							},
							False: &Sentence{
								Words: []interface{}{
									"display", NewText("no match!"),
								},
							},
						},
						&Sentence{
							Words: []interface{}{
								"display", NewText("done"),
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
						NewText("foo"), NewText("bar"), NewText("match!"), NewText("no match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorNotEqual,
							True:     3,
							False:    1,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&JumpInstruction{
							Forward: 2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{4},
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
			compiler := NewCompiler(test.program)
			cf := compiler.Compile()

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
