package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	for testName, test := range map[string]struct {
		bento    string
		expected *Program
	}{
		"Empty": {
			bento: "",
			expected: &Program{
				Functions: map[string]*Function{},
			},
		},
		"EmptyStart": {
			bento: "start:",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
					},
				},
			},
		},
		"Display": {
			bento: `start: Display "Hello, World!"`,
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", NewText("Hello, World!"),
								},
							},
						},
					},
				},
			},
		},
		"Display2": {
			bento: "start:\nDisplay \"Hello, World!\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", NewText("Hello, World!"),
								},
							},
						},
					},
				},
			},
		},
		"DisplayTwice": {
			bento: "start: Display \"hello\"\ndisplay \"twice!\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", NewText("hello"),
								},
							},
							&Sentence{
								Words: []interface{}{
									"display", NewText("twice!"),
								},
							},
						},
					},
				},
			},
		},
		"Declare1": {
			bento: "start: declare some-variable is text",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "some-variable",
								Type:       "text",
								LocalScope: true,
							},
						},
					},
				},
			},
		},
		"Declare2": {
			bento: "start: declare foo is text\ndisplay foo",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", VariableReference("foo"),
								},
							},
						},
					},
				},
			},
		},
		"Function1": {
			bento: "start:  display \"hi\"\ndo something:\ndisplay \"ok\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", NewText("hi"),
								},
							},
						},
					},
					"do something": {
						Definition: &Sentence{Words: []interface{}{"do", "something"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", NewText("ok"),
								},
							},
						},
					},
				},
			},
		},
		"Function2": {
			bento: "start:do something\ndo something:\ndisplay \"ok\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"do", "something",
								},
							},
						},
					},
					"do something": {
						Definition: &Sentence{Words: []interface{}{"do", "something"}},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", NewText("ok"),
								},
							},
						},
					},
				},
			},
		},
		"FunctionWithArgument": {
			bento: "greet persons-name now (persons-name is text):",
			expected: &Program{
				Functions: map[string]*Function{
					"greet ? now": {
						Definition: &Sentence{Words: []interface{}{"greet", VariableReference("persons-name"), "now"}},
						Variables: []*VariableDefinition{
							{
								Name:       "persons-name",
								Type:       "text",
								LocalScope: false,
							},
						},
					},
				},
			},
		},
		"CallWithArgument": {
			bento: "greet persons-name now (persons-name is text):\ndisplay persons-name",
			expected: &Program{
				Functions: map[string]*Function{
					"greet ? now": {
						Definition: &Sentence{Words: []interface{}{"greet", VariableReference("persons-name"), "now"}},
						Variables: []*VariableDefinition{
							{
								Name:       "persons-name",
								Type:       "text",
								LocalScope: false,
							},
						},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", VariableReference("persons-name"),
								},
							},
						},
					},
				},
			},
		},
		"FunctionWithArguments": {
			bento: "say greeting to persons-name (persons-name is text, greeting is text):",
			expected: &Program{
				Functions: map[string]*Function{
					"say ? to ?": {
						Definition: &Sentence{Words: []interface{}{
							"say",
							VariableReference("greeting"),
							"to",
							VariableReference("persons-name"),
						}},
						Variables: []*VariableDefinition{
							{
								Name:       "greeting",
								Type:       "text",
								LocalScope: false,
							},
							{
								Name:       "persons-name",
								Type:       "text",
								LocalScope: false,
							},
						},
					},
				},
			},
		},
		"DeclareNumber": {
			bento: "start: declare foo is number",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "number",
								LocalScope: true,
								Precision:  6,
							},
						},
					},
				},
			},
		},
		"SetNegativeNumber": {
			bento: "start: declare foo is number\nset foo to -1.23",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "number",
								LocalScope: true,
								Precision:  6,
							},
						},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"set", VariableReference("foo"), "to", NewNumber("-1.23", 6),
								},
							},
						},
					},
				},
			},
		},
		"InlineIf": {
			bento: "start: declare foo is text\nif foo = \"qux\", quux 1.234\ncorge",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&If{
								Condition: &Condition{
									Left:     VariableReference("foo"),
									Right:    NewText("qux"),
									Operator: OperatorEqual,
								},
								True: &Sentence{
									Words: []interface{}{
										"quux", NewNumber("1.234", 6),
									},
								},
							},
							&Sentence{
								Words: []interface{}{
									"corge",
								},
							},
						},
					},
				},
			},
		},
		"InlineIfElse": {
			bento: "start: declare foo is text\nif foo = \"qux\", quux 1.234, otherwise corge\ndisplay",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&If{
								Condition: &Condition{
									Left:     VariableReference("foo"),
									Right:    NewText("qux"),
									Operator: OperatorEqual,
								},
								True: &Sentence{
									Words: []interface{}{
										"quux", NewNumber("1.234", 6),
									},
								},
								False: &Sentence{
									Words: []interface{}{
										"corge",
									},
								},
							},
							&Sentence{
								Words: []interface{}{
									"display",
								},
							},
						},
					},
				},
			},
		},
		"InlineUnless": {
			bento: "start: declare foo is text\nunless foo = \"qux\", quux 1.234\ncorge",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&If{
								Unless: true,
								Condition: &Condition{
									Left:     VariableReference("foo"),
									Right:    NewText("qux"),
									Operator: OperatorEqual,
								},
								True: &Sentence{
									Words: []interface{}{
										"quux", NewNumber("1.234", 6),
									},
								},
							},
							&Sentence{
								Words: []interface{}{
									"corge",
								},
							},
						},
					},
				},
			},
		},
		"InlineUnlessElse": {
			bento: "start: declare foo is text\nunless foo = \"qux\", quux 1.234, otherwise corge\ndisplay",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&If{
								Unless: true,
								Condition: &Condition{
									Left:     VariableReference("foo"),
									Right:    NewText("qux"),
									Operator: OperatorEqual,
								},
								True: &Sentence{
									Words: []interface{}{
										"quux", NewNumber("1.234", 6),
									},
								},
								False: &Sentence{
									Words: []interface{}{
										"corge",
									},
								},
							},
							&Sentence{
								Words: []interface{}{
									"display",
								},
							},
						},
					},
				},
			},
		},
		"InlineWhile": {
			bento: "start: declare foo is text\nwhile foo = \"qux\", quux 1.234\ncorge",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&While{
								Condition: &Condition{
									Left:     VariableReference("foo"),
									Right:    NewText("qux"),
									Operator: OperatorEqual,
								},
								True: &Sentence{
									Words: []interface{}{
										"quux", NewNumber("1.234", 6),
									},
								},
							},
							&Sentence{
								Words: []interface{}{
									"corge",
								},
							},
						},
					},
				},
			},
		},
		"InlineUntil": {
			bento: "start: declare foo is text\nuntil foo = \"qux\", quux 1.234\ncorge",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Statements: []Statement{
							&While{
								Until: true,
								Condition: &Condition{
									Left:     VariableReference("foo"),
									Right:    NewText("qux"),
									Operator: OperatorEqual,
								},
								True: &Sentence{
									Words: []interface{}{
										"quux", NewNumber("1.234", 6),
									},
								},
							},
							&Sentence{
								Words: []interface{}{
									"corge",
								},
							},
						},
					},
				},
			},
		},
		"NumberWithPrecision": {
			bento: "start: declare foo is number with 2 decimal places\ndisplay foo",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "number",
								LocalScope: true,
								Precision:  2,
							},
						},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", VariableReference("foo"),
								},
							},
						},
					},
				},
			},
		},
		"NumberWithPrecision1": {
			bento: "start: declare foo is number with 1 decimal place\ndisplay foo",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Words: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "number",
								LocalScope: true,
								Precision:  1,
							},
						},
						Statements: []Statement{
							&Sentence{
								Words: []interface{}{
									"display", VariableReference("foo"),
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			parser := NewParser(strings.NewReader(test.bento))
			actual, err := parser.Parse()
			require.NoError(t, err)

			diff := cmp.Diff(test.expected, actual,
				cmpopts.AcyclicTransformer("NumberToString",
					func(number *Number) string {
						return number.String()
					}))

			assert.Empty(t, diff)
		})
	}
}
